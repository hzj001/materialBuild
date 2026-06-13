package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	pwd := "admin123"
	if len(os.Args) > 1 {
		pwd = os.Args[1]
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	fmt.Println(string(hash))
}
