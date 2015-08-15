package web

import (
	"database/sql"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/cloudfoundry-incubator/notifications/v2/models"
	"github.com/cloudfoundry-incubator/notifications/v2/web/campaigntypes"
	"github.com/cloudfoundry-incubator/notifications/v2/web/info"
	"github.com/cloudfoundry-incubator/notifications/v2/web/middleware"
	"github.com/cloudfoundry-incubator/notifications/v2/web/senders"
	"github.com/cloudfoundry-incubator/notifications/v2/web/templates"
	"github.com/gorilla/mux"
	"github.com/nu7hatch/gouuid"
	"github.com/pivotal-golang/lager"
	"github.com/ryanmoran/stack"
)

type muxer interface {
	Handle(method, path string, handler stack.Handler, middleware ...stack.Middleware)
	GetRouter() *mux.Router
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

type mother interface {
	SQLDatabase()
}

type Config struct {
	DBLoggingEnabled bool
	SQLDB            *sql.DB
	Logger           lager.Logger
	UAAPublicKey     string
}

func NewRouter(mx muxer, config Config) http.Handler {
	logging := middleware.NewRequestLogging(config.Logger)
	notificationsWriteAuthenticator := middleware.NewAuthenticator(config.UAAPublicKey, "notifications.write")
	databaseAllocator := middleware.NewDatabaseAllocator(config.SQLDB, config.DBLoggingEnabled)

	sendersRepository := models.NewSendersRepository(uuid.NewV4)
	campaignTypesRepository := models.NewCampaignTypesRepository(uuid.NewV4)
	templatesRepository := models.NewTemplatesRepository(uuid.NewV4)

	sendersCollection := collections.NewSendersCollection(sendersRepository)
	campaignTypesCollection := collections.NewCampaignTypesCollection(campaignTypesRepository, sendersRepository)
	templatesCollection := collections.NewTemplatesCollection(templatesRepository)

	info.Routes{
		RequestLogging: logging,
	}.Register(mx)

	senders.Routes{
		RequestLogging:    logging,
		Authenticator:     notificationsWriteAuthenticator,
		DatabaseAllocator: databaseAllocator,
		SendersCollection: sendersCollection,
	}.Register(mx)

	campaigntypes.Routes{
		RequestLogging:          logging,
		Authenticator:           notificationsWriteAuthenticator,
		DatabaseAllocator:       databaseAllocator,
		CampaignTypesCollection: campaignTypesCollection,
	}.Register(mx)

	templates.Routes{
		RequestLogging:      logging,
		Authenticator:       notificationsWriteAuthenticator,
		DatabaseAllocator:   databaseAllocator,
		TemplatesCollection: templatesCollection,
	}.Register(mx)

	return mx
}
