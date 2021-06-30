package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Client struct {
	ID      int    `json:"idCliente"`
	Nombre  string `json:"Nombre"`
	Paterno string `json:"Paterno"`
	Materno string `json:"Materno"`
	Edad    int    `json:"Edad"`
}

var db *sql.DB
var err error

func main() {
	db, err = sql.Open("mysql", "root:Qwerty123451*@tcp(127.0.0.1:3306)/Clientes")
	//db, err = sql.Open("mysql", "usrtest:Qwerty@tcp(godockerDB)/Clientes")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	router := mux.NewRouter()
	router.HandleFunc("/clientes", getClients).Methods("GET")
	router.HandleFunc("/clientes", createClient).Methods("POST")
	router.HandleFunc("/clientes/{id}", getClient).Methods("GET")
	router.HandleFunc("/clientes/{id}", updateClient).Methods("PUT")
	router.HandleFunc("/clientes/{id}", deleteClient).Methods("DELETE")
	http.ListenAndServe(":8080", router)
}
func getClients(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var posts []Client
	result, err := db.Query("SELECT idCliente, Nombre, Paterno, Materno, Edad from Cliente")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	for result.Next() {
		var post Client
		err := result.Scan(&post.ID, &post.Nombre, &post.Paterno, &post.Materno, &post.Edad)
		if err != nil {
			panic(err.Error())
		}
		posts = append(posts, post)
	}
	json.NewEncoder(w).Encode(posts)
}
func createClient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	stmt, err := db.Prepare("INSERT INTO Cliente(idCliente,Nombre,Paterno,Materno,Edad) VALUES(?,?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	var client Client
	json.Unmarshal(body, &client)

	_, err = stmt.Exec(client.ID, client.Nombre, client.Paterno, client.Materno, client.Edad)

	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "New user was created")
}
func getClient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	result, err := db.Query("select idCliente, Nombre, Paterno, Materno, Edad from Cliente WHERE idCliente = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	var post Client
	for result.Next() {
		err := result.Scan(&post.ID, &post.Nombre, &post.Paterno, &post.Materno, &post.Edad)
		if err != nil {
			panic(err.Error())
		}
	}
	json.NewEncoder(w).Encode(post)
}
func updateClient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	stmt, err := db.Prepare("UPDATE Cliente SET Nombre = ?, Paterno = ?,Materno=?,Edad=? WHERE idCliente = ?")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	var client Client
	json.Unmarshal(body, &client)

	_, err = stmt.Exec(client.Nombre, client.Paterno, client.Materno, client.Edad, params["id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Client with ID = %s was updated", params["id"])
}
func deleteClient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM Cliente WHERE idCliente = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Client with ID = %s was deleted", params["id"])
}
