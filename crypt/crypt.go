package crypt

const (
	defaultKeyLength = 4096
	BlockType        = "Payload"
)

var (
	keyOutputPrefix = "example"
	keyOutputDir    = "."

	privateKey = keyOutputPrefix + ".privkey"
	publicKey  = keyOutputPrefix + ".pubkey"
)

func init() {
	if keysMissing() {
		generateKeys()
	}
}
