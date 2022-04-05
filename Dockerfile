# build stage
FROM golang:1.18-buster as builder

ARG GOPATH=/tmp/go
RUN apt-get update -y \
    && apt-get install -y upx libcap2-bin \
    && go install github.com/goreleaser/goreleaser@latest

WORKDIR /root/snd
COPY . /root/snd/
RUN goreleaser build --single-target --id "snd" --output "dist/snd" --snapshot --rm-dist

# production stage
FROM debian:buster-slim
LABEL maintainer="docker@public.swineson.me"

# Import the user and group files from the builder.
COPY --from=builder /etc/passwd /etc/group /etc/

COPY --from=builder /root/snd/dist/snd /usr/local/bin/
COPY --from=builder /root/snd/examples/config.toml /etc/snd/

# nope
# See: https://github.com/moby/moby/issues/8460
# USER nobody:nogroup

EXPOSE 53/tcp 53/udp
ENTRYPOINT [ "/usr/local/bin/snd" ]
CMD [ "-config",  "/etc/snd/config.toml" ]
