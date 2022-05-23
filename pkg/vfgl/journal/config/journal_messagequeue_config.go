package config

const (
	MQ_TOPIC_JOURNAL_CREATED      string = "when-journal-created"
	MQ_TOPIC_JOURNAL_UPDATED      string = "when-journal-updated"
	MQ_TOPIC_JOURNAL_DELETED      string = "when-journal-deleted"
	MQ_TOPIC_JOURNAL_BULK_CREATED string = "when-journal-bulk-created"
)

type JournalMessageQueueConfig struct{}

func (JournalMessageQueueConfig) TopicCreated() string {
	return MQ_TOPIC_JOURNAL_CREATED
}

func (JournalMessageQueueConfig) TopicUpdated() string {
	return MQ_TOPIC_JOURNAL_UPDATED
}

func (JournalMessageQueueConfig) TopicDeleted() string {
	return MQ_TOPIC_JOURNAL_DELETED
}

func (JournalMessageQueueConfig) TopicBulkCreated() string {
	return MQ_TOPIC_JOURNAL_BULK_CREATED
}
