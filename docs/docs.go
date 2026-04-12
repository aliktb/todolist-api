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
      "get": {"summary": "Liveness probe", "tags": ["health"], "responses": {"200": {"description": "OK"}}}
    },
    "/health/ready": {
      "get": {"summary": "Readiness probe", "tags": ["health"], "responses": {"200": {"description": "Ready"}, "503": {"description": "Database unavailable"}}}
    },
    "/api/v1/tasks": {
      "get": {
        "summary": "List tasks",
        "tags": ["tasks"],
        "produces": ["application/json"],
        "responses": {"200": {"description": "OK", "schema": {"type": "array", "items": {"$ref": "#/definitions/tasks.Task"}}}}
      },
      "post": {
        "summary": "Create a task",
        "tags": ["tasks"],
        "consumes": ["application/json"],
        "produces": ["application/json"],
        "parameters": [{"in": "body", "name": "payload", "required": true, "schema": {"$ref": "#/definitions/tasks.CreateTaskInput"}}],
        "responses": {"201": {"description": "Created", "schema": {"$ref": "#/definitions/tasks.Task"}}}
      }
    },
    "/api/v1/tasks/search": {
      "get": {
        "summary": "Search tasks",
        "tags": ["tasks"],
        "produces": ["application/json"],
        "parameters": [{"type": "string", "in": "query", "name": "q", "required": true}],
        "responses": {"200": {"description": "OK", "schema": {"type": "array", "items": {"$ref": "#/definitions/tasks.TaskSummary"}}}}
      }
    },
    "/api/v1/tasks/{taskID}": {
      "get": {
        "summary": "Get a task",
        "tags": ["tasks"],
        "produces": ["application/json"],
        "parameters": [{"type": "integer", "in": "path", "name": "taskID", "required": true}],
        "responses": {"200": {"description": "OK", "schema": {"$ref": "#/definitions/tasks.Task"}}}
      },
      "put": {
        "summary": "Update a task",
        "tags": ["tasks"],
        "consumes": ["application/json"],
        "produces": ["application/json"],
        "parameters": [
          {"type": "integer", "in": "path", "name": "taskID", "required": true},
          {"in": "body", "name": "payload", "required": true, "schema": {"$ref": "#/definitions/tasks.UpdateTaskInput"}}
        ],
        "responses": {"200": {"description": "OK", "schema": {"$ref": "#/definitions/tasks.Task"}}}
      },
      "delete": {
        "summary": "Delete a task",
        "tags": ["tasks"],
        "parameters": [{"type": "integer", "in": "path", "name": "taskID", "required": true}],
        "responses": {"204": {"description": "No Content"}}
      }
    },
    "/api/v1/tasks/{taskID}/history": {
      "get": {
        "summary": "Get task history",
        "tags": ["tasks"],
        "produces": ["application/json"],
        "parameters": [{"type": "integer", "in": "path", "name": "taskID", "required": true}],
        "responses": {"200": {"description": "OK", "schema": {"type": "array", "items": {"$ref": "#/definitions/tasks.TaskEvent"}}}}
      }
    },
    "/api/v1/tasks/{taskID}/links": {
      "get": {
        "summary": "List task links",
        "tags": ["tasks"],
        "produces": ["application/json"],
        "parameters": [{"type": "integer", "in": "path", "name": "taskID", "required": true}],
        "responses": {"200": {"description": "OK", "schema": {"type": "array", "items": {"$ref": "#/definitions/tasks.TaskLink"}}}}
      },
      "post": {
        "summary": "Create a task link",
        "tags": ["tasks"],
        "consumes": ["application/json"],
        "produces": ["application/json"],
        "parameters": [
          {"type": "integer", "in": "path", "name": "taskID", "required": true},
          {"in": "body", "name": "payload", "required": true, "schema": {"$ref": "#/definitions/tasks.CreateTaskLinkInput"}}
        ],
        "responses": {"201": {"description": "Created", "schema": {"$ref": "#/definitions/tasks.TaskLink"}}}
      }
    },
    "/api/v1/tasks/{taskID}/links/{linkID}": {
      "delete": {
        "summary": "Delete a task link",
        "tags": ["tasks"],
        "parameters": [
          {"type": "integer", "in": "path", "name": "taskID", "required": true},
          {"type": "integer", "in": "path", "name": "linkID", "required": true}
        ],
        "responses": {"204": {"description": "No Content"}}
      }
    }
  },
  "definitions": {
    "tasks.CreateTaskInput": {
      "type": "object",
      "required": ["title"],
      "properties": {
        "title": {"type": "string"},
        "description": {"type": "string"},
        "status": {"type": "string"},
        "priority": {"type": "string"},
        "due_at": {"type": "string", "format": "date-time"},
        "parent_task_id": {"type": "integer"}
      }
    },
    "tasks.UpdateTaskInput": {
      "type": "object",
      "required": ["title"],
      "properties": {
        "title": {"type": "string"},
        "description": {"type": "string"},
        "status": {"type": "string"},
        "priority": {"type": "string"},
        "due_at": {"type": "string", "format": "date-time"},
        "parent_task_id": {"type": "integer"}
      }
    },
    "tasks.Task": {
      "type": "object",
      "properties": {
        "id": {"type": "integer"},
        "title": {"type": "string"},
        "description": {"type": "string"},
        "status": {"type": "string"},
        "priority": {"type": "string"},
        "due_at": {"type": "string", "format": "date-time"},
        "completed_at": {"type": "string", "format": "date-time"},
        "parent_task_id": {"type": "integer"},
        "created_at": {"type": "string", "format": "date-time"},
        "updated_at": {"type": "string", "format": "date-time"}
      }
    },
    "tasks.TaskSummary": {
      "type": "object",
      "properties": {
        "id": {"type": "integer"},
        "title": {"type": "string"},
        "status": {"type": "string"},
        "priority": {"type": "string"},
        "due_at": {"type": "string", "format": "date-time"}
      }
    },
    "tasks.TaskEvent": {
      "type": "object",
      "properties": {
        "id": {"type": "integer"},
        "task_id": {"type": "integer"},
        "event_type": {"type": "string"},
        "field_name": {"type": "string"},
        "old_value": {"type": "object"},
        "new_value": {"type": "object"},
        "created_at": {"type": "string", "format": "date-time"}
      }
    },
    "tasks.CreateTaskLinkInput": {
      "type": "object",
      "required": ["target_task_id"],
      "properties": {
        "target_task_id": {"type": "integer"},
        "link_type": {"type": "string"}
      }
    },
    "tasks.TaskLink": {
      "type": "object",
      "properties": {
        "id": {"type": "integer"},
        "source_task_id": {"type": "integer"},
        "target_task_id": {"type": "integer"},
        "link_type": {"type": "string"},
        "created_at": {"type": "string", "format": "date-time"},
        "linked_task": {"$ref": "#/definitions/tasks.TaskSummary"}
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
