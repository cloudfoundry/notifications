package stack

import (
    "log"
    "net/http"
    "os"
)

type Logging struct {
    logger *log.Logger
}

func NewLogging(logger *log.Logger) Logging {
    return Logging{
        logger: logger,
    }
}

func (ware Logging) ServeHTTP(w http.ResponseWriter, req *http.Request, context Context) bool {
    if os.Getenv("HTTP_LOGGING_ENABLED") == "true" {
        ware.logger.Printf("%s %s\n", req.Method, req.URL.String())
    }
    return true
}
