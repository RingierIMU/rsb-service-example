package crypt

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"os"
	"path/filepath"
)

func generateKeys() {
	key, err := rsa.GenerateKey(rand.Reader, defaultKeyLength)
	if err != nil {
		fmt.Println("Unable to generate key: " + err.Error())
	}

	privKeyFile, err := os.Create(filepath.Join(keyOutputDir, keyOutputPrefix+".privkey"))
	if err != nil {
		fmt.Println("Unable to save private key: " + err.Error())
	}
	defer privKeyFile.Close()

	pubKeyfile, err := os.Create(filepath.Join(keyOutputDir, keyOutputPrefix+".pubkey"))
	if err != nil {
		fmt.Println("Unable to save public key: " + err.Error())
	}
	defer pubKeyfile.Close()

	encodePrivateKey(privKeyFile, key)
	encodePublicKey(pubKeyfile, key)
}
