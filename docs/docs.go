package docs

import "github.com/swaggo/swag"

const doc = `{
  "swagger": "2.0",
  "info": {
    "title": "Todo List API",
    "description": "REST API for managing todo tasks.",
    "version": "1.0"
  },
  "basePath": "/",
  "schemes": ["http", "https"],
  "paths": {
    "/health/live": {
      "get": {
        "summary": "Liveness probe",
        "tags": ["health"],
        "responses": {
          "200": {"description": "OK"}
        }
      }
    },
    "/health/ready": {
      "get": {
        "summary": "Readiness probe",
        "tags": ["health"],
        "responses": {
          "200": {"description": "Ready"},
          "503": {"description": "Database unavailable"}
        }
      }
    },
    "/api/v1/tasks": {
      "get": {
        "summary": "List tasks",
        "tags": ["tasks"],
        "produces": ["application/json"],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "type": "array",
              "items": {"$ref": "#/definitions/tasks.Task"}
            }
          }
        }
      },
      "post": {
        "summary": "Create a task",
        "tags": ["tasks"],
        "consumes": ["application/json"],
        "produces": ["application/json"],
        "parameters": [
          {
            "in": "body",
            "name": "payload",
            "required": true,
            "schema": {"$ref": "#/definitions/tasks.CreateTaskInput"}
          }
        ],
        "responses": {
          "201": {
            "description": "Created",
            "schema": {"$ref": "#/definitions/tasks.Task"}
          }
        }
      }
    },
    "/api/v1/tasks/{taskID}": {
      "get": {
        "summary": "Get a task",
        "tags": ["tasks"],
        "produces": ["application/json"],
        "parameters": [
          {
            "type": "integer",
            "in": "path",
            "name": "taskID",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {"$ref": "#/definitions/tasks.Task"}
          }
        }
      },
      "put": {
        "summary": "Update a task",
        "tags": ["tasks"],
        "consumes": ["application/json"],
        "produces": ["application/json"],
        "parameters": [
          {
            "type": "integer",
            "in": "path",
            "name": "taskID",
            "required": true
          },
          {
            "in": "body",
            "name": "payload",
            "required": true,
            "schema": {"$ref": "#/definitions/tasks.UpdateTaskInput"}
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {"$ref": "#/definitions/tasks.Task"}
          }
        }
      },
      "delete": {
        "summary": "Delete a task",
        "tags": ["tasks"],
        "parameters": [
          {
            "type": "integer",
            "in": "path",
            "name": "taskID",
            "required": true
          }
        ],
        "responses": {
          "204": {"description": "No Content"}
        }
      }
    }
  },
  "definitions": {
    "tasks.CreateTaskInput": {
      "type": "object",
      "required": ["title"],
      "properties": {
        "title": {"type": "string", "example": "Write React integration"},
        "description": {"type": "string", "example": "Hook the frontend up to the API"},
        "status": {"type": "string", "example": "todo"},
        "priority": {"type": "string", "example": "high"},
        "due_at": {"type": "string", "format": "date-time", "example": "2026-04-15T09:00:00Z"},
        "parent_task_id": {"type": "integer", "example": 1}
      }
    },
    "tasks.UpdateTaskInput": {
      "type": "object",
      "required": ["title"],
      "properties": {
        "title": {"type": "string", "example": "Write React integration"},
        "description": {"type": "string", "example": "Frontend now calls the API"},
        "status": {"type": "string", "example": "in_progress"},
        "priority": {"type": "string", "example": "medium"},
        "due_at": {"type": "string", "format": "date-time", "example": "2026-04-16T17:00:00Z"},
        "parent_task_id": {"type": "integer", "example": 1}
      }
    },
    "tasks.Task": {
      "type": "object",
      "properties": {
        "id": {"type": "integer", "example": 2},
        "title": {"type": "string", "example": "Write React integration"},
        "description": {"type": "string", "example": "Frontend now calls the API"},
        "status": {"type": "string", "example": "done"},
        "priority": {"type": "string", "example": "high"},
        "due_at": {"type": "string", "format": "date-time", "example": "2026-04-16T17:00:00Z"},
        "completed_at": {"type": "string", "format": "date-time", "example": "2026-04-14T12:30:00Z"},
        "parent_task_id": {"type": "integer", "example": 1},
        "created_at": {"type": "string", "format": "date-time", "example": "2026-04-12T10:00:00Z"},
        "updated_at": {"type": "string", "format": "date-time", "example": "2026-04-14T12:30:00Z"}
      }
    }
  }
}`

type swaggerDoc struct{}

func (swaggerDoc) ReadDoc() string {
	return doc
}

func init() {
	swag.Register(swag.Name, swaggerDoc{})
}
