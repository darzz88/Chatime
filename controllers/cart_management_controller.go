package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func InsertCart(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		return
	}

	id_drink := r.Form.Get("id_drink")
	quantity := r.Form.Get("quantity")

	var detailedCart DetailedCart
	var detailedCarts []DetailedCart

	rows, _ := db.Query(`SELECT detailed_carts.id_detailed_cart, detailed_carts.id_cart, detailed_carts.id_drink, detailed_carts.quantity 
	FROM detailed_carts 
	JOIN carts 
	ON detailed_carts.id_cart=carts.id_cart 
	WHERE carts.id_user =?`, GetOnlineUserId(r))

	for rows.Next() {
		if err := rows.Scan(&detailedCart.Id_Detailed_Cart, &detailedCart.Id_Cart, &detailedCart.Id_Drink, &detailedCart.Quantity); err != nil {
			log.Fatal(err.Error())
			PrintError(400, "No Product Data Inserted To []Product", w)
		} else {
			detailedCarts = append(detailedCarts, detailedCart)
		}
	}

	id_drink_int, _ := strconv.Atoi(id_drink)
	isFound := false

	for i := 0; i < len(detailedCarts); i++ {
		if detailedCarts[i].Id_Drink == id_drink_int {
			_, errQuery := db.Exec("UPDATE detailed_carts SET quantity = quantity + "+quantity+" WHERE id_detailed_cart = ? ", detailedCarts[i].Id_Detailed_Cart)
			isFound = true
			if errQuery == nil {
				PrintSuccess(200, "Added To Cart", w)
			} else {
				PrintError(400, "Failed", w)
			}
			return
		}
	}
	if isFound == false {
		_, errQuery := db.Exec("INSERT INTO detailed_carts(id_cart, id_drink, quantity) VALUES (?,?,?)", GetCartID(r), id_drink, quantity)

		if errQuery == nil {
			PrintSuccess(200, "Added To Cart", w)
		} else {
			PrintError(400, "Failed", w)
		}
	}

}
func SeeCart(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		return
	}
	Id_Drink := r.Form.Get("Id_Drink")

	query := `SELECT drinks.name, drinks.price, detailed_carts.quantity 
	FROM detailed_carts JOIN drinks ON detailed_carts.id_drink=drinks.id_drink
	JOIN carts ON carts.id_cart=detailed_carts.id_cart
	JOIN users ON users.id_user=carts.id_user
	WHERE users.id_user = ?`

	if len(Id_Drink) > 0 {
		query += ` WHERE detailed_carts.Id_Drink = ` + Id_Drink
	}

	rows, err := db.Query(query, GetOnlineUserId(r))

	if err != nil {
		log.Println(err)
		PrintError(400, "Rows Are Empty - Carts", w)
		return
	}

	var detailCart DetailedCartDrink
	var detailCarts []DetailedCartDrink
	for rows.Next() {
		if err := rows.Scan(&detailCart.DrinkData.Name,
			&detailCart.DrinkData.Price,
			&detailCart.Quantity); err != nil {
			log.Fatal(err.Error())
			PrintError(400, "No Product Data Inserted To []Product", w)
		} else {
			detailCarts = append(detailCarts, detailCart)
		}
	}

	var response DetailedCartDrinkResponse

	if len(detailCarts) > 0 {
		response.Data = detailCarts
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		PrintError(400, "No Detail Cart In []detailCart", w)
		return
	}
}
func UpdateQuantity(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		return
	}

	id_drink := r.Form.Get("id_drink")
	quantity := r.Form.Get("quantity")
	_, errQuery := db.Exec(`UPDATE detailed_carts 
	JOIN carts 
	ON detailed_carts.id_cart = carts.id_cart 
	SET detailed_carts.quantity =?   
	WHERE detailed_carts.id_drink =? AND carts.id_user =?`, quantity, id_drink, GetOnlineUserId(r))

	if errQuery == nil {
		PrintSuccess(200, "Updated Quantity", w)
	} else {
		PrintError(400, "Update Failed", w)
	}
}
func DeleteCart(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	id_drink := r.URL.Query()["id_drink"]

	_, errQuery := db.Exec(`DELETE detailed_carts
	FROM detailed_carts 
	JOIN carts
	ON detailed_carts.id_cart = carts.id_cart 
	WHERE detailed_carts.id_drink = ?
	AND carts.id_user =?`, id_drink[0], GetOnlineUserId(r))

	if errQuery == nil {
		PrintSuccess(200, "Delete Success", w)
	} else {
		PrintError(400, "Delete Failed", w)
	}

}

func GetCartID(r *http.Request) int {
	db := connect()
	defer db.Close()

	var idCart int
	err := db.QueryRow("SELECT id_cart FROM carts WHERE id_user = ?", GetOnlineUserId(r)).Scan(&idCart)
	if err != nil {
		log.Println(err)
	}
	return idCart
}
