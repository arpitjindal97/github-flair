FROM mongo:latest

RUN mkdir /arpit

COPY main /arpit
COPY crt-bundle.pem /arpit
COPY ssl-private.key /arpit
COPY entrypoint.sh /arpit
COPY arialbd.ttf /arpit

#COPY cron_file /arpit/tor_restart.sh
#RUN chmod +x /arpit/tor_restart.sh
#RUN echo "* * * * * root /arpit/tor_restart.sh > /dev/fd/1" > /etc/cron.d/tor_cron
#RUN crontab /etc/cron.d/tor_cron

EXPOSE 443
RUN chmod +x /arpit/entrypoint.sh
ENTRYPOINT ["/arpit/entrypoint.sh"]
