package main

import (
	"fmt"
	"os"
	"strings"
	"net/http"
	"html/template"
)
// Match retourne une chaîne de caractères correspondant à la rune spécifiée dans la map ASCII.
func Match(r rune, i int, ascii map[byte][]string) string {
	var str string
	for ind, v := range ascii {
		if rune(ind) == r {
			str += v[i]
		}
	}
	return str
}

// NewLine vérifie si le tableau de chaînes de caractères contient uniquement des lignes vides.
func NewLine(tab []string) bool {
	for i := 0; i < len(tab); i++ {
		if tab[i] != "" {
			return false
		}
	}
	return true
}

// Printable vérifie si tous les caractères du tableau de runes sont imprimables.
func Printable(tab []rune) bool {
	for i := 0; i < len(tab); i++ {
		if tab[i] < 32 || tab[i] > 126 {
			return false
		}
	}
	return true
}

// Banner retourne le chemin d'accès au fichier ASCII correspondant à la bannière spécifiée.
func Banner(s string) string {
	return "ascii-art/" + s + ".txt"
}

// cooking génère l'art ASCII à partir du texte et de la bannière spécifiés.
func cooking(s string, option string) (string, string) {
	//largeur, _, _ := term.GetSize(0)
	ascii := make(map[byte][]string)
	var index byte = 32
	Erreur := ""
	banner := Banner(option)

	file, err := os.ReadFile(banner)
	if err != nil {
		Erreur = "500"
		// log.Fatal("Error : Not a ascci file in the repertory")
	}
	if option == "thinkertoy" {
		Split := strings.Split(string(file), "\r\n")
		for i := 1; i+8 < len(Split); i += 9 {
			ascii[index] = Split[i : i+8]
			index++
		}
	} else {
		Split := strings.Split(string(file), "\n")
		for i := 1; i+8 < len(Split); i += 9 {
			ascii[index] = Split[i : i+8]
			index++
		}
	}

	tabascii := ascii
	var affiche string
	var split []string
	split = strings.Split(s, "\r\n")
	if NewLine(split) {
		split = split[:len(split)-1]
	}
	for _, v := range split {
		tabrune := []rune(v)
		if Printable(tabrune) {
			for j := 0; j < 8; j++ {
				for i := 0; i < len(tabrune); i++ {
					affiche += Match(tabrune[i], j, tabascii)
				}
				if len(tabrune) != 0 {
					affiche += "\n"
				} else {
					affiche += "\n"
					break
				}
			}
		} else {
			Erreur = "400"
			fmt.Println("Error : Non-displayable character !!!")
			break
		}

	}

	return affiche, Erreur
}

// error500Handler gère la réponse HTTP avec une erreur 500 (Internal Server Error).
func error500Handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	tmpl := template.Must(template.ParseFiles("template/500.html"))
	tmpl.Execute(w, nil)
}

// error404Handler gère la réponse HTTP avec une erreur 404 (Not Found).
func error404Handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	tmpl := template.Must(template.ParseFiles("template/404.html"))
	tmpl.Execute(w, nil)
}

// error400Handler gère la réponse HTTP avec une erreur 400 (Bad Request).
func error400Handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	tmpl := template.Must(template.ParseFiles("template/400.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		handleError(w, r, http.StatusInternalServerError)
		return
	}
}

// indexHandler gère la réponse HTTP pour la page d'accueil.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("template/index.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		handleError(w, r, http.StatusInternalServerError)
		return
	}
}

// asciiArtHandler gère la réponse HTTP pour la génération de l'art ASCII.
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

// exportHandler gère la réponse HTTP pour l'export du fichier ASCII.
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

// handleError gère les erreurs HTTP et redirige vers les gestionnaires d'erreurs appropriés.
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
			// Gère l'erreur 500 (Internal Server Error)
			error500Handler(w, r)
		default:
			// Gère les autres erreurs internes
			http.Error(w, "Autre erreur interne", errorCode)
		}
	case http.StatusBadRequest:
		error400Handler(w, r)
	default:
		// Gère les erreurs inattendues
		error500Handler(w, r)
	}
}

// errorHandler est un middleware pour gérer les erreurs lors du traitement des requêtes HTTP.
func errorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Erreur interne du serveur:", r)
				error500Handler(w, nil)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
