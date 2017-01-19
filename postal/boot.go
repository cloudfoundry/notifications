package postal

import (
	"crypto/rand"
	"database/sql"
	"os"
	"time"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/mail"
	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/postal/common"
	"github.com/cloudfoundry-incubator/notifications/postal/v1"
	"github.com/cloudfoundry-incubator/notifications/postal/v2"
	"github.com/cloudfoundry-incubator/notifications/uaa"
	"github.com/cloudfoundry-incubator/notifications/util"
	v1models "github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/cloudfoundry-incubator/notifications/v1/web/notify"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/cloudfoundry-incubator/notifications/v2/horde"
	v2models "github.com/cloudfoundry-incubator/notifications/v2/models"
	"github.com/cloudfoundry-incubator/notifications/v2/queue"
	"github.com/dgrijalva/jwt-go"
	"github.com/pivotal-golang/conceal"
	"github.com/pivotal-golang/lager"
)

type mother interface {
	SQLDatabase() *sql.DB
	Database() db.DatabaseInterface
	MailClient() *mail.Client
}

type uaaTokenValidator interface {
	Parse(token string) (*jwt.Token, error)
}

type Config struct {
	UAAClientID          string
	UAAClientSecret      string
	UAATokenValidator    *uaa.TokenValidator
	UAAHost              string
	VerifySSL            bool
	InstanceIndex        int
	WorkerCount          int
	EncryptionKey        []byte
	DBLoggingEnabled     bool
	Sender               string
	Domain               string
	QueueWaitMaxDuration int
	CCHost               string
}

func Boot(mom mother, config Config) {
	uaaClient := uaa.NewZonedUAAClient(config.UAAClientID, config.UAAClientSecret, config.VerifySSL, config.UAATokenValidator)

	logger := lager.NewLogger("notifications")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))

	clock := util.NewClock()

	sqlDatabase := mom.SQLDatabase()
	database := mom.Database()

	gobbleDatabase := gobble.NewDatabase(sqlDatabase)
	gobbleQueue := gobble.NewQueue(gobbleDatabase, clock, gobble.Config{
		WaitMaxDuration: time.Duration(config.QueueWaitMaxDuration) * time.Millisecond,
	})

	cloak, err := conceal.NewCloak(config.EncryptionKey)
	if err != nil {
		panic(err)
	}

	guidGenerator := util.NewIDGenerator(rand.Reader)

	// V1
	receiptsRepo := v1models.NewReceiptsRepo()
	unsubscribesRepo := v1models.NewUnsubscribesRepo()
	globalUnsubscribesRepo := v1models.NewGlobalUnsubscribesRepo()
	messagesRepo := v1models.NewMessagesRepo(guidGenerator.Generate)
	clientsRepo := v1models.NewClientsRepo()
	kindsRepo := v1models.NewKindsRepo()
	templatesRepo := v1models.NewTemplatesRepo()
	v1TemplateLoader := v1.NewTemplatesLoader(database, clientsRepo, kindsRepo, templatesRepo)
	deliveryFailureHandler := common.NewDeliveryFailureHandler()
	messageStatusUpdater := v1.NewMessageStatusUpdater(messagesRepo)
	userLoader := common.NewUserLoader(uaaClient)
	tokenLoader := uaa.NewTokenLoader(uaaClient)
	packager := common.NewPackager(v1TemplateLoader, cloak)

	// V2
	metricsEmitter := metrics.NewEmitter(metrics.DefaultLogger)

	messagesRepository := v2models.NewMessagesRepository(clock, guidGenerator.Generate)
	gobbleInitializer := gobble.Initializer{}

	v2enqueuer := queue.NewJobEnqueuer(gobbleQueue, messagesRepository, gobbleInitializer)

	cloudController := cf.NewCloudController(config.CCHost, !config.VerifySSL)
	spaceLoader := services.NewSpaceLoader(cloudController)
	organizationLoader := services.NewOrganizationLoader(cloudController)
	findsUserIDs := services.NewFindsUserIDs(cloudController, uaaClient)

	usersAudienceGenerator := horde.NewUsers()
	emailsAudienceGenerator := horde.NewEmails()
	spacesAudienceGenerator := horde.NewSpaces(findsUserIDs, organizationLoader, spaceLoader, tokenLoader, config.UAAHost)
	orgsAudienceGenerator := horde.NewOrganizations(findsUserIDs, organizationLoader, tokenLoader, config.UAAHost)

	v2database := v2models.NewDatabase(sqlDatabase, v2models.Config{})
	v2messageStatusUpdater := v2.NewV2MessageStatusUpdater(messagesRepository)
	unsubscribersRepository := v2models.NewUnsubscribersRepository(guidGenerator.Generate)
	campaignsRepository := v2models.NewCampaignsRepository(guidGenerator.Generate, clock)
	v2templatesRepo := v2models.NewTemplatesRepository(guidGenerator.Generate)
	templatesCollection := collections.NewTemplatesCollection(v2templatesRepo)
	v2TemplateLoader := v2.NewTemplatesLoader(v2database, templatesCollection)
	v2deliveryFailureHandler := common.NewDeliveryFailureHandler()
	campaignJobProcessor := v2.NewCampaignJobProcessor(notify.EmailFormatter{}, notify.HTMLExtractor{},
		emailsAudienceGenerator, spacesAudienceGenerator, orgsAudienceGenerator, usersAudienceGenerator, v2enqueuer)

	WorkerGenerator{
		InstanceIndex: config.InstanceIndex,
		Count:         config.WorkerCount,
	}.Work(func(index int) Worker {

		mailClient := mom.MailClient()

		v1DeliveryJobProcessor := v1.NewDeliveryJobProcessor(v1.DeliveryJobProcessorConfig{
			DBTrace: config.DBLoggingEnabled,
			UAAHost: config.UAAHost,
			Sender:  config.Sender,
			Domain:  config.Domain,

			Packager:    packager,
			MailClient:  mailClient,
			Database:    database,
			TokenLoader: tokenLoader,
			UserLoader:  userLoader,

			KindsRepo:              kindsRepo,
			ReceiptsRepo:           receiptsRepo,
			UnsubscribesRepo:       unsubscribesRepo,
			GlobalUnsubscribesRepo: globalUnsubscribesRepo,
			MessageStatusUpdater:   messageStatusUpdater,
			DeliveryFailureHandler: deliveryFailureHandler,
		})

		v2mailClient := mom.MailClient()

		v2DeliveryJobProcessor := v2.NewDeliveryJobProcessor(v2mailClient, common.NewPackager(v2TemplateLoader, cloak),
			common.NewUserLoader(uaaClient), uaa.NewTokenLoader(uaaClient), v2messageStatusUpdater, v2database,
			unsubscribersRepository, campaignsRepository, config.Sender, config.Domain, config.UAAHost, metricsEmitter)

		worker := NewDeliveryWorker(v1DeliveryJobProcessor, v2DeliveryJobProcessor, DeliveryWorkerConfig{
			ID:      index,
			UAAHost: config.UAAHost,
			DBTrace: config.DBLoggingEnabled,

			Logger: logger.Session("worker", lager.Data{"worker_id": index}),
			Queue:  gobbleQueue,

			Database:               v2database,
			CampaignJobProcessor:   campaignJobProcessor,
			DeliveryFailureHandler: v2deliveryFailureHandler,
			MessageStatusUpdater:   v2messageStatusUpdater,
		})

		return &worker
	})
}
