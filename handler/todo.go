package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

// ServeHTTP inplements http.Handler interface
func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var reqBody model.CreateTODORequest

		err := json.NewDecoder(r.Body).Decode(&reqBody)

		if err != nil {
			log.Println("json Decode Error:", err)
			http.Error(w, "リクエストが不正です", http.StatusBadRequest)
			return
		}

		if reqBody.Subject == "" {
			http.Error(w, "SUbjectが空です", http.StatusBadRequest)
			return
		} else {
			todo, err := h.svc.CreateTODO(r.Context(), reqBody.Subject, reqBody.Description)
			fmt.Println("todoの値:", todo)

			if err != nil {
				log.Println("Create TODO Error:", err)
				http.Error(w, "TODO作成に失敗しました", http.StatusInternalServerError)
				return
			}

			resBody := map[string]interface{}{
				"todo": todo,
			}

			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(resBody)

			if err != nil {
				log.Println("json Encode Error:", err)
				http.Error(w, "レスポンスのエンコードに失敗しました", http.StatusInternalServerError)
				return
			}

		}

	} else if r.Method == "PUT" {
		var reqBody model.UpdateTODORequest

		err := json.NewDecoder(r.Body).Decode(&reqBody)

		if err != nil {
			log.Println("json Decode ERror:", err)
			http.Error(w, "リクエストが不正です", http.StatusBadRequest)
			return
		}

		if reqBody.ID == 0 || reqBody.Subject == "" {
			http.Error(w, "IDが0またはSubjectが空です", http.StatusBadRequest)
			return
		}

		todo, err := h.svc.UpdateTODO(r.Context(), reqBody.ID, reqBody.Subject, reqBody.Description)

		if err != nil {
			if e, ok := err.(*model.ErrNotFound); ok {
				http.Error(w, e.Error(), http.StatusNotFound)
				return
			}
		}

		resBody := map[string]interface{}{
			"todo": todo,
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(resBody)

		if err != nil {
			log.Println("json Encode Error:", err)
			http.Error(w, "レスポンスのエンコードに失敗しました", http.StatusInternalServerError)
			return
		}

	} else if r.Method == "GET" {
		var reqParams model.ReadTODORequest

		queryParams := r.URL.Query()

		if prevIDStr := queryParams.Get("prev_id"); prevIDStr == "" {
			reqParams.PrevID = 0
		} else {
			prevID, err := strconv.ParseInt(prevIDStr, 10, 64)
			if err != nil {
				fmt.Println("Invalid prev_id:", err)
				http.Error(w, "invalid prev_id", http.StatusBadRequest)
				return
			} else {
				reqParams.PrevID = prevID
			}
		}

		if sizeStr := queryParams.Get("size"); sizeStr == "" {
			reqParams.Size = 5
		} else {
			size, err := strconv.ParseInt(sizeStr, 10, 64)
			if err != nil {
				fmt.Println("Invalid size:", err)
				http.Error(w, "Invalid size", http.StatusBadRequest)
				return
			} else {
				reqParams.Size = size
			}

		}

		todos, err := h.svc.ReadTODO(r.Context(), reqParams.PrevID, reqParams.Size)

		if err != nil {
			fmt.Println("TODOの読み込みに失敗しました:", err)
			http.Error(w, "ReadTODO error", http.StatusInternalServerError)
			return
		}

		resBody := map[string]interface{}{
			"todos": todos,
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(resBody)

		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}

	} else if r.Method == "DELETE" {
		var reqBody model.DeleteTODORequest
		err := json.NewDecoder(r.Body).Decode(&reqBody)

		if err != nil {
			fmt.Println("Decode Error:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if len(reqBody.IDs) == 0 {
			fmt.Println("IDが指定されていません")
			http.Error(w, "IDが指定されていません", http.StatusBadRequest)
			return
		}

		err = h.svc.DeleteTODO(r.Context(), reqBody.IDs)

		if err != nil {
			if e, ok := err.(*model.ErrNotFound); ok {
				fmt.Println("Not found Error:", e)
				http.Error(w, "Not found Error", http.StatusNotFound)
				return
			} else {
				fmt.Println("DeleteTODO Error:", e)
				http.Error(w, "Delete TODO Error", http.StatusInternalServerError)
				return
			}
		}

		resBody := map[string]interface{}{}

		err = json.NewEncoder(w).Encode(resBody)

		if err != nil {
			fmt.Println("Encode Error", err)
			http.Error(w, "Encode Error", http.StatusInternalServerError)
			return
		}

		return

	}

}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	_, _ = h.svc.CreateTODO(ctx, "", "")
	return &model.CreateTODOResponse{}, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	_, _ = h.svc.ReadTODO(ctx, 0, 0)
	return &model.ReadTODOResponse{}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	_, _ = h.svc.UpdateTODO(ctx, 0, "", "")
	return &model.UpdateTODOResponse{}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	_ = h.svc.DeleteTODO(ctx, nil)
	return &model.DeleteTODOResponse{}, nil
}
