package models

// Task represents a task in the system
type Task struct {
	ID          string `json:"id" bson:"_id,omitempty"`
	Title       string `json:"title" bson:"title" binding:"required"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
	DueDate     string `json:"due_date,omitempty" bson:"due_date,omitempty"` // e.g. "2025-11-25"
	Status      string `json:"status,omitempty" bson:"status,omitempty"`     // e.g. "pending", "in-progress", "done"
}
