FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN apk add --no-cache build-base && go build -o forum cmd/web/main.go

FROM alpine:3.14
WORKDIR /app
COPY --from=builder /app .

EXPOSE 8080
LABEL name="FORUM" \
      authors="asundeto, atemerzh" \
      release_date="01.02.2024"
CMD ["./forum"]