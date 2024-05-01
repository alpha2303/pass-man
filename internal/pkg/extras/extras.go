package extras

import "fmt"

// func Input(prompt string, strref ...any) {
// 	fmt.Printf("%s", prompt)
// 	if _, err := fmt.Scanln(strref); err != nil {
// 		fmt.Println(err.Error())
// 	}

// }

func Input(prompt string) string {
	var strref string
	fmt.Printf("%s", prompt)
	if _, err := fmt.Scanln(&strref); err != nil {
		fmt.Println(err.Error())
	}
	return strref
}

func SilentInput(prompt string) string {
	var strref string

	fmt.Printf("%s", prompt)
	fmt.Print("\033[8m")
	if _, err := fmt.Scanln(&strref); err != nil {
		fmt.Println(err.Error())
	}
	fmt.Print("\033[28m")

	return strref
}
