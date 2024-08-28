package models

type Base struct {
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
	Deleted   bool   `db:"deleted"`
}

type IDResponse struct {
	ID int64 `json:"id"`
}

type SuccessResponse struct {
	Success bool `json:"success"`
}
