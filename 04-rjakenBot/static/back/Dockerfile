FROM golang:1.18.2 as build

ENV CGO_ENABLED=0

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . .

RUN go build -o /rjaken-bot-app ./cmd/rjakenBot/main.go

################################################################

FROM scratch

WORKDIR /

COPY --from=build /rjaken-bot-app /rjaken-bot-app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY static/ static/

USER 1000

EXPOSE 13337

CMD ["/rjaken-bot-app"]
