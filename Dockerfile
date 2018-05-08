FROM debian:jessie-slim

RUN apt-get update && apt-get install -y ca-certificates --no-install-recommends
RUN mkdir /arpit

COPY output/flair* /arpit/flair-bin

ENTRYPOINT ["/arpit/flair-bin"]
