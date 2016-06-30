FROM debian:jessie

ADD bin/premkit /premkit

RUN apt-get -y update && \
    apt-get -y --no-install-recommends install ca-certificates && \
    rm -rf /var/lib/apt/lists/*

EXPOSE 80 443

ENTRYPOINT ["/premkit"]
CMD ["daemon"]
