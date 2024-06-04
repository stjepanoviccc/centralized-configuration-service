# build
FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# run
RUN go build -o app .
FROM alpine:latest
WORKDIR /app
COPY --from=build /app/app .
COPY swagger.yaml /app/swagger.yaml
EXPOSE 8000
CMD ["./app"]