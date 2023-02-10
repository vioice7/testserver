package dbtools

import (
	"database/sql"
	"fmt"
	"log"
	"testserver/database/model"

	_ "github.com/go-sql-driver/mysql"
)

//
// Conection and initialase mysql
//

var driverName string
var dataSourceName string

func DBInitilize(dn, dsn string) {

	driverName = dn
	dataSourceName = dsn

}

func connect() (db *sql.DB) {

	db, err := sql.Open(driverName, dataSourceName)

	if err != nil {
		log.Fatal(err.Error())
	}

	return db
}

// -------------------------
// Website table connection
// -------------------------

func SelectAllWebsites() []model.Website {

	db := connect()

	rows, err := db.Query("select * from websites")

	if err != nil {
		log.Fatal(err.Error())
	}

	defer db.Close()

	websites := []model.Website{}

	for rows.Next() {

		website := model.Website{}

		err = rows.Scan(&website.ID, &website.Name, &website.Adress)

		if err != nil {
			log.Fatal(err.Error())
			continue
		}

		websites = append(websites, website)
	}

	return websites

}

func SelectWebsiteBasedName(name string) (model.Website, error) {

	db := connect()

	rows := db.QueryRow("select * from websites where name = ?", name)

	defer db.Close()

	website := model.Website{}

	err := rows.Scan(&website.ID, &website.Name, &website.Adress)

	return website, err

}

func SelectWebsiteBasedId(id string) (model.Website, error) {

	db := connect()

	rows := db.QueryRow("select * from websites where id = ?", id)

	defer db.Close()

	website := model.Website{}

	err := rows.Scan(&website.ID, &website.Name, &website.Adress)

	return website, err

}

func SelectWebsiteBasedAdress(adress string) (model.Website, error) {

	db := connect()

	defer db.Close()

	rows := db.QueryRow("select * from websites where adress = ?", adress)

	website := model.Website{}

	err := rows.Scan(&website.ID, &website.Name, &website.Adress)

	return website, err

}

func SaveWebsite(website model.Website) int64 {

	db := connect()

	defer db.Close()

	websiteCheck := model.Website{}

	rows := db.QueryRow("select adress from websites where adress = ?", website.Adress)

	errDupl := rows.Scan(&websiteCheck.Adress)

	// if there is no record don't add a record and return

	if website.Adress == "" {
		return 0
	}

	// if there is a duplicate entry don't add a record and return

	if errDupl == nil {
		fmt.Println("Duplicate entry in database! No records added!")
		return 0
	}

	save, err := db.Prepare("insert into websites(id,name,adress) values(?,?,?)")

	if err != nil {
		log.Fatal(err.Error())
	}

	result, err := save.Exec(website.ID, website.Name, website.Adress)

	if err != nil {
		log.Fatal(err.Error())
	}

	websiteID, err := result.LastInsertId()

	if err != nil {
		log.Fatal(err.Error())
	}

	return websiteID
}

func UpdateWebsite(website model.Website) int64 {

	db := connect()

	defer db.Close()

	// make new model for comparing

	websiteCheck := model.Website{}

	// check if there are records to compare

	rowsCheck := db.QueryRow("select id from websites where id = ?", website.ID)

	errCheck := rowsCheck.Scan(&websiteCheck.ID)

	if errCheck != nil {
		fmt.Println("Error, no records with that id available!")
		return 0
	}

	// check for unique adress

	checkAdressID := ""

	rowsDuplAdressCheck := db.QueryRow("select adress from websites where adress = ? and id = ?", website.Adress, website.ID)

	errDuplAdressID := rowsDuplAdressCheck.Scan(&websiteCheck.Adress)

	if errDuplAdressID != nil {

		fmt.Println(errDuplAdressID.Error())

		checkAdressID = errDuplAdressID.Error()

	}

	if checkAdressID == "sql: no rows in result set" {

		rowsDuplAdressCheck := db.QueryRow("select adress from websites where adress = ?", website.Adress)

		errDuplAdress := rowsDuplAdressCheck.Scan(&websiteCheck.Adress)

		if errDuplAdress != nil {
			fmt.Println(errDuplAdress.Error())
		} else {
			return 0
		}

	}

	// create a new model to compare

	websiteData := model.Website{}

	rows := db.QueryRow("select * from websites where id = ?", website.ID)

	err := rows.Scan(&websiteData.ID, &websiteData.Name, &websiteData.Adress)

	if err != nil {
		log.Fatal(err.Error())
	}

	// compare the data between models in order to remove blank data

	if website.Name == "" {
		website.Name = websiteData.Name
	}

	if website.Adress == "" {
		website.Adress = websiteData.Adress
	}

	update, err := db.Prepare("update websites set name=?, adress=? where id=?")

	if err != nil {
		log.Fatal(err.Error())
	}

	result, err := update.Exec(website.Name, website.Adress, website.ID)

	if err != nil {
		log.Fatal(err.Error())
	}

	rowsEffected, err := result.RowsAffected()

	if err != nil {
		log.Fatal(err.Error())
	}

	return rowsEffected

}

func DeleteWebsiteId(id int) int64 {

	db := connect()

	defer db.Close()

	delete, err := db.Prepare("delete from websites where id = ?")

	if err != nil {
		log.Fatal(err.Error())
	}

	result, err := delete.Exec(id)

	if err != nil {
		log.Fatal(err.Error())
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		log.Fatal(err.Error())
	}

	return rowsAffected
}

func DeleteAllWebsites() int64 {

	db := connect()

	defer db.Close()

	// count records before delete table

	var count int64

	err := db.QueryRow("select count(*) from websites").Scan(&count)

	if err != nil {
		log.Fatal(err.Error())
	}

	// delete table to delete all records

	deleteAll, err := db.Prepare("delete from websites where id > 0")

	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = deleteAll.Exec()

	if err != nil {
		log.Fatal(err.Error())
	}

	alterTable, err := db.Prepare("alter table websites auto_increment = 1")

	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = alterTable.Exec()

	if err != nil {
		log.Fatal(err.Error())
	}

	return count

}

func SaveMultipleWebsites(websites []model.Website) int64 {

	db := connect()

	defer db.Close()

	var nrAddedWebsites int64

	for _, website := range websites {

		// check for unique adress

		rowsDupl := db.QueryRow("select adress from websites where adress = ?", website.Adress)

		errDupl := rowsDupl.Scan(&website.Adress)

		// if there is a duplicate entry don't add a record and return

		if errDupl == nil {
			fmt.Println("Duplicate entry in database! No records added!")
			continue
		}

		save, err := db.Prepare("insert into websites(id,name,adress) values(?,?,?)")

		if err != nil {
			log.Fatal(err.Error())
		}

		result, err := save.Exec(website.ID, website.Name, website.Adress)

		if err != nil {
			log.Fatal(err.Error())
		}

		websiteID, err := result.LastInsertId()

		if err != nil {
			log.Fatal(err.Error())
		}

		fmt.Println("Website ID ", websiteID, " inserted!")

	}

	return nrAddedWebsites

}
