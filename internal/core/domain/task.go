package domain

import "github.com/uptrace/bun"

type Task struct {
	bun.BaseModel `bun:"task"`
	ID            uint   `json:"id" bun:"id,pk,autoincrement"`
	Name          string `json:"name" bun:"name,notnull"`
	Desc          string `json:"description" bun:"description"`
}
