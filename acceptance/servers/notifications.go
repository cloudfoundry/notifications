package servers

import (
    "net/http"
    "os/exec"
    "time"

    "github.com/cloudfoundry-incubator/notifications/config"
)

type Notifications struct {
    cmd *exec.Cmd
    env config.Environment
}

func NewNotifications() Notifications {
    env := config.NewEnvironment()
    cmd := exec.Cmd{
        Path: env.RootPath + "/bin/notifications",
        Dir:  env.RootPath,
        //Stdout: os.Stdout, // Uncomment to get server output for debugging
        //Stderr: os.Stderr,
    }

    return Notifications{
        cmd: &cmd,
        env: config.NewEnvironment(),
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

func (s Notifications) SpacesPath(space string) string {
    return "http://localhost:" + s.env.Port + "/spaces/" + space
}

func (s Notifications) UsersPath(user string) string {
    return "http://localhost:" + s.env.Port + "/users/" + user
}

func (s Notifications) EmailPath() string {
    return "http://localhost:" + s.env.Port + "/emails"
}

func (s Notifications) RegistrationPath() string {
    return "http://localhost:" + s.env.Port + "/registration"
}

func (s Notifications) UserPreferencesPath() string {
    return "http://localhost:" + s.env.Port + "/user_preferences"
}

func (s Notifications) SpecificUserPreferencesPath(userGUID string) string {
    return "http://localhost:" + s.env.Port + "/user_preferences/" + userGUID
}

func (s Notifications) SpaceTemplatePath() string {
    return "http://localhost:" + s.env.Port + "/templates/space_body"
}

func (s Notifications) UserTemplatePath() string {
    return "http://localhost:" + s.env.Port + "/templates/user_body"
}

func (s Notifications) UserTemplateForClientPath(clientID string) string {
    return "http://localhost:" + s.env.Port + "/templates/" + clientID + ".user_body"
}

func (s Notifications) UserTemplateForClientAndKindPath(clientID, kindID string) string {
    return "http://localhost:" + s.env.Port + "/templates/" + clientID + "." + kindID + ".user_body"
}
