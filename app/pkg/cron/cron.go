package cron2

import (
	pc "avito/app/pkg/context"
	e "avito/app/pkg/errors"
	"context"
	"github.com/go-co-op/gocron/v2"
	"time"
)

type Cron struct {
	Interval  time.Duration
	Scheduler *gocron.Scheduler
}

func NewCron(ctx context.Context) (*Cron, error) {

	cfg := pc.GetConfig(ctx)

	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, e.Wrap(err)
	}

	return &Cron{cfg.Cron.Interval, &s}, nil

}
