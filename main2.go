package main2

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
	username  = "paulo"
	password  = "pinkel"
	hostname  = "127.0.0.1:3306"
	dbname    = "neueDatabase"
	tablename = "postsTable17"
)

func main() {

	db, err = sql.Open("mysql", get_dsn(""))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	log.Printf("Opened mysql DSN successfully\n")

	connectToDB()
	//createTableIfNotExists()

	/*
		query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s.%s (`%s` int(6) unsigned NOT NULL AUTO_INCREMENT, `%s` varchar(30) NOT NULL, PRIMARY KEY (`%s`)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;", dbname, tablename, "id", "title", "id")
		insert, err := db.Query(query)
		if err != nil {
			panic(err.Error())
		}
		//	defer insert.Close()
		log.Printf("Executed query %s to table %s with %s successfully\n", query, tablename, insert)

		//*/

	//*
	stmt, err := db.Prepare("INSERT INTO " + dbname + "." + tablename + "(title) VALUES(?)")
	if err != nil {
		//		panic(err.Error())
		log.Printf("Error %s at Insert into table name %s at stmt %s", err, tablename, stmt)
		//		createTable(db, dbname, tablename)
		log.Printf("Table %s created DB\n", tablename)
	}
	//*/

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
	http.ListenAndServe(":8000", router)
	log.Printf("API open and running\n")

}

// func connectToDB(targetDbName string) {
func connectToDB() {
	/*
			db, err = sql.Open("mysql", get_dsn(""))
			if err != nil {
				log.Printf("Error %s when opening DB %S", err, targetDbName)
				return
			}
			contexti, cnclfnc := context.WithTimeout(context.Background(), 5*time.Second)
			defer cnclfnc()
			//	perms, err := db.ExecContext(ctx, "GRANT ALL ON '"+targetDbName+"' TO '"+username+"'@'localhost';")
			perms, err := db.ExecContext(contexti, "GRANT ALL PRIVILEGES ON *.* TO '"+username+"'@'localhost' IDENTIFIED BY '"+password+"';")
			if err != nil {
				log.Printf("Error %s when granting user permissions to %s on localhost\n", err, username)
				return
			}
			log.Printf("permissions granted for %d\n", perms)
			db.Close()

		//	db, err = sql.Open("mysql", get_dsn(targetDbName))
		db, err = sql.Open("mysql", get_dsn(""))
		if err != nil {
			log.Printf("Error %s when opening mysql dsn", err)
			return
		}
	*/

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	res, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+dbname)
	if err != nil {
		log.Printf("Error %s when cueureating DB\n", err)
		return
	}

	no, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when fetching rows", err)
		//		return
	}
	log.Printf("rows affected %d\n", no)

	db.Close()

	//	db, err = sql.Open("mysql", get_dsn(targetDbName))
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
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dsnDbName)
}

func createTableIfNotExists() {

	//	query := fmt.Sprintf("CREATE TABLE %s (`%s` int(6) unsigned NOT NULL AUTO_INCREMENT, `%s` varchar(30) NOT NULL, PRIMARY KEY (`%s`)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;", tableName, "json:\"id\"", "json:\"title\"", "json:\"id\"")
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s.%s (`%s` int(6) unsigned NOT NULL AUTO_INCREMENT, `%s` varchar(30) NOT NULL, PRIMARY KEY (`%s`)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;", dbname, tablename, "id", "title", "id")
	insert, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
	log.Printf("Executed query %s to table %s successfully\n", query, tablename)
}

func insertTestValuesIntoTable(theDB *sql.DB, targetDbName string, newTableName string) {

	query := fmt.Sprintf("INSERT INTO %s VALUES('1', 'Merkel', 'Ghostbusters')", newTableName)
	insert, err := theDB.Query(query)
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
	log.Printf("Executed query %s to table %s successfully\n", query, newTableName)
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
	w.Header().Set("Content-Type", "application/json")
	//	log.Printf("Preparing DB\n")
	stmt, err := db.Prepare("INSERT INTO " + tablename + "(title) VALUES(?)")
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
	//	fmt.Fprintf(w, "New post was created")
	var newPostReport = "new post " + strconv.Itoa(99) + " was created"
	fmt.Fprintf(w, newPostReport)
	log.Printf("created post\n")
}
func getPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	result, err := db.Query("SELECT id, title FROM "+tablename+" WHERE id = ?", params["id"])
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
