package microservice

import (
	"context"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// IMQ is interface to manage Kafka
type IMQ interface {
	CreateTopic(topic string, partitions int, replications int) error
	CreateTopicR(topic string, partitions int, replications int, retentionPeriod time.Duration) error
}

// IMQConfig is mq configuration interface
type IMQConfig interface {
	URI() string
}

// MQ is message queue
type MQ struct {
	ms      *Microservice
	servers string
}

// NewMQ return new MQ
func NewMQ(mqConfig IMQConfig, ms *Microservice) *MQ {
	return &MQ{
		ms:      ms,
		servers: mqConfig.URI(),
	}
}

func (q *MQ) getAdminClient() (*kafka.AdminClient, error) {
	admin, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": q.servers})
	if err != nil {
		q.ms.Logger.WithError(err).Error("Failed to Connect to Kafka")
		return nil, err
	}
	return admin, nil
}

// CreateTopicR create topic with retention period
func (q *MQ) CreateTopicR(topic string, partitions int, replications int, retentionPeriod time.Duration) error {
	return q.createTopic(topic, partitions, replications, retentionPeriod)
}

// CreateTopic create the topic
func (q *MQ) CreateTopic(topic string, partitions int, replications int) error {
	return q.createTopic(topic, partitions, replications, 0)
}

func (q *MQ) createTopic(topic string, partitions int, replications int, retentionPeriod time.Duration) error {
	if retentionPeriod <= 0 {
		retentionPeriod = 7 * (time.Hour * 24) // default = 7 days (Message will keep 7 days)
	}

	admin, err := q.getAdminClient()
	if err != nil {
		return err
	}

	defer admin.Close()

	// Operation timeout for create topic = 5 minutes
	timeout, err := time.ParseDuration("5m")
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	retentionPeriodMillisec := fmt.Sprintf("%d", int64(retentionPeriod/time.Millisecond))

	var results []kafka.TopicResult
	if timeout > 0 {
		results, err = admin.CreateTopics(
			ctx,
			[]kafka.TopicSpecification{{
				Topic:             topic,
				NumPartitions:     partitions,
				ReplicationFactor: replications,
				Config: map[string]string{
					"retention.ms": retentionPeriodMillisec,
				}}},
			kafka.SetAdminOperationTimeout(timeout))
	} else {
		results, err = admin.CreateTopics(
			ctx,
			[]kafka.TopicSpecification{{
				Topic:             topic,
				NumPartitions:     partitions,
				ReplicationFactor: replications,
				Config: map[string]string{
					"retention.ms": retentionPeriodMillisec,
				}}})
	}
	if err != nil {
		return err
	}

	for _, result := range results {
		q.ms.Logger.Debugf("Create Topic \"%s\" Result: %s", topic, result.String())
	}

	return nil
}
