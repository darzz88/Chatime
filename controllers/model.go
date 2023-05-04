package controllers

// Error Success
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// Admin & User
type Admin struct {
	Id_Admin int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	Id_User  int    `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    User   `json:"data"`
}

type UsersResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []User `json:"data"`
}

// Drink
type Drink struct {
	Id_Drink    int    `json:"id"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Description string `json:"description"`
}

type DrinksResponse struct {
	Status  int     `json:"status"`
	Message string  `json:"message"`
	Data    []Drink `json:"data"`
}

// Cart
type Cart struct {
	IdCart int `json:"idCart"`
	IdUser int `json:"idUser"`
}

type DetailedCartDrink struct {
	IdDetailedCart int   `json:"ID"`
	DrinkData      Drink `json:"User"`
	Quantity       int   `json:"Quantity"`
}

type DetailedCartDrinkResponse struct {
	Data []DetailedCartDrink `json:"data"`
}

type DetailedCart struct {
	Id_Detailed_Cart int `json:"id_detailed_cart"`
	Id_Cart          int `json:"id_cart"`
	Id_Drink         int `json:"id_drink"`
	Quantity         int `json:"quantity"`
}

// Transaction
type Transaction struct {
	Id_Transaction int    `json:"id_transaction"`
	Id_User        int    `json:"id_user"`
	Status         string `json:"status"`
	Total          int    `json:"total"`
	Date           string `json:"date"`
}

type TransactionResponse struct {
	Status  int           `json:"status"`
	Message string        `json:"message"`
	Data    []Transaction `json:"data"`
}

type DetailedTransaction struct {
	Drink_Name string `json:"drink_name"`
	Quantity   int    `json:"quantity"`
}

type DetailedTransactionResponse struct {
	Status  int                   `json:"status"`
	Message string                `json:"message"`
	Data    []DetailedTransaction `json:"data"`
}
