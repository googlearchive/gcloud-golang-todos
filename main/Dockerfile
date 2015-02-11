# Dockerfile extending the generic Go image with application files for a
# single application.
FROM google/appengine-go

ADD . /app
RUN /bin/bash /app/_ah/build.sh
