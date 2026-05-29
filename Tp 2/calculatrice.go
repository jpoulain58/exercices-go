package main

import "fmt"

// creerOperation retourne une closure pour l'opération demandée.
// Renvoie nil si l'opérateur est inconnu.
func creerOperation(op string) func(float64, float64) float64 {
	switch op {
	case "+":
		return func(a, b float64) float64 { return a + b }
	case "-":
		return func(a, b float64) float64 { return a - b }
	case "*":
		return func(a, b float64) float64 { return a * b }
	case "/":
		return func(a, b float64) float64 { return a / b }
	default:
		return nil
	}
}

// operer applique l'opération avec validation : erreur si division par
// zéro ou opérateur inconnu.
func operer(a, b float64, op string) (float64, error) {
	if op == "/" && b == 0 {
		return 0, fmt.Errorf("division par zéro : impossible de diviser %v par 0", a)
	}
	calcul := creerOperation(op)
	if calcul == nil {
		return 0, fmt.Errorf("opération inconnue : %q", op)
	}
	return calcul(a, b), nil
}

func main() {
	// Boucle infinie : lit nombre1 nombre2 opérateur, 'quit' pour sortir.
	for {
		var a, b float64
		var op string
		fmt.Print("Entrez : nombre1 nombre2 opérateur (ex. 10 5 + — ou 0 0 quit) : ")
		fmt.Scan(&a, &b, &op)

		if op == "quit" {
			fmt.Println("Au revoir !")
			break
		}

		resultat, err := operer(a, b, op)
		if err != nil {
			fmt.Println("Erreur:", err)
			continue
		}
		fmt.Printf("%.2f %s %.2f = %.2f\n", a, op, b, resultat)
	}
}
