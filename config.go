package rabbitmq

type Config struct {
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
