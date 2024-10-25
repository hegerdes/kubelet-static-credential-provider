FROM ghcr.io/siderolabs/ecr-credential-provider:v1.31.1 AS TMP

FROM debian:bookworm

COPY --from=TMP / /root/

ENTRYPOINT ["sleep", "7d"]
