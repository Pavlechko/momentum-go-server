package handlers

import (
	"net/http"

	"momentum-go-server/internal/utils"
)

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserID(r)

	utils.InfoLogger.Println("userID:", userID)
	Response := h.Service.GetData(userID)

	err := h.WriteJSON(w, http.StatusOK, Response)
	if err != nil {
		utils.ErrorLogger.Printf("Error write JSON %s", err.Error())
	}
}
