FROM golang:1.17.6-alpine3.15 as builder

EXPOSE 2701/tcp

WORKDIR /workspace
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd/test_tcp/main.go main.go
COPY internal/ internal/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o test_tcp main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/test_tcp .

ENTRYPOINT ["/test_tcp"]


