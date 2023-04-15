package v1

import (
	"shortener/db"
)

type TaskServerV1 struct {
	PgContext db.PgCaller
}

// NewTaskServerV1 creates a new TaskServerV1 with created database connection
func NewTaskServerV1(pgContext db.PgCaller) *TaskServerV1 {
	return &TaskServerV1{
		PgContext: pgContext,
	}
}
