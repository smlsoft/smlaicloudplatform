package config

const (
	MQ_TOPIC_CHARTOFACCOUNT_CREATED      string = "when-chartofaccount-created"
	MQ_TOPIC_CHARTOFACCOUNT_UPDATED      string = "when-chartofaccount-updated"
	MQ_TOPIC_CHARTOFACCOUNT_DELETED      string = "when-chartofaccount-deleted"
	MQ_TOPIC_CHARTOFACCOUNT_BULK_CREATED string = "when-chartofaccount-bulk-created"
)

type ChartOfAccountMessageQueueConfig struct{}

func (ChartOfAccountMessageQueueConfig) TopicCreated() string {
	return MQ_TOPIC_CHARTOFACCOUNT_CREATED
}

func (ChartOfAccountMessageQueueConfig) TopicUpdated() string {
	return MQ_TOPIC_CHARTOFACCOUNT_UPDATED
}

func (ChartOfAccountMessageQueueConfig) TopicDeleted() string {
	return MQ_TOPIC_CHARTOFACCOUNT_DELETED
}

func (ChartOfAccountMessageQueueConfig) TopicBulkCreated() string {
	return MQ_TOPIC_CHARTOFACCOUNT_BULK_CREATED
}
