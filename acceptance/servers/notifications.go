package servers

import (
	"net/http"
	"os/exec"
	"time"

	"github.com/cloudfoundry-incubator/notifications/application"
)

type Notifications struct {
	cmd *exec.Cmd
	env application.Environment
}

func NewNotifications() Notifications {
	env := application.NewEnvironment()
	cmd := exec.Cmd{
		Path: env.RootPath + "/bin/notifications",
		Dir:  env.RootPath,
		//Stdout: os.Stdout, // Uncomment to get server output for debugging
		//Stderr: os.Stderr,
	}

	return Notifications{
		cmd: &cmd,
		env: env,
	}
}

func (s Notifications) Boot() {
	err := s.cmd.Start()
	if err != nil {
		panic(err)
	}
	s.Ping()
}

func (s Notifications) Ping() {
	timer := time.After(0 * time.Second)
	timeout := time.After(10 * time.Second)
	for {
		select {
		case <-timeout:
			panic("Failed to boot!")
		case <-timer:
			_, err := http.Get("http://localhost:" + s.env.Port + "/info")
			if err == nil {
				return
			}

			timer = time.After(1 * time.Second)
		}
	}
}

func (s Notifications) Close() {
	err := s.cmd.Process.Kill()
	if err != nil {
		panic(err)
	}
}

func (s Notifications) RootPath() string {
	return "http://localhost:" + s.env.Port
}

func (s Notifications) SpacesPath(space string) string {
	return s.RootPath() + "/spaces/" + space
}

func (s Notifications) OrganizationsPath(organization string) string {
	return s.RootPath() + "/organizations/" + organization
}

func (s Notifications) EveryonePath() string {
	return s.RootPath() + "/everyone"
}

func (s Notifications) ScopesPath(scope string) string {
	return s.RootPath() + "/uaa_scopes/" + scope
}

func (s Notifications) UsersPath(user string) string {
	return s.RootPath() + "/users/" + user
}

func (s Notifications) EmailPath() string {
	return s.RootPath() + "/emails"
}

func (s Notifications) NotificationsPath() string {
	return s.RootPath() + "/notifications"
}

func (s Notifications) RegistrationPath() string {
	return s.RootPath() + "/registration"
}

func (s Notifications) UserPreferencesPath() string {
	return s.RootPath() + "/user_preferences"
}

func (s Notifications) SpecificUserPreferencesPath(userGUID string) string {
	return s.RootPath() + "/user_preferences/" + userGUID
}

func (s Notifications) DefaultTemplatePath() string {
	return s.RootPath() + "/default_template"
}

func (s Notifications) TemplatesBasePath() string {
	return s.RootPath() + "/templates"
}

func (s Notifications) TemplatePath(templateID string) string {
	return s.RootPath() + "/templates/" + templateID
}

func (s Notifications) DeprecatedTemplatePath(templateName string) string {
	return s.RootPath() + "/deprecated_templates/" + templateName
}

func (s Notifications) ClientsTemplatePath(clientID string) string {
	return s.RootPath() + "/clients/" + clientID + "/template"
}

func (s Notifications) ClientsNotificationsTemplatePath(clientID, notificationID string) string {
	return s.RootPath() + "/clients/" + clientID + "/notifications/" + notificationID + "/template"
}
