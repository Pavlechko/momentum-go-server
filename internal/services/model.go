package services

import (
	"sync"

	"momentum-go-server/internal/store"
)

type Service struct {
	DB      *store.Database
	Mu      sync.Mutex
	Counter int
	Quit    chan bool
}
