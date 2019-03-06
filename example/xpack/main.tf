provider "elastic" {
  kibana_url = "http://127.0.0.1:5601"
  elasticsearch_url = "http://127.0.0.1:9200"
}

resource "elastic_elasticsearch_watcher" "watcherA" {
  active = true
  name = "awesome"
  trigger = <<-EOF
	{
	  "schedule" : { "cron" : "0 0/1 * * * ?" }
	}
	EOF
  input = <<-EOF
	  {
	    "search" : {
	      "request" : {
	        "indices" : [
	          "logstash*"
	        ],
	        "body" : {
	          "query" : {
	            "bool" : {
	              "must" : {
	                "match": {
	                   "response": 404
	                }
	              },
	              "filter" : {
	                "range": {
	                  "@timestamp": {
	                    "from": "{{ctx.trigger.scheduled_time}}||-5m",
	                    "to": "{{ctx.trigger.triggered_time}}"
	                  }
	                }
	              }
	            }
	          }
	        }
	      }
	    }
	  }
	EOF
  condition = <<-EOF
	{
	    "compare" : { "ctx.payload.hits.total" : { "gt" : 0 }}
	}
	EOF
}
