package crlmgr

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Parse the targets from the configuration.
func Parse(logger *zap.Logger, cfg *viper.Viper) []*Target {
	o := []*Target{}

	targets := make(map[string]bool)

	for _, key := range cfg.AllKeys() {
		spkey := strings.SplitN(key, ".", 2) //nolint:gomnd // 2 is more than 1, just need the first part.
		targets[spkey[0]] = true
	}

	for target := range targets {
		switch target {
		case "watchdog", "debug":
			// skip these categories
			continue
		}

		logger.Debug("adding target from config", zap.String("target", target))

		i := &Target{
			source:  cfg.GetString(fmt.Sprintf("%s.source", target)),
			target:  cfg.GetString(fmt.Sprintf("%s.target", target)),
			workdir: cfg.GetString(fmt.Sprintf("%s.workdir", target)),
			actions: make(map[EventType]*Action),
		}

		for _, et := range eventTypes {
			if v := cfg.GetStringSlice(fmt.Sprintf("%s.actions.%s", target, et)); len(v) > 0 {
				i.actions[et] = ActionFromString(v)
				if i.workdir != "" {
					i.actions[et].workdir = i.workdir
				}
			}
		}

		o = append(o, i)
	}

	return o
}
