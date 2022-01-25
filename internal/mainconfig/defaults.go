package mainconfig

import "github.com/spf13/viper"

// ConfigInit is the common config initialisation for the commands.
func ConfigInit() {
	viper.SetConfigName("crl-updater")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./artifacts")
	viper.AddConfigPath("./test")
	viper.AddConfigPath("$HOME/.crl-updater")
	viper.AddConfigPath("$HOME/.config")
	viper.AddConfigPath("/run/secrets")
	viper.AddConfigPath("/etc")
	viper.AddConfigPath("/etc/crl-updater")
	viper.AddConfigPath("/usr/local/etc")
	viper.AddConfigPath("/usr/local/crl-updater/etc")
	viper.AddConfigPath("/etc/ssl")
	viper.AddConfigPath(".")

	_ = viper.ReadInConfig()
}
