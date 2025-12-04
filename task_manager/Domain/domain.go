package Domain

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type Task struct {
	ID          string `json:"id" bson:"_id,omitempty"`
	Title       string `json:"title" bson:"title"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
	DueDate     string `json:"due_date,omitempty" bson:"due_date,omitempty"`
	Status      string `json:"status,omitempty" bson:"status,omitempty"`
}

type User struct {
	ID           string `json:"id" bson:"_id,omitempty"`
	Username     string `json:"username" bson:"username"`
	PasswordHash string `json:"-" bson:"password_hash"`
	Role         Role   `json:"role" bson:"role"`
}
