FROM golang:1.25-alpine AS builder
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata git
ENV GOPROXY=https://proxy.golang.org,direct GOSUMDB=sum.golang.org GIT_TERMINAL_PROMPT=0
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o prinxelioback ./cmd/server

FROM alpine:3.19
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /app/prinxelioback /app/prinxelioback
COPY public /app/public
EXPOSE 3020
ENV PORT=3020
CMD ["/app/prinxelioback"]
