package fakes

import (
	"database/sql"

	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/cloudfoundry-incubator/notifications/web/webutil"
)

type Mother struct{}

func NewMother() Mother {
	return Mother{}
}

func (mother Mother) Registrar() services.Registrar {
	return services.Registrar{}
}

func (mother Mother) EmailStrategy() services.EmailStrategy {
	return services.EmailStrategy{}
}

func (mother Mother) UserStrategy() services.UserStrategy {
	return services.UserStrategy{}
}

func (mother Mother) SpaceStrategy() services.SpaceStrategy {
	return services.SpaceStrategy{}
}

func (mother Mother) OrganizationStrategy() services.OrganizationStrategy {
	return services.OrganizationStrategy{}
}

func (mother Mother) EveryoneStrategy() services.EveryoneStrategy {
	return services.EveryoneStrategy{}
}

func (mother Mother) UAAScopeStrategy() services.UAAScopeStrategy {
	return services.UAAScopeStrategy{}
}

func (mother Mother) NotificationsFinder() services.NotificationsFinder {
	return services.NotificationsFinder{}
}

func (mother Mother) NotificationsUpdater() services.NotificationsUpdater {
	return services.NotificationsUpdater{}
}

func (mother Mother) PreferencesFinder() *services.PreferencesFinder {
	return &services.PreferencesFinder{}
}

func (mother Mother) PreferenceUpdater() services.PreferenceUpdater {
	return services.PreferenceUpdater{}
}

func (mother Mother) MessageFinder() services.MessageFinder {
	return services.MessageFinder{}
}

func (mother Mother) TemplateServiceObjects() (services.TemplateCreator, services.TemplateFinder,
	services.TemplateUpdater, services.TemplateDeleter, services.TemplateLister,
	services.TemplateAssigner, services.TemplateAssociationLister) {
	return services.TemplateCreator{}, services.TemplateFinder{}, services.TemplateUpdater{}, services.TemplateDeleter{}, services.TemplateLister{}, services.TemplateAssigner{}, services.TemplateAssociationLister{}
}

func (mother Mother) Database() db.DatabaseInterface {
	return NewDatabase()
}

func (mother Mother) SQLDatabase() *sql.DB {
	return &sql.DB{}
}

func (mother Mother) Logging() middleware.RequestLogging {
	return middleware.RequestLogging{}
}

func (mother Mother) ErrorWriter() webutil.ErrorWriter {
	return webutil.ErrorWriter{}
}

func (mother Mother) Authenticator(scopes ...string) middleware.Authenticator {
	return middleware.Authenticator{
		Scopes: scopes,
	}
}

func (mother Mother) CORS() middleware.CORS {
	return middleware.CORS{}
}
