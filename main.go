package main

import (
	"fmt"

	app "github.com/alpha2303/pass-man/internal/app"
)

func main() {
	var choice int = 1

	for choice != 0 {
		fmt.Println("\n\n** Welcome to Pass-Man! **")
		fmt.Println("\nChoose your option:")
		fmt.Println("1. Create a Vault")
		fmt.Println("2. Sign In to a Vault")
		fmt.Println("3. Delete Vault")
		fmt.Println("0. Exit")

		fmt.Print("\nChoose your option: ")

		_, err := fmt.Scanln(&choice)

		if err != nil {
			fmt.Printf("error occured: %s\n", err.Error())
			continue
		}

		switch choice {
		case 1:
			app.HandleVaultCreate()
		case 2:
			app.HandleVaultSignIn()
		case 3:
			app.HandleVaultDelete()
		default:
			choice = 0
		}
	}
	fmt.Println("\n** Thanks for using Pass-Man! **")
}
