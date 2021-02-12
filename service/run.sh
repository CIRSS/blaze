#!/bin/bash

# avoid error message on ctrl-c
cleanup() {
    echo
    exit 0
}
trap cleanup EXIT

# run the service
echo
echo "--------------------------------------------------------------------------"
echo "The Blazegraph service will now start in the REPRO."
echo "Connect to it by navigating in a web browser to http://localhost:9999 "
echo
echo "Terminate the service by typing CTRL-C in this terminal."
echo "--------------------------------------------------------------------------"
cd ${BLAZEGRAPH_DOT_DIR}
${BLAZEGRAPH_CMD} 2>&1 > `eval echo ${BLAZEGRAPH_LOG}`
