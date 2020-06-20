FROM alpine:latest

WORKDIR /
ADD graphql-api-demo graphql-api-demo
ADD config.json      config.json

EXPOSE 8080

ENTRYPOINT ["/graphql-api-demo"]
