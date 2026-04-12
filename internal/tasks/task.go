package tasks

import "time"

const (
	StatusTodo       = "todo"
	StatusInProgress = "in_progress"
	StatusDone       = "done"
	StatusArchived   = "archived"

	PriorityLow    = "low"
	PriorityMedium = "medium"
	PriorityHigh   = "high"
)

type Task struct {
	ID           int64      `json:"id"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Status       string     `json:"status"`
	Priority     string     `json:"priority"`
	DueAt        *time.Time `json:"due_at"`
	CompletedAt  *time.Time `json:"completed_at"`
	ParentTaskID *int64     `json:"parent_task_id"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type CreateTaskInput struct {
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Status       string     `json:"status"`
	Priority     string     `json:"priority"`
	DueAt        *time.Time `json:"due_at"`
	ParentTaskID *int64     `json:"parent_task_id"`
}

type UpdateTaskInput struct {
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Status       string     `json:"status"`
	Priority     string     `json:"priority"`
	DueAt        *time.Time `json:"due_at"`
	ParentTaskID *int64     `json:"parent_task_id"`
}

func NormalizeCreateInput(input *CreateTaskInput) {
	if input.Status == "" {
		input.Status = StatusTodo
	}
	if input.Priority == "" {
		input.Priority = PriorityMedium
	}
}

func NormalizeUpdateInput(input *UpdateTaskInput) {
	if input.Status == "" {
		input.Status = StatusTodo
	}
	if input.Priority == "" {
		input.Priority = PriorityMedium
	}
}

func ValidStatus(status string) bool {
	switch status {
	case StatusTodo, StatusInProgress, StatusDone, StatusArchived:
		return true
	default:
		return false
	}
}

func ValidPriority(priority string) bool {
	switch priority {
	case PriorityLow, PriorityMedium, PriorityHigh:
		return true
	default:
		return false
	}
}
