package rabbitmq

type Config struct {
	RabbitMQ `mapstructure:"rabbitmq"`
	Delivery `mapstructure:"delivery"`
}

type RabbitMQ struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type Delivery struct {
	Queue    `mapstructure:"queue"`
	Channel  `mapstructure:"channel"`
	Publish  `mapstructure:"publish"`
	Exchange `mapstructure:"exchange"`
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
	Mandatory bool `mapstructure:"mandatory"`
	Immediate bool `mapstructure:"immediate"`
}

type Exchange struct {
	Type     string `mapstructure:"type"`
	Durable  bool   `mapstructure:"durable"`
	AutoDel  bool   `mapstructure:"auto_delete"`
	Internal bool   `mapstructure:"internal"`
	NoWait   bool   `mapstructure:"no_wait"`
}
