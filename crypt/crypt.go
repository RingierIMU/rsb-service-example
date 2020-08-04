package crypt

const (
	defaultKeyLength = 4096
)

var (
	keyOutputPrefix = "example"
	keyOutputDir    = "."

	privateKey = keyOutputPrefix + ".privkey"
	publicKey  = keyOutputPrefix + ".pubkey"
)

func init() {
	//generateKeys()

	encryptToFile("/Users/zebroc/go/src/github.com/RingierIMU/rsb-service-example/file.txt")

	decryptFromFile("/Users/zebroc/go/src/github.com/RingierIMU/rsb-service-example/file.txt")
}
