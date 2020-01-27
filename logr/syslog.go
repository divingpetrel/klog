// +build go1.12
// +build syslog

package logr

import (
	"fmt"
	"log/syslog"
)

type Syslog struct {
	*syslog.Writer
}

func New(daemon, version string) (*Syslog, error) {
	sys, err := syslog.New(syslog.LOG_ERR, "")
	defer sys.Close()
	if err != nil {
		fmt.Printf("error %s", err.Error())
	}

	return &Syslog{sys}, nil
}

func (s *Syslog) Infof(format string, args ...interface{}) {
	s.Info(fmt.Sprintf(format, args...))
}

func (s *Syslog) Debugf(format string, args ...interface{}) {
	s.Debug(fmt.Sprintf(format, args...))
}

func (s *Syslog) Errorf(format string, args ...interface{}) {
	s.Err(fmt.Sprintf(format, args...))
}

func (s *Syslog) Criticalf(format string, args ...interface{}) {
	s.Crit(fmt.Sprintf(format, args...))
}
