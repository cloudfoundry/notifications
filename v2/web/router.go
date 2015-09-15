package web

import (
	"database/sql"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/uaa"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/cloudfoundry-incubator/notifications/v2/models"
	"github.com/cloudfoundry-incubator/notifications/v2/queue"
	"github.com/cloudfoundry-incubator/notifications/v2/util"
	"github.com/cloudfoundry-incubator/notifications/v2/web/campaigns"
	"github.com/cloudfoundry-incubator/notifications/v2/web/campaigntypes"
	"github.com/cloudfoundry-incubator/notifications/v2/web/info"
	"github.com/cloudfoundry-incubator/notifications/v2/web/middleware"
	"github.com/cloudfoundry-incubator/notifications/v2/web/senders"
	"github.com/cloudfoundry-incubator/notifications/v2/web/templates"
	"github.com/cloudfoundry-incubator/notifications/v2/web/unsubscribers"
	"github.com/gorilla/mux"
	"github.com/nu7hatch/gouuid"
	"github.com/pivotal-cf-experimental/rainmaker"
	"github.com/pivotal-cf-experimental/warrant"
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
	SkipVerifySSL    bool
	SQLDB            *sql.DB
	Logger           lager.Logger
	Queue            queue.Enqueuer

	UAAPublicKey    string
	UAAHost         string
	UAAClientID     string
	UAAClientSecret string
	CCHost          string
}

func NewRouter(mx muxer, config Config) http.Handler {
	requestCounter := middleware.NewRequestCounter(mx.GetRouter(), metrics.DefaultLogger)
	logging := middleware.NewRequestLogging(config.Logger)
	notificationsWriteAuthenticator := middleware.NewAuthenticator(config.UAAPublicKey, "notifications.write")
	notificationsPreferencesAdminAuthenticator := middleware.NewAuthenticator(config.UAAPublicKey, "notification_preferences.admin")
	databaseAllocator := middleware.NewDatabaseAllocator(config.SQLDB, config.DBLoggingEnabled)

	warrantConfig := warrant.Config{
		Host:          config.UAAHost,
		SkipVerifySSL: config.SkipVerifySSL,
	}
	warrantUsersService := warrant.NewUsersService(warrantConfig)
	warrantClientsService := warrant.NewClientsService(warrantConfig)

	rainmakerConfig := rainmaker.Config{
		Host:          config.CCHost,
		SkipVerifySSL: config.SkipVerifySSL,
	}
	rainmakerSpacesService := rainmaker.NewSpacesService(rainmakerConfig)
	rainmakerOrganizationsService := rainmaker.NewOrganizationsService(rainmakerConfig)

	userFinder := uaa.NewUserFinder(config.UAAClientID, config.UAAClientSecret, warrantUsersService, warrantClientsService)
	spaceFinder := cf.NewSpaceFinder(config.UAAClientID, config.UAAClientSecret, warrantClientsService, rainmakerSpacesService)
	orgFinder := cf.NewOrgFinder(config.UAAClientID, config.UAAClientSecret, warrantClientsService, rainmakerOrganizationsService)

	campaignEnqueuer := queue.NewCampaignEnqueuer(config.Queue)

	sendersRepository := models.NewSendersRepository(uuid.NewV4)
	campaignTypesRepository := models.NewCampaignTypesRepository(uuid.NewV4)
	templatesRepository := models.NewTemplatesRepository(uuid.NewV4)
	campaignsRepository := models.NewCampaignsRepository(uuid.NewV4)
	messagesRepository := models.NewMessagesRepository(util.NewClock())
	unsubscribersRepository := models.NewUnsubscribersRepository(uuid.NewV4)

	sendersCollection := collections.NewSendersCollection(sendersRepository, campaignTypesRepository)
	templatesCollection := collections.NewTemplatesCollection(templatesRepository)
	campaignTypesCollection := collections.NewCampaignTypesCollection(campaignTypesRepository, sendersRepository, templatesRepository)
	campaignsCollection := collections.NewCampaignsCollection(campaignEnqueuer, campaignsRepository, campaignTypesRepository, templatesRepository, sendersRepository, userFinder, spaceFinder, orgFinder)
	campaignStatusesCollection := collections.NewCampaignStatusesCollection(campaignsRepository, sendersRepository, messagesRepository)
	unsubscribersCollection := collections.NewUnsubscribersCollection(unsubscribersRepository, campaignTypesRepository, userFinder)

	info.Routes{
		RequestCounter: requestCounter,
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

	campaigns.Routes{
		Clock:                      util.NewClock(),
		RequestLogging:             logging,
		Authenticator:              notificationsWriteAuthenticator,
		DatabaseAllocator:          databaseAllocator,
		CampaignsCollection:        campaignsCollection,
		CampaignStatusesCollection: campaignStatusesCollection,
	}.Register(mx)

	unsubscribers.Routes{
		RequestLogging:          logging,
		Authenticator:           notificationsPreferencesAdminAuthenticator,
		DatabaseAllocator:       databaseAllocator,
		UnsubscribersCollection: unsubscribersCollection,
	}.Register(mx)

	return mx
}
