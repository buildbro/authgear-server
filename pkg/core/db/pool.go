package db

import (
	"fmt"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/skygeario/skygear-server/pkg/core/config"
)

type Pool interface {
	Open(tConfig config.TenantConfiguration) (*sqlx.DB, error)
	Close() error
}

type poolImpl struct {
	closed     bool
	closeMutex sync.RWMutex

	cache      map[string]*sqlx.DB
	cacheMutex sync.RWMutex
}

func NewPool() Pool {
	p := &poolImpl{cache: map[string]*sqlx.DB{}}
	return p
}

var errPoolClosed = fmt.Errorf("database pool is closed")

func (p *poolImpl) Open(tConfig config.TenantConfiguration) (db *sqlx.DB, err error) {
	p.closeMutex.RLock()
	defer func() { p.closeMutex.RUnlock() }()
	if p.closed {
		return nil, errPoolClosed
	}

	source := tConfig.AppConfig.DatabaseURL

	p.cacheMutex.RLock()
	db, exists := p.cache[source]
	p.cacheMutex.RUnlock()

	if !exists {
		p.cacheMutex.Lock()
		db, exists = p.cache[source]
		if !exists {
			db, err = sqlx.Open("postgres", source)
			if err == nil {
				p.cache[source] = db
			}
		}
		p.cacheMutex.Unlock()
	}

	return
}

func (p *poolImpl) Close() (err error) {
	p.closeMutex.Lock()
	defer func() { p.closeMutex.Unlock() }()

	p.closed = true
	for _, db := range p.cache {
		if closeErr := db.Close(); closeErr != nil {
			err = closeErr
		}
	}

	return
}
