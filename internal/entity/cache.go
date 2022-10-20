package entity

import (
	"sync"
	"time"
)

type Cache struct {
	sync.RWMutex
	DefaultExpiration time.Duration
	Data              ResultSetT
	Created           time.Time
	Expiration        int64
}
