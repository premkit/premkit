FROM debian:jessie

RUN apt-get -y update && \
    apt-get -y --no-install-recommends install ca-certificates && \
    rm -rf /var/lib/apt/lists/*

EXPOSE 80 443

VOLUME ["/data"]

ADD bin/premkit /premkit

ENTRYPOINT ["/premkit"]
CMD ["daemon"]
