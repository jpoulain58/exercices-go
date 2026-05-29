package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// =============================================================================
// TP1 — TechShop : gestionnaire de catalogue d'une boutique de matériel info.
// =============================================================================

// --- DONNÉES ---

type Produit struct {
	ID        int
	Nom       string
	Marque    string
	Prix      float64
	Stock     int
	Categorie string
	Actif     bool
}

type Catalogue struct {
	produits []Produit
}

// --- FONCTIONNALITÉS (méthodes sur Catalogue) ---

// AjouterProduit ajoute un produit après validation ; erreur si données
// invalides ou ID déjà utilisé. Pointer receiver : modifie le catalogue.
func (c *Catalogue) AjouterProduit(p Produit) error {
	if p.ID <= 0 {
		return fmt.Errorf("ID invalide : %d (doit être > 0)", p.ID)
	}
	if strings.TrimSpace(p.Nom) == "" {
		return fmt.Errorf("le nom du produit est obligatoire")
	}
	if p.Prix < 0 {
		return fmt.Errorf("prix invalide : %.2f (ne peut pas être négatif)", p.Prix)
	}
	if p.Stock < 0 {
		return fmt.Errorf("stock invalide : %d (ne peut pas être négatif)", p.Stock)
	}
	for _, existant := range c.produits {
		if existant.ID == p.ID {
			return fmt.Errorf("ID %d déjà utilisé (produit « %s »)", p.ID, existant.Nom)
		}
	}
	c.produits = append(c.produits, p)
	return nil
}

// TrouverParID retourne le produit correspondant ou une erreur.
// Value receiver : lecture seule.
func (c Catalogue) TrouverParID(id int) (Produit, error) {
	for _, p := range c.produits {
		if p.ID == id {
			return p, nil
		}
	}
	return Produit{}, fmt.Errorf("aucun produit avec l'ID %d", id)
}

// TrouverParCategorie retourne tous les produits d'une catégorie
// (comparaison insensible à la casse).
func (c Catalogue) TrouverParCategorie(cat string) []Produit {
	var resultats []Produit
	for _, p := range c.produits {
		if strings.EqualFold(p.Categorie, cat) {
			resultats = append(resultats, p)
		}
	}
	return resultats
}

// AppliquerReduction applique un % de réduction à toute une catégorie
// et retourne le nombre de produits modifiés.
func (c *Catalogue) AppliquerReduction(categorie string, pct float64) int {
	if pct <= 0 || pct > 100 {
		return 0
	}
	modifies := 0
	for i := range c.produits {
		if strings.EqualFold(c.produits[i].Categorie, categorie) {
			nouveauPrix := c.produits[i].Prix * (1 - pct/100)
			c.produits[i].Prix = math.Round(nouveauPrix*100) / 100 // arrondi à 2 décimales
			modifies++
		}
	}
	return modifies
}

// Vendre réduit le stock d'un produit ; erreur si stock insuffisant.
func (c *Catalogue) Vendre(id int, qte int) error {
	if qte <= 0 {
		return fmt.Errorf("quantité invalide : %d", qte)
	}
	for i := range c.produits {
		if c.produits[i].ID == id {
			if qte > c.produits[i].Stock {
				return fmt.Errorf("stock insuffisant pour « %s » : %d demandé, %d disponible",
					c.produits[i].Nom, qte, c.produits[i].Stock)
			}
			c.produits[i].Stock -= qte
			return nil
		}
	}
	return fmt.Errorf("aucun produit avec l'ID %d", id)
}

// Rapport retourne un résumé : nb de produits et valeur totale du stock.
func (c Catalogue) Rapport() string {
	var valeurTotale float64
	for _, p := range c.produits {
		valeurTotale += p.Prix * float64(p.Stock)
	}
	return fmt.Sprintf("Catalogue : %d produit(s) | Valeur totale du stock : %.2f €",
		len(c.produits), valeurTotale)
}

// --- Affichage ---

func (p Produit) String() string {
	statut := "indisponible"
	if p.Actif && p.Stock > 0 {
		statut = "disponible"
	}
	return fmt.Sprintf("#%d %s %s | %.2f € | stock: %d | %s | %s",
		p.ID, p.Marque, p.Nom, p.Prix, p.Stock, p.Categorie, statut)
}

// --- Lecture d'entrée robuste ---

func lireLigne(r *bufio.Reader, invite string) string {
	fmt.Print(invite)
	texte, _ := r.ReadString('\n')
	return strings.TrimSpace(texte)
}

func lireInt(r *bufio.Reader, invite string) (int, error) {
	return strconv.Atoi(lireLigne(r, invite))
}

func lireFloat(r *bufio.Reader, invite string) (float64, error) {
	return strconv.ParseFloat(lireLigne(r, invite), 64)
}

// --- MAIN : menu CLI interactif ---

func main() {
	cat := &Catalogue{}

	// Pré-remplissage : 5 produits high-tech réalistes
	produitsInitiaux := []Produit{
		{ID: 1, Nom: "iPhone 15", Marque: "Apple", Prix: 999.99, Stock: 25, Categorie: "Smartphone", Actif: true},
		{ID: 2, Nom: "MacBook Pro 14", Marque: "Apple", Prix: 2499.00, Stock: 10, Categorie: "Ordinateur", Actif: true},
		{ID: 3, Nom: "Galaxy S24", Marque: "Samsung", Prix: 899.00, Stock: 18, Categorie: "Smartphone", Actif: true},
		{ID: 4, Nom: "ThinkPad X1", Marque: "Lenovo", Prix: 1799.00, Stock: 7, Categorie: "Ordinateur", Actif: true},
		{ID: 5, Nom: "AirPods Pro", Marque: "Apple", Prix: 279.00, Stock: 50, Categorie: "Audio", Actif: true},
	}
	for _, p := range produitsInitiaux {
		_ = cat.AjouterProduit(p)
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n===== TechShop =====")
		fmt.Println("[1] Ajouter  [2] Chercher  [3] Soldes  [4] Vendre  [5] Rapport  [0] Quitter")
		choix := lireLigne(reader, "Votre choix : ")

		switch choix {
		case "1": // Ajouter
			id, err := lireInt(reader, "ID : ")
			if err != nil {
				fmt.Println("Erreur : ID invalide")
				continue
			}
			nom := lireLigne(reader, "Nom : ")
			marque := lireLigne(reader, "Marque : ")
			prix, err := lireFloat(reader, "Prix : ")
			if err != nil {
				fmt.Println("Erreur : prix invalide")
				continue
			}
			stock, err := lireInt(reader, "Stock : ")
			if err != nil {
				fmt.Println("Erreur : stock invalide")
				continue
			}
			categorie := lireLigne(reader, "Catégorie : ")
			p := Produit{ID: id, Nom: nom, Marque: marque, Prix: prix, Stock: stock, Categorie: categorie, Actif: true}
			if err := cat.AjouterProduit(p); err != nil {
				fmt.Println("Erreur :", err)
			} else {
				fmt.Println("Produit ajouté :", p)
			}

		case "2": // Chercher
			fmt.Println("[a] par ID   [b] par catégorie")
			sous := lireLigne(reader, "Recherche : ")
			if sous == "a" {
				id, err := lireInt(reader, "ID recherché : ")
				if err != nil {
					fmt.Println("Erreur : ID invalide")
					continue
				}
				p, err := cat.TrouverParID(id)
				if err != nil {
					fmt.Println("Erreur :", err)
				} else {
					fmt.Println(p)
				}
			} else {
				c := lireLigne(reader, "Catégorie recherchée : ")
				resultats := cat.TrouverParCategorie(c)
				if len(resultats) == 0 {
					fmt.Println("Aucun produit dans cette catégorie.")
				}
				for _, p := range resultats {
					fmt.Println(p)
				}
			}

		case "3": // Soldes
			categorie := lireLigne(reader, "Catégorie à solder : ")
			pct, err := lireFloat(reader, "Pourcentage de réduction : ")
			if err != nil {
				fmt.Println("Erreur : pourcentage invalide")
				continue
			}
			n := cat.AppliquerReduction(categorie, pct)
			fmt.Printf("%d produit(s) modifié(s).\n", n)

		case "4": // Vendre
			id, err := lireInt(reader, "ID du produit : ")
			if err != nil {
				fmt.Println("Erreur : ID invalide")
				continue
			}
			qte, err := lireInt(reader, "Quantité : ")
			if err != nil {
				fmt.Println("Erreur : quantité invalide")
				continue
			}
			if err := cat.Vendre(id, qte); err != nil {
				fmt.Println("Erreur :", err)
			} else {
				fmt.Printf("Vente enregistrée : %d unité(s) du produit #%d\n", qte, id)
			}

		case "5": // Rapport
			fmt.Println(cat.Rapport())

		case "0": // Quitter
			fmt.Println("Au revoir !")
			return

		default:
			fmt.Println("Choix invalide.")
		}
	}
}
