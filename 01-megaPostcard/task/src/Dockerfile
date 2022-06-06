FROM golang:1.18.2 as build

ENV CGO_ENABLED=0

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . .

RUN go build -o /megapostcard-app ./cmd/megapostcard/main.go

################################################################

FROM scratch

WORKDIR /

COPY --from=build /megapostcard-app /megapostcard-app

COPY static/ static/

USER 1000

EXPOSE 13337

CMD ["/megapostcard-app"]
