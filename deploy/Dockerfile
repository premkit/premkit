FROM alpine@sha256:f006ecbb824d87947d0b51ab8488634bf69fe4094959d935c0c103f4820a417d

RUN apk add --no-cache \
  ca-certificates

ADD ./deploy/bin/premkit /usr/bin/premkit

EXPOSE 80 443 2080 2443

# Set up required directories with permission
RUN mkdir -p /data
RUN chmod -R a+rw /data

VOLUME /data

ENTRYPOINT ["/usr/bin/premkit"]

CMD ["daemon"]
