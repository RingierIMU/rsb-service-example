FROM golang:1.14-alpine

RUN apk add --update git

ARG GITHUBOAUTHTOKEN
ENV GITHUBOAUTHTOKEN=$GITHUBOAUTHTOKEN
RUN printf 'machine github.com\nlogin git\npassword %s' $GITHUBOAUTHTOKEN > ~/.netrc

ENV GOPROXY "https://proxy.golang.org,direct"
ENV CGO_ENABLED 0
ENV GOPRIVATE "github.com/RingierIMU/rsb-go-lib"

WORKDIR $HOME/github.com/RingierIMU/rsb-service-example
COPY . .

RUN go build -a -ldflags '-extldflags "-static"' -v -o /rsb-service-example .

FROM gcr.io/distroless/static
COPY --from=0 /rsb-service-example /
CMD ["/rsb-service-example"]
