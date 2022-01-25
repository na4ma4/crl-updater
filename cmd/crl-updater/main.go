package main

import (
	"context"
	"errors"

	"github.com/na4ma4/config"
	"github.com/na4ma4/crl-updater/internal/crlmgr"
	"github.com/na4ma4/crl-updater/internal/mainconfig"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

//nolint:gochecknoglobals // cobra uses globals in main
var rootCmd = &cobra.Command{
	Use:  "crl-updater",
	RunE: mainCommand,
}

//nolint:gochecknoinits // init is used in main for cobra
func init() {
	cobra.OnInitialize(mainconfig.ConfigInit)

	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Debug output")
	_ = viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	_ = viper.BindEnv("debug", "DEBUG")
}

func main() {
	_ = rootCmd.Execute()
}

func mainCommand(cmd *cobra.Command, args []string) error {
	cfg := config.NewViperConfDFromViper(viper.GetViper(), "/etc/crl-updater/conf.d/", "crl-updater")

	logcfg := cfg.ZapConfig()
	logcfg.OutputPaths = []string{"stdout"}

	logger, _ := logcfg.Build()
	defer logger.Sync() //nolint:errcheck

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	targets := crlmgr.Parse(logger, viper.GetViper())

	for _, target := range targets {
		targetLogger := logger.With(
			zap.String("source", target.Source()),
			zap.String("target", target.Target()),
		)

		err := target.Run(ctx, targetLogger)

		switch {
		case errors.Is(err, crlmgr.ErrNotModified):
			targetLogger.Debug("source not modified, skipping")
		case err != nil:
			targetLogger.Error("unable to execute target", zap.Error(err))
		default:
			targetLogger.Info("target updated")
		}
	}

	return nil
}
