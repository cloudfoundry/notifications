package mail

import (
    "fmt"
    "net"
    "net/smtp"
)

type Client struct {
    Host   string
    Port   string
    user   string
    pass   string
    client *smtp.Client
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
    client, err := smtp.Dial(fmt.Sprintf("%s:%s", c.Host, c.Port))
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

    if ok, _ := c.client.Extension("AUTH"); ok {
        err = c.client.Auth(smtp.PlainAuth("", c.user, c.pass, c.Host))
        if err != nil {
            return err
        }
    }

    err = c.client.Mail(msg.From)
    if err != nil {
        return err
    }

    err = c.client.Rcpt(msg.To)
    if err != nil {
        return err
    }

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

    err = c.client.Quit()
    if err != nil {
        return err
    }

    return nil
}
