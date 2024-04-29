package extras

import "fmt"

func Input(prompt string, strref ...any) {
	fmt.Printf("%s", prompt)
	fmt.Scanln(strref)
}

func SilentInput(prompt string, strref *string) {
	fmt.Printf("%s", prompt)
	fmt.Print("\033[8m")
	fmt.Scanln(strref)
	fmt.Print("\033[28m")
}
