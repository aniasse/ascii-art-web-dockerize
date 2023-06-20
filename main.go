package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Serveur en cours d'exécution sur http://localhost:8080")

	// Gestion des fichiers statiques
	staticHandler := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", staticHandler))
	// Gestion des fichiers CSS et images
	http.Handle("/template/style.css", http.FileServer(http.Dir(".")))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))
	http.Handle("/template/", http.StripPrefix("/template/", http.FileServer(http.Dir("./template"))))

	// Routes avec gestion d'erreur
	http.Handle("/", errorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/", "/index.html", "/template":
			indexHandler(w, r)
		case "/template/Art.html/", "/ascii-art":
			asciiArtHandler(w, r)
		case "/ascii-art/download":
			exportHandler(w, r)
		default:
			error404Handler(w, r)
		}

	})))

	// Démarrage du serveur
	http.ListenAndServe(":8080", nil)
}
