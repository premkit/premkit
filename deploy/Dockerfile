FROM alpine:3.4

RUN apk add --no-cache \
  ca-certificates

ADD ./deploy/bin/premkit /usr/bin/premkit

VOLUME ["/data"]
EXPOSE 80 443

ENTRYPOINT ["/usr/bin/premkit"]
CMD ["daemon"]
