package models

type Card struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Column   string `json:"column"` // 'todo', 'doing', 'done'
	EditedAt string `json:"edited_at"`
}
