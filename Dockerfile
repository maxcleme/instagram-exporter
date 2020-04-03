FROM golang:1.13 AS builder
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build main.go

FROM gcr.io/distroless/base
COPY --from=builder /app .
ENV HOME=/app
CMD ["./main"]