package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"taskflow/internal/domain"
	"taskflow/internal/middleware"
	"taskflow/internal/service"
)

type Handler struct {
	service *service.TaskService
	limiter *middleware.RateLimiter
}

func NewHandler(s *service.TaskService, limiter *middleware.RateLimiter) *Handler {
	return &Handler{
		service: s,
		limiter: limiter,
	}
}

func parseID(path string) (int, error) {
	pathSlice := strings.Split(path, "/")
	var taskIDstr string
	taskErr := errors.New("incorrect path")
	if len(pathSlice) == 3 {
		taskIDstr = pathSlice[2]
	} else {
		return 0, taskErr
	}

	taskID, err := strconv.Atoi(taskIDstr)
	if err != nil {
		return 0, taskErr
	}
	return taskID, nil
}

func (h *Handler) HandleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodPost:
		apiKey, ok := middleware.GetAPIKeyFromContext(r.Context())
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		allowed, err := h.limiter.Allow(r.Context(), "rate:key"+apiKey+":task_create")
		if err != nil {
			http.Error(w, "rate limiter unavailable", http.StatusServiceUnavailable)
			return
		}
		if !allowed {
			http.Error(w, "too many requests", http.StatusTooManyRequests)
			return
		}

		var req CreateTaskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid JSON")
			return
		}

		if req.Title == "" {
			writeJSONError(w, http.StatusBadRequest, "empty title")
			return
		}

		task, err := h.service.CreateTask(r.Context(), req.Title)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "server error")
			return
		}

		if err := writeJSONResponse(w, http.StatusCreated, task); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "server error")
			return
		}

	case http.MethodGet:
		tasks, err := h.service.GetTasks(r.Context())
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "server error")
			return
		}

		if err := writeJSONResponseAll(w, http.StatusOK, tasks); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "server error")
			return
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) HandleTasksByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch r.Method {

	case http.MethodGet:
		task, err := h.service.GetTaskByID(r.Context(), id)
		if errors.Is(err, domain.ErrNotFound) {
			writeJSONError(w, http.StatusNotFound, "task not found")
			return
		}
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "server error")
			return
		}
		if err := writeJSONResponse(w, http.StatusCreated, task); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "server error")
			return
		}

	case http.MethodPut:
		var req CompleteTaskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid JSON")
			return
		}

		task, err := h.service.UpdateTask(r.Context(), id, req.Title, req.Completed)
		if errors.Is(err, domain.ErrNotFound) {
			writeJSONError(w, http.StatusNotFound, "task not found")
			return
		}
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "server error")
			return
		}

		if err := writeJSONResponse(w, http.StatusCreated, task); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "server error")
			return
		}

	case http.MethodDelete:
		err := h.service.DeleteTask(r.Context(), id)
		if errors.Is(err, domain.ErrNotFound) {
			writeJSONError(w, http.StatusNotFound, "task not found")
			return
		}
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "server error")
			return
		}

		w.WriteHeader(http.StatusNoContent)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
