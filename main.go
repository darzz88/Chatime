package main

import (
	c "Kokutime/controllers"
	"fmt"
	"os"

	"log"
	"net/http"

	"github.com/claudiu/gocron"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql" // Connection
	"github.com/gorilla/mux"           // Router
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	loadEnv()

	// 0 = Admin ; 1 = User

	// Processing, Delivering, Received

	router := mux.NewRouter()

	// Register Endpoint
	router.HandleFunc("/register", c.UserRegister).Methods("POST")

	// Login Endpoint
	router.HandleFunc("/login-user", c.UserLogin).Methods("POST")
	router.HandleFunc("/login-admin", c.AdminLogin).Methods("POST")

	// Logout Endpoint
	router.HandleFunc("/logout", c.Logout).Methods("POST")

	// User Endpoint
	router.HandleFunc("/user", c.Authenticate(c.UpdateUser, 1)).Methods("PUT")
	router.HandleFunc("/user", c.Authenticate(c.GetUserProfile, 1)).Methods("GET")
	router.HandleFunc("/user-admin", c.Authenticate(c.SeeUsers, 0)).Methods("GET")
	router.HandleFunc("/user-admin", c.Authenticate(c.DeleteUser, 0)).Methods("DELETE")

	// Drink Endpoint
	router.HandleFunc("/drink", c.GetDrinks).Methods("GET")
	router.HandleFunc("/drink", c.Authenticate(c.AddDrinks, 0)).Methods("POST")
	router.HandleFunc("/drink", c.Authenticate(c.DeleteDrink, 0)).Methods("DELETE")
	router.HandleFunc("/drink", c.Authenticate(c.UpdateDrink, 0)).Methods("PUT")

	// Transaction Endpoint
	router.HandleFunc("/transaction-admin", c.Authenticate(c.SalesReport, 0)).Methods("GET")
	router.HandleFunc("/transaction-admin", c.Authenticate(c.StatusManagement, 0)).Methods("PUT")
	router.HandleFunc("/transaction-basic", c.Authenticate(c.SeeOrder, 1)).Methods("GET")
	router.HandleFunc("/detail-transaction-basic", c.Authenticate(c.SeeDetailOrder, 1)).Methods("GET")

	// Cart Endpoint
	router.HandleFunc("/cart", c.Authenticate(c.SeeCart, 1)).Methods("GET")
	router.HandleFunc("/cart", c.Authenticate(c.InsertCart, 1)).Methods("POST")
	router.HandleFunc("/cart", c.Authenticate(c.UpdateQuantity, 1)).Methods("PUT")
	router.HandleFunc("/cart", c.Authenticate(c.DeleteCart, 1)).Methods("DELETE")

	// Checkout Endpoint
	router.HandleFunc("/checkout", c.Authenticate(c.GetTotalPrice, 1)).Methods("GET")
	router.HandleFunc("/checkout", c.Authenticate(c.Checkout, 1)).Methods("POST")

	// Tools
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	c.SetRedis(rdb, "eng", "Initialization", 0)
	c.SetRedis(rdb, "idn", "Inisialisasi", 0)

	gocron.Start()
	gocron.Every(3600).Seconds().Do(c.Task)

	//cors
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowCredentials: true,
	})

	Handler := corsHandler.Handler(router)

	// Connection Notif
	http.Handle("/", router)
	log.Println("Connected to port 8080")
	log.Fatal(http.ListenAndServe(":8080", Handler))
}

func loadEnv() {
	err := godotenv.Load()
	c.CheckError(err)

	appName := os.Getenv("APP_NAME")
	fmt.Println(appName)
}
