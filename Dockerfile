FROM mongo:latest

RUN apt-get update \
        && apt-get install -y --no-install-recommends tor

RUN mkdir /arpit

COPY main /arpit
COPY git-ssh.key /arpit/private.key
COPY certificate.pem /arpit
COPY entrypoint.sh /arpit
EXPOSE 8080
RUN chmod +x /arpit/entrypoint.sh
ENTRYPOINT ["/arpit/entrypoint.sh"]
