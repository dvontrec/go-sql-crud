package main

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// sets a db variable
var db *sql.DB
var err error

func main() {
	// "user:password@tcp(localhost:555)/dbname?charset=utf8"
	db, err = sql.Open("mysql", "root:password@(localhost)/homies?charset=utf8") // connects to local host using local credentials
	check(err)
	// defer the close
	defer db.Close()

	// Pings the db
	err = db.Ping()
	check(err)

	// Setup routes
	http.HandleFunc("/", index)
	http.HandleFunc("/homies", homies)
	http.HandleFunc("/create", create)
	http.HandleFunc("/insert", insert)
	http.HandleFunc("/read", read)
	http.HandleFunc("/update", update)
	http.HandleFunc("/delete", del)
	http.HandleFunc("/drop", drop)
	http.Handle("/faveicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
	check(err)
}

// index route code
func index(w http.ResponseWriter, req *http.Request) {
	_, err := io.WriteString(w, "Successfully connected")
	check(err)
}

// route to show all homies
func homies(w http.ResponseWriter, req *http.Request) {
	// runs a query to pull data from the database
	rows, err := db.Query(`SELECT hName FROM homies`)
	check(err)

	// variables for data to print
	var s, name string
	s = "The Homies:\n"

	// loop to end of rows
	for rows.Next() {
		// sets the name var to be the content of the row
		err = rows.Scan(&name)
		check(err)
		s += name + "\n"
	}
	fmt.Fprintln(w, s)
}

func create(w http.ResponseWriter, req *http.Request) {
	// prepares a query to create the customers table
	stms, err := db.Prepare(`CREATE TABLE customers (name VARCHAR(20));`)
	check(err)
	// executes the query
	r, err := stms.Exec()
	check(err)

	// if the query is successful there will be no errors and returns the number of rows effected
	n, err := r.RowsAffected()
	check(err)

	// prints the success to browser
	fmt.Fprintln(w, "CREATED TABLE customers", n)
}

// route to add customer james to customers table
func insert(w http.ResponseWriter, req *http.Request) {
	stmt, err := db.Prepare(`INSERT INTO customers VALUES("james");`)
	check(err)

	r, err := stmt.Exec()
	check(err)

	n, err := r.RowsAffected()
	check(err)

	fmt.Fprintln(w, "INSERTED RECORD", n)
}

// route to read content from customers table
func read(w http.ResponseWriter, req *http.Request) {
	// query to get all customers
	rows, err := db.Query(`SELECT * FROM customers`)
	check(err)

	var name string
	for rows.Next() {
		// sets the value of name to be the content of the row
		err = rows.Scan(&name)
		check(err)
		fmt.Println(name)
		fmt.Fprintln(w, "RETREIVED RECORD: ", name)
	}
}

// route to update rows in db
func update(w http.ResponseWriter, req *http.Request) {
	stmt, err := db.Prepare(`UPDATE customers SET name="Jimmy" WHERE name="James";`)
	check(err)

	r, err := stmt.Exec()
	check(err)

	n, err := r.RowsAffected()
	check(err)

	fmt.Fprintln(w, "UPDATED RECORD", n)
}

// route to delete elements from table
func del(w http.ResponseWriter, req *http.Request) {
	stmt, err := db.Prepare(`DELETE FROM customers WHERE name="Jimmy";`)
	check(err)

	r, err := stmt.Exec()
	check(err)

	n, err := r.RowsAffected()
	check(err)

	fmt.Fprintln(w, "Deleted customer", n)
}

// route for droping a table
func drop(w http.ResponseWriter, req *http.Request) {
	stmt, err := db.Prepare(`DROP TABLE customers;`)
	check(err)

	_, err = stmt.Exec()
	check(err)

	fmt.Fprintln(w, "customers TABLE DROPPED")
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
