package flags

import (
	"fmt"
	"github.com/lightstep/lightstep-partner-sdk/collector/generatorreceiver/internal/cron"
	"github.com/lightstep/lightstep-partner-sdk/collector/generatorreceiver/internal/incidents"
	"go.uber.org/zap"
	"time"
)

const (
	DisabledState = 0.0
	EnabledState  = 1.0
)

// TODO: separate config types from code types generally

type IncidentConfig struct {
	Name  string
	Start time.Duration
	End   time.Duration
}

type CronConfig struct {
	Start string
	End   string
}

type Flag struct {
	Name     string          `json:"name" yaml:"name"`
	Incident *IncidentConfig `json:"incident" yaml:"incident"`
	Cron     *CronConfig     `json:"cron" yaml:"cron"`

	state float64
}

func (f *Flag) Enabled() bool {
	return f.GetState() > DisabledState
}

func (f *Flag) GetState() float64 {
	if f.Incident == nil {
		return f.state
	}

	incident := incidents.Manager.GetIncident(f.Incident.Name)
	// TODO: where does this logic belong?
	if f.Name == fmt.Sprintf("%s.default", incident.Name) {
		if !incident.Active() {
			return EnabledState
		}
		return DisabledState
	}
	duration := incident.CurrentDuration()
	if duration > f.Incident.Start {
		if f.Incident.End == 0 || duration < f.Incident.End {
			return EnabledState
		}
	}
	return DisabledState
}

func (f *Flag) Enable() {
	f.SetState(EnabledState)
}

func (f *Flag) Disable() {
	f.SetState(DisabledState)
}

func (f *Flag) SetState(state float64) {
	f.state = state
}

func (f *Flag) Setup(logger *zap.Logger) {
	// TODO: add validation to disallow having cron and incident both configured?
	if f.Cron != nil {
		f.SetupCron(logger)
	}
}

func (f *Flag) SetupCron(logger *zap.Logger) {
	_, err := cron.Add(f.Cron.Start, func() {
		logger.Info("toggling flag on", zap.String("flag", f.Name))
		f.Enable()
	})
	if err != nil {
		logger.Error("error adding flag start schedule", zap.Error(err))
	}

	_, err = cron.Add(f.Cron.End, func() {
		logger.Info("toggling flag off", zap.String("flag", f.Name))
		f.Disable()
	})
	if err != nil {
		logger.Error("error adding flag stop schedule", zap.Error(err))
	}
}
