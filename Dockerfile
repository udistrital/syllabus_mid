# Build stage for installing libmcrypt
FROM debian:12 as builder

RUN apt-get update && \
    apt-get install -y libmcrypt4 libmcrypt-dev && \
    rm -rf /var/lib/apt/lists/*

# Final stage with distroless
FROM gcr.io/distroless/base-debian12

COPY --from=builder /usr/lib/x86_64-linux-gnu/libmcrypt.so* /usr/lib/x86_64-linux-gnu/

COPY main main
COPY conf/app.conf conf/app.conf

ENTRYPOINT ["/main"]
