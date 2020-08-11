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
)

func Encrypt(data []byte) ([]byte, error) {
	pubKey := decodePublicKey(publicKey)
	privKey := decodePrivateKey(privateKey)

	to := createEntityFromKeys(pubKey, privKey)

	ecryptedBuffer := new(bytes.Buffer)

	encryptedBufferWriter, err := armor.Encode(ecryptedBuffer, BlockType, make(map[string]string))
	if err != nil {
		fmt.Println("Error creating armor: " + err.Error())
		return nil, err
	}

	plain, err := openpgp.Encrypt(encryptedBufferWriter, []*openpgp.Entity{to}, nil, nil, nil)
	if err != nil {
		fmt.Println("Error creating entity for encryption: " + err.Error())
		return nil, err
	}

	compressed, err := gzip.NewWriterLevel(plain, gzip.BestCompression)
	if err != nil {
		fmt.Println("Invalid compression level: " + err.Error())
		return nil, err
	}

	_, err = io.Copy(compressed, bytes.NewReader(data))
	if err != nil {
		fmt.Println("Error writing encrypted filed: " + err.Error())
		return nil, err
	}

	compressed.Close()
	plain.Close()
	encryptedBufferWriter.Close()

	return ecryptedBuffer.Bytes(), nil
}
