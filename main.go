package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
)

func error500Handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	tmpl := template.Must(template.ParseFiles("template/500.html"))
	tmpl.Execute(w, nil)
}

func error404Handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	tmpl := template.Must(template.ParseFiles("template/404.html"))
	tmpl.Execute(w, nil)
}

func error400Handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	tmpl := template.Must(template.ParseFiles("template/400.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		handleError(w, r, http.StatusInternalServerError)
		return
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("template/index.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		handleError(w, r, http.StatusInternalServerError)
		return
	}
}

func asciiArtHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	text := r.FormValue("text")
	banner := r.FormValue("banner")
	asciiArt := ""
	asciiArt, erreur := cooking(text, banner)
	if erreur == "400" {
		error400Handler(w, r)
		return
	}
	if erreur == "500" {
		handleError(w, r, http.StatusInternalServerError)
		return
	}
	fmt.Println(asciiArt)

	tmpl2 := template.Must(template.ParseFiles("template/Art2.html"))
	err1 := tmpl2.Execute(w, struct {
		Text     string
		Banner   string
		AsciiArt string
	}{
		Text:     text,
		Banner:   banner,
		AsciiArt: asciiArt,
	})
	if err1 != nil {
		handleError(w, r, http.StatusInternalServerError)
		return
	}
	file, err2 := os.Create("template/ascii.txt")
	if err2 != nil {
		handleError(w, r, http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Écrire le contenu dans le fichier texte
	_, err := file.WriteString(asciiArt)
	if err != nil {
		handleError(w, r, http.StatusInternalServerError)
		return
	}
}

func exportHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.ReadFile("template/ascii.txt")
	if err != nil {
		handleError(w, r, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Disposition", "attachment; filename=ascii.txt")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(file)))

	// Transférer le fichier
	w.Write(file)
}

func handleError(w http.ResponseWriter, r *http.Request, statusCode int) {
	switch statusCode {
	case http.StatusNotFound:
		error404Handler(w, r)
	case http.StatusInternalServerError:
		// Récupérer le code d'erreur interne du serveur
		errorCode := http.StatusInternalServerError

		// Votre traitement en fonction du code d'erreur interne
		switch errorCode {
		case http.StatusInternalServerError:
			// Faites quelque chose en cas d'erreur interne spécifique
			error500Handler(w, r)
		default:
			// Faites autre chose en cas d'erreur interne différente
			http.Error(w, "Autre erreur interne", errorCode)
		}
	case http.StatusBadRequest:
		error400Handler(w, r)
	// Ajoutez d'autres cas pour les autres codes d'erreur si nécessaire
	default:
		// Faites quelque chose en cas d'erreur inattendue
		error500Handler(w, r)
		// /http.Error(w, "Erreur inattendue", http.StatusInternalServerError)
	}
}

func errorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Erreur interne du serveur:", r)
				//handleError(w, r.(*http.Request), http.StatusInternalServerError)
				error500Handler(w, nil)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

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
