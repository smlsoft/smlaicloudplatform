package transaction

// rebuild transaction to postgresql
// fetch data from mongodb and send to kafka

type ITransactionHTTP interface{}

type TransactionHTTP struct{}

func NewTransactionHttp() *TransactionHTTP {
	return &TransactionHTTP{}
}
