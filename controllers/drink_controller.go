package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-redis/redis/v8"
)

func GetDrinks(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	idDrink := r.URL.Query()["id_drink"]
	drinkName := r.URL.Query()["name"]

	var drink Drink
	var drinks []Drink

	query := "SELECT * FROM drinks"

	if len(drinkName[0]) > 0 {
		query += " WHERE name LIKE '%" + drinkName[0] + "%'"
	} else if len(idDrink[0]) > 0 {
		query += " WHERE id_drink = " + idDrink[0]
	}

	rows, err := db.Query(query)
	CheckError(err)

	for rows.Next() {
		if err := rows.Scan(&drink.Id_Drink, &drink.Name, &drink.Price, &drink.Description); err != nil {
			PrintError(400, "No User Data Inserted To []User", w)
			log.Fatal(err)
			return
		} else {
			drinks = append(drinks, drink)
		}
	}

	var response DrinksResponse
	if err == nil {
		response.Status = 200
		response.Message = "Get Drinks Success"
		response.Data = drinks
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		PrintError(400, "No Data Found", w)
		return
	}
}
func AddDrinks(w http.ResponseWriter, r *http.Request) {

	db := connectGorm()

	err := r.ParseForm()
	CheckError(err)

	name := r.Form.Get("name")
	price := r.Form.Get("price")
	description := r.Form.Get("description")

	var drink Drink

	drink.Name = name
	drink.Price, _ = strconv.Atoi(price)
	drink.Description = description

	query := db.Select("name", "price", "description").Create(&drink)

	if query != nil {
		rdb := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})
		SetRedis(rdb, "eng", "New Drink: "+name, 0)
		SetRedis(rdb, "idn", "Minuman Baru: "+name, 0)
		PrintSuccess(200, "Drink Inserted", w)
	} else {
		PrintError(400, "insert Drinks Failed", w)
	}
}
func DeleteDrink(w http.ResponseWriter, r *http.Request) {

	db := connectGorm()

	idDrink := r.URL.Query()["id_drink"]

	var drink Drink

	err := db.Delete(&drink, idDrink)

	if err != nil {
		PrintSuccess(200, "Drink Deleted", w)
	} else {
		PrintError(400, "Delete Drink Failed", w)
		return
	}

}
func UpdateDrink(w http.ResponseWriter, r *http.Request) {

	db := connectGorm()

	err := r.ParseForm()
	CheckError(err)

	idDrink := r.Form.Get("id_drink")
	name := r.Form.Get("name")
	price := r.Form.Get("price")
	description := r.Form.Get("description")

	var drink Drink

	drink.Id_Drink, _ = strconv.Atoi(idDrink)
	priceint, _ := strconv.Atoi(price)

	query := db.Model(&drink).Where("id_drink = ?", idDrink).Updates(Drink{Name: name, Price: priceint, Description: description})

	if query != nil {
		PrintSuccess(200, "Drink Updated", w)
	} else {
		PrintError(400, "Update Drinks Failed", w)
	}

}
