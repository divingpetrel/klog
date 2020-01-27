// +build go1.12
// +build stackdriver

package logr

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"
	"time"
)

// ServiceContext is required by the Stackdriver Error format.
type ServiceContext struct {
	Service string `json:"service,omitempty"`
	Version string `json:"version,omitempty"`
}

// ReportLocation is required by the Stackdriver Error format
type ReportLocation struct {
	FilePath     string `json:"filePath"`
	FunctionName string `json:"functionName"`
	LineNumber   int    `json:"lineNumber"`
}

// Context is required by the Stackdriver Error format
type Context struct {
	ReportLocation *ReportLocation `json:"reportLocation,omitempty"`
}

// cloud.google.com/error-reporting/docs/formatting-error-messages#formatting_requirements
type Stackdriver struct {
	Severity       string          `json:"severity"`
	EventTime      string          `json:"eventTime"`
	Message        string          `json:"message"`
	ServiceContext *ServiceContext `json:"serviceContext,omitempty"`
	Context        *Context        `json:"context,omitempty"`
	Stacktrace     string          `json:"stacktrace,omitempty"`
}

func New(daemon, version string) (*Stackdriver, error) {
	return &Stackdriver{"", "", "", &ServiceContext{Service: daemon, Version: version}, &Context{ReportLocation: &ReportLocation{}}, ""}, nil
}

// INFO logs a message with INFO severity level
func (s *Stackdriver) Infof(format string, args ...interface{}) {
	s.basicFormat("INFO", fmt.Sprintf(format, args...))
}

// DEBUG logs a message with DEBUG severity level
func (s *Stackdriver) Debugf(format string, args ...interface{}) {
	s.basicFormat("DEBUG", fmt.Sprintf(format, args...))
}

// ERROR prints out a message with ERROR severity level + picked up by stackdriver error reporting
func (s *Stackdriver) Errorf(format string, args ...interface{}) {
	s.reportFormat("ERROR", fmt.Sprintf(format, args...))
}

// CRITICAL is system failure
func (s *Stackdriver) Criticalf(format string, args ...interface{}) {
	s.reportFormat("CRITICAL", fmt.Sprintf(format, args...))
}

// Emits to stdout with correct format allowing fluentd to submit a logEntry to stackdriver, no alerting/reporting with INFO/DEBUG levels
// cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry
func (s *Stackdriver) basicFormat(severity, message string) {

	s.Severity = severity
	s.EventTime = time.Now().Format(time.RFC3339)
	s.Message = message

	payloadJson, ok := json.Marshal(&s)
	if ok != nil {
		fmt.Printf("logging ERROR: cannot marshal payload: %s", ok.Error())
		payloadJson = []byte{}
	}

	fmt.Println(string(payloadJson))
}

// ....  Error and Critical levels will result in stackdriver error reporting being intitiated
// As error/critical logs are relatively rare we can take the performance hit of retrieving stack information etc...
func (s *Stackdriver) reportFormat(severity, message string) {

	buffer := make([]byte, 1024)
	buffer = buffer[:runtime.Stack(buffer, false)]

	// Who called me?
	// Caller reports file and line number information about function invocations on the calling goroutine's stack
	fpc, file, line, _ := runtime.Caller(2)

	// FuncForPC returns a *Func describing the function that contains the given program counter address, or else nil.
	funcName := "unknown"
	fun := runtime.FuncForPC(fpc)
	if fun != nil {
		_, funcName = filepath.Split(fun.Name())
	}

	s.Severity = severity
	s.EventTime = time.Now().Format(time.RFC3339)
	s.Message = message
	s.Context.ReportLocation.FilePath = file
	s.Context.ReportLocation.FunctionName = funcName
	s.Context.ReportLocation.LineNumber = line
	s.Stacktrace = string(buffer)

	payloadJson, ok := json.Marshal(&s)
	if ok != nil {
		fmt.Printf("logging ERROR: cannot marshal payload: %s", ok.Error())
		payloadJson = []byte{}
	}

	fmt.Println(string(payloadJson))
}
