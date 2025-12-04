package cli

import "fmt"

// color codes
const (
	green    = "\033[32m"
	yellow   = "\033[33m"
	red      = "\033[31m"
	cyanBold = "\033[1;36m"
	reset    = "\033[0m"
)

func printSection(title string) {
	fmt.Printf("\n%s%s%s\n", cyanBold, title, reset)
}

func printStep(title string) {
	fmt.Printf("%s==>%s %s\n", cyanBold, reset, title)
}

func printInfo(msg string) {
	fmt.Printf("  %s•%s %s\n", green, reset, msg)
}

func printWarn(msg string) {
	fmt.Printf("  %s!%s %s\n", yellow, reset, msg)
}

func printError(msg string) {
	fmt.Printf("  %s✖%s %s\n", red, reset, msg)
}

func printSuccess(msg string) {
	fmt.Printf("%s✔%s %s\n", green, reset, msg)
}

func colorGreen(msg string) string {
	return green + msg + reset
}
