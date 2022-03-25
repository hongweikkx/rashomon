package service

import (
	"github.com/newrelic/go-agent/v3/newrelic"
)

var MonitorNewRelic *newrelic.Application

func init() {
	var err error
	MonitorNewRelic, err = newrelic.NewApplication(
		newrelic.ConfigAppName("rashomon"),
		newrelic.ConfigLicense("8e361fcb912dbcc9af2883af61b26774e527NRAL"),
		newrelic.ConfigDistributedTracerEnabled(true),
	)
	if err != nil {
		panic(err)
	}
}
