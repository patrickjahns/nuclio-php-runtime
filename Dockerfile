ARG NUCLIO_LABEL=0.7.1
ARG NUCLIO_ARCH=amd64
ARG NUCLIO_BASE_IMAGE=php:7.3-fpm-alpine
ARG NUCLIO_ONBUILD_IMAGE=nuclio/handler-builder-golang-onbuild:${NUCLIO_LABEL}-${NUCLIO_ARCH}-alpine

# Supplies processor uhttpc, used for healthcheck
FROM nuclio/uhttpc:0.0.1-amd64 as uhttpc

# Builds source, supplies processor binary and handler plugin
FROM ${NUCLIO_ONBUILD_IMAGE} as builder

# From the base image
FROM ${NUCLIO_BASE_IMAGE}


# Copy required objects from the suppliers
COPY --from=builder /home/nuclio/bin/processor /usr/local/bin/processor
COPY --from=builder /home/nuclio/bin/handler.so /opt/nuclio/handler.so
COPY --from=uhttpc /home/nuclio/bin/uhttpc /usr/local/bin/uhttpc

# Readiness probe
HEALTHCHECK --interval=1s --timeout=3s CMD /usr/local/bin/uhttpc --url http://127.0.0.1:8082/ready || exit 1


# add php src
ADD test.php /var/task/src/test.php
ADD php-fpm.conf /var/task

# set script
ENV PHP_FPM_BIN=/usr/local/sbin/php-fpm
ENV PHP_SCRIPT=/var/task/src/test.php

# ensure permissions are valid
RUN chown -R www-data:www-data /var/task

# Run processor with configuration and platform configuration
CMD [ "processor" ]