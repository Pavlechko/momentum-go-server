package handlers

import (
	"momentum-go-server/internal/services"
	"momentum-go-server/internal/utils"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	userId := utils.GetUserId(r)

	Response := services.GetData(userId)

	WriteJSON(w, http.StatusOK, Response)
}
