package mail

import (
    "crypto/tls"
    "errors"
    "fmt"
    "log"
    "net"
    "net/smtp"
    "strings"
    "time"
)

type Client struct {
    config Config
    client *smtp.Client
    logger *log.Logger
}

type Config struct {
    Host           string
    Port           string
    User           string
    Pass           string
    TestMode       bool
    SkipVerifySSL  bool
    DisableTLS     bool
    ConnectTimeout time.Duration
    LoggingEnabled bool
}

type ClientInterface interface {
    Connect() error
    Send(Message) error
}

type connection struct {
    client *smtp.Client
    err    error
}

func NewClient(config Config, logger *log.Logger) (*Client, error) {
    client := &Client{
        config: config,
        logger: logger,
    }

    if client.config.ConnectTimeout == 0 {
        client.config.ConnectTimeout = 15 * time.Second
    }

    return client, nil
}

func (c *Client) Connect() error {
    c.Log("Connecting...")
    if c.config.TestMode {
        c.Log("Test Mode enabled, not connected")
        return nil
    }

    if c.client != nil {
        c.Log("Already connected.")
        return nil
    }

    select {
    case connection := <-c.connect():
        c.Log("Connected")
        if connection.err != nil {
            return connection.err
        }

        c.client = connection.client
    case <-time.After(c.config.ConnectTimeout):
        c.Log("Timed out after %v", c.config.ConnectTimeout)
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

func (c *Client) Send(msg Message) error {
    if c.config.TestMode {
        c.logger.Println("TEST_MODE is enabled, emails not being sent")
        return nil
    }

    err := c.Connect()
    if err != nil {
        return c.Error(err)
    }

    c.Log("Initiating hello...")
    err = c.Hello()
    if err != nil {
        return c.Error(err)
    }
    c.Log("Hello complete.")

    if !c.config.DisableTLS {
        c.Log("Starting TLS...")
        err = c.StartTLS()
        if err != nil {
            return c.Error(err)
        }
        c.Log("TLS connection opened.")

        c.Log("Starting authentication...")
        err = c.Auth()
        if err != nil {
            return c.Error(err)
        }
        c.Log("Authenticated.")
    }

    c.Log("Sending mail from: %s", msg.From)
    err = c.client.Mail(msg.From)
    if err != nil {
        return c.Error(err)
    }

    c.Log("Sending mail to: %s", msg.To)
    err = c.client.Rcpt(msg.To)
    if err != nil {
        return c.Error(err)
    }

    c.Log("Sending mail data...")
    err = c.Data(msg)
    if err != nil {
        return c.Error(err)
    }
    c.Log("Mail data sent.")

    c.Log("Quitting...")
    err = c.Quit()
    if err != nil {
        return c.Error(err)
    }
    c.Log("Goodbye.")

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

func (c *Client) Auth() error {
    if ok, _ := c.Extension("AUTH"); ok {
        err := c.client.Auth(smtp.PlainAuth("", c.config.User, c.config.Pass, c.config.Host))
        if err != nil {
            return err
        }
    }

    return nil
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
    if err != nil {
        return err
    }

    c.client = nil

    return nil
}

func (c *Client) Error(err error) error {
    c.logger.Printf("SMTP Error: %s", err.Error())
    return err
}

func (c *Client) Log(format string, arguments ...interface{}) {
    if c.config.LoggingEnabled {
        c.logger.Printf("[SMTP] "+format, arguments...)
    }
}
