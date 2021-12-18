#!/bin/sh

EXTRA_ARGS=$EXTRA_ARGS
if [ $LISTENPORT ]; then
    EXTRA_ARGS=$EXTRA_ARGS' -ListenPort='$LISTENPORT
fi

if [ $ENDPOINTNAME ]; then
    EXTRA_ARGS=$EXTRA_ARGS' -EndpointName='$ENDPOINTNAME
fi

if [ $DEV ]; then
    echo $EXTRA_ARGS
fi

/var/app/magicLite $EXTRA_ARGS "$@"
