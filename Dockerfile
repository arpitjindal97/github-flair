FROM mongo:latest

RUN apt-get update \
        && apt-get install -y --no-install-recommends tor varnish

RUN mkdir /arpit

COPY main /arpit
COPY https_server /arpit
COPY secrets/git-ssh.key /arpit/private.key
COPY certificate.pem /arpit
COPY entrypoint.sh /arpit
COPY default.vcl /arpit
EXPOSE 443
RUN chmod +x /arpit/entrypoint.sh
ENTRYPOINT ["/arpit/entrypoint.sh"]
