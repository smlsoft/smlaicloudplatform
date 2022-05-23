package microservice_test

import (
	"fmt"
	"smlcloudplatform/internal/microservice"
	"testing"
)

type ConfigPostgresqlDBTest struct{}

func (cfg *ConfigPostgresqlDBTest) Host() string {
	return "localhost"
}

func (cfg *ConfigPostgresqlDBTest) Port() string {
	return "5432"
}

func (cfg *ConfigPostgresqlDBTest) DB() string {
	return "dev"
}

func (cfg *ConfigPostgresqlDBTest) Username() string {
	return "postgres"
}

func (cfg *ConfigPostgresqlDBTest) Password() string {
	return "sml"
}

func (cfg *ConfigPostgresqlDBTest) SSLMode() string {
	return "disable"
}

func (cfg *ConfigPostgresqlDBTest) TimeZone() string {
	return "Asia/Bangkok"
}

func (cfg *ConfigPostgresqlDBTest) LoggerLevel() string {
	return "debug"
}

type User struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"column:name"`
}

func (u User) TableName() string {
	return "users"
}

func TestTransaction(t *testing.T) {
	pgConfig := &ConfigPostgresqlDBTest{}

	pst := microservice.NewPersister(pgConfig)

	err := pst.Transaction(func(pst *microservice.Persister) error {
		u := User{
			Name: "namex",
		}
		err := pst.Create(&u)
		if err != nil {
			t.Log("create:: ", err)
			return err
		}

		fmt.Println(u.ID)

		return nil
	})

	if err != nil {
		t.Log("trans :: ", err)
	}
}

func TestUpdate(t *testing.T) {
	pgConfig := &ConfigPostgresqlDBTest{}

	pst := microservice.NewPersister(pgConfig)
	u := User{
		Name: "name1 modifyx",
	}

	pst.Update(&u, map[string]interface{}{
		"id":   2,
		"name": "name edited",
	})
}

func TestDelete(t *testing.T) {
	pgConfig := &ConfigPostgresqlDBTest{}

	pst := microservice.NewPersister(pgConfig)
	pst.Delete(&User{}, map[string]interface{}{
		"id":   3,
		"name": "namex",
	})
}
