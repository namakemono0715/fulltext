FROM golang:1.24.5
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

# air をインストール
RUN go install github.com/air-verse/air@latest

COPY . .

EXPOSE 8080
CMD ["air"]

# RUN go build -o server

# EXPOSE 8080

# CMD ["./server"]