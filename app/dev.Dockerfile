FROM node:20 as client
WORKDIR /build
COPY client .
RUN npm ci && npm run build

FROM golang:1.20 as server
WORKDIR /build
COPY --from=client /build/dist static
COPY server .

# build the server with debug information
RUN go generate ./...
RUN CGO_ENABLED=0 GOOS=linux go build -gcflags "all=-N -l" -o server



FROM golang:1.20
WORKDIR /app
COPY --from=server /build .

RUN go install github.com/go-delve/delve/cmd/dlv
RUN go install github.com/cosmtrek/air

CMD air -c .air.toml