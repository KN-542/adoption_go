package service

import (
	"crypto/rand"
	"log"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

func GenerateHash(minLength, maxLength int) (*string, *string, error) {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	length, err := rand.Int(rand.Reader, big.NewInt(int64(maxLength-minLength+1)))
	if err != nil {
		return nil, nil, err
	}
	strLength := minLength + int(length.Int64())

	buffer := make([]byte, strLength)
	_, err2 := rand.Read(buffer)
	if err2 != nil {
		return nil, nil, err2
	}
	for i := 0; i < strLength; i++ {
		buffer[i] = chars[int(buffer[i])%len(chars)]
	}
	str := string(buffer)

	buffer2, err3 := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err3 != nil {
		log.Printf("%v", err3)
		return nil, nil, err3
	}
	hash := string(buffer2)

	return &str, &hash, nil
}
