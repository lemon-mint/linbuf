package main

import (
	"fmt"

	"github.com/lemon-mint/linbuf"
)

func main() {
	ve := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	fmt.Printf("original: %#v\n", ve)

	e := linbuf.NewEncoder().
		Bytes(ve).
		Finalize()

	data := e.Bytes()
	fmt.Printf("encoded: %#v\n", data)

	var vd []byte
	err := linbuf.NewDecoder(data).
		Bytes(&vd).
		Finalize()
	fmt.Printf("err: %v\n", err)
	fmt.Printf("decoded: %#v\n", vd)
}
