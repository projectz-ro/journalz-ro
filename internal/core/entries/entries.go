package core

import "time"

type Entry struct {
	ID       int64     `json:"id"`
	Title    string    `json:"title"`
	Tags     []string  `json:"tags"`
	Path     string    `json:"path"`
	LastRead time.Time `json:"last_read"`
}
