package main

import (
	"fmt"
	"sort"
)

// =============================================================================
// Exercice perso — Système de ligue de joueurs
// Démontre : groupe const + iota, slices, la boucle for unique, fallthrough
// =============================================================================

// --- Groupe const avec iota n°1 : les rangs (type nommé + méthode String) ---
type Rang int

const (
	Bronze  Rang = iota // 0
	Argent              // 1
	Or                  // 2
	Platine             // 3
	Diamant             // 4
)

func (r Rang) String() string {
	return [...]string{"Bronze", "Argent", "Or", "Platine", "Diamant"}[r]
}

// --- Groupe const avec iota n°2 : paliers d'XP (iota avec calcul) ---
const (
	XPBronze  = iota * 1000 // 0
	XPArgent                // 1000
	XPOr                    // 2000
	XPPlatine               // 3000
	XPDiamant               // 4000
)

type Joueur struct {
	Nom string
	XP  int
}

// rangDepuisXP utilise la boucle for en mode "while" : on monte de palier
// tant que l'XP est suffisante.
func rangDepuisXP(xp int) Rang {
	paliers := []int{XPBronze, XPArgent, XPOr, XPPlatine, XPDiamant}
	niveau := 0
	for niveau+1 < len(paliers) && xp >= paliers[niveau+1] {
		niveau++
	}
	return Rang(niveau)
}

// recompenses utilise fallthrough : chaque rang hérite des récompenses de
// tous les rangs inférieurs (effet "cumulatif").
func recompenses(r Rang) []string {
	var perks []string
	switch r {
	case Diamant:
		perks = append(perks, "Skin légendaire")
		fallthrough
	case Platine:
		perks = append(perks, "Badge animé")
		fallthrough
	case Or:
		perks = append(perks, "Pseudo doré")
		fallthrough
	case Argent:
		perks = append(perks, "Accès au salon privé")
		fallthrough
	case Bronze:
		perks = append(perks, "Bienvenue dans la ligue !")
	}
	return perks
}

func ligne(n int) {
	// for classique (à la C) : pour tracer un séparateur
	for i := 0; i < n; i++ {
		fmt.Print("=")
	}
	fmt.Println()
}

func main() {
	// --- Slices : construction dynamique avec append ---
	joueurs := []Joueur{
		{"Alice", 3500},
		{"Bob", 800},
		{"Chloé", 4200},
	}
	joueurs = append(joueurs, Joueur{"David", 2100}, Joueur{"Emma", 50})

	// for range : itérer sur le slice, index + valeur
	fmt.Println("CLASSEMENT & RÉCOMPENSES")
	ligne(40)
	for i, j := range joueurs {
		r := rangDepuisXP(j.XP)
		fmt.Printf("%d. %-6s | %5d XP | %s\n", i+1, j.Nom, j.XP, r)
		for _, perk := range recompenses(r) {
			fmt.Printf("      - %s\n", perk)
		}
	}
	ligne(40)

	// --- Slicing : podium (top 3) après tri décroissant par XP ---
	sort.Slice(joueurs, func(a, b int) bool {
		return joueurs[a].XP > joueurs[b].XP
	})
	podium := joueurs[:3] // sous-slice des 3 premiers
	places := []string{"1er", "2e", "3e"}
	fmt.Println("PODIUM")
	for i, j := range podium {
		fmt.Printf("%s : %s (%d XP, %s)\n", places[i], j.Nom, j.XP, rangDepuisXP(j.XP))
	}
	ligne(40)

	// --- for infini + break : simulation de grind jusqu'au rang suivant ---
	debutant := Joueur{"Newbie", 50}
	rangInitial := rangDepuisXP(debutant.XP)
	sessions := 0
	for {
		debutant.XP += 250 // +250 XP par session
		sessions++
		if rangDepuisXP(debutant.XP) > rangInitial {
			break
		}
	}
	fmt.Printf("%s est passé de %s à %s en %d sessions (XP=%d)\n",
		debutant.Nom, rangInitial, rangDepuisXP(debutant.XP), sessions, debutant.XP)
}
