package crypt

import (
	"bytes"
	"compress/gzip"
	_ "crypto/sha256"
	"fmt"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	_ "golang.org/x/crypto/ripemd160"
	"io"
	"os"
)

func encryptToFile(filename string) {
	pubKey := decodePublicKey(publicKey)
	privKey := decodePrivateKey(privateKey)

	to := createEntityFromKeys(pubKey, privKey)

	ecryptedFile, err := os.Create(filename)
	if err != nil {
		fmt.Println("Unable to create file: " + err.Error())
	}
	defer ecryptedFile.Close()

	w, err := armor.Encode(ecryptedFile, "Message", make(map[string]string))
	if err != nil {
		fmt.Println("Error creating armor: " + err.Error())
		return
	}
	defer w.Close()

	plain, err := openpgp.Encrypt(w, []*openpgp.Entity{to}, nil, nil, nil)
	if err != nil {
		fmt.Println("Error creating entity for encryption: " + err.Error())
		return
	}
	defer plain.Close()

	compressed, err := gzip.NewWriterLevel(plain, gzip.BestCompression)
	if err != nil {
		fmt.Println("Invalid compression level: " + err.Error())
		return
	}

	n, err := io.Copy(compressed, bytes.NewReader([]byte("Hello World!\n")))
	if err != nil {
		fmt.Println("Error writing encrypted filed: " + err.Error())
		return
	}
	fmt.Printf("Wrote %d bytes", n)

	compressed.Close()
}
