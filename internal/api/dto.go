package api

import (
	"encoding/json"
	"net/http"
	"taskflow/internal/domain"
	"time"
)

type CreateTaskRequest struct {
	Title string `json:"title"`
}

type CompleteTaskRequest struct {
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type TaskResponse struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	Completed  bool      `json:"completed"`
	Created_at time.Time `json:"created_at"`
}

func toTaskResponse(t domain.Task) TaskResponse {
	return TaskResponse{
		ID:         t.ID,
		Title:      t.Title,
		Completed:  t.Completed,
		Created_at: t.CreatedAt,
	}
}

func writeJSONResponse(w http.ResponseWriter, status int, t domain.Task) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(toTaskResponse(t))

	if err != nil {
		return err
	}
	return nil
}

func writeJSONResponseAll(w http.ResponseWriter, status int, t []domain.Task) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := make([]TaskResponse, 0, len(t))
	for _, t := range t {
		resp = append(resp, toTaskResponse(t))
	}

	err := json.NewEncoder(w).Encode(resp)

	if err != nil {
		return err
	}
	return nil
}
