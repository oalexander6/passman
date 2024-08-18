package models

import "time"

type base struct {
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Deleted   bool      `db:"deleted"`
}

type IDResponse struct {
	ID int `json:"id"`
}

type SuccessResponse struct {
	Success bool `json:"success"`
}
