package repository

import (
	"crypto/rand"
	"log"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

func GenerateHash(minLength, maxLength int) (*string, *string, error) {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	length, lengthErr := rand.Int(rand.Reader, big.NewInt(int64(maxLength-minLength+1)))
	if lengthErr != nil {
		return nil, nil, lengthErr
	}
	strLength := minLength + int(length.Int64())

	buffer := make([]byte, strLength)
	_, bufferErr := rand.Read(buffer)
	if bufferErr != nil {
		return nil, nil, bufferErr
	}
	for i := 0; i < strLength; i++ {
		buffer[i] = chars[int(buffer[i])%len(chars)]
	}
	str := string(buffer)

	buffer2, buffer2Err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if buffer2Err != nil {
		log.Printf("%v", buffer2Err)
		return nil, nil, buffer2Err
	}
	hash := string(buffer2)

	return &str, &hash, nil
}
