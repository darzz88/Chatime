package controllers

import (
	"encoding/json"
	"log"
	"net/http"
)

// General error checking
func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
		// fmt.Println(err)
	}
}

// Postman output
func PrintSuccess(status int, message string, w http.ResponseWriter) {
	var succResponse SuccessResponse
	succResponse.Status = status
	succResponse.Message = message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(succResponse)
}

func PrintError(status int, message string, w http.ResponseWriter) {
	var errResponse ErrorResponse
	errResponse.Status = status
	errResponse.Message = message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(errResponse)
}

// Register
func UserRegister(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	CheckError(err)

	Name := r.Form.Get("name")
	Address := r.Form.Get("address")
	Email := r.Form.Get("email")
	Password := r.Form.Get("password")

	_, errQuery := db.Exec("INSERT INTO users(Name, Address, Email, Password) VALUES(?, ?, ?, ?)",
		Name,
		Address,
		Email,
		Password)

	if errQuery == nil {
		// Auto create new cart for new user
		rows, err := db.Query("SELECT id_user FROM users WHERE email = ? AND password = ?", Email, Password)
		CheckError(err)

		var id_user int

		for rows.Next() {
			rows.Scan(&id_user)
		}

		db.Exec("INSERT INTO carts(id_user) VALUES(?)", id_user)

		PrintSuccess(200, "Registered", w)
	} else {
		PrintError(400, "Registration Failed", w)
		return
	}
}

// Login
func AdminLogin(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	email := r.URL.Query()["email"]
	password := r.URL.Query()["password"]

	row := db.QueryRow("SELECT * FROM admins WHERE email = ? AND password = ?", email[0], password[0])

	var admin Admin
	err := row.Scan(&admin.Id_Admin, &admin.Name, &admin.Email, &admin.Password)

	if err != nil {
		PrintError(400, "Admin Not Found", w)
	} else {
		userType := 0
		generateToken(w, admin.Id_Admin, admin.Name, userType)
		PrintSuccess(200, "Logged In", w)
	}
}

func UserLogin(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	email := r.URL.Query()["email"]
	password := r.URL.Query()["password"]

	row := db.QueryRow("SELECT * FROM users WHERE email = ? AND password = ?", email[0], password[0])

	var user User
	err := row.Scan(&user.Id_User, &user.Name, &user.Address, &user.Email, &user.Password)

	if err != nil {
		PrintError(400, "User Not Found", w)
	} else {
		userType := 1
		generateToken(w, user.Id_User, user.Name, userType)
		PrintSuccess(200, "Logged In", w)
	}
}

// Logout
func Logout(w http.ResponseWriter, r *http.Request) {
	resetToken(w)
	PrintSuccess(200, "Logged Out", w)
}
