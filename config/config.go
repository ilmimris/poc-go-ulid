package config

import (
	"os"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/joho/godotenv"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

type Source string

const (

	// default values by convention
	DefaultType = "json"
	DefaultFile = "config"

	// environment variables
	EnvConsulHostKey = "GOCONF_CONSUL"
	EnvTypeKey       = "GOCONF_TYPE"
	EnvFileNameKey   = "GOCONF_FILENAME"
	EnvPrefixKey     = "GOCONF_ENV_PREFIX"

	SourceEnv    Source = "env"
	SourceFile   Source = "file"
	SourceConsul Source = "consul"
)

var (
	// Config is the global configuration
	dirs = []string{"./", "./config/", "./configs/", "/etc/", "/app/config/", "/app/configs/"}

	typ       = DefaultType
	fname     = DefaultFile
	GlobalCfg Config
)

func LoadConfig(cfgFile string) {
	GlobalCfg = NewConfig()
	setConfigFile(cfgFile)
	// readFileRemote()

	// attempt to load from file in configured directories
	for _, dir := range dirs {
		viper.AddConfigPath(dir)
	}

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&GlobalCfg, decoderTime())
	if err != nil {
		panic(err)
	}
}

func GetConfig() Config {
	return GlobalCfg
}

func decoderTime() viper.DecoderConfigOption {
	return func(m *mapstructure.DecoderConfig) {
		m.DecodeHook = mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
			mapstructure.StringToTimeHookFunc(time.RFC3339Nano),
		)
	}
}

func setConfigFile(cfgFile string) {
	var prefix string
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		err := godotenv.Load()
		if err != nil {
			log.WithError(err).Error()
		}

		if v := os.Getenv(EnvTypeKey); len(v) > 0 {
			typ = v
		}
		if v := os.Getenv(EnvFileNameKey); len(v) > 0 {
			fname = v
		}
		if v := os.Getenv(EnvPrefixKey); len(v) > 0 {
			prefix = v
		}

		// setup and configure viper instance
		viper.SetConfigType(typ)
		viper.SetConfigName(fname)
		if len(prefix) > 0 {
			viper.SetEnvPrefix(prefix)
		}
	}
	viper.AutomaticEnv()
}

func readFileRemote() {
	var err error
	// next we load from consul; only if consul host defined
	if ch := os.Getenv(EnvConsulHostKey); ch != "" {
		err = viper.AddRemoteProvider(string(SourceConsul), ch, fname)
		if err != nil {
			log.WithError(err).Error()
		}

		connect := func() error { return viper.ReadRemoteConfig() }
		notify := func(err error, t time.Duration) { log.Println("[goconf]", err.Error(), t) }
		b := backoff.NewExponentialBackOff()
		b.MaxElapsedTime = 2 * time.Minute

		err = backoff.RetryNotify(connect, b, notify)
		if err != nil {
			log.WithError(err).Error("[goconf] giving up connecting to remote config ")
		}
	} else {
		log.Error("failed loading remote source; ENV not defined")
	}
}
