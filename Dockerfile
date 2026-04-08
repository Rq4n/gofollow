# build the binary
FROM golang:1.26.1-alpine AS builder

# set the working directory
WORKDIR /app

RUN apk add --no-cache git

# copy dependency files first
COPY go.mod go.sum ./
RUN go mod download

# copy the rest of the source code
COPY . .

# cria um link estatico (binario -> CGO_ENABLED = 0) para a plataforma alvo
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/.

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/ 

COPY --from=builder /app/main .

EXPOSE 3030

CMD ["./main"]
