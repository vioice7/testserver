package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testserver/database/dbtools"
	"testserver/restlayer"
)

type Configuration struct {
	DriverName     string `json:"driverName"`
	DataSourceName string `json:"dataSourceName"`
}

type AdressConfig struct {
	Adress string `json:"adress"`
}

func main() {

	fileDataBase, err := os.Open("database/configuration/config.json")

	if err != nil {
		log.Fatal(err.Error())
	}

	defer fileDataBase.Close()

	conf := new(Configuration)

	json.NewDecoder(fileDataBase).Decode(conf)

	dbtools.DBInitilize(conf.DriverName, conf.DataSourceName)

	fileAdress, err := os.Open("configHtmlServer.json")

	if err != nil {
		log.Fatal(err.Error())
	}

	defer fileAdress.Close()

	adress := new(AdressConfig)

	json.NewDecoder(fileAdress).Decode(adress)

	fmt.Println("Server started at adress: " + adress.Adress)

	restlayer.RestStart(adress.Adress)

}
