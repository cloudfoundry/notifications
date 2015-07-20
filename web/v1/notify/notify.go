package notify

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/web/webutil"
	"github.com/dgrijalva/jwt-go"
	"github.com/ryanmoran/stack"
)

type NotifyInterface interface {
	Execute(models.ConnectionInterface, *http.Request, stack.Context, string, services.StrategyInterface, ValidatorInterface, string) ([]byte, error)
}

type Notify struct {
	finder    services.NotificationsFinderInterface
	registrar services.RegistrarInterface
}

func NewNotify(finder services.NotificationsFinderInterface, registrar services.RegistrarInterface) Notify {
	return Notify{
		finder:    finder,
		registrar: registrar,
	}
}

type ValidatorInterface interface {
	Validate(*NotifyParams) bool
}

func (h Notify) Execute(connection models.ConnectionInterface, req *http.Request, context stack.Context,
	guid string, strategy services.StrategyInterface, validator ValidatorInterface, vcapRequestID string) ([]byte, error) {
	parameters, err := NewNotifyParams(req.Body)
	if err != nil {
		return []byte{}, err
	}

	if !validator.Validate(&parameters) {
		return []byte{}, webutil.ValidationError(parameters.Errors)
	}

	requestReceivedTime, ok := context.Get(RequestReceivedTime).(time.Time)
	if !ok {
		panic("programmer error: missing RequestReceivedTime in http context")
	}
	token := context.Get("token").(*jwt.Token) // TODO: (rm) get rid of the context object, just pass in the token
	clientID := token.Claims["client_id"].(string)

	tokenIssuerURL, err := url.Parse(token.Claims["iss"].(string))
	if err != nil {
		return []byte{}, errors.New("Token issuer URL invalid")
	}
	uaaHost := tokenIssuerURL.Scheme + "://" + tokenIssuerURL.Host

	client, kind, err := h.finder.ClientAndKind(context.Get("database").(models.DatabaseInterface), clientID, parameters.KindID)
	if err != nil {
		return []byte{}, err
	}

	if kind.Critical && !h.hasCriticalNotificationsWriteScope(token.Claims["scope"]) {
		return []byte{}, postal.NewCriticalNotificationError(kind.ID)
	}

	err = h.registrar.Register(connection, client, []models.Kind{kind})
	if err != nil {
		return []byte{}, err
	}

	var responses []services.Response

	responses, err = strategy.Dispatch(services.Dispatch{
		GUID:       guid,
		Connection: connection,
		Role:       parameters.Role,
		Client: services.DispatchClient{
			ID:          clientID,
			Description: client.Description,
		},
		Kind: services.DispatchKind{
			ID:          parameters.KindID,
			Description: kind.Description,
		},
		UAAHost: uaaHost,
		VCAPRequest: services.DispatchVCAPRequest{
			ID:          vcapRequestID,
			ReceiptTime: requestReceivedTime,
		},
		Message: services.DispatchMessage{
			To:      parameters.To,
			ReplyTo: parameters.ReplyTo,
			Subject: parameters.Subject,
			Text:    parameters.Text,
			HTML: services.HTML{
				BodyContent:    parameters.ParsedHTML.BodyContent,
				BodyAttributes: parameters.ParsedHTML.BodyAttributes,
				Head:           parameters.ParsedHTML.Head,
				Doctype:        parameters.ParsedHTML.Doctype,
			},
		},
	})
	if err != nil {
		return []byte{}, err
	}

	output, err := json.Marshal(responses)
	if err != nil {
		panic(err)
	}

	return output, nil
}

func (h Notify) hasCriticalNotificationsWriteScope(elements interface{}) bool {
	for _, elem := range elements.([]interface{}) {
		if elem.(string) == "critical_notifications.write" {
			return true
		}
	}
	return false
}
