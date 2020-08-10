package rsb

import "encoding/json"

type Event struct {
	Event   string          `json:"event"`
	Payload json.RawMessage `json:"payload"`
}

type EncryptedPayload struct {
	Payload string `json:"encrypted_payload"`
}
