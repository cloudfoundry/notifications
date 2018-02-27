package common

import (
	"fmt"
	"net/http"
)

const (
	JSONErrorMessageTemplate = `{"error_description":%q,"error":%q}`
	HTMLErrorMessageTemplate = `<html><head><title>Apache Tomcat/7.0.61 - Error report</title><style><!--H1 {font-family:Tahoma,Arial,sans-serif;color:white;background-color:#525D76;font-size:22px;} H2 {font-family:Tahoma,Arial,sans-serif;color:white;background-color:#525D76;font-size:16px;} H3 {font-family:Tahoma,Arial,sans-serif;color:white;background-color:#525D76;font-size:14px;} BODY {font-family:Tahoma,Arial,sans-serif;color:black;background-color:white;} B {font-family:Tahoma,Arial,sans-serif;color:white;background-color:#525D76;} P {font-family:Tahoma,Arial,sans-serif;background:white;color:black;font-size:12px;}A {color : black;}A.name {color : black;}HR {color : #525D76;}--></style> </head><body><h1>HTTP Status %d - </h1><HR size="1" noshade="noshade"><p><b>type</b> Status report</p><p><b>message</b> <u>%s</u></p><p><b>description</b> <u>%s</u></p><HR size="1" noshade="noshade"><h3>Apache Tomcat/7.0.61</h3></body></html>`
)

func JSONError(w http.ResponseWriter, status int, message, errorType string) {
	output := fmt.Sprintf(JSONErrorMessageTemplate, message, errorType)

	w.WriteHeader(status)
	w.Write([]byte(output))
}

func HTMLError(w http.ResponseWriter, status int, message, errorType string) {
	output := fmt.Sprintf(HTMLErrorMessageTemplate, status, message, errorType)

	w.WriteHeader(status)
	w.Write([]byte(output))
}
