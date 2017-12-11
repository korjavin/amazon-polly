FROM korjavin/korjavin-base
RUN apt-get update
#RUN DEBIAN_FRONTEND=noninteractive apt-get install -y python3-pip python-pip
#RUN /usr/bin/pip install --upgrade --user awscli
#RUN ln -s /root/.local/bin/aws /bin/aws
RUN DEBIAN_FRONTEND=noninteractive apt-get install -y mp3wrap
RUN mkdir /site
ADD index.html /site/index.html
#ADD aws /root/.aws
RUN apt-get install -y ca-certificates
ADD webPolly /site/webPolly

WORKDIR /site
ENTRYPOINT ["/site/webPolly"]
