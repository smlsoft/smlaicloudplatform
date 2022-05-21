package vfgl

type ChartOfAccount struct {
	// รหัสผังบัญชี
	AccountCode string
	// ชื่อบัญชี
	AccountName string
	// หมวดบัญชี 1=สินทรัพย์, 2=หนี้สิน, 3=ทุน, 4=รายได้, 5=ค่าใช้จ่าย
	AccountCategory int
	// ด้านบัญชี เดบิต,เครดิต
	AccountBalanceType string
	// กลุ่มบัญชี
	AccountGroup string
}

type ChartOfAccountDoc struct {
}
