package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type JSONB []NameX

// Value Marshal
func (a JSONB) Value() (driver.Value, error) {

	j, err := json.Marshal(a)
	return j, err
	//return json.Marshal(a)
}

// Scan Unmarshal
// ✅ รองรับ JSONB ทั้งแบบ Object `{}` และ Array `[]`
func (a *JSONB) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	// ✅ ตรวจสอบว่าเป็น JSON Object หรือ JSON Array
	if string(b) == "null" {
		*a = nil
		return nil
	}

	// ✅ ลอง Unmarshal เป็น Array
	var temp []NameX
	if err := json.Unmarshal(b, &temp); err == nil {
		*a = temp
		return nil
	}

	// ❌ ถ้าไม่ใช่ Array ลอง Unmarshal เป็น Object แล้วแปลงเป็น Array
	var tempObj NameX
	if err := json.Unmarshal(b, &tempObj); err == nil {
		*a = []NameX{tempObj} // ✅ แปลง Object เป็น Array
		return nil
	}

	return errors.New("failed to unmarshal JSONB into JSONB type")
}
