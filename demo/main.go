package main

import (
	"fmt"
	"html/template"
	"net/http"

	ja "github.com/io84team/jazzavatar-go"

	"github.com/gorilla/mux"
)

func main() {
	port := 3030

	r := mux.NewRouter()

	r.HandleFunc("/{name}/{size:[0-9]+}{radius:[/0-9]*}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		size := vars["size"]
		radius := vars["radius"]

		if len(radius) > 1 {
			radius = radius[1:]
		}

		ja, err := new(ja.Jazzavatar).Init(name, size, radius)
		if err != nil {
			fmt.Fprintln(w, "500 error:", err)
		} else {
			requestAvatar(ja, w, r)
		}

	})

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "<h1>Hello Jazzavatar!</h1><br/>Demo: <img src=\"/name/36\"/> <img src=\"/name2/72\"/> <img src=\"/name2/72/18\"/> <img src=\"/name2/72/36\"/>")
	})

	fmt.Println("Server is running on port:", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}

func requestAvatar(ja *ja.Jazzavatar, w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("avatar.xml"))

	w.Header().Add("Content-Type", "image/svg+xml")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Add("Access-Control-Allow-Methods", "GET")

	tmpl.Execute(w, ja)
}
