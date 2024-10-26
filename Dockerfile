ARG BUILD_VERSION=1.23-bookworm
FROM golang:$BUILD_VERSION AS build

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN go vet -v
RUN go test -v
RUN CGO_ENABLED=0 go build -o /go/bin/static-credential-provider

FROM gcr.io/distroless/static-debian12

COPY --from=build /go/bin/* /
CMD ["/static-credential-provider"]
