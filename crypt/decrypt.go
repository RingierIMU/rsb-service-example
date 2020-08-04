package crypt

import (
	"bufio"
	"compress/gzip"
	"crypto"
	_ "crypto/sha256"
	"fmt"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
	_ "golang.org/x/crypto/ripemd160"
	"io"
	"os"
)

func decryptFromFile(filename string) {
	pubKey := decodePublicKey(publicKey)
	privKey := decodePrivateKey(privateKey)

	entity := createEntityFromKeys(pubKey, privKey)

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Unable to open file: " + err.Error())
		return
	}

	block, err := armor.Decode(bufio.NewReader(file))
	if err != nil {
		fmt.Println("Unable to decode file: " + err.Error())
		return
	}

	if block.Type != "Message" {
		fmt.Println("Invalid message type")
		return
	}

	var entityList openpgp.EntityList
	entityList = append(entityList, entity)

	md, err := openpgp.ReadMessage(block.Body, entityList, nil, nil)
	if err != nil {
		fmt.Println("Error reading message: " + err.Error())
		return
	}

	compressed, err := gzip.NewReader(md.UnverifiedBody)
	if err != nil {
		fmt.Println("Invalid compression level: " + err.Error())
		return
	}
	defer compressed.Close()

	n, err := io.Copy(os.Stdout, compressed)
	if err != nil {
		fmt.Println("Error reading encrypted file:" + err.Error())
		return
	}
	fmt.Printf("Decrypted %d bytes", n)
}

func createEntityFromKeys(pubKey *packet.PublicKey, privKey *packet.PrivateKey) *openpgp.Entity {
	config := packet.Config{
		DefaultHash:            crypto.SHA256,
		DefaultCipher:          packet.CipherAES256,
		DefaultCompressionAlgo: packet.CompressionZLIB,
		CompressionConfig: &packet.CompressionConfig{
			Level: 9,
		},
		RSABits: defaultKeyLength,
	}
	currentTime := config.Now()
	uid := packet.NewUserId("", "", "")

	e := openpgp.Entity{
		PrimaryKey: pubKey,
		PrivateKey: privKey,
		Identities: make(map[string]*openpgp.Identity),
	}
	isPrimaryId := false

	e.Identities[uid.Id] = &openpgp.Identity{
		Name:   uid.Name,
		UserId: uid,
		SelfSignature: &packet.Signature{
			CreationTime: currentTime,
			SigType:      packet.SigTypePositiveCert,
			PubKeyAlgo:   packet.PubKeyAlgoRSA,
			Hash:         config.Hash(),
			IsPrimaryId:  &isPrimaryId,
			FlagsValid:   true,
			FlagSign:     true,
			FlagCertify:  true,
			IssuerKeyId:  &e.PrimaryKey.KeyId,
		},
	}

	keyLifetimeSecs := uint32(86400 * 365)

	e.Subkeys = make([]openpgp.Subkey, 1)
	e.Subkeys[0] = openpgp.Subkey{
		PublicKey:  pubKey,
		PrivateKey: privKey,
		Sig: &packet.Signature{
			CreationTime:              currentTime,
			SigType:                   packet.SigTypeSubkeyBinding,
			PubKeyAlgo:                packet.PubKeyAlgoRSA,
			Hash:                      config.Hash(),
			PreferredHash:             []uint8{8}, // SHA-256
			FlagsValid:                true,
			FlagEncryptStorage:        true,
			FlagEncryptCommunications: true,
			IssuerKeyId:               &e.PrimaryKey.KeyId,
			KeyLifetimeSecs:           &keyLifetimeSecs,
		},
	}
	return &e
}

func decodePublicKey(filename string) *packet.PublicKey {
	in, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening public key: " + err.Error())
	}
	defer in.Close()

	block, err := armor.Decode(in)
	if err != nil {
		fmt.Println("Error decoding OpenPGP Armor: " + err.Error())
	}

	if block.Type != openpgp.PublicKeyType {
		fmt.Println("Invalid public key file")
	}

	reader := packet.NewReader(block.Body)
	pkt, err := reader.Next()
	if err != nil {
		fmt.Println("Error reading private key: " + err.Error())
	}

	key, ok := pkt.(*packet.PublicKey)
	if !ok {
		fmt.Println("Invalid public key")
	}
	return key
}

func decodePrivateKey(filename string) *packet.PrivateKey {
	in, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening private key: " + err.Error())
	}
	defer in.Close()

	block, err := armor.Decode(in)
	if err != nil {
		fmt.Println("Error decoding OpenPGP Armor: " + err.Error())
	}

	if block.Type != openpgp.PrivateKeyType {
		fmt.Println("Invalid private key file")
	}

	reader := packet.NewReader(block.Body)
	pkt, err := reader.Next()
	if err != nil {
		fmt.Println("Error reading private key: " + err.Error())
	}

	key, ok := pkt.(*packet.PrivateKey)
	if !ok {
		fmt.Println("Invalid private key")
	}

	return key
}
