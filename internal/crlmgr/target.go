package crlmgr

import (
	"context"
	"errors"
	"io"
	"net/url"
	"os"
	"time"

	"go.uber.org/zap"
)

// Target is a matching target and source for syncronisation.
type Target struct {
	source  string
	target  string
	workdir string
	actions map[EventType]*Action
}

// Source is the source of the CRL.
func (t *Target) Source() string {
	return t.source
}

// Target is the on-disk location for the CRL file.
func (t *Target) Target() string {
	return t.target
}

// Retrieve the CRL source.
func (t *Target) Retrieve(ctx context.Context, ts time.Time) (io.ReadCloser, time.Time, error) {
	u, _ := url.Parse(t.Source())

	return Retrieve(ctx, u, ts)
}

// Run the target sync.
func (t *Target) Run(ctx context.Context, logger *zap.Logger) error {
	defer t.ExecAction(ctx, logger, PostAction)
	t.ExecAction(ctx, logger, PreCheckAction)

	modTime, err := Stat(t.Target())
	if errors.Is(err, os.ErrNotExist) {
		// do nothing
	} else if err != nil {
		return err
	}

	body, ts, err := t.Retrieve(ctx, modTime)
	if errors.Is(err, ErrNotModified) {
		return err
	} else if err != nil {
		return err
	}

	t.ExecAction(ctx, logger, PreInstallAction)

	if err := Store(t.Target(), ts, body); err != nil {
		return err
	}

	t.ExecAction(ctx, logger, PostInstallAction)

	return nil
}

// ExecAction will execute the specified action and log anything interesting, returning nothing.
func (t *Target) ExecAction(ctx context.Context, logger *zap.Logger, event EventType) {
	if a := t.Action(event); a != nil {
		output, err := a.Exec(ctx)
		if err != nil {
			logger.Error("unable to execute event action", zap.String("event", event.String()), zap.String("output", output))
		}

		logger.Info("action executed", zap.String("event", event.String()), zap.String("output", output))
	}
}

// Action returns the action specified for the event, or nil if it isn't valid.
func (t *Target) Action(event EventType) *Action {
	if v, ok := t.actions[event]; ok {
		return v
	}

	return nil
}
