package config

// IMQConfig is mq configuration interface
type IMQConfig interface {
	URI() string
}

// MQ and Producer configuration
type MQConfig struct{}

func NewMQConfig() *MQConfig {
	return &MQConfig{}
}

func (MQConfig) URI() string {
	return getEnv("KAFKA_SERVER_URL", "") // localhost:9094
}

func (*Config) MQConfig() IMQConfig {
	return NewMQConfig()
}
