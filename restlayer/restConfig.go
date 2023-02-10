package restlayer

import (
	"net/http"

	"github.com/gorilla/mux"
)

func restConfig(router *mux.Router) {

	restRouter := router.PathPrefix("/testserver/api").Subrouter()

	httpRouter := router.PathPrefix("").Subrouter()

	// ------------
	// Websites API
	// ------------

	// localhost:8080/testserver/api/websites
	restRouter.Methods("GET").Path("/websites").HandlerFunc(SelectAllWebsites)

	//localhost:8080/testserver/api/website/{name}
	restRouter.Methods("GET").Path("/website/{name}").HandlerFunc(SelectWebsiteBasedName)

	//localhost:8080/testserver/api/website/id/{id}
	restRouter.Methods("GET").Path("/website/id/{id}").HandlerFunc(SelectWebsiteBasedId)

	//localhost:8080/testserver/api/website/adress/{adress}
	restRouter.Methods("GET").Path("/website/adress/{adress}").HandlerFunc(SelectWebsiteBasedAdress)

	//localhost:8080/testserver/api/website/add
	restRouter.Methods("PUT").Path("/website/add").HandlerFunc(SaveWebsite)

	//localhost:8080/testserver/api/website/edit
	restRouter.Methods("POST").Path("/website/edit").HandlerFunc(UpdateWebsite)

	//localhost:8080/testserver/api/website/deleteid/{id}
	restRouter.Methods("DELETE").Path("/website/deleteid/{id}").HandlerFunc(DeleteWebsiteId)

	//localhost:8080/testserver/api/website/deleteall
	restRouter.Methods("DELETE").Path("/website/deleteall").HandlerFunc(DeleteAllWebsites)

	//localhost:8080/testserver/api/website/addmultiple
	restRouter.Methods("POST").Path("/website/addmultiple").HandlerFunc(SaveMultipleWebsites)

	//localhost:8080/testserver/api/websites/checkall
	restRouter.Methods("GET").Path("/websites/checkall").HandlerFunc(CheckAllWebsites)

	// ---
	// File Show System
	// ---

	//localhost:8080/configHtmlServer.json
	httpRouter.Methods("GET").Path("/configHtmlServer.json").HandlerFunc(ShowHtmlFile)

	// ---
	// HTML Files Show
	// ---

	//localhost:8080/index.html
	httpRouter.HandleFunc("/", indexTemplateHandling)

	// ---
	// HTML Files Show websites
	// ---

	//localhost:8080/create_website.html
	httpRouter.HandleFunc("/create_website.html", createWebsiteTemplateHandling)

	//localhost:8080/delete_all_website.html
	httpRouter.HandleFunc("/delete_all_website.html", deleteAllWebsitesTemplateHandling)

	//localhost:8080/delete_id_website.html
	httpRouter.HandleFunc("/delete_id_website.html", deleteIdWebsiteTemplateHandling)

	//localhost:8080/edit_id_website.html
	httpRouter.HandleFunc("/edit_id_website.html", editIdWebsiteTemplateHandling)

	//localhost:8080/show_all_website.html
	httpRouter.HandleFunc("/show_all_website.html", showAllWebsitesTemplateHandling)

	//localhost:8080/check_website.html
	httpRouter.HandleFunc("/check_website.html", checkWebsiteTemplateHandling)
}

func RestStart(endpoint string) error {

	router := mux.NewRouter()

	restConfig(router)

	return http.ListenAndServe(endpoint, router)
}
