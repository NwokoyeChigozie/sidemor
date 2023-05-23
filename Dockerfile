# Build stage
FROM golang:1.20.1-alpine3.17 as build

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN if test -e app.env; then echo 'found app.env'; else mv app-sample.env app.env; fi; \
    go build -v -o /dist/vesicash-mor-api

# Deployment stage
FROM alpine:3.17

WORKDIR /usr/src/app

COPY --from=build /usr/src/app ./

COPY --from=build /dist/vesicash-mor-api /usr/local/bin/vesicash-mor-api

CMD ["vesicash-mor-api"]