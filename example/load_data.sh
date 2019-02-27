#!/bin/bash

OUTPUT=/tmp/shakespeare.json

if [ ! -e ${OUTPUT} ];then
	curl https://download.elastic.co/demos/kibana/gettingstarted/shakespeare_6.0.json > ${OUTPUT}
fi

echo curl -H 'Content-Type: application/x-ndjson' -XPOST 'localhost:9200/shakespeare/doc/_bulk?pretty' --data-binary @${OUTPUT}
curl -H 'Content-Type: application/x-ndjson' -XPOST 'localhost:9200/shakespeare/doc/_bulk?pretty' --data-binary @${OUTPUT}
