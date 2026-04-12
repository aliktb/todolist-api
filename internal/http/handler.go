package http

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "github.com/ali/todolist-api/docs"
	"github.com/ali/todolist-api/internal/config"
	"github.com/ali/todolist-api/internal/http/middleware"
	"github.com/ali/todolist-api/internal/response"
	"github.com/ali/todolist-api/internal/tasks"
)

type Handler struct {
	cfg    config.Config
	logger *slog.Logger
	repo   *tasks.Repository
	db     *pgxpool.Pool
}

func NewHandler(cfg config.Config, logger *slog.Logger, repo *tasks.Repository, db *pgxpool.Pool) *Handler {
	return &Handler{
		cfg:    cfg,
		logger: logger,
		repo:   repo,
		db:     db,
	}
}

func (h *Handler) Router() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestLogger(h.logger))
	r.Use(middleware.CORS(h.cfg.AllowedOrigin))

	r.Get("/health/live", h.liveness)
	r.Get("/health/ready", h.readiness)
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	r.Route("/api/v1/tasks", func(r chi.Router) {
		r.Get("/", h.listTasks)
		r.Post("/", h.createTask)
		r.Get("/{taskID}", h.getTask)
		r.Put("/{taskID}", h.updateTask)
		r.Delete("/{taskID}", h.deleteTask)
	})

	return r
}

// liveness godoc
// @Summary Liveness probe
// @Tags health
// @Success 200 {object} map[string]string
// @Router /health/live [get]
func (h *Handler) liveness(w http.ResponseWriter, _ *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// readiness godoc
// @Summary Readiness probe
// @Tags health
// @Success 200 {object} map[string]string
// @Failure 503 {object} map[string]string
// @Router /health/ready [get]
func (h *Handler) readiness(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.ReadinessTimeout)
	defer cancel()

	if err := h.db.Ping(ctx); err != nil {
		response.Error(w, http.StatusServiceUnavailable, "database unavailable")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

// listTasks godoc
// @Summary List tasks
// @Tags tasks
// @Produce json
// @Success 200 {array} tasks.Task
// @Failure 500 {object} map[string]string
// @Router /api/v1/tasks [get]
func (h *Handler) listTasks(w http.ResponseWriter, r *http.Request) {
	allTasks, err := h.repo.List(r.Context())
	if err != nil {
		h.logger.Error("failed to list tasks", "error", err)
		response.Error(w, http.StatusInternalServerError, "failed to list tasks")
		return
	}

	response.JSON(w, http.StatusOK, allTasks)
}

// getTask godoc
// @Summary Get a task
// @Tags tasks
// @Produce json
// @Param taskID path int true "Task ID"
// @Success 200 {object} tasks.Task
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/tasks/{taskID} [get]
func (h *Handler) getTask(w http.ResponseWriter, r *http.Request) {
	id, ok := parseTaskID(w, r)
	if !ok {
		return
	}

	task, err := h.repo.Get(r.Context(), id)
	if errors.Is(err, tasks.ErrNotFound) {
		response.Error(w, http.StatusNotFound, "task not found")
		return
	}
	if err != nil {
		h.logger.Error("failed to get task", "task_id", id, "error", err)
		response.Error(w, http.StatusInternalServerError, "failed to get task")
		return
	}

	response.JSON(w, http.StatusOK, task)
}

// createTask godoc
// @Summary Create a task
// @Tags tasks
// @Accept json
// @Produce json
// @Param payload body tasks.CreateTaskInput true "Task payload"
// @Success 201 {object} tasks.Task
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/tasks [post]
func (h *Handler) createTask(w http.ResponseWriter, r *http.Request) {
	var input tasks.CreateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	tasks.NormalizeCreateInput(&input)
	if message := validateCreateInput(input); message != "" {
		response.Error(w, http.StatusBadRequest, message)
		return
	}

	task, err := h.repo.Create(r.Context(), input)
	if errors.Is(err, tasks.ErrInvalidParent) {
		response.Error(w, http.StatusBadRequest, "parent task does not exist")
		return
	}
	if err != nil {
		h.logger.Error("failed to create task", "error", err)
		response.Error(w, http.StatusInternalServerError, "failed to create task")
		return
	}

	response.JSON(w, http.StatusCreated, task)
}

// updateTask godoc
// @Summary Update a task
// @Tags tasks
// @Accept json
// @Produce json
// @Param taskID path int true "Task ID"
// @Param payload body tasks.UpdateTaskInput true "Task payload"
// @Success 200 {object} tasks.Task
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/tasks/{taskID} [put]
func (h *Handler) updateTask(w http.ResponseWriter, r *http.Request) {
	id, ok := parseTaskID(w, r)
	if !ok {
		return
	}

	var input tasks.UpdateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	tasks.NormalizeUpdateInput(&input)
	if message := validateUpdateInput(id, input); message != "" {
		response.Error(w, http.StatusBadRequest, message)
		return
	}

	task, err := h.repo.Update(r.Context(), id, input)
	if errors.Is(err, tasks.ErrInvalidParent) {
		response.Error(w, http.StatusBadRequest, "parent task does not exist")
		return
	}
	if errors.Is(err, tasks.ErrNotFound) {
		response.Error(w, http.StatusNotFound, "task not found")
		return
	}
	if err != nil {
		h.logger.Error("failed to update task", "task_id", id, "error", err)
		response.Error(w, http.StatusInternalServerError, "failed to update task")
		return
	}

	response.JSON(w, http.StatusOK, task)
}

// deleteTask godoc
// @Summary Delete a task
// @Tags tasks
// @Param taskID path int true "Task ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/tasks/{taskID} [delete]
func (h *Handler) deleteTask(w http.ResponseWriter, r *http.Request) {
	id, ok := parseTaskID(w, r)
	if !ok {
		return
	}

	err := h.repo.Delete(r.Context(), id)
	if errors.Is(err, tasks.ErrNotFound) {
		response.Error(w, http.StatusNotFound, "task not found")
		return
	}
	if err != nil {
		h.logger.Error("failed to delete task", "task_id", id, "error", err)
		response.Error(w, http.StatusInternalServerError, "failed to delete task")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseTaskID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(chi.URLParam(r, "taskID"), 10, 64)
	if err != nil || id < 1 {
		response.Error(w, http.StatusBadRequest, "invalid task id")
		return 0, false
	}
	return id, true
}

func validateCreateInput(input tasks.CreateTaskInput) string {
	if input.Title == "" {
		return "title is required"
	}
	if !tasks.ValidStatus(input.Status) {
		return "status must be one of: todo, in_progress, done, archived"
	}
	if !tasks.ValidPriority(input.Priority) {
		return "priority must be one of: low, medium, high"
	}
	if input.ParentTaskID != nil && *input.ParentTaskID < 1 {
		return "parent_task_id must be a positive integer"
	}
	return ""
}

func validateUpdateInput(taskID int64, input tasks.UpdateTaskInput) string {
	if input.Title == "" {
		return "title is required"
	}
	if !tasks.ValidStatus(input.Status) {
		return "status must be one of: todo, in_progress, done, archived"
	}
	if !tasks.ValidPriority(input.Priority) {
		return "priority must be one of: low, medium, high"
	}
	if input.ParentTaskID != nil {
		if *input.ParentTaskID < 1 {
			return "parent_task_id must be a positive integer"
		}
		if *input.ParentTaskID == taskID {
			return "parent_task_id cannot reference the task itself"
		}
	}
	return ""
}
