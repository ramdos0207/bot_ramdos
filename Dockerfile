FROM golang:1.23-bookworm AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/bot-ramdos .

FROM gcr.io/distroless/base-debian12:nonroot

WORKDIR /app
COPY --from=builder /out/bot-ramdos /usr/local/bin/bot-ramdos
COPY help.txt ./help.txt

USER nonroot:nonroot
ENTRYPOINT ["/usr/local/bin/bot-ramdos"]
