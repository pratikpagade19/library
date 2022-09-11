package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

const (
	API_PATH = "/apis/v1/books"
)

type library struct {
	dbHost, dbPass, dbName string
}

type Book struct {
	Id, Name, Isbn string
}

func main() {

	//DB_HOST = host:port
	dbHost := os.Getenv("DB_HOST")

	if dbHost == "" {
		dbHost = "remotemysql.com:3306"
	}

	dbPass := os.Getenv("DB_PASS")

	if dbPass == "" {
		dbPass = "shsfVD1mBf"
	}

	apiPath := os.Getenv("API_PATH")

	if apiPath == "" {
		apiPath = API_PATH
	}

	dbName := os.Getenv("DB_NAME")

	if dbName == "" {
		dbName = "L9t7YYAl66"
	}

	l := library{
		dbHost: dbHost,
		dbPass: dbPass,
		dbName: dbName,
	}
	r := mux.NewRouter()

	r.HandleFunc("/", HomeHandler)
	r.HandleFunc(apiPath, l.getBooksHandler).Methods("GET")
	r.HandleFunc(apiPath, l.postBookHandler).Methods("POST")
	http.Handle("/", r)

	http.ListenAndServe("localhost:8080", r)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("Home route called")
}

func (l library) postBookHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Post Book route called")

	book := Book{}

	json.NewDecoder(r.Body).Decode(&book)

	db := l.openConnection()

	insertQuery, err := db.Prepare("insert into books values (?, ?, ?)")

	if err != nil {
		log.Fatalf("Preparing insertQuery %s \n", err.Error())
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf(" Beginning transaction %s \n", err.Error())
	}

	_, err = tx.Stmt(insertQuery).Exec(book.Id, book.Name, book.Isbn)
	if err != nil {
		log.Fatalf(" Execing the insert command %s \n", err.Error())
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf(" while commiting the transaction %s \n", err.Error())
	}

	l.closeConnection(db)
}

func (l library) getBooksHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("Get Books route called")
	db := l.openConnection()
	rows, err := db.Query("select * from books")
	if err != nil {
		log.Fatalf("Error while querying the books table %s \n", err.Error())
	}

	books := []Book{}
	for rows.Next() {
		var id, name, isbn string
		err := rows.Scan(&id, &name, &isbn)

		if err != nil {
			log.Fatalf("Error while scanning book rows %s \n", err.Error())
		}

		aBook := Book{

			Id:   id,
			Name: name,
			Isbn: isbn,
		}

		books = append(books, aBook)

	}

	json.NewEncoder(w).Encode(books)

	l.closeConnection(db)

}

func (l library) openConnection() *sql.DB {

	// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s)/%s", "L9t7YYAl66", l.dbPass, l.dbHost, l.dbName))

	if err != nil {
		log.Fatalf("opening the connection to database %s \n", err.Error())
	}

	return db
}

//close connection
func (l library) closeConnection(db *sql.DB) {

	err := db.Close()

	if err != nil {
		log.Fatalf("Closing connection %s \n", err.Error())
	}

}
