package streamreceiver

import (
	"context"
	"fmt"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenterror"
	"go.opentelemetry.io/collector/consumer"
	"go.uber.org/zap"
	"net/url"
	"time"
)

type streamReceiver struct {
	logger     *zap.Logger
	traceConsumer   consumer.Traces
	ticker    *time.Ticker
	client Client
	stop chan struct{}
}

func (s streamReceiver) Start(ctx context.Context, host component.Host) error {
	s.ticker = time.NewTicker(5 * time.Second)
	s.stop = make(chan struct{})
	go func() {
		for {
			select {
			case <- s.ticker.C:
				s.getTraces()
			case <- s.stop:
				s.ticker.Stop()
				return
			}
		}
	}()

	return nil
}

func (s streamReceiver) getTraces() {
	s.logger.Info("Getting traces...")
	resp, err := s.client.GetStreamTraces()
	if err != nil {
		s.logger.Info(fmt.Sprintf("Could not get traces: %v", err))
		return
	}
	s.logger.Info(fmt.Sprintf("found exemplars: %v", len(resp.Data.Attributes.Exemplars)))
}

func (s streamReceiver) Shutdown(ctx context.Context) error {
	close(s.stop)
	return nil
}

var snapReceiver = streamReceiver{}

func newTraceReceiver(config *Config,
	consumer consumer.Traces,
	logger *zap.Logger) (component.TracesReceiver, error) {

	if consumer == nil {
		return nil, componenterror.ErrNilNextConsumer
	}
	u, _ := url.Parse("https://api.lightstep.com/public/v0.2/")
	snapReceiver.logger = logger
	c := NewClientProvider(*u, config.Organization, config.Project, config.ApiKey, logger).BuildClient()
	snapReceiver.client = c
	return &snapReceiver, nil
}
