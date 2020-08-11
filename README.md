# rsb-service-example

A service to show conventions and best practices

## Functionality

This services serves no real-world purpose, it is just an example for coding. It returns a 200 OK on / for the
load balancer and an informational message when an event is sent to /api/events.

It comes with all the necessary files for deployment, Dockerfile, buildspec.yml and BaseTaskDefinition.json.

## Build / compile

You can use the go toolchain to compile and run the service:

go build / go run

You can build and run it from Docker:

docker build -t  rsb-service-example
ducker run

## Documentation

Besides this Readme, you can find a postman collection in /api/postman

There is no further documentation, but if there were, this would be the place to link to it

## Environment

Variables that can / should be set:

* PORT - The HTTP port to listen on
* RSB_ENV - The name of the environment that the service is deployed in

## Encryption

Please see [End to end encryption on Ringier Event Bus](https://docs.google.com/document/d/1YfHAeQeU_N1zjeX5TKFSatqHgl7_EvUymaZH32gKudI) 
for an overview of how the encryption works in principle.

In this service, a keypair is automatically generated if none is present. Incoming events are decrypted, only if they
are encrypted. If not, the plaintext payload is used.

The service prints the (decrypted) payload to Stdout, along with an encrypted payload wrapped in a callback event.

In the real world, the service would need to have the public key of the receiver
(which it can get from the servicerepository e.g.).

To test the decryption, you can copy & paste the callback event into Postman and send it to the service.
