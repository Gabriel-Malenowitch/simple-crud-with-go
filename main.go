package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Names struct {
	Id    int
	Name  string
	Email string
}

func throwErrorIfNecessary(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func dbConnect() (db *sql.DB) {
	const dbDriver, dbUser, dbPass, dbName string = "mysql", "male", "", "crudgo"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	throwErrorIfNecessary(err)

	return db
}

func safeQuery(query string) *sql.Rows {
	db := dbConnect()
	result, err := db.Query(query)
	throwErrorIfNecessary(err)

	return result
}

func getList(writer http.ResponseWriter, request *http.Request) {
	namesRows := safeQuery("select * from names")
	var names []Names
	for namesRows.Next() {
		var id int
		var name, email string

		err := namesRows.Scan(&id, &name, &email)
		throwErrorIfNecessary(err)
		names = append(names, Names{id, name, email})
	}

	bytesNames, err := json.Marshal(names)
	throwErrorIfNecessary(err)
	writer.Write(bytesNames)

	defer dbConnect().Close()
}

func getOne(writer http.ResponseWriter, request *http.Request) {
	var id string = strings.Split(request.URL.String(), "/")[2]
	nameRow := safeQuery("select * from names where id = " + id)
	var names []Names
	for nameRow.Next() {
		var id int
		var name, email string

		err := nameRow.Scan(&id, &name, &email)
		throwErrorIfNecessary(err)
		names = append(names, Names{id, name, email})
	}

	bytesNames, err := json.Marshal(names)
	throwErrorIfNecessary(err)
	writer.Write(bytesNames)

	defer dbConnect().Close()
}

func create(writer http.ResponseWriter, request *http.Request) {
	var names Names
	err := json.NewDecoder(request.Body).Decode(&names)
	throwErrorIfNecessary(err)

	name := names.Name
	email := names.Email

	db := dbConnect()
	insForm, err := db.Prepare("INSERT INTO names(name, email) VALUES(?,?)")
	throwErrorIfNecessary(err)
	insForm.Exec(name, email)

	defer dbConnect().Close()
}

func edit(writer http.ResponseWriter, request *http.Request) {
	var id string = strings.Split(request.URL.String(), "/")[2]

	var names Names
	err := json.NewDecoder(request.Body).Decode(&names)
	throwErrorIfNecessary(err)

	name := names.Name
	email := names.Email

	db := dbConnect()
	insForm, err := db.Prepare("UPDATE names SET name=?, email=? WHERE id=?")
	throwErrorIfNecessary(err)
	insForm.Exec(name, email, id)

	defer dbConnect().Close()
}

func delete(writer http.ResponseWriter, request *http.Request) {
	var id string = strings.Split(request.URL.String(), "/")[2]

	db := dbConnect()
	insForm, err := db.Prepare("DELETE FROM names WHERE id=?")
	throwErrorIfNecessary(err)

	insForm.Exec(id)
	defer dbConnect().Close()
}

func main() {
	log.Println("Server started on: http://localhost:9000")

	http.HandleFunc("/names", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == "GET" {
			getList(writer, request)
		} else if request.Method == "POST" {
			create(writer, request)
		}
	})
	http.HandleFunc("/names/", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == "GET" {
			getOne(writer, request)
		} else if request.Method == "PUT" {
			edit(writer, request)
		} else if request.Method == "DELETE" {
			delete(writer, request)
		}
	})

	http.ListenAndServe(":9000", nil)

}
