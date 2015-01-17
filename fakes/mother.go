package fakes

import (
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal/strategies"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/cloudfoundry-incubator/notifications/web/services"
	"github.com/ryanmoran/stack"
)

type Mother struct{}

func NewMother() Mother {
	return Mother{}
}

func (mother Mother) Registrar() services.Registrar {
	return services.Registrar{}
}

func (mother Mother) EmailStrategy() strategies.EmailStrategy {
	return strategies.EmailStrategy{}
}

func (mother Mother) UserStrategy() strategies.UserStrategy {
	return strategies.UserStrategy{}
}

func (mother Mother) SpaceStrategy() strategies.SpaceStrategy {
	return strategies.SpaceStrategy{}
}

func (mother Mother) OrganizationStrategy() strategies.OrganizationStrategy {
	return strategies.OrganizationStrategy{}
}

func (mother Mother) EveryoneStrategy() strategies.EveryoneStrategy {
	return strategies.EveryoneStrategy{}
}

func (mother Mother) UAAScopeStrategy() strategies.UAAScopeStrategy {
	return strategies.UAAScopeStrategy{}
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

func (mother Mother) Database() models.DatabaseInterface {
	return NewDatabase()
}

func (mother Mother) Logging() stack.Middleware {
	return stack.Logging{}
}

func (mother Mother) ErrorWriter() handlers.ErrorWriter {
	return handlers.ErrorWriter{}
}

func (mother Mother) Authenticator(scopes ...string) middleware.Authenticator {
	return middleware.Authenticator{
		Scopes: scopes,
	}
}

func (mother Mother) CORS() middleware.CORS {
	return middleware.CORS{}
}
