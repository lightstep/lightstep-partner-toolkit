package incidents

import (
	"github.com/lightstep/lightstep-partner-sdk/collector/generatorreceiver/internal/cron"
	"go.uber.org/zap"
	"time"
)

type Incident struct {
	CronStart string `json:"start" yaml:"start"`
	CronEnd   string `json:"end" yaml:"end"`
	Name      string

	started time.Time
}

func (i *Incident) Setup(logger *zap.Logger) {
	if i.CronStart != "" {
		_, err := cron.Add(i.CronStart, func() {
			logger.Info("starting incident", zap.String("incident", i.Name))
			i.Start()
		})
		if err != nil {
			logger.Error("error adding incident start schedule", zap.Error(err))
		}
	}
	if i.CronEnd != "" {
		_, err := cron.Add(i.CronEnd, func() {
			logger.Info("ending incident", zap.String("incident", i.Name))
			i.End()
		})
		if err != nil {
			logger.Error("error adding incident end schedule", zap.Error(err))
		}
	}
}

func (i *Incident) Active() bool {
	return i.CurrentDuration() > 0
}

func (i *Incident) CurrentDuration() time.Duration {
	if i.started.IsZero() {
		return time.Duration(0)
	}
	return time.Now().Sub(i.started)
}

func (i *Incident) Start() {
	if i.started.IsZero() {
		i.started = time.Now()
	}
}

func (i *Incident) End() {
	i.started = time.Time{}
}
