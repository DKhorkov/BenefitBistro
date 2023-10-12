package mycrypto

import (
	"testing"

	"github.com/stretchr/testify/require"

	"mycrypto/testdata"
)

func TestHashPassword(test *testing.T) {
	hashed_password, err := HashPassword(testdata.Password)
	if err != nil {
		test.Errorf("Failed to hash password=%v. Error: %v\n", testdata.Password, err)
	}

	hashed_password_length := len(hashed_password)
	if hashed_password_length != testdata.PasswordExpectedLength {
		test.Errorf(
			"Actual hashed password length (%v) is not euqal to expected length (%v)!\n", 
			hashed_password_length, 
			testdata.PasswordExpectedLength)
	}
}

func TestEncrypt(test *testing.T) {
	encrypted, err := Encrypt(testdata.EncryptAndDecryptData)
	if err != nil {
		test.Errorf("Failed to encrypt data=%v. Error: %v\n", testdata.EncryptAndDecryptData, err)
	}

	const isBase64 = "^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{4})$"
	require.Nil(test, err)
	require.Regexp(test, isBase64, encrypted)
}

func TestDecrypt(test *testing.T) {
	encrypted, err := Encrypt(testdata.EncryptAndDecryptData)
	if err != nil {
		test.Errorf("Failed to encrypt data=%v. Error: %v\n", testdata.EncryptAndDecryptData, err)
	}

	var decrypted string
	decrypted, err = Decrypt(encrypted)
	if err != nil {
		test.Errorf("Failed to decrypt data=%v. Error: %v\n", encrypted, err)
	}

	if decrypted != testdata.EncryptAndDecryptData {
		test.Errorf(
			"Decrypted data (%v) is not equal to base data before encryption (%v)\n", 
			decrypted, 
			testdata.EncryptAndDecryptData)
	}
}
