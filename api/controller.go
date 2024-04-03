package api

import (
	"encoding/json"
	"fmt"
	"groupie-tracker/domain"
	"groupie-tracker/internal"
	"html/template"
	"io/ioutil"
	"net/http"
)

type Controller struct {
	Usecase internal.Usecase
}

var templates = template.Must(template.ParseGlob("static/*"))

func (cl *Controller) GetOneArtist(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")

	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists/" + id)
	if err != nil {
		fmt.Println("No response from request")
		InternalError(w, r)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var art domain.Artist

	if err := json.Unmarshal(body, &art); err != nil {
		fmt.Println("Can not unmarshal JSON", err.Error())
		InternalError(w, r)
		return
	}
	if art.ID == 0 {
		NotFound(w, r)
		return
	}
	resp, err = http.Get("https://groupietrackers.herokuapp.com/api/relation/" + id)
	if err != nil {
		fmt.Println("No response from request")
		InternalError(w, r)
		return
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	body, err = ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("Can not unmarshal JSON", err.Error())
		InternalError(w, r)
		return
	}
	fmt.Println(data["datesLocations"])

	err = templates.ExecuteTemplate(w, "artist.html", struct {
		Artists  domain.Artist
		Relation interface{}
	}{art, data["datesLocations"]})
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		InternalError(w, r)
		return
	}

}

func (cl *Controller) GetAllArtist(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		fmt.Println("No response from request")
		InternalError(w, r)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var arts []domain.Artist
	if err := json.Unmarshal(body, &arts); err != nil {
		fmt.Println("Can not unmarshal JSON", err.Error())
		InternalError(w, r)
		return
	}

	err = templates.ExecuteTemplate(w, "index.html", struct{ Artists []domain.Artist }{arts})
	if err != nil {
		InternalError(w, r)
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	err := templates.ExecuteTemplate(w, "404.html", nil)
	if err != nil {
		http.Error(w, "Error rendering HTML", http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}
}

func BadRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	err := templates.ExecuteTemplate(w, "400.html", nil)
	if err != nil {
		http.Error(w, "Error rendering HTML", http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}
}
func InternalError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	err := templates.ExecuteTemplate(w, "500.html", nil)
	if err != nil {
		http.Error(w, "Error rendering HTML", http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}
}
