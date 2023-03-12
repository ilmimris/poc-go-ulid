package config

import (
	"github.com/ilmimris/poc-go-ulid/shared/constant"
)

type Config struct {
	Rest struct {
		Port              int  `mapstructure:"port"`
		BodyLimit         int  `mapstructure:"body_limit"`
		ReduceMemoryUsage bool `mapstructure:"reduce_memory_usage"`
		ReadTimeOut       int  `mapstructure:"read_timeout"`
		WriteTimeOut      int  `mapstructure:"write_timeout"`
	} `mapstructure:"rest"`
	Logger struct {
		Level string `mapstructure:"level"`
	} `mapstructure:"logger"`
	Database struct {
		Postgres struct {
			Host     string `mapstructure:"host"`
			Port     int    `mapstructure:"port"`
			Username string `mapstructure:"username"`
			Password string `mapstructure:"password"`
			Database string `mapstructure:"database"`
		} `mapstructure:"postgresdb"`
		Maria struct {
			Dsn string `mapstructure:"dsn"`
		} `mapstructure:"mariadb"`
		Mongo struct {
			URI               string `mapstructure:"uri"`
			Db                string `mapstructure:"db"`
			ConnectionTimeOut int64  `mapstructure:"connection_timeout"`
			PingTimeOut       int64  `mapstructure:"ping_timeout"`
		} `mapstructure:"mongodb"`
	} `mapstructure:"database"`
	TracerConfig TracerConfig `mapstructure:"tracer_config"`
}

type TracerConfig struct {
	Tracer             string  `mapstructure:"tracer" json:"tracer"`
	JaegerCollectorURL string  `mapstructure:"jaeger_url" json:"jaeger_url"`
	JaegerMode         string  `mapstructure:"jaeger_mode" json:"jaeger_mode"`
	SamplingRate       float32 `mapstructure:"tracer_sample_rate" json:"tracer_sample_rate"`
}

func NewConfig() (cfg Config) {
	cfg.Rest.Port = constant.DefaultPort
	cfg.Rest.BodyLimit = constant.MaxBodyLimit
	cfg.Rest.ReadTimeOut = constant.MaxReadTimeOut
	cfg.Rest.WriteTimeOut = constant.MaxWriteTimeOut
	cfg.Logger.Level = constant.DefaultLogLevel

	return
}
