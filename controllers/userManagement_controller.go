package controllers

import (
	"encoding/json"
	"net/http"
)

func SeeUsers(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	var response UsersResponse
	var user User
	var users []User

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		PrintError(400, "Table Not Found", w)
		return
	}

	for rows.Next() {
		if err := rows.Scan(&user.Id_User, &user.Name, &user.Address, &user.Email, &user.Password); err != nil {
			PrintError(400, "Field Undefined", w)
			return
		} else {
			users = append(users, user)

		}
	}

	if len(users) != 0 {
		response.Status = 200
		response.Message = "Success"
		response.Data = users
	} else {
		PrintError(400, "No User", w)
	}
	json.NewEncoder(w).Encode(response)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	var response UsersResponse

	err := r.ParseForm()
	if err != nil {
		return
	}

	idUser := r.Form.Get("id_user")
	name := r.Form.Get("name")
	address := r.Form.Get("address")
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	if len(idUser) <= 0 || len(name) <= 0 || len(address) <= 0 || len(email) <= 0 || len(password) <= 0 {
		PrintError(400, "Please Input All Field", w)
		return
	}

	result, errQuery := db.Exec("UPDATE users SET name=?, address=?, email=?, password=? WHERE id_user=?", name, address, email, password, idUser)
	if errQuery != nil {
		PrintError(400, "No Table Found", w)
		return
	} else {
		num, _ := result.RowsAffected()
		if num == 0 {
			PrintError(400, "User Not Found", w)
		} else {
			w.WriteHeader(http.StatusOK)
			response.Status = 200
			response.Message = "Data Updated!"
			json.NewEncoder(w).Encode(response)
		}
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	idUser := r.URL.Query()["id_user"]

	db.Exec("DELETE FROM carts WHERE id_user=?", idUser[0])
	result, errQuery := db.Exec("DELETE FROM users WHERE id_user=?", idUser[0])

	if errQuery != nil {
		PrintError(400, "No Table Found", w)
		return
	} else {
		num, _ := result.RowsAffected()
		if num == 0 {
			PrintError(400, "User Not Found", w)
		} else {
			PrintSuccess(200, "User Deleted", w)
		}
	}

}
