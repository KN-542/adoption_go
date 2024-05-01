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
	_, err = rand.Read(buffer)
	if err != nil {
		return nil, nil, err
	}
	for i := 0; i < strLength; i++ {
		buffer[i] = chars[int(buffer[i])%len(chars)]
	}
	str := string(buffer)

	buffer2, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("%v", err)
		return nil, nil, err
	}
	hash := string(buffer2)

	return &str, &hash, nil
}
