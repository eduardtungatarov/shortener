package pkg1

func someFuncA() {
	panic("error") // want "usage of panic"
}

func main() {
	panic("error") // want "usage of panic"
}
