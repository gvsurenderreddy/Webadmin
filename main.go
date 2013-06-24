package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"io"
	"net/http"
	"os"
)

var db *sql.DB
var config map[string]interface{}

type Page struct {
	Title string
	Body  []byte
}

func main() {
	loadConfig()
	db = connectToSql()
	defer db.Close()
	loadServer()
}

func connectToSql() *sql.DB {
	var err error
	db, err = sql.Open("sqlite3", "Schumix.db3")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return db
}

func loadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var dat []byte
	data := make([]byte, 100)
	for {
		count, err := file.Read(data)
		if err == io.EOF {
			break
		}
		dat = append(dat, data[:count]...)
	}

	if err := json.Unmarshal(dat, &config); err != nil {
		panic(err)
	}
}

func loadServer() {
	fmt.Print("Starting web server on localhost...")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.ParseFiles("www/index.html")
		rows, err := db.Query("SELECT name FROM admins")
		if err != nil {
			panic(err)
		}
		defer rows.Close()
		var name string
		for rows.Next() {
			rows.Scan(&name)
		}
		rows.Close()
		p := &Page{Title: "Schumix WebAdmin"}
		t.Execute(w, p)
	})
	fmt.Print("Done. Serving...")
	http.ListenAndServe(":45987", nil)
}
