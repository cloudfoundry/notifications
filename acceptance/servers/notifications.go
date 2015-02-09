package servers

import (
	"net/http"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/cloudfoundry-incubator/notifications/application"
)

type Notifications struct {
	cmd *exec.Cmd
	env application.Environment
}

func NewNotifications() Notifications {
	return Notifications{
		env: application.NewEnvironment(),
	}
}

func (s Notifications) Compile() {
	path, err := exec.LookPath("go")
	if err != nil {
		panic(err)
	}

	cmd := exec.Cmd{
		Path:   path,
		Args:   []string{"go", "build", "-o", "bin/notifications", "main.go"},
		Dir:    s.env.RootPath,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}

func (s Notifications) Destroy() {
	err := os.Remove(path.Join(s.env.RootPath, "bin", "notifications"))
	if err != nil {
		panic(err)
	}
}

func (s *Notifications) Boot() {
	cmd := exec.Cmd{
		Path: path.Join(s.env.RootPath, "bin", "notifications"),
		Dir:  s.env.RootPath,
	}
	if os.Getenv("TRACE") != "" {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	s.cmd = &cmd

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

func (s *Notifications) Restart() {
	s.Close()
	s.Boot()
}

func (s Notifications) URL() string {
	return "http://localhost:" + s.env.Port
}
