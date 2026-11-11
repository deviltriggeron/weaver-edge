package handler

import (
	"encoding/json"
	"net/http"

	"weaver-edge/internal/application"
	tgClient "weaver-edge/internal/infrastructure/connectors/telegram/dto"
	"weaver-edge/internal/transport/dto"
)

type Handler struct {
	application application.Application
}

func NewHandler(application application.Application) Handler {
	return Handler{
		application: application,
	}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, dto.StatusResponse{
		Status: "ok",
	})
}

func (h *Handler) TelegramWebhook(w http.ResponseWriter, r *http.Request) {
	var req tgClient.TelegramUpdate

	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid request",
		})
		return
	}

	if req.Message == nil {
		writeJSON(w, http.StatusOK, dto.StatusResponse{
			Status: "ignored",
		})
		return
	}

	msg := tgClient.MapMessage(req.Message)

	err := h.application.Handle(r.Context(), msg)
	if err != nil {
		// add error handling
		writeJSON(w, http.StatusInternalServerError, dto.ErrorResponse{
			Error: "internal_error",
		})
		return
	}

	writeJSON(w, http.StatusOK, dto.StatusResponse{
		Status: "OK",
	})
}

func (h *Handler) Context(w http.ResponseWriter, r *http.Request) {
	var req dto.ContextRequest

	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid request",
		})
		return
	}

	contextRequest := dto.ContextReqToDomain(req)

	resp, err := h.application.GetContext(r.Context(), contextRequest)
	if err != nil {
		// add error handling
		writeJSON(w, http.StatusInternalServerError, dto.ErrorResponse{
			Error: "internal_error",
		})
		return
	}

	response := dto.ContextRespDomainToResp(resp)

	writeJSON(w, http.StatusOK, response)
}

func decodeJSON(r *http.Request, dst any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(dst)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
