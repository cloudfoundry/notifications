package stack

type Context interface {
    Keys() []string
    Get(string) interface{}
    Set(string, interface{})
}

type context struct {
    values map[string]interface{}
}

func NewContext() *context {
    return &context{
        values: make(map[string]interface{}),
    }
}

func (c *context) Keys() []string {
    var keys []string

    for key, _ := range c.values {
        keys = append(keys, key)
    }

    return keys
}

func (c *context) Get(key string) interface{} {
    value, _ := c.values[key]
    return value
}

func (c *context) Set(key string, value interface{}) {
    c.values[key] = value
}
