package rabbitmq

type Config struct {
	Queue `mapstructure:"queue"`
	Channel  `mapstructure:"channel"`
	Publish  `mapstructure:"publish"`
}

func NewConfig(queue Queue, channel Channel, publish Publish) *Config {
	return &Config{Queue: queue, Channel: channel, Publish: publish}
}

type Queue struct {
	Durable   bool `mapstructure:"durable"`
	AutoDel   bool `mapstructure:"auto_delete"`
	Exclusive bool `mapstructure:"exclusive"`
	NoWait    bool `mapstructure:"no_wait"`
}

func NewQueue(durable bool, autoDel bool, exclusive bool, noWait bool) *Queue {
	return &Queue{Durable: durable, AutoDel: autoDel, Exclusive: exclusive, NoWait: noWait}
}

type Channel struct {
	AutoAck   bool `mapstructure:"auto_ack"`
	Exclusive bool `mapstructure:"exclusive"`
	NoLocal   bool `mapstructure:"no_local"`
	NoWait    bool `mapstructure:"no_wait"`
}

func NewChannel(autoAck bool, exclusive bool, noLocal bool, noWait bool) *Channel {
	return &Channel{AutoAck: autoAck, Exclusive: exclusive, NoLocal: noLocal, NoWait: noWait}
}

type Publish struct {
	Mandatory   bool `mapstructure:"mandatory"`
	Immediate   bool `mapstructure:"immediate"`
}

func NewPublish(mandatory bool, immediate bool) *Publish {
	return &Publish{Mandatory: mandatory, Immediate: immediate}
}
