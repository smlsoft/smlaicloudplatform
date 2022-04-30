package models

type Payment struct {
	Cash              float64             `json:"cash" bson:"cash" `
	CreditCard        float64             `json:"creditcard" bson:"creditcard" `
	CreditCardDetails []CreditCardPayment `json:"creditcarddetails" bson:"creditcarddetails" `
}

type CreditCardPayment struct {
	CardType     string  `json:"cardtype" bson:"cardtype" `
	CardNumber   string  `json:"cardnumber" bson:"cardnumber" `
	Amount       float64 `json:"amount" bson:"amount" `
	ApprovedCode string  `json:"approvedcode" bson:"approvedcode" `
	Remark       string  `json:"remark" bson:"remark" `
}
