FROM alpine
MAINTAINER "Gildas Le Nadan"

EXPOSE 8080

COPY glock /usr/local/bin

CMD /usr/local/bin/glock
