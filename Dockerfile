
FROM golang:1.6

RUN go get -u github.com/kardianos/govendor # 149dead46a0b5ec0d44f36ab95b7077a049eb68b

ENV PROJECTPATH=/go/src/github.com/premkit/premkit
ENV PATH $PATH:$PROJECTPATH/go/bin

ENV LOG_LEVEL DEBUG

EXPOSE 80 443

WORKDIR $PROJECTPATH

VOLUME /data

CMD ["/bin/bash"]