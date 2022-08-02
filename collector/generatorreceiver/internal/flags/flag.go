package flags

import (
	"github.com/lightstep/lightstep-partner-sdk/collector/generatorreceiver/internal/cron"
	"go.uber.org/zap"
	"time"
)

const (
	DisabledState = 0.0
	EnabledState  = 1.0
)

// TODO: separate config types from code types generally

type IncidentConfig struct {
	ParentFlag string        `json:"parentFlag" yaml:"parentFlag"`
	Start      time.Duration `json:"start" yaml:"start"`
	End        time.Duration `json:"end" yaml:"end"`
}

type CronConfig struct {
	Start string `json:"start" yaml:"start"`
	End   string `json:"end" yaml:"end"`
}

type Flag struct {
	Name     string          `json:"name" yaml:"name"`
	Incident *IncidentConfig `json:"incident" yaml:"incident"`
	Cron     *CronConfig     `json:"cron" yaml:"cron"`

	started time.Time
}

func (f *Flag) Active() bool {
	f.update()
	return f.active()
}

func (f *Flag) active() bool {
	return !f.started.IsZero()
}

// update checks if the given flag f has a parent flag ("Incident"); if so,
// updates f's state based on its start and end times relative to the parent.
func (f *Flag) update() {
	if f.Incident == nil {
		// managed by cron or manual-only
		return
	}
	// Managed by parent
	parent := Manager.GetFlag(f.Incident.ParentFlag)
	if parent == nil {
		// TODO: with validation, this should never happen
		return
	}

	incidentDuration := parent.CurrentDuration()
	afterStart := incidentDuration > f.Incident.Start
	beforeEnd := f.Incident.End == 0 || incidentDuration < f.Incident.End
	// shouldBeActive will be true if and only if both of the following are true:
	// - parent has been Active for at least f's incident start time
	// - either f has no end *or* parent has been active for less than f's incident end time
	shouldBeActive := afterStart && beforeEnd

	if f.active() != shouldBeActive {
		f.Toggle()
	}
}

func (f *Flag) CurrentDuration() time.Duration {
	if !f.Active() {
		return 0
	}
	return time.Now().Sub(f.started)
}

func (f *Flag) Enable() {
	if f.started.IsZero() {
		f.started = time.Now()
	}
}

func (f *Flag) Disable() {
	f.started = time.Time{}
}

func (f *Flag) Toggle() {
	if f.active() {
		f.Disable()
	} else {
		f.Enable()
	}
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
