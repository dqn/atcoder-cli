package atcoder

import "fmt"

func colorize(s string, c string) string {
	var cd string
	switch c {
	case "green":
		cd = "\x1b[32m"
	case "yellow":
		cd = "\x1b[33m"
	default:
		panic(fmt.Sprintf("unknown color: %s", c))
	}
	return cd + s + "\x1b[0m"
}

func printDivider() {
	fmt.Println("-----------------------")
}

func printHeader(i int, result string) {
	printDivider()
	fmt.Printf("case %d: %s\n", i+1, result)
}

func printAC(i int) {
	printHeader(i, colorize("AC", "green"))
}

func printWA(i int, expect string, actual string) {
	printHeader(i, colorize("WA", "yellow"))
	fmt.Printf("%s\n", colorize(expect, "green"))
	fmt.Printf("%s\n", colorize(actual, "yellow"))
}
