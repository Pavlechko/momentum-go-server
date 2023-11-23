package handlers

import (
	"net/http"

	"momentum-go-server/internal/services"
	"momentum-go-server/internal/utils"
)

func Home(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserID(r)
	res := &services.Res{
		Counter: 0,
		Quit:    make(chan bool),
	}

	Response := res.GetData(userID)

	err := WriteJSON(w, http.StatusOK, Response)
	if err != nil {
		utils.ErrorLogger.Printf("Error write JSON %s", err.Error())
	}
}
