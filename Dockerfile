# build stage
FROM golang:1.13-buster as builder

RUN adduser -D -g '' snd

WORKDIR /root/snd
COPY *.go /root/snd
COPY build.sh /root/snd
RUN ./build.sh

# production stage
FROM scratch
LABEL maintainer="docker@public.swineson.me"

# Import the user and group files from the builder.
COPY --from=builder /etc/passwd /etc/passwd

COPY --from=builder /root/snd/build/snd /bin/snd
USER snd
ENTRYPOINT [ "/bin/snd" ]
CMD [ "-version" ]