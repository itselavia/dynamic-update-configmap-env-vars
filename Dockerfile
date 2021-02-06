FROM alpine:3.12

RUN apk update && APP_NAME=update_env && addgroup -g 2000 $APP_NAME && \
    adduser -s /bin/bash -u 2000 -G $APP_NAME \
    -h /opt/$APP_NAME -H -D $APP_NAME

USER $APP_NAME

ADD bin/$APP_NAME /opt/

ENTRYPOINT ["/opt/update_env"]