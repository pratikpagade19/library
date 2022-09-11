FROM golang:1.16.15

WORKDIR /home

COPY ./pkg /home

RUN cd /home && go build -o library

CMD ["/home/library"]