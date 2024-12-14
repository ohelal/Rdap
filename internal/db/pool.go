package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"time"
)

type DBPool struct {
	master  *sql.DB
	slaves  []*sql.DB
	current int
}

func NewDBPool(masterDSN string, slaveDSNs []string) (*DBPool, error) {
	master, err := sql.Open("postgres", masterDSN)
	if err != nil {
		return nil, err
	}

	master.SetMaxOpenConns(100)
	master.SetMaxIdleConns(10)
	master.SetConnMaxLifetime(time.Hour)

	slaves := make([]*sql.DB, len(slaveDSNs))
	for i, dsn := range slaveDSNs {
		slave, err := sql.Open("postgres", dsn)
		if err != nil {
			return nil, err
		}
		slave.SetMaxOpenConns(100)
		slave.SetMaxIdleConns(10)
		slave.SetConnMaxLifetime(time.Hour)
		slaves[i] = slave
	}

	return &DBPool{
		master: master,
		slaves: slaves,
	}, nil
}

func (p *DBPool) GetReader() *sql.DB {
	if len(p.slaves) == 0 {
		return p.master
	}
	p.current = (p.current + 1) % len(p.slaves)
	return p.slaves[p.current]
}

func (p *DBPool) GetWriter() *sql.DB {
	return p.master
}
