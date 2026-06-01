FROM golang:1.26

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
