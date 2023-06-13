package config

const (
	MQ_TOPIC_CREATED      string = "when-purchase-created"
	MQ_TOPIC_UPDATED      string = "when-purchase-updated"
	MQ_TOPIC_DELETED      string = "when-purchase-deleted"
	MQ_TOPIC_BULK_CREATED string = "when-purchase-bulk-created"
	MQ_TOPIC_BULK_UPDATED string = "when-purchase-bulk-updated"
	MQ_TOPIC_BULK_DELETED string = "when-purchase-bulk-deleted"
)

type PurchaseQueueConfig struct{}

func (PurchaseQueueConfig) TopicCreated() string {
	return MQ_TOPIC_CREATED
}

func (PurchaseQueueConfig) TopicUpdated() string {
	return MQ_TOPIC_UPDATED
}

func (PurchaseQueueConfig) TopicDeleted() string {
	return MQ_TOPIC_DELETED
}

func (PurchaseQueueConfig) TopicBulkCreated() string {
	return MQ_TOPIC_BULK_CREATED
}

func (PurchaseQueueConfig) TopicBulkUpdated() string {
	return MQ_TOPIC_BULK_UPDATED
}

func (PurchaseQueueConfig) TopicBulkDeleted() string {
	return MQ_TOPIC_BULK_DELETED
}
