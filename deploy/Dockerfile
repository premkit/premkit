FROM debian:bookworm-slim

RUN apt-get update && apt-get upgrade -y && apt-get install -y --no-install-recommends \
    ca-certificates \
    \
    \
  && rm -rf /var/lib/apt/lists/*

ADD ./bin/premkit /usr/bin/premkit

EXPOSE 80 443 2080 2443

# Set up required directories with permission
RUN mkdir -p /data
RUN chmod -R a+rw /data

VOLUME /data

ENTRYPOINT ["/usr/bin/premkit"]

CMD ["daemon"]
