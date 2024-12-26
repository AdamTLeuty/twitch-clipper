package main

import (
	"html/template"
	"net/http"
)

// Define a struct for passing data to the template
type PageData struct {
	Username string
}

func home(w http.ResponseWriter, r *http.Request) {
	// Parse the HTML template
	tmpl := template.Must(template.ParseFiles("template.html"))

	// Data to be passed into the template
	data := PageData{
		Username: "Testy",
	}

	// Execute the template with data
	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", home)
	// Handle form submission
	http.HandleFunc("/twitchURL", streamPickerHandler)

	// Serve static files if needed (e.g., CSS, images)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		println("Error starting server:", err)
	}
}
