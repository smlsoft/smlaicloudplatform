package config

// IMQConfig is mq configuration interface
type IMQConfig interface {
	URI() string
	SecurityProtocol() string
	SSLKeyFile() string
	SSLCAFile() string
	SSLCertFile() string
}

// MQ and Producer configuration
type MQConfig struct{}

func NewMQConfig() *MQConfig {
	return &MQConfig{}
}

func (MQConfig) URI() string {
	return getEnv("KAFKA_SERVER_URL", "") // localhost:9094
}

func (MQConfig) SecurityProtocol() string {
	return getEnv("KAFKA_SECURITY_PROTOCOL", "plaintext") // SASL_SSL
}

func (MQConfig) SSLKeyFile() string {
	return getEnv("KAFKA_SSL_KEY_FILE", "") // /path/to/key
}

func (MQConfig) SSLCAFile() string {
	return getEnv("KAFKA_SSL_CA_FILE", "") // /path/to/ca
}

func (MQConfig) SSLCertFile() string {
	return getEnv("KAFKA_SSL_CERT_FILE", "") // /path/to/cert
}

func (*Config) MQConfig() IMQConfig {
	return NewMQConfig()
}
