// +build go1.12
// +build console

package logr

import (
	"fmt"
)

type Console struct{}

func New(daemon, version string) (*Console, error) {
	return &Console{}, nil
}

// INFO logs a message with INFO severity level
func (c *Console) Infof(format string, args ...interface{}) {
	fmt.Printf("INFO: %s", fmt.Sprintf(format, args...))
}

// DEBUG ...
func (c *Console) Debugf(format string, args ...interface{}) {
	fmt.Printf("DEBUG: %s", fmt.Sprintf(format, args...))
}

// ERROR ....
func (c *Console) Errorf(format string, args ...interface{}) {
	fmt.Printf("ERROR: %s", fmt.Sprintf(format, args...))
}

// CRITICAL is system failure
func (c *Console) Criticalf(format string, args ...interface{}) {
	fmt.Printf("CRITICAL: %s", fmt.Sprintf(format, args...))
}
