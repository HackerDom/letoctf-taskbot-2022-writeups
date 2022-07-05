FROM node:14.19-alpine3.16 AS front_build_stage

WORKDIR /src

COPY ./webui/package.json package.json
COPY ./webui/yarn.lock yarn.lock
RUN yarn install

COPY ./webui/scripts scripts
COPY ./webui/config config
COPY ./webui/public public
COPY ./webui/src src
COPY ./webui/tsconfig.json tsconfig.json
RUN yarn run build

FROM golang:1.18.3-alpine3.16 AS build_stage

WORKDIR /src

COPY ./go.mod go.mod
COPY ./go.sum go.sum
RUN go mod download

RUN mkdir build

COPY ./pkg pkg
RUN go build -o /src/build/filestore /src/pkg

FROM alpine:3.16.0

WORKDIR /app

COPY --from=build_stage /src/build/filestore filestore
COPY --from=front_build_stage /src/build /webui

CMD ["./filestore"]
