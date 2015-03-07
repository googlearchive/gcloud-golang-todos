# Dockerfile extending the generic Go image with application files for a
# single application.
FROM google/appengine-go

# TODO: Remove once base image updates to 1.4.2 toolchain.
RUN rm -rf /goroot && mkdir /goroot && curl https://storage.googleapis.com/golang/go1.4.2.linux-amd64.tar.gz | tar xvzf - -C /goroot --strip-components=1

ADD . /app
RUN /bin/bash /app/_ah/build.sh
