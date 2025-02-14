FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /bin/hyper-todo cmd/api/main.go

FROM gcr.io/distroless/base-debian12

WORKDIR /

COPY --from=builder /bin/hyper-todo /bin/hyper-todo

ENV GIN_MODE=release

CMD ["/bin/hyper-todo"]
