package tasks

import (
	"encoding/json"
	"time"
)

const (
	StatusTodo       = "todo"
	StatusInProgress = "in_progress"
	StatusDone       = "done"
	StatusArchived   = "archived"

	PriorityLow    = "low"
	PriorityMedium = "medium"
	PriorityHigh   = "high"

	LinkTypeRelatesTo  = "relates_to"
	LinkTypeBlocks     = "blocks"
	LinkTypeDuplicates = "duplicates"
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

type TaskSummary struct {
	ID       int64      `json:"id"`
	Title    string     `json:"title"`
	Status   string     `json:"status"`
	Priority string     `json:"priority"`
	DueAt    *time.Time `json:"due_at"`
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

type TaskEvent struct {
	ID        int64           `json:"id"`
	TaskID    int64           `json:"task_id"`
	EventType string          `json:"event_type"`
	FieldName string          `json:"field_name,omitempty"`
	OldValue  json.RawMessage `json:"old_value,omitempty"`
	NewValue  json.RawMessage `json:"new_value,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
}

type CreateTaskLinkInput struct {
	TargetTaskID int64  `json:"target_task_id"`
	LinkType     string `json:"link_type"`
}

type TaskLink struct {
	ID           int64       `json:"id"`
	SourceTaskID int64       `json:"source_task_id"`
	TargetTaskID int64       `json:"target_task_id"`
	LinkType     string      `json:"link_type"`
	CreatedAt    time.Time   `json:"created_at"`
	LinkedTask   TaskSummary `json:"linked_task"`
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

func NormalizeCreateTaskLinkInput(input *CreateTaskLinkInput) {
	if input.LinkType == "" {
		input.LinkType = LinkTypeRelatesTo
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

func ValidLinkType(linkType string) bool {
	switch linkType {
	case LinkTypeRelatesTo, LinkTypeBlocks, LinkTypeDuplicates:
		return true
	default:
		return false
	}
}
