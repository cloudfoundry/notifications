package mail

import "time"

func (c *Client) ConnectTimeout() time.Duration {
	return c.config.ConnectTimeout
}
