package main

import "fmt"

// Bonus : prénom personnalisé
const Nom = "Jeremy"

// Seuils de catégorie d'IMC
const (
	IMCMaigreur = 18.5
	IMCNormal   = 25.0
	IMCSurpoids = 30.0
)

func main() {
	// 1. Données : poids (kg) et taille (m)
	var poids float64 = 82
	var taille float64 = 1.89

	// 3. Calcul de l'IMC
	imc := poids / (taille * taille)

	// 4. Affichage de l'IMC avec 2 décimales
	fmt.Printf("%s, ton IMC est : %.2f\n", Nom, imc)

	// 5. Affichage de la catégorie
	var categorie string
	switch {
	case imc < IMCMaigreur:
		categorie = "Maigreur"
	case imc < IMCNormal:
		categorie = "Normal"
	case imc < IMCSurpoids:
		categorie = "Surpoids"
	default:
		categorie = "Obésité"
	}

	fmt.Printf("Catégorie : %s\n", categorie)
}
