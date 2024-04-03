package api

import (
	"fmt"
	"groupie-tracker/internal"
	"net/http"
)

func SetUpRouter() {
	port := "8080"
	fmt.Println("Server is running on http://localhost:" + port)

	NewRouter()

	http.ListenAndServe(":"+port, nil)
}

func NewRouter() {
	u := internal.Usecase{}
	c := Controller{Usecase: u}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/artist", c.GetAllArtist)
	http.HandleFunc("/artist/one", c.GetOneArtist)
	http.HandleFunc("/", NotFound)

}
