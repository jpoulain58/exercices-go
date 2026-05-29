package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// =============================================================================
// Exercice noté 3 — Système de contacts (Personne, Employé, Étudiant)
// Bonus : struct tags JSON + maps
// =============================================================================

// --- 1. Personne (avec tags JSON) ---
type Personne struct {
	Prenom string `json:"prenom"`
	Nom    string `json:"nom"`
	Age    int    `json:"age"`
	Email  string `json:"email,omitempty"` // omis si vide
}

func (p Personne) NomComplet() string {
	return p.Prenom + " " + p.Nom
}

func (p Personne) Presentation() string {
	return fmt.Sprintf("%s, %d ans (%s)", p.NomComplet(), p.Age, p.Email)
}

// --- 2. Adresse ---
type Adresse struct {
	Rue        string `json:"rue"`
	Ville      string `json:"ville"`
	CodePostal string `json:"code_postal"`
}

func (a Adresse) Format() string {
	return fmt.Sprintf("%s, %s %s", a.Rue, a.CodePostal, a.Ville)
}

// --- 3. Employe : embedding de Personne + Adresse ---
type Employe struct {
	Personne          // méthodes promues : NomComplet(), Presentation()
	Adresse           // méthode promue : Format()
	Poste    string  `json:"poste"`
	Salaire  float64 `json:"salaire"`
}

func (e Employe) FicheEmploye() string {
	return fmt.Sprintf(
		"[EMPLOYÉ] %s\n  Poste   : %s\n  Salaire : %.2f €\n  Adresse : %s\n  Email   : %s",
		e.NomComplet(), e.Poste, e.Salaire, e.Format(), e.Email,
	)
}

// Pointer receiver : modifie le vrai objet
func (e *Employe) AugmenterSalaire(pct float64) {
	if pct <= 0 {
		return
	}
	e.Salaire = e.Salaire * (1 + pct/100)
}

// --- 4. Etudiant : embedding de Personne ---
type Etudiant struct {
	Personne         // méthodes promues
	Promo    string  `json:"promo"`
	Moyenne  float64 `json:"moyenne"`
}

func (e Etudiant) MentionObtenue() string {
	switch {
	case e.Moyenne >= 16:
		return "Très bien"
	case e.Moyenne >= 14:
		return "Bien"
	case e.Moyenne >= 12:
		return "Assez bien"
	case e.Moyenne >= 10:
		return "Passable"
	default:
		return "Insuffisant"
	}
}

func (e Etudiant) FicheEtudiant() string {
	return fmt.Sprintf(
		"[ÉTUDIANT] %s\n  Promo   : %s\n  Moyenne : %.2f/20\n  Mention : %s",
		e.NomComplet(), e.Promo, e.Moyenne, e.MentionObtenue(),
	)
}

func main() {
	// --- 5. Création de 2 employés et 2 étudiants ---
	employes := []Employe{
		{
			Personne: Personne{Prenom: "Alice", Nom: "Martin", Age: 34, Email: "alice@corp.fr"},
			Adresse:  Adresse{Rue: "12 rue de la Paix", Ville: "Paris", CodePostal: "75001"},
			Poste:    "Développeuse",
			Salaire:  45000,
		},
		{
			Personne: Personne{Prenom: "Bruno", Nom: "Lopez", Age: 41, Email: "bruno@corp.fr"},
			Adresse:  Adresse{Rue: "5 avenue des Tilleuls", Ville: "Lyon", CodePostal: "69003"},
			Poste:    "Développeur",
			Salaire:  52000,
		},
	}

	etudiants := []Etudiant{
		{Personne: Personne{Prenom: "Chloé", Nom: "Petit", Age: 21}, Promo: "M2", Moyenne: 16.4},
		{Personne: Personne{Prenom: "David", Nom: "Roux", Age: 22}, Promo: "M1", Moyenne: 11.8},
	}

	// Démo pointer receiver : on augmente le salaire d'Alice de 10 %
	employes[0].AugmenterSalaire(10)

	// Affichage des fiches
	fmt.Println(strings.Repeat("=", 50))
	for _, e := range employes {
		fmt.Println(e.FicheEmploye())
		fmt.Println(strings.Repeat("-", 50))
	}
	for _, e := range etudiants {
		fmt.Println(e.FicheEtudiant())
		fmt.Println(strings.Repeat("-", 50))
	}

	// --- MAPS : statistiques d'annuaire ---
	// Masse salariale par poste
	masseParPoste := map[string]float64{}
	for _, e := range employes {
		masseParPoste[e.Poste] += e.Salaire
	}
	fmt.Println("MASSE SALARIALE PAR POSTE")
	for poste, total := range masseParPoste {
		fmt.Printf("  %-14s : %.2f €\n", poste, total)
	}

	// Nombre d'étudiants par mention
	parMention := map[string]int{}
	for _, e := range etudiants {
		parMention[e.MentionObtenue()]++
	}
	fmt.Println("ÉTUDIANTS PAR MENTION")
	for mention, n := range parMention {
		fmt.Printf("  %-12s : %d\n", mention, n)
	}

	// Annuaire indexé par email (map clé → struct)
	annuaire := map[string]Personne{}
	for _, e := range employes {
		annuaire[e.Email] = e.Personne
	}
	if p, ok := annuaire["alice@corp.fr"]; ok {
		fmt.Println(strings.Repeat("=", 50))
		fmt.Println("Recherche annuaire (alice@corp.fr) :", p.Presentation())
	}

	// --- TAGS JSON : sérialisation d'un employé ---
	data, _ := json.MarshalIndent(employes[0], "", "  ")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("EXPORT JSON (struct tags) :")
	fmt.Println(string(data))
}
