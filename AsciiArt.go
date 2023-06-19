package main

import (
	"fmt"
//	"log"
	"os"
	"strings"
)

func Match(r rune, i int, ascii map[byte][]string) string {
	var str string
	for ind, v := range ascii {
		if rune(ind) == r {
			str += v[i]
		}
	}
	return str
}

func NewLine(tab []string) bool {
	for i := 0; i < len(tab); i++ {
		if tab[i] != "" {
			return false
		}
	}
	return true
}

func Printable(tab []rune) bool {
	for i := 0; i < len(tab); i++ {
		if tab[i] < 32 || tab[i] > 126 {
			return false
		}
	}
	return true
}
func Banner(s string) string {
	return "ascii-art/" + s + ".txt"
}

func cooking(s string, option string) (string, string) {
	//largeur, _, _ := term.GetSize(0)
	ascii := make(map[byte][]string)
	var index byte = 32
	Erreur:=""
	banner := Banner(option)

	file, err := os.ReadFile(banner)
	if err != nil {
		Erreur="500"
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
			Erreur="400"
			fmt.Println("Error : Non-displayable character !!!")
			break
		}

	}

	return affiche, Erreur
}
