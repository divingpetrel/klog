// +build go1.12

package klog

type Logger interface {
	Infof(message string, args ...interface{})
	Debugf(message string, args ...interface{})
	Errorf(message string, args ...interface{})
	Criticalf(message string, args ...interface{})
}
