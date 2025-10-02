FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

# Instalar swag para gerar documentação Swagger
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .

# Gerar documentação Swagger
RUN swag init -g cmd/server/main.go -o docs

# Compilar a aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server cmd/server/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

ENV TZ=America/Sao_Paulo

WORKDIR /root/

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]
