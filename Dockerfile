ARG BUILD_VERSION=1.24-bookworm
FROM golang:$BUILD_VERSION AS build

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/static-credential-provider

# Test
ENV KSCP_REGISTRY_USERNAME=my-user
ENV KSCP_REGISTRY_PASSWORD=my-password
RUN echo '{"apiVersion":"credentialprovider.kubelet.k8s.io/v1","kind":"CredentialProviderRequest","image":"your.registry.example.org/org/image:version"}' | /go/bin/static-credential-provider

FROM gcr.io/distroless/static-debian12

COPY --from=build /go/bin/* /
ENTRYPOINT ["/static-credential-provider"]
