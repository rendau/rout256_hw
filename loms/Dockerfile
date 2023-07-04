FROM ubuntu:22.04

RUN apt-get update && apt-get install -y curl

ADD ./migrations /migrations
ADD ./bin/app /app

CMD ["/app"]
