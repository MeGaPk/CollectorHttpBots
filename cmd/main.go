package main

import (
	"fmt"
	"net/http"

	"../database"
	"encoding/json"
	"io/ioutil"
	"os"
	"github.com/gorilla/handlers"
	"flag"
)

var db *database.DatabaseConnection

func handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	link := r.URL.String()
	remote_ip := r.RemoteAddr
	fmt.Fprintf(w, "My url: %s", link)

	header, _ := json.Marshal(r.Header)
	body, _ := ioutil.ReadAll(r.Body)
	form, _ := json.Marshal(r.Form)
	postForm, _ := json.Marshal(r.PostForm)
	db.AddUrl(&database.Bot{
		Link: link,
		Header: string(header),
		Body: string(body),
		Form: string(form),
		PostForm: string(postForm),
		RemoteIp: remote_ip,
	})
}

func GetUrls(w http.ResponseWriter, r *http.Request) {
	js, err := json.MarshalIndent(db.GetUrls(), "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

var type_database string
var port int
var hostname string
var port_database int
var login string
var password string
var dbname string

func main() {
	flag.StringVar(&type_database, "type_database", "sqlite", "type (mysql/sqlite)")
	flag.StringVar(&hostname, "hostname", "127.0.0.1", "hostname")
	flag.StringVar(&login, "login", "root", "login")
	flag.StringVar(&password, "password", "root", "password")
	flag.StringVar(&dbname, "dbname", "db", "database name")
	flag.IntVar(&port_database, "port_database", 3306, "port number for mysql server")
	flag.IntVar(&port, "port", 8080, "port number for http server")
	flag.Parse()
	db = nil;
	if (type_database == "mysql") {
		db = database.NewMySQL(hostname, port_database, login, password, dbname)
	} else {
		db = database.NewSqlite3("db.sqlite3")
	}
	db.AutoMigrate(&database.Bot{})

	r := http.NewServeMux()
	r.HandleFunc("/get_urls", GetUrls)
	r.Handle("/", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler)))
	http.ListenAndServe(fmt.Sprintf(":%d", port), handlers.CompressHandler(r))
}