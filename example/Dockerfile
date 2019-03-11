FROM patrickjahns/nucleo-php-runtime

# add php src
ADD example/test.php /var/task/src/test.php

# configure the script to be the default entrypoint
ENV PHP_SCRIPT=/var/task/src/test.php