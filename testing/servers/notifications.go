package servers

import (
	"database/sql"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/onsi/ginkgo/v2"
)

type Notifications struct {
	cmd  *exec.Cmd
	env  application.Environment
	port string
}

func NewNotifications() Notifications {
	env, err := application.NewEnvironment()
	if err != nil {
		panic(err)
	}

	return Notifications{
		env: env,
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
	s.port = freePort()
	environment := append(os.Environ(), fmt.Sprintf("PORT=%s", s.port))

	cmd := exec.Cmd{
		Path: path.Join(s.env.RootPath, "bin", "notifications"),
		Dir:  s.env.RootPath,
		Env:  environment,
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
	s.waitForStart()
}

func (s Notifications) waitForStart() {
	timer := time.Tick(1 * time.Second)
	timeout := time.After(10 * time.Second)
	for {
		select {
		case <-timeout:
			panic("Failed to boot!")
		case <-timer:
			resp, err := http.Get("http://localhost:" + s.port + "/info")
			if err == nil && resp.StatusCode == http.StatusOK {
				return
			}
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
	return "http://localhost:" + s.port
}

func (s Notifications) MigrateDatabase() {
	env, err := application.NewEnvironment()
	if err != nil {
		panic(err)
	}

	sqlDB := openDatabase(env)
	defer sqlDB.Close()
	database, gobbleDB := fetchDatabases(env, sqlDB)

	migrator := models.DatabaseMigrator{}
	migrator.Migrate(database.RawConnection(), env.ModelMigrationsPath)

	gobbleDB.Migrate(env.GobbleMigrationsPath)
}

func (s Notifications) ResetDatabase() {
	env, err := application.NewEnvironment()
	if err != nil {
		panic(err)
	}

	sqlDB := openDatabase(env)
	defer sqlDB.Close()
	database, gobbleDB := fetchDatabases(env, sqlDB)

	models.Setup(database)

	err = database.Connection().(*db.Connection).TruncateTables()
	if err != nil {
		panic(err)
	}

	migrator := models.DatabaseMigrator{}
	migrator.Seed(database, path.Join(env.RootPath, "templates", "default.json"))

	gobbleDB.Connection.TruncateTables()
}

func (s Notifications) WaitForJobsQueueToEmpty() error {
	env, err := application.NewEnvironment()
	if err != nil {
		return err
	}

	sqlDB := openDatabase(env)
	defer sqlDB.Close()
	_, gobbleDB := fetchDatabases(env, sqlDB)

	timer := time.After(10 * time.Second)
	for {
		select {
		case <-timer:
			return errors.New("timed out waiting for jobs queue to empty")
		default:
			count, err := gobbleDB.Connection.SelectInt("SELECT COUNT(*) FROM `jobs`")
			if err != nil {
				return err
			}
			if count == 0 {
				return nil
			}
		}
	}
}

func freePort() string {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		ginkgo.Fail(err.Error(), 1)
	}
	defer listener.Close()

	address := listener.Addr().String()
	addressParts := strings.SplitN(address, ":", 2)
	return addressParts[1]
}

func openDatabase(env application.Environment) *sql.DB {
	sqlDB, err := sql.Open("mysql", env.DatabaseURL)
	if err != nil {
		ginkgo.Fail(err.Error(), 1)
	}
	return sqlDB
}

func fetchDatabases(env application.Environment, sqlDB *sql.DB) (*db.DB, *gobble.DB) {
	database := db.NewDatabase(sqlDB, db.Config{DefaultTemplatePath: path.Join(env.RootPath, "templates", "default.json")})
	gobbleDB := gobble.NewDatabase(sqlDB)

	return database, gobbleDB
}
