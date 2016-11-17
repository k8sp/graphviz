FROM golang:wheezy

ENV DIST /go/src/github.com/k8sp/graphviz
RUN sed -i "s#wheezy#jessie#g" /etc/apt/sources.list
RUN apt-get update
RUN apt-get install -y graphviz && dot -v

COPY . $DIST
RUN cd $DIST && go get ./... && go get .

EXPOSE 9090
VOLUME ["/cache"]
ENTRYPOINT ["graphviz"]
CMD ["-addr=:9090", "-dir=/cache"]
