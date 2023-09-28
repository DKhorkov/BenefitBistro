package mycrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"config"
	"logging"
)

/*
	This code is working according to https://gist.github.com/fracasula/38aa1a4e7481f9cedfa78a0cdd5f1865 with some changes.

	For hashing passwords it's better to use in future SHA256 only to hash password and not decrypting it back
	https://stackoverflow.com/questions/10701874/generating-the-sha-hash-of-a-string-using-golang
*/

func EncryptMessage(message string) (string, error) {
	byteMsg := []byte(message)
	block, err := aes.NewCipher(config.CryptKey)
	if err != nil {
		err := fmt.Errorf("could not create new cipher: %v", err)
		logging.Log.Println(err)
		return "", err
	}

	cipherText := make([]byte, aes.BlockSize+len(byteMsg))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		err := fmt.Errorf("could not encrypt: %v", err)
		logging.Log.Println(err)
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], byteMsg)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func DecryptMessage(message string) (string, error) {
	cipherText, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		err := fmt.Errorf("could not base64 decode: %v", err)
		logging.Log.Println(err)
		return "", err
	}

	block, err := aes.NewCipher(config.CryptKey)
	if err != nil {
		err := fmt.Errorf("could not create new cipher: %v", err)
		logging.Log.Println(err)
		return "", err
	}

	if len(cipherText) < aes.BlockSize {
		err := fmt.Errorf("invalid ciphertext block size")
		logging.Log.Println(err)
		return "", err
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}