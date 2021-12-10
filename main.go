package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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
	username = "pedro"
	password = "polinesis"
	//	hostname  = "127.0.0.1"
	hostname  = "localhost"
	port      = "3306"
	dbname    = "newDatabase"
	tablename = "postsTable11"
)

func main() {

	connectToDB()
	createTableIfNotExists()

	stmt, err := db.Prepare("INSERT INTO " + dbname + "." + tablename + "(title) VALUES(?)")
	if err != nil {
		log.Printf("Error %s at Insert into table name %s at stmt %s", err, tablename, stmt)
		log.Printf("Table %s created DB\n", tablename)
	}

	//*
	runAPI()
	//defer db.Close()
	/*/
	// insertValuesIntoTable(db, dbname)
	//*/
}

func runAPI() {
	log.Printf("starting API\n")
	router := mux.NewRouter()
	router.HandleFunc("/posts", getPosts).Methods("GET")
	router.HandleFunc("/posts", createPost).Methods("POST")
	router.HandleFunc("/posts/{id}", getPost).Methods("GET")
	router.HandleFunc("/posts/{id}", updatePost).Methods("PUT")
	router.HandleFunc("/posts/{id}", deletePost).Methods("DELETE")
	http.ListenAndServe(":"+string(port), router)
	log.Printf("API open and running\n")

}

func connectToDB() {
	db, err = sql.Open("mysql", get_dsn(dbname))
	if err != nil {
		log.Printf("Error %s when opening mysql dsn", err)
		return
	}

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	// create db
	log.Printf("create db %s if not existing\n", dbname)
	res, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+dbname)
	if err != nil {
		log.Printf("Error %s when cueureating db %s\n", err, dbname)
		return
	}

	no, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when fetching rows", err)
		return
	}
	log.Printf("CREATE: rows affected %d\n", no)

	// use db
	log.Printf("use db %s\n", dbname)
	res2, err := db.ExecContext(ctx, "USE "+dbname+";")
	if err != nil {
		log.Printf("Error %s when setting USE DB "+dbname+"\n", err)
		return
	}

	no2, err := res2.RowsAffected()
	if err != nil {
		log.Printf("Error %s when fetching rows", err)
	}
	log.Printf("USE: rows affected %d\n", no2)

	// mysql privileges
	log.Printf("granting privileges on db %s table %s\n", dbname, tablename)
	/*
		res3, err := db.ExecContext(ctx, "SHOW GRANTS;")
		/*/
	res3, err := db.ExecContext(ctx, "GRANT ALL ON "+dbname+"."+tablename+" TO '"+username+"'@'"+hostname+"' WITH GRANT OPTION;")
	//*/
	if err != nil {
		log.Printf("Error %s when granting privileges on %s\n", err, dbname)
		return
	}
	//	log.Printf("mysql result: %s", res3.LastInsertId().)

	//	log.Printf("hello\n")

	//lastInsertId, err := res3.LastInsertId()
	//log.Printf("last command (id %s)\n", lastInsertId)

	//	log.Printf("hello\n")

	no3, err := res3.RowsAffected()
	if err != nil {
		log.Printf("Error %s when fetching rows", err)
	}
	log.Printf("GRANT: rows affected %d\n", no3)

	db.Close()

	log.Printf("open connection to db " + dbname + "\n")
	db, err = sql.Open("mysql", get_dsn(dbname))
	if err != nil {
		log.Printf("Error %s when opening mysql dsn", err)
		return
	}

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
	log.Printf("Connected to DB %s successfully\n", dbname)

}

func get_dsn(dsnDbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, hostname, port, dsnDbName)
}

func createTableIfNotExists() {

	//	query := fmt.Sprintf("CREATE TABLE %s (`%s` int(6) unsigned NOT NULL AUTO_INCREMENT, `%s` varchar(30) NOT NULL, PRIMARY KEY (`%s`)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;", tableName, "json:\"id\"", "json:\"title\"", "json:\"id\"")
	// query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s.%s (`%s` int(6) unsigned NOT NULL AUTO_INCREMENT, `%s` varchar(30) NOT NULL, PRIMARY KEY (`%s`)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;", dbname, tablename, "json:\"id\"", "json:\"title\"", "json:\"id\"")
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s.%s (`%s` int(6) unsigned NOT NULL AUTO_INCREMENT, `%s` varchar(30) NOT NULL, PRIMARY KEY (`%s`)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;", dbname, tablename, "id", "title", "id")
	insert, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
	log.Printf("Executed query %s to table %s successfully\n", query, tablename)
}

func insertTestValuesIntoTable(theDB *sql.DB) {

	query := fmt.Sprintf("INSERT INTO %s VALUES('1', 'Ghostbusters')", tablename)
	insert, err := theDB.Query(query)
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
	log.Printf("Executed query %s to table %s successfully\n", query, tablename)
}

func getPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var posts []Post
	result, err := db.Query("SELECT id, title from " + tablename)
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
	log.Printf("Creating post\n")
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	// title := "'" + keyVal["title"] + "'"
	title := keyVal["title"]
	log.Printf("body = " + string(body) + "\n")
	log.Printf("id = " + keyVal["id"] + "\n")
	log.Printf("title = " + title + "\n")
	if err != nil {
		panic(err.Error())
	}

	log.Printf("Preparing statement\n")
	// stmt, err := db.Prepare("INSERT INTO " + tablename + "(title) VALUES(" + title + ");")
	// _, err = stmt.Exec()
	stmt, err := db.Prepare("INSERT INTO " + tablename + "(title) VALUES(?);")
	_, err = stmt.Exec(title)
	if err != nil {
		panic(err.Error())
	}

	log.Printf("Prepared statement successfully\n")

	// var responseText = "New post with title " + string(keyVal["title"]) + " was created"
	var responseText = "New post with title " + title + " was created in database " + dbname + ", table " + tablename
	fmt.Fprintf(w, responseText)
	log.Printf(responseText)
}

func createPost3(w http.ResponseWriter, r *http.Request) {
	log.Printf("Creating post\n")
	w.Header().Set("Content-Type", "application/json")
	log.Printf("Preparing statement\n")
	stmt, err := db.Prepare("INSERT INTO " + tablename + "(title) VALUES(?)")
	if err != nil {
		panic(err.Error())
	}
	log.Printf("Prepared statement successfully\n")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	title := keyVal["title"]
	_, err = stmt.Exec(title)
	log.Printf("body = " + string(body) + "\n")
	log.Printf("id = " + keyVal["id"] + "\n")
	log.Printf("title = " + keyVal["title"] + "\n")
	if err != nil {
		panic(err.Error())
	}
	var responseText = "New post with title " + string(keyVal["title"]) + " was created"
	fmt.Fprintf(w, responseText)
}

func createPost2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//	log.Printf("Preparing DB\n")
	// stmt, err := db.Prepare(fmt.Sprintf("INSERT INTO "+tablename+"(id, title) VALUES(%s, %s);", "1", "trustme"))
	// stmt, err := db.Prepare(fmt.Sprintf("INSERT INTO "+tablename+" VALUES(%s, '%s');", "17", "jeans"))
	stmt, err := db.Prepare("INSERT INTO " + tablename + " VALUES(9, 'patata');")
	if err != nil {
		panic(err.Error())
	}
	log.Printf("posted to DB successfully stmt %s\n", stmt)
	//*
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
	//	fmt.Fprintf(w, "New post was created")
	var newPostReport = "new post " + strconv.Itoa(99) + " was created"
	fmt.Fprintf(w, newPostReport)
	log.Printf("created post\n")
	//*/
}
func getPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	result, err := db.Query("SELECT id, title FROM "+tablename+" WHERE id = ?", params["id"])
	// result, err := db.Query("SELECT id, title FROM "+tablename+" WHERE id = %s", params["id"])
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
	log.Printf("got post %s/%s (id %s)\n", post.Title, params["title"], params["id"])
}
func updatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	stmt, err := db.Prepare("UPDATE " + tablename + " SET title = ? WHERE id = ?")
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
	stmt, err := db.Prepare("DELETE FROM " + tablename + " WHERE id = ?")
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
