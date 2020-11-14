FROM golang:1.15.2 AS build

ENV GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

WORKDIR /go/src/inventoryservice
COPY . .
RUN go mod download
RUN go build -v 

FROM alpine:3.12.0
COPY --from=build /go/src/inventoryservice/inventoryservice /usr/local/lib/inventoryservice/inventoryservice
COPY --from=build /go/src/inventoryservice/products.json /usr/local/lib/inventoryservice/products.json
WORKDIR /usr/local/lib/inventoryservice
CMD ["/usr/local/lib/inventoryservice/inventoryservice"]