package extras

import (
	"bufio"
	"fmt"
	"os"
)

func Input(prompt string) string {
	var strref string
	fmt.Printf("%s", prompt)
	if err := readLine(&strref); err != nil {
		fmt.Println(err.Error())
	}
	return strref
}

func SilentInput(prompt string) string {
	var strref string

	fmt.Printf("%s", prompt)
	fmt.Print("\033[8m")
	if err := readLine(&strref); err != nil {
		fmt.Println(err.Error())
	}
	fmt.Print("\033[28m")

	return strref
}

func readLine(strref *string) error {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	*strref = scanner.Text()
	return scanner.Err()
}
