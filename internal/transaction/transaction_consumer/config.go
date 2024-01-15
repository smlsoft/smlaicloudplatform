package transactionconsumer

type ITransactionConsumerConfig interface {
	PurchaseTopicCreate() string
	PurchaseTopicUpdate() string
	PurchaseTopicDelete() string
}

type TransactionConsumerConfig struct{}
