package main

import (
	"encoding/json"
	"fmt"
	"github.com/RingierIMU/rsb-service-example/crypt"
	_ "github.com/RingierIMU/rsb-service-example/crypt"
	"github.com/RingierIMU/rsb-service-example/rsb"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var (
	port string
	env  string
	ok   bool
)

func main() {
	getEnv()

	// Returns a 200 OK for the load balancer
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "OK")
	})

	// Receives the event
	http.HandleFunc("/api/events", eventsHandler)

	// Opens a webserver on $port
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		panic(fmt.Sprintf("Unable to listen on port %s: %v", port, err))
	}
}

// Checks / sets environment variables, fallback to sane defaults
func getEnv() {
	if port, ok = os.LookupEnv("PORT"); !ok {
		port = "8080"
	}

	if env, ok = os.LookupEnv("RSB_ENV"); !ok {
		env = "NA"
	}
}

// Receives an (encrypted) event, sends an encrypted callback, returns an informational response
func eventsHandler(w http.ResponseWriter, req *http.Request) {
	// Read event from HTTP body, convert into Event struct
	body, errReadBody := ioutil.ReadAll(req.Body)
	if errReadBody != nil {
		fmt.Println("Unable to read request body: " + errReadBody.Error())
	}

	var event rsb.Event
	errUnmarshal := json.Unmarshal(body, &event)
	if errUnmarshal != nil {
		fmt.Println("Unable to unmarshal data: " + errUnmarshal.Error())
	}

	// Decrypt is payload is encrypted, otherwise show plaintext payload
	if strings.Contains(string(event.Payload), fmt.Sprintf("-----BEGIN %s-----", crypt.BlockType)) {
		payload, errDecrypt := crypt.Decrypt(event.Payload)
		if errDecrypt != nil {
			fmt.Println("Error decrypting payload: " + errDecrypt.Error())
		} else {
			fmt.Println("Received this payload (encrypted):\n" + string(payload))
		}
	} else {
		fmt.Println("Received this payload (plaintext):\n" + string(event.Payload))
	}

	// Encrypt payload of callback, wrap into event
	message := []byte(`{"code":400,"message":"Could not handle this payload"}`)
	encryptedPayload, errEncrypt := crypt.Encrypt(message)
	if errEncrypt != nil {
		fmt.Println("Error encrypting payload: " + errEncrypt.Error())
	} else {
		fmt.Printf("Encrypted %s to:\n%s\n", message, string(encryptedPayload))
	}
	pl, errMarshalPl := json.Marshal(rsb.EncryptedPayload{
		Payload: string(encryptedPayload),
	})
	if errMarshalPl != nil {
		fmt.Println("Error marshalling payload: " + errMarshalPl.Error())
	} else {
		callbackEvent, _ := json.Marshal(
			rsb.Event{
				Event:   "CallbackEvent",
				Payload: pl,
			},
		)
		fmt.Println("Here's the callback event:\n" + string(callbackEvent))
	}

	// Send response to the caller of the service (usually the bus)
	_, err := fmt.Fprintf(w, fmt.Sprintf("Received an event from %s on environment %s", req.RemoteAddr, env))
	if err != nil {
		fmt.Println("Unable to return data: " + err.Error())
	}
}
