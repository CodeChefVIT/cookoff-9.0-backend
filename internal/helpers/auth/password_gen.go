package auth

import (
	"math/rand"
	"time"
)

func PasswordGenerator(passwordLength int) string {
	upperCase := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers := "0123456789"
	specialChar := "!@#$%^&*()_-+={}[/?]"

	password := ""

	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	for n := 0; n < passwordLength; n++ {
		randNum := rng.Intn(4)

		switch randNum {
		case 0:
			randCharNum := rng.Intn(len(upperCase))
			password += string(upperCase[randCharNum])
		case 1:
			randCharNum := rng.Intn(len(upperCase))
			password += string(upperCase[randCharNum])
		case 2:
			randCharNum := rng.Intn(len(numbers))
			password += string(numbers[randCharNum])
		case 3:
			randCharNum := rng.Intn(len(specialChar))
			password += string(specialChar[randCharNum])
		}
	}

	return password
}
