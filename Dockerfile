FROM  golang:1.10-alpine3.8  AS builder
RUN go get github.com/muesli/cache2go 
RUN go get github.com/prometheus/client_golang/prometheus
RUN go get github.com/prometheus/client_golang/prometheus/promhttp
RUN go get github.com/sirupsen/logrus
RUN go get gopkg.in/yaml.v2

ADD . /go/src/apiexporter/
RUN cd src/apiexporter && go build -o apiexporter




From alpine:3.8
COPY --from=builder /go/src/apiexporter/apiexporter /usr/bin/apiexporter
CMD ["/usr/bin/apiexporter"]   

#ENTRYPOINT ["/apiexporter"]
