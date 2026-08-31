package crud

import (
	"gorm.io/gorm"
)

type CRUDController[T any, S any, R any] struct {
	DB *gorm.DB
	Preloads []string
}

func NewCRUDController[T any, S any, R any](db *gorm.DB, preloads ...string) *CRUDController[T, S, R] {
	return &CRUDController[T, S, R]{DB: db, Preloads: preloads,}
}