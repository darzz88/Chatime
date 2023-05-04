package controllers

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

func GetTotalPrice(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	row := db.QueryRow(`
	SELECT SUM(drinks.price * detailed_carts.quantity) 
	FROM drinks
	JOIN detailed_carts ON drinks.id_drink = detailed_carts.id_drink
	JOIN carts ON detailed_carts.id_cart = carts.id_cart
	WHERE carts.id_user = ?`,
		GetOnlineUserId(r))

	var total int
	err := row.Scan(&total)

	if err == nil {
		PrintSuccess(200, "Total Price: "+strconv.Itoa(total), w)
	} else {
		PrintError(400, "User Not Found", w)
		return
	}
}

func CalculateTotalPrice(id_user int) int {
	db := connect()
	defer db.Close()

	row := db.QueryRow(`
	SELECT SUM(drinks.price * detailed_carts.quantity) 
	FROM drinks
	JOIN detailed_carts ON drinks.id_drink = detailed_carts.id_drink
	JOIN carts ON detailed_carts.id_cart = carts.id_cart
	WHERE carts.id_user = ?`, id_user)

	var total int
	err := row.Scan(&total)

	if err == nil {
		return total
	} else {
		return -1
	}
}

func Checkout(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	currentTime := time.Now()

	err := r.ParseForm()
	CheckError(err)

	// Get data from detailed_carts
	rows, err := db.Query(`
	SELECT id_detailed_cart, detailed_carts.id_cart, id_drink, quantity 
	FROM detailed_carts
	JOIN carts ON detailed_carts.id_cart = carts.id_cart
	WHERE carts.id_user = ?`, GetOnlineUserId(r))
	CheckError(err)

	var detailed_cart DetailedCart
	var detailed_carts []DetailedCart
	for rows.Next() {
		if err := rows.Scan(
			&detailed_cart.Id_Detailed_Cart,
			&detailed_cart.Id_Cart,
			&detailed_cart.Id_Drink,
			&detailed_cart.Quantity); err != nil {
			PrintError(400, "No Item In Cart", w)
			log.Fatal(err)
			return
		} else {
			detailed_carts = append(detailed_carts, detailed_cart)
		}
	}

	totalPrice := CalculateTotalPrice(GetOnlineUserId(r))

	if totalPrice == -1 {
		PrintError(400, "No Item In Cart", w)
		return
	} else {
		_, err := db.Exec(`
		INSERT INTO transactions(id_user, status, date) 
		VALUES( ?, ?, ?)`, GetOnlineUserId(r), "Processing", currentTime.Format("2006-01-02"))
		CheckError(err)
	}

	// Get newest id_transaction for the person
	rows, err = db.Query(`
	SELECT id_transaction
	FROM transactions
	WHERE id_user = ?`, GetOnlineUserId(r))
	CheckError(err)

	var id_transaction int
	var id_transactions []int

	for rows.Next() {
		if err := rows.Scan(
			&id_transaction); err != nil {
			PrintError(400, "Error in Inserting Transaction", w)
			log.Fatal(err)
			return
		} else {
			id_transactions = append(id_transactions, id_transaction)
		}
	}

	// Insert into detailed_transactions
	for i := range detailed_carts {
		db.Exec(`INSERT INTO detailed_transactions(id_transaction, id_drink, quantity)
		VALUES(?, ?, ?)`, id_transactions[len(id_transactions)-1], detailed_carts[i].Id_Drink, detailed_carts[i].Quantity)
	}

	// Delete from detailed_cart
	db.Exec(`DELETE detailed_carts
		FROM detailed_carts
		JOIN carts ON detailed_carts.id_cart = carts.id_cart
		WHERE carts.id_user = ?`, GetOnlineUserId(r))

	text := "Received payment: Rp." + strconv.Itoa(totalPrice)
	SendMail(GetEmailUser(r), "Thanks For Ordering", text)

	PrintSuccess(200, "Checked out", w)
}

func GetEmailUser(r *http.Request) string {
	db := connect()
	defer db.Close()

	var email string
	err := db.QueryRow("SELECT email FROM users WHERE id_user = ?", GetOnlineUserId(r)).Scan(&email)
	if err != nil {
		log.Println(err)
	}
	return email
}
