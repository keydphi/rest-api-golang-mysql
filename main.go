package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
)

type Post struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

var db *sql.DB
var err error

const (
	username  = "paulo"
	password  = "pinkel"
	hostname  = "127.0.0.1:3306"
	dbname    = "newDatabase"
	tableName = "postsTable6"
)

func dsn(dsnDbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dsnDbName)
	//	return fmt.Sprintf("%s:%s@tcp(%s)", username, password, hostname)
	//	return fmt.Sprintf("user:password@tcp(127.0.0.1:3306)/database-name")
}

func main() {
	db, err = sql.Open("mysql", dsn(""))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	log.Printf("Opened mysql DSN successfully\n")

	connectToDB(dbname)

	//*
	runAPI()
	/*/
	queryDBnStuff(db, dbname)
	//*/
}
func runAPI() {
	router := mux.NewRouter()
	router.HandleFunc("/posts", getPosts).Methods("GET")
	router.HandleFunc("/posts", createPost).Methods("POST")
	router.HandleFunc("/posts/{id}", getPost).Methods("GET")
	router.HandleFunc("/posts/{id}", updatePost).Methods("PUT")
	router.HandleFunc("/posts/{id}", deletePost).Methods("DELETE")
	http.ListenAndServe(":8000", router)
	log.Printf("API open and running\n")

}
func connectToDB(targetDbName string) {
	//*
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	res, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+targetDbName)
	if err != nil {
		log.Printf("Error %s when creating DB\n", err)
		return
	}
	no, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when fetching rows", err)
		return
	}
	log.Printf("rows affected %d\n", no)
	//*/

	db.Close()
	db, err = sql.Open("mysql", dsn(targetDbName))
	if err != nil {
		log.Printf("Error %s when opening DB", err)
		return
	}
	//	defer db.Close()

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Errors %s pinging DB %s", err, connectToDB)
		return
	}
	log.Printf("Connected to DB %s successfully\n", targetDbName)

}
func queryDBnStuff(theDB *sql.DB, targetDbName string) {

	//*
	//	query := fmt.Sprintf("CREATE TABLE %s (`%s` int(6) unsigned NOT NULL AUTO_INCREMENT, `%s` varchar(30) NOT NULL, PRIMARY KEY (`%s`)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;", tableName, "json:\"id\"", "json:\"title\"", "json:\"id\"")
	query := fmt.Sprintf("CREATE TABLE %s (`%s` int(6) unsigned NOT NULL AUTO_INCREMENT, `%s` varchar(30) NOT NULL, PRIMARY KEY (`%s`)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;", tableName, "id", "title", "id")
	/*/
	query := fmt.Sprintf("INSERT INTO %s VALUES('1', 'Merkel', 'Ghostbusters')", tableName)
	//*/
	insert, err := theDB.Query(query)
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
	log.Printf("Executed query %s to table %s successfully\n", query, tableName)
}

func getPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var posts []Post
	result, err := db.Query("SELECT id, title from " + tableName)
	//result, err := db.Query("SELECT json:\"id\", title from " + tableName)
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	for result.Next() {
		var post Post
		err := result.Scan(&post.ID, &post.Title)
		if err != nil {
			panic(err.Error())
		}
		posts = append(posts, post)
	}
	json.NewEncoder(w).Encode(posts)
	log.Printf("got posts\n")
}
func createPost(w http.ResponseWriter, r *http.Request) {
	//	log.Printf("Creating post\n")
	w.Header().Set("Content-Type", "application/json")
	//	log.Printf("Preparing DB\n")
	stmt, err := db.Prepare("INSERT INTO " + tableName + "(title) VALUES(?)")
	if err != nil {
		panic(err.Error())
	}
	log.Printf("Prepared DB successfully\n")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	title := keyVal["title"]
	_, err = stmt.Exec(title)
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "New post was created")
	log.Printf("created post\n")
}
func getPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	result, err := db.Query("SELECT id, title FROM "+tableName+" WHERE id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	var post Post
	for result.Next() {
		err := result.Scan(&post.ID, &post.Title)
		if err != nil {
			panic(err.Error())
		}
	}
	json.NewEncoder(w).Encode(post)
	log.Printf("got post\n")
}
func updatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	stmt, err := db.Prepare("UPDATE " + tableName + " SET title = ? WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	newTitle := keyVal["title"]
	_, err = stmt.Exec(newTitle, params["id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Post with ID = %s was updated", params["id"])
	log.Printf("updated post\n")
}
func deletePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM " + tableName + " WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Post with ID = %s was deleted", params["id"])
	log.Printf("deleted post\n")
}
