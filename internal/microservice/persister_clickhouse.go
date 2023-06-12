package microservice

import (
	"context"
	"fmt"
	"smlcloudplatform/pkg/config"
	"sync"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type IPersisterClickHouse interface {
	Conn() driver.Conn
}

type ChModel interface {
	TableName() string
}

type PersisterClickHouse struct {
	cfg     config.IPersisterClickHouseConfig
	ctx     context.Context
	dbMutex sync.Mutex
	db      driver.Conn
}

func NewPersisterClickHouse(cfg config.IPersisterClickHouseConfig) *PersisterClickHouse {

	ctx := context.Background()

	pst := PersisterClickHouse{
		cfg: cfg,
		ctx: ctx,
	}

	_, err := pst.getClient()
	if err != nil {
		panic(err)
	}

	return &pst
}

func (pst *PersisterClickHouse) getClient() (driver.Conn, error) {
	if pst.db != nil {
		return pst.db, nil
	}

	// connect

	connection, err := clickhouse.Open(&clickhouse.Options{
		Addr: pst.cfg.ServerAddress(),
		Auth: clickhouse.Auth{
			Database: pst.cfg.DatabaseName(),
			Username: pst.cfg.Username(),
			Password: pst.cfg.Password(),
		},
		Debug: true,
	})

	if err != nil {
		return nil, err
	}
	pst.dbMutex.Lock()
	pst.db = connection
	pst.dbMutex.Unlock()
	return pst.db, nil
}

func (pst *PersisterClickHouse) GetTableName(model interface{}) (string, error) {

	chmodel, ok := model.(ChModel)

	if ok {
		return chmodel.TableName(), nil
	}
	return "", fmt.Errorf("struct is not implement Clickhouse Model")
}

func (pst *PersisterClickHouse) Conn() driver.Conn {
	return pst.db
}
