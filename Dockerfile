FROM debian

MAINTAINER roman.saveljev@haltian.com

RUN apt-get update
RUN apt-get -y install rsync make markdown lynx bash

RUN mkdir /opt/symbols
ADD docker/opt /opt
ADD docker/bin /bin
ADD help /opt/help

CMD ["usage-help"]
