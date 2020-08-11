package crypt

import (
	"crypto"
	"crypto/rsa"
	"fmt"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
	"io"
	"os"
	"time"
)

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

func keysMissing() bool {
	if _, err := os.Stat(privateKey); os.IsNotExist(err) {
		return true
	}

	if _, err := os.Stat(privateKey); os.IsNotExist(err) {
		return true
	}

	return false
}
