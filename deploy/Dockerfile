FROM alpine:3.4

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
