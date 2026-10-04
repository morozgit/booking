package handlers

import (
	"booking/internal/dto"
	"booking/internal/service"
	"encoding/json"
	"net/http"
)

type HotelsHandler struct {
	service *service.HotelsService
}

func NewHotelsHandler(service *service.HotelsService) *HotelsHandler {
	return &HotelsHandler{service: service}
}

// CreateHotel @Summary      Create Hotel
// @Description  Add Hotel
// @Tags         hotels
// @Produce      json
// @Success      200 {object} dto.StatusResponse
// @Param request body dto.HotelAdd true "Registration data"
// @Router       /hotels/ [post]
func (h *HotelsHandler) CreteHotel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req dto.HotelAdd

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := h.service.AddHotel(ctx, req); err != nil {

		http.Error(w, "Hotel already exists", http.StatusConflict)
		return

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]string{
		"status": "OK",
	})
}
