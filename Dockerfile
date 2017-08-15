FROM golang:1.8.3-alpine3.6

# build arguments (get passed down from docker-compose)
ARG debug=0
ARG hotswap=0
ENV debugenv=${debug}
ENV hotswapenv=${hotswap}

#install git since alpine doesn't come with it
RUN apk update && apk add git

# install graphql packages
RUN go get github.com/graphql-go/graphql
RUN go get github.com/graphql-go/handler
RUN go get github.com/mnmtanish/go-graphiql
RUN go get github.com/rs/cors

# install neo4j package
RUN go get github.com/johnnadratowski/golang-neo4j-bolt-driver

# expose the remote debugging port
EXPOSE 2345

EXPOSE 8080

# copy source code from current directory to /go/src/graphql-server (this should be done for production)
# COPY . /go/src/graphql-server

# change working directoy to /go/src
WORKDIR /go/bin

# copy in init script to run the project taking into account the env vars
COPY ./docker-entrypoint.sh ./docker-entrypoint.sh

# copy in fresh library config
COPY ./runner.conf ./runner.conf

# set execute permissions for the script
RUN chmod 755 ./docker-entrypoint.sh

CMD ["${debugenv}", "${hotswapenv}"]

# run the init script
ENTRYPOINT ["/go/bin/docker-entrypoint.sh"]
