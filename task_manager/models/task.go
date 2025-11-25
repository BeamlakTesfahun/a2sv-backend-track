package models

// Task represents a task in the system
type Task struct {
	ID          int64  `json:"id"`
	Title       string `json:"title" binding:"required"` // required field
	Description string `json:"description,omitempty"`
	DueDate     string `json:"due_date,omitempty"` // 2025-11-25
	Status      string `json:"status,omitempty"`   // "pending", "in-progress", "done"
}
