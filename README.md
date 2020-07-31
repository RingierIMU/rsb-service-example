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
