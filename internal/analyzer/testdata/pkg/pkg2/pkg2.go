package pkg2

import "log"

func someFuncA() {
	log.Fatal("error") // want "usage of log.Fatal()"
}

func main() {
	log.Fatal("error") // want "usage of log.Fatal()"
}
