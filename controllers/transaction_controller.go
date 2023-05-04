package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Admin Status Management
func StatusManagement(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	var response TransactionResponse

	err := r.ParseForm()
	if err != nil {
		return
	}

	idTransaction := r.Form.Get("id_transaction")
	status := r.Form.Get("status")

	if len(idTransaction) <= 0 || len(status) <= 0 {
		PrintError(400, "Please Input All Field", w)
		return
	}

	result, errQuery := db.Exec("UPDATE transactions SET status=? WHERE id_transaction=?", status, idTransaction)
	if errQuery != nil {
		PrintError(400, "Table Not Found", w)
		return
	} else {
		num, _ := result.RowsAffected()
		if num == 0 {
			PrintError(400, "Transaction With This ID Are Not Found", w)
		} else {
			w.WriteHeader(http.StatusOK)
			response.Status = 200
			response.Message = "Status Transaction Updated!"
			json.NewEncoder(w).Encode(response)
		}
	}
}

func SeeOrder(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	var response TransactionResponse
	var transaction Transaction
	var transactions []Transaction
	err := r.ParseForm()
	if err != nil {
		return
	}

	rows, err := db.Query("SELECT * FROM transactions WHERE id_user=?", GetOnlineUserId(r))
	if err != nil {
		PrintError(404, "Table Not Found", w)
		return
	}

	for rows.Next() {
		if err := rows.Scan(&transaction.Id_Transaction, &transaction.Id_User, &transaction.Status, &transaction.Date); err != nil {
			PrintError(400, "Print Undefined", w)
			return
		} else {
			transaction.Total = TotalPrice(transaction.Id_Transaction)
			transactions = append(transactions, transaction)
		}
	}
	if len(transactions) != 0 {
		response.Status = 200
		response.Message = "Success"
		response.Data = transactions
	} else {
		PrintError(400, "Array Size Not Found", w)
	}
	json.NewEncoder(w).Encode(response)
}

func SeeDetailOrder(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	var response DetailedTransactionResponse
	var detail DetailedTransaction
	var details []DetailedTransaction

	err := r.ParseForm()
	if err != nil {
		return
	}

	id_transaction := r.URL.Query()["id_transaction"]

	rows, err := db.Query(`SELECT drinks.name, detailed_transactions.quantity
	FROM detailed_transactions
	JOIN drinks ON drinks.id_drink = detailed_transactions.id_drink
	WHERE detailed_transactions.id_transaction=?`, id_transaction[0])
	if err != nil {
		PrintError(400, "Table Not Found", w)
		return
	}

	for rows.Next() {
		if err := rows.Scan(&detail.Drink_Name, &detail.Quantity); err != nil {
			PrintError(400, "Print Undefined", w)
			return
		} else {
			details = append(details, detail)
		}
	}

	if len(details) != 0 {
		response.Status = 200
		response.Message = "Success"
		response.Data = details
	} else {
		PrintError(400, "No Detail Transaction Found", w)
	}
	json.NewEncoder(w).Encode(response)

}

func SalesReport(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	var response TransactionResponse
	var transaction Transaction
	var transactions []Transaction

	err := r.ParseForm()
	if err != nil {
		return
	}

	rows, err := db.Query(`SELECT * FROM transactions`)
	if err != nil {
		fmt.Println(err)
		PrintError(404, "Table Not Found", w)
		return
	}

	for rows.Next() {
		if err := rows.Scan(&transaction.Id_Transaction, &transaction.Id_User, &transaction.Status, &transaction.Date); err != nil {
			PrintError(400, "Print Undefined", w)
			return
		} else {
			transaction.Total = TotalPrice(transaction.Id_Transaction)
			transactions = append(transactions, transaction)
		}
	}

	if len(transactions) != 0 {
		response.Status = 200
		response.Message = "Success"
		response.Data = transactions
	} else {
		PrintError(400, "No Sales Found", w)
	}
	json.NewEncoder(w).Encode(response)
}

func TotalPrice(id_transaction int) int {
	db := connect()
	defer db.Close()

	row := db.QueryRow(`SELECT SUM(c.price * a.quantity)
	FROM detailed_transactions a
    JOIN transactions b ON a.id_transaction=b.id_transaction
    JOIN drinks c ON a.id_drink=c.id_drink
    WHERE b.id_transaction = ?`, id_transaction)

	var total int
	err := row.Scan(&total)

	if err == nil {
		return total
	} else {
		return -1
	}
}
