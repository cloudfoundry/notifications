package mail

import (
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"

	"github.com/pivotal-golang/lager"
)

const (
	AuthNone AuthMechanism = iota
	AuthPlain
	AuthCRAMMD5
)

type AuthMechanism int

type Client struct {
	config Config
	client *smtp.Client
}

type Config struct {
	Host           string
	Port           string
	User           string
	Pass           string
	Secret         string
	AuthMechanism  AuthMechanism
	TestMode       bool
	SkipVerifySSL  bool
	DisableTLS     bool
	ConnectTimeout time.Duration
	LoggingEnabled bool
}

type ClientInterface interface {
	Connect(lager.Logger) error
	Send(Message, lager.Logger) error
}

type connection struct {
	client *smtp.Client
	err    error
}

func NewClient(config Config) *Client {
	client := &Client{config: config}

	if client.config.ConnectTimeout == 0 {
		client.config.ConnectTimeout = 15 * time.Second
	}

	return client
}

func (c *Client) Connect(logger lager.Logger) error {
	c.PrintLog(logger, "Connecting...")
	if c.config.TestMode {
		c.PrintLog(logger, "Test Mode enabled, not connected")
		return nil
	}

	if c.client != nil {
		c.PrintLog(logger, "Already connected.")
		return nil
	}

	select {
	case connection := <-c.connect():
		c.PrintLog(logger, "Connected")
		if connection.err != nil {
			return connection.err
		}

		c.client = connection.client
	case <-time.After(c.config.ConnectTimeout):
		c.PrintLog(logger, "Timed out after %+v", c.config.ConnectTimeout)
		return errors.New("server timeout")
	}

	return nil
}

func (c *Client) connect() chan connection {
	channel := make(chan connection)

	go func() {
		client, err := smtp.Dial(net.JoinHostPort(c.config.Host, c.config.Port))
		channel <- connection{
			client: client,
			err:    err,
		}
	}()

	return channel
}

func (c *Client) Send(msg Message, logger lager.Logger) error {
	if c.config.TestMode {
		logger.Info("TEST_MODE is enabled, emails not being sent")
		return nil
	}

	err := c.Connect(logger)
	if err != nil {
		return c.Error(logger, err)
	}

	c.PrintLog(logger, "Initiating hello...")
	err = c.Hello()
	if err != nil {
		return c.Error(logger, err)
	}
	c.PrintLog(logger, "Hello complete.")

	if !c.config.DisableTLS {
		c.PrintLog(logger, "Starting TLS...")
		err = c.StartTLS()
		if err != nil {
			return c.Error(logger, err)
		}
		c.PrintLog(logger, "TLS connection opened.")

		c.PrintLog(logger, "Starting authentication...")
		err = c.Auth(logger)
		if err != nil {
			return c.Error(logger, err)
		}
		c.PrintLog(logger, "Authenticated.")
	}

	c.PrintLog(logger, "Sending mail from: %s", msg.From)
	err = c.client.Mail(msg.From)
	if err != nil {
		return c.Error(logger, err)
	}

	c.PrintLog(logger, "Sending mail to: %s", msg.To)
	err = c.client.Rcpt(msg.To)
	if err != nil {
		return c.Error(logger, err)
	}

	c.PrintLog(logger, "Sending mail data...")
	c.PrintLog(logger, "Message Data: %s", base64.StdEncoding.EncodeToString([]byte(msg.Data())))
	err = c.Data(msg)
	if err != nil {
		return c.Error(logger, err)
	}
	c.PrintLog(logger, "Mail data sent.")

	c.PrintLog(logger, "Quitting...")
	err = c.Quit()
	if err != nil {
		return c.Error(logger, err)
	}
	c.PrintLog(logger, "Goodbye.")

	return nil
}

func (c *Client) Hello() error {
	err := c.client.Hello("localhost")
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Extension(name string) (bool, string) {
	return c.client.Extension(name)
}

func (c *Client) StartTLS() error {
	if ok, _ := c.Extension("STARTTLS"); ok {
		err := c.client.StartTLS(&tls.Config{
			ServerName:         c.config.Host,
			InsecureSkipVerify: c.config.SkipVerifySSL,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) Auth(logger lager.Logger) error {
	if ok, _ := c.Extension("AUTH"); ok {
		if mechanism := c.AuthMechanism(logger); mechanism != nil {
			err := c.client.Auth(mechanism)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Client) AuthMechanism(logger lager.Logger) smtp.Auth {
	switch c.config.AuthMechanism {
	case AuthCRAMMD5:
		c.PrintLog(logger, "Using CRAMMD5 to authenticate")
		return smtp.CRAMMD5Auth(c.config.User, c.config.Secret)
	case AuthPlain:
		c.PrintLog(logger, "Using PLAIN to authenticate")
		return smtp.PlainAuth("", c.config.User, c.config.Pass, c.config.Host)
	default:
		c.PrintLog(logger, "No authentication mechanism configured")
		return nil
	}
}

func (c *Client) Data(msg Message) error {
	wc, err := c.client.Data()
	if err != nil {
		return err
	}

	data := strings.Replace(string(msg.Data()), "%", "%%", -1)
	_, err = fmt.Fprintf(wc, data)
	if err != nil {
		return err
	}

	err = wc.Close()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Quit() error {
	err := c.client.Quit()
	c.client = nil
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Error(logger lager.Logger, err error) error {
	if c.client != nil {
		failure := c.Quit()
		if failure != nil {
			return failure
		}
	}

	logger.Error("SMTP", err)

	return err
}

func (c *Client) PrintLog(logger lager.Logger, format string, arguments ...interface{}) {
	if c.config.LoggingEnabled {
		logger.Info(fmt.Sprintf("[SMTP] "+format, arguments...))
	}
}
