package config

const (
	MQ_TOPIC_ACCOUNTGROUP_CREATED      string = "when-journalbook-created"
	MQ_TOPIC_ACCOUNTGROUP_UPDATED      string = "when-journalbook-updated"
	MQ_TOPIC_ACCOUNTGROUP_DELETED      string = "when-journalbook-deleted"
	MQ_TOPIC_ACCOUNTGROUP_BULK_CREATED string = "when-journalbook-bulk-created"
)

type JournalBookMessageQueueConfig struct{}

func (JournalBookMessageQueueConfig) TopicCreated() string {
	return MQ_TOPIC_ACCOUNTGROUP_CREATED
}

func (JournalBookMessageQueueConfig) TopicUpdated() string {
	return MQ_TOPIC_ACCOUNTGROUP_UPDATED
}

func (JournalBookMessageQueueConfig) TopicDeleted() string {
	return MQ_TOPIC_ACCOUNTGROUP_DELETED
}

func (JournalBookMessageQueueConfig) TopicBulkCreated() string {
	return MQ_TOPIC_ACCOUNTGROUP_BULK_CREATED
}
