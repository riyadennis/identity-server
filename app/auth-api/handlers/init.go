package handlers

import (
	"time"

	"github.com/riyadennis/identity-server/business/store"
)

var (
	Idb store.Store
)

const (
	tokenTTL = 120 * time.Hour
)

func dataSource() store.Store {
	return Idb
}
