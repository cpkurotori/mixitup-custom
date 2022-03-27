package main

import (
	"fmt"
	"mixitup-custom/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprint(os.Stdout, "ERROR")
	}
}
