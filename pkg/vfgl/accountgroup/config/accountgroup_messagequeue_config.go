package config

const (
	MQ_TOPIC_ACCOUNTGROUP_CREATED      string = "when-accountgroup-created"
	MQ_TOPIC_ACCOUNTGROUP_UPDATED      string = "when-accountgroup-updated"
	MQ_TOPIC_ACCOUNTGROUP_DELETED      string = "when-accountgroup-deleted"
	MQ_TOPIC_ACCOUNTGROUP_BULK_CREATED string = "when-accountgroup-bulk-created"
)

type AccountGroupMessageQueueConfig struct{}

func (AccountGroupMessageQueueConfig) TopicCreated() string {
	return MQ_TOPIC_ACCOUNTGROUP_CREATED
}

func (AccountGroupMessageQueueConfig) TopicUpdated() string {
	return MQ_TOPIC_ACCOUNTGROUP_UPDATED
}

func (AccountGroupMessageQueueConfig) TopicDeleted() string {
	return MQ_TOPIC_ACCOUNTGROUP_DELETED
}

func (AccountGroupMessageQueueConfig) TopicBulkCreated() string {
	return MQ_TOPIC_ACCOUNTGROUP_BULK_CREATED
}
