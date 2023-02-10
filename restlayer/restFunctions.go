package restlayer

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"testserver/database/dbtools"
	"testserver/database/model"

	"github.com/gorilla/mux"
)

// ---
// Rest Layer API Websites
// ---

func SelectWebsiteBasedName(response http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)

	name, ok := vars["name"]

	if !ok {
		response.WriteHeader(http.StatusBadRequest)

		fmt.Fprintln(response, "Website not found.")
	}

	website, err := dbtools.SelectWebsiteBasedName(name)

	if err != nil {
		json.NewEncoder(response).Encode("Website not found.")
	} else {
		json.NewEncoder(response).Encode(website)
	}

}

func SelectWebsiteBasedId(response http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)

	id, ok := vars["id"]

	if !ok {
		response.WriteHeader(http.StatusBadRequest)

		fmt.Fprintln(response, "Website not found.")
	}

	website, err := dbtools.SelectWebsiteBasedId(id)

	if err != nil {
		json.NewEncoder(response).Encode("Website not found.")
	} else {
		json.NewEncoder(response).Encode(website)
	}
}

func SelectWebsiteBasedAdress(response http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)

	adress, ok := vars["adress"]

	if !ok {
		response.WriteHeader(http.StatusBadRequest)

		fmt.Fprintln(response, "Website not found.")
	}

	website, err := dbtools.SelectWebsiteBasedAdress(adress)

	if err != nil {
		json.NewEncoder(response).Encode("Website not found.")
	} else {
		json.NewEncoder(response).Encode(website)
	}

}

func SelectAllWebsites(response http.ResponseWriter, request *http.Request) {

	websites := dbtools.SelectAllWebsites()

	json.NewEncoder(response).Encode(websites)
}

func SaveWebsite(response http.ResponseWriter, request *http.Request) {

	var website model.Website

	err := json.NewDecoder(request.Body).Decode(&website)

	if err != nil {
		fmt.Println(err)

		response.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(response, "Could not add new Website by error: %v.", err)
		return
	}

	websiteCheck := dbtools.SaveWebsite(website)

	json.NewEncoder(response).Encode(websiteCheck)

}

func UpdateWebsite(response http.ResponseWriter, request *http.Request) {

	var website model.Website

	err := json.NewDecoder(request.Body).Decode(&website)

	if err != nil {
		fmt.Println(err)

		response.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(response, "Could not update website by error: %v.", err)
		return
	}

	websiteCheck := dbtools.UpdateWebsite(website)

	json.NewEncoder(response).Encode(websiteCheck)

}

func DeleteWebsiteId(response http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)

	id, ok := vars["id"]

	if !ok {
		response.WriteHeader(http.StatusBadRequest)

		fmt.Fprintln(response, "Id website not found.")
	}

	// convert string to int
	idWebsite, err := strconv.Atoi(id)

	if err != nil {
		fmt.Println("Cannot convert string to int.")
	}

	website := dbtools.DeleteWebsiteId(idWebsite)

	json.NewEncoder(response).Encode(website)

}

func DeleteAllWebsites(response http.ResponseWriter, request *http.Request) {

	website := dbtools.DeleteAllWebsites()

	json.NewEncoder(response).Encode(website)

}

func SaveMultipleWebsites(response http.ResponseWriter, request *http.Request) {

	var websites []model.Website

	err := json.NewDecoder(request.Body).Decode(&websites)

	if err != nil {
		fmt.Println(err)

		response.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(response, "Could not update website by error: %v.", err)
		return
	}

	dbtools.SaveMultipleWebsites(websites)

}

func CheckAllWebsites(response http.ResponseWriter, request *http.Request) {

	websites := dbtools.SelectAllWebsites()

	for i := 0; i < len(websites); i++ {
		_, err := http.Get("http://" + websites[i].Adress)
		if err != nil {
			websites[i].Check = false
		} else {
			websites[i].Check = true
		}
	}

	json.NewEncoder(response).Encode(websites)

}

// ---
// HTML Files Template Handling
// ---

// --- These functions allow to show the data in HTML format on the same domain (CORS policy avoid) ---

func ShowHtmlFile(response http.ResponseWriter, request *http.Request) {

	http.ServeFile(response, request, request.URL.Path[1:])

}

// ---
// HTML Template Handling
// ---

func indexTemplateHandling(response http.ResponseWriter, request *http.Request) {

	// Initialize a slice containing the paths to the two files. It's important
	// to note that the file containing our base template must be the *first*
	// file in the slice.
	files := []string{
		"templates/base_template.html",
		"templates/index_template.html",
	}

	// Use the template.ParseFiles() function to read the files and store the
	// templates in a template set. Notice that we can pass the slice of file
	// paths as a variadic parameter?
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(response, "Internal Server Error", 500)
		return
	}

	// Use the ExecuteTemplate() method to write the content of the "base"
	// template as the response body.
	err = ts.ExecuteTemplate(response, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(response, "Internal Server Error", 500)
	}
}

// ---
// Website HTML Template Handling
// ---

func createWebsiteTemplateHandling(response http.ResponseWriter, request *http.Request) {

	files := []string{
		"templates/base_template.html",
		"templates/website/create_website_template.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(response, "Internal Server Error", 500)
		return
	}

	err = ts.ExecuteTemplate(response, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(response, "Internal Server Error", 500)
	}
}

func deleteAllWebsitesTemplateHandling(response http.ResponseWriter, request *http.Request) {

	files := []string{
		"templates/base_template.html",
		"templates/website/delete_all_website_template.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(response, "Internal Server Error", 500)
		return
	}

	err = ts.ExecuteTemplate(response, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(response, "Internal Server Error", 500)
	}
}

func deleteIdWebsiteTemplateHandling(response http.ResponseWriter, request *http.Request) {

	files := []string{
		"templates/base_template.html",
		"templates/website/delete_id_website_template.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(response, "Internal Server Error", 500)
		return
	}

	err = ts.ExecuteTemplate(response, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(response, "Internal Server Error", 500)
	}
}

func editIdWebsiteTemplateHandling(response http.ResponseWriter, request *http.Request) {

	files := []string{
		"templates/base_template.html",
		"templates/website/edit_website_template.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(response, "Internal Server Error", 500)
		return
	}

	err = ts.ExecuteTemplate(response, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(response, "Internal Server Error", 500)
	}
}

func showAllWebsitesTemplateHandling(response http.ResponseWriter, request *http.Request) {

	files := []string{
		"templates/base_template.html",
		"templates/website/show_all_website_template.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(response, "Internal Server Error", 500)
		return
	}

	err = ts.ExecuteTemplate(response, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(response, "Internal Server Error", 500)
	}
}

func checkWebsiteTemplateHandling(response http.ResponseWriter, request *http.Request) {

	files := []string{
		"templates/base_template.html",
		"templates/website/check_website_template.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(response, "Internal Server Error", 500)
		return
	}

	err = ts.ExecuteTemplate(response, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(response, "Internal Server Error", 500)
	}
}
