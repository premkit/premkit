FROM golang:1.11

RUN go get github.com/kardianos/govendor
RUN go get github.com/go-swagger/go-swagger/cmd/swagger

ENV PROJECTPATH=/go/src/github.com/premkit/premkit
ENV PATH $PATH:$PROJECTPATH/go/bin

ENV LOG_LEVEL DEBUG

EXPOSE 80 443 2080 2443

WORKDIR $PROJECTPATH

# Set up required directories with permissions
RUN mkdir -p /data
RUN chmod -R a+rw /data

VOLUME /data

CMD ["/bin/bash"]
