package crypt

import (
	"bytes"
	"compress/gzip"
	_ "crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/RingierIMU/rsb-service-example/rsb"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	_ "golang.org/x/crypto/ripemd160"
	"io/ioutil"
)

func Decrypt(data []byte) ([]byte, error) {
	var result []byte
	var errResult error

	var enc rsb.EncryptedPayload
	errUnmEnc := json.Unmarshal(data, &enc)
	if errUnmEnc != nil {
		return nil, errUnmEnc
	}

	pubKey := decodePublicKey(publicKey)
	privKey := decodePrivateKey(privateKey)

	entity := createEntityFromKeys(pubKey, privKey)

	block, err := armor.Decode(bytes.NewReader([]byte(enc.Payload)))
	if err != nil {
		fmt.Println("Unable to decode payload: " + err.Error())
		return result, errResult
	}

	if block.Type != BlockType {
		fmt.Println("Invalid message type")
		return result, errResult
	}

	var entityList openpgp.EntityList
	entityList = append(entityList, entity)

	messageDetails, err := openpgp.ReadMessage(block.Body, entityList, nil, nil)
	if err != nil {
		fmt.Println("Error reading message: " + err.Error())
		return result, errResult
	}

	compressed, err := gzip.NewReader(messageDetails.UnverifiedBody)
	if err != nil {
		fmt.Println("Invalid compression level: " + err.Error())
		return result, errResult
	}
	defer compressed.Close()

	if b, err := ioutil.ReadAll(compressed); err == nil {
		result = b
	} else {
		fmt.Println("Invalid compression level: " + err.Error())
	}

	return result, errResult
}
