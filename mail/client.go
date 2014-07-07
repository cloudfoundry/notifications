package mail

import (
    "crypto/tls"
    "fmt"
    "net"
    "net/smtp"

    "github.com/cloudfoundry-incubator/notifications/config"
)

type Client struct {
    Host     string
    Port     string
    user     string
    pass     string
    client   *smtp.Client
    Insecure bool
}

type ClientInterface interface {
    Connect() error
    Send(Message) error
}

func NewClient(user, pass, url string) (Client, error) {
    client := Client{
        user: user,
        pass: pass,
    }
    host, port, err := net.SplitHostPort(url)
    if err != nil {
        return client, err
    }
    client.Host = host
    client.Port = port
    return client, nil
}

func (c *Client) Connect() error {
    if c.client != nil {
        return nil
    }

    client, err := smtp.Dial(net.JoinHostPort(c.Host, c.Port))
    if err != nil {
        return err
    }

    c.client = client
    return nil
}

func (c *Client) Send(msg Message) error {
    err := c.Connect()
    if err != nil {
        return err
    }

    err = c.client.Hello("localhost")
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

    err = c.client.Quit()
    if err != nil {
        return err
    }

    c.client = nil

    return nil
}

func (c *Client) StartTLS() error {
    env := config.NewEnvironment()
    if env.SMTPTLS {
        if ok, _ := c.client.Extension("STARTTLS"); ok {
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
    if ok, _ := c.client.Extension("AUTH"); ok {
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
