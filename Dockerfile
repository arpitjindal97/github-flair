FROM mongo:latest

RUN mkdir /arpit

COPY main /arpit
COPY entrypoint.sh /arpit
COPY cron_file /arpit
COPY update /arpit

#COPY cron_file /arpit/tor_restart.sh
#RUN chmod +x /arpit/tor_restart.sh
#RUN echo "* * * * * root /arpit/tor_restart.sh > /dev/fd/1" > /etc/cron.d/tor_cron
#RUN crontab /etc/cron.d/tor_cron

EXPOSE 443
RUN chmod +x /arpit/entrypoint.sh
ENTRYPOINT ["/arpit/entrypoint.sh"]
