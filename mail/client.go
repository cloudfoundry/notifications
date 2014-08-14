package mail

import (
    "crypto/tls"
    "errors"
    "fmt"
    "log"
    "net"
    "net/smtp"
    "time"

    "github.com/cloudfoundry-incubator/notifications/config"
)

type Client struct {
    Host           string
    Port           string
    user           string
    pass           string
    client         *smtp.Client
    Insecure       bool
    ConnectTimeout time.Duration
}

type ClientInterface interface {
    Connect() error
    Send(Message) error
}

type connection struct {
    client *smtp.Client
    err    error
}

func NewClient(user, pass, url string) (*Client, error) {
    client := &Client{
        user: user,
        pass: pass,
    }

    host, port, err := net.SplitHostPort(url)
    if err != nil {
        return client, err
    }

    client.Host = host
    client.Port = port
    client.ConnectTimeout = 15 * time.Second

    return client, nil
}

func (c *Client) Connect() error {

    env := config.NewEnvironment()
    if env.TestMode {
        return nil
    }

    if c.client != nil {
        return nil
    }

    select {
    case connection := <-c.connect():
        if connection.err != nil {
            return connection.err
        }

        c.client = connection.client
    case <-time.After(c.ConnectTimeout):
        return errors.New("dial tcp: i/o timeout")
    }

    return nil
}

func (c *Client) connect() chan connection {
    channel := make(chan connection)

    go func() {
        client, err := smtp.Dial(net.JoinHostPort(c.Host, c.Port))
        channel <- connection{
            client: client,
            err:    err,
        }
    }()

    return channel
}

func (c *Client) Send(msg Message) error {
    env := config.NewEnvironment()
    if env.TestMode {
        log.Println("TEST_MODE is true, emails not being sent")
        return nil
    }

    err := c.Connect()
    if err != nil {
        return err
    }

    err = c.Hello()
    if err != nil {
        return err
    }

    err = c.StartTLS()
    if err != nil {
        return err
    }

    err = c.Auth()
    if err != nil {
        return err
    }

    err = c.client.Mail(msg.From)
    if err != nil {
        return err
    }

    err = c.client.Rcpt(msg.To)
    if err != nil {
        return err
    }

    err = c.Data(msg)
    if err != nil {
        return err
    }

    err = c.Quit()
    if err != nil {
        return err
    }

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
    env := config.NewEnvironment()
    if env.SMTPTLS {
        if ok, _ := c.Extension("STARTTLS"); ok {
            err := c.client.StartTLS(&tls.Config{
                ServerName:         c.Host,
                InsecureSkipVerify: c.Insecure,
            })
            if err != nil {
                return err
            }
        }
    }

    return nil
}

func (c *Client) Auth() error {
    if ok, _ := c.Extension("AUTH"); ok {
        err := c.client.Auth(smtp.PlainAuth("", c.user, c.pass, c.Host))
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

    _, err = fmt.Fprintf(wc, msg.Data())
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
