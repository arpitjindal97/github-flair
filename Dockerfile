FROM debian:jessie-slim

RUN apt-get update && apt-get install -y ca-certificates --no-install-recommends
RUN mkdir /arpit

COPY output/flair* /arpit/flair-bin
COPY secrets /arpit
COPY entrypoint.sh /arpit
COPY arialbd.ttf /arpit

RUN chmod +x /arpit/entrypoint.sh
ENTRYPOINT ["/arpit/entrypoint.sh"]
