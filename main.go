package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	go startForm()

	connectToDB()
	createTableIfNotExists()

	/*
		stmt, err := db.Prepare("INSERT INTO " + dbname + "." + tablename + "(title) VALUES(?)")
		if err != nil {
			log.Printf("Table %s created DB\n", tablename)
			} else {
			log.Printf("Error %s at Insert into table name %s at stmt %s", err, tablename, stmt)
		}
		//*/

	//*
	runAPI()
	//defer db.Close()
	/*/
	// insertValuesIntoTable(db, dbname)
	//*/

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
