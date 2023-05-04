package controllers

import (
	"encoding/json"
	"net/http"
)

func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	row := db.QueryRow("SELECT * FROM users WHERE id_user = ?", GetOnlineUserId(r))

	var user User
	err := row.Scan(&user.Id_User, &user.Name, &user.Address, &user.Email, &user.Password)

	var response UserResponse

	if err == nil {
		response.Status = 200
		response.Message = "Success"
		response.Data = user
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		PrintError(400, "User Not Found", w)
		return
	}
}
