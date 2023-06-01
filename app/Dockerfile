FROM node:20 as client
WORKDIR /build
COPY client .
RUN npm ci && npm run build

FROM golang:1.20 as server
WORKDIR /build
COPY --from=client /build/dist static
COPY server .

# build
RUN go generate ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o server

FROM alpine:3.18
WORKDIR /app
COPY --from=server /build/server ./server
COPY --from=server /static ./server/static
CMD ["./server"]