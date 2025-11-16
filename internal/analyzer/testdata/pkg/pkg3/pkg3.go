package pkg3

import (
	"os"
)

func someFuncA() {
	os.Exit(1) // want "usage of os.Exit()"
}

func main() {
	os.Exit(1) // want "usage of os.Exit()"
}
