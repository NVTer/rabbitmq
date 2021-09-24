package rabbitmq

import (
	"github.com/spf13/viper"
)


type config struct {
	Queue `mapstructure:"queue"`
	Channel  `mapstructure:"channel"`
	Publish  `mapstructure:"publish"`
}

type Queue struct {
	Durable   bool `mapstructure:"durable"`
	AutoDel   bool `mapstructure:"auto_delete"`
	Exclusive bool `mapstructure:"exclusive"`
	NoWait    bool `mapstructure:"no_wait"`
}

type Channel struct {
	AutoAck   bool `mapstructure:"auto_ack"`
	Exclusive bool `mapstructure:"exclusive"`
	NoLocal   bool `mapstructure:"no_local"`
	NoWait    bool `mapstructure:"no_wait"`
}

type Publish struct {
	Mandatory   bool `mapstructure:"mandatory"`
	Immediate   bool `mapstructure:"immediate"`
}

func newConfig() (*config, error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	var serviceConfig config
	err = viper.Unmarshal(&serviceConfig)
	if err != nil {
		return nil, err
	}
	return &serviceConfig, nil
}
