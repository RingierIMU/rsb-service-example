package crypt

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
	"io"
	"os"
	"path/filepath"
	"time"
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

func encodePrivateKey(out io.Writer, key *rsa.PrivateKey) {
	w, err := armor.Encode(out, openpgp.PrivateKeyType, make(map[string]string))
	if err != nil {
		fmt.Println("Unable to save private key: " + err.Error())
	}
	defer w.Close()

	pgpKey := packet.NewRSAPrivateKey(time.Now(), key)
	err = pgpKey.Serialize(w)
	if err != nil {
		fmt.Println("Unable to serialize private key: " + err.Error())
	}
}

func encodePublicKey(out io.Writer, key *rsa.PrivateKey) {
	w, err := armor.Encode(out, openpgp.PublicKeyType, make(map[string]string))
	if err != nil {
		fmt.Println("Unable to save public key: " + err.Error())
	}
	defer w.Close()

	pgpKey := packet.NewRSAPublicKey(time.Now(), &key.PublicKey)
	err = pgpKey.Serialize(w)
	if err != nil {
		fmt.Println("Unable to serialize public key: " + err.Error())
	}
}
