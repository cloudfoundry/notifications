package notify

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/cloudfoundry-incubator/notifications/v1/web/webutil"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ryanmoran/stack"
)

type clientAndKindFinder interface {
	ClientAndKind(database services.DatabaseInterface, clientID, kindID string) (models.Client, models.Kind, error)
}

type registrar interface {
	Register(services.ConnectionInterface, models.Client, []models.Kind) error
	Prune(services.ConnectionInterface, models.Client, []models.Kind) error
}

type Notify struct {
	finder    clientAndKindFinder
	registrar registrar
}

func NewNotify(finder clientAndKindFinder, registrar registrar) Notify {
	return Notify{
		finder:    finder,
		registrar: registrar,
	}
}

type ValidatorInterface interface {
	Validate(*NotifyParams) bool
}

func (h Notify) Execute(connection ConnectionInterface, req *http.Request, context stack.Context,
	guid string, strategy Dispatcher, validator ValidatorInterface, vcapRequestID string) ([]byte, error) {

	parameters, err := NewNotifyParams(req.Body)
	if err != nil {
		return []byte{}, err
	}

	if !validator.Validate(&parameters) {
		return []byte{}, webutil.ValidationError{Err: errors.New(strings.Join(parameters.Errors, ","))}
	}

	requestReceivedTime, ok := context.Get(RequestReceivedTime).(time.Time)
	if !ok {
		panic("programmer error: missing RequestReceivedTime in http context")
	}
	token := context.Get("token").(*jwt.Token) // TODO: (rm) get rid of the context object, just pass in the token
	claims := token.Claims.(jwt.MapClaims)
	clientID := claims["client_id"].(string)

	tokenIssuerURL, err := url.Parse(claims["iss"].(string))
	if err != nil {
		return []byte{}, errors.New("Token issuer URL invalid")
	}
	uaaHost := tokenIssuerURL.Scheme + "://" + tokenIssuerURL.Host

	client, kind, err := h.finder.ClientAndKind(context.Get("database").(DatabaseInterface), clientID, parameters.KindID)
	if err != nil {
		return []byte{}, err
	}

	if kind.Critical && !h.hasCriticalNotificationsWriteScope(claims["scope"]) {
		return []byte{}, webutil.NewCriticalNotificationError(kind.ID)
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
