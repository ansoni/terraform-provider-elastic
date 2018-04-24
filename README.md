# terraform-provider-elastic

A Terraform Provider for Elastic.co products (Elasticsearch, Kibana, etc)

# Build

	go get
	go build

# Install

	cp terraform-provider-elastic ~/.terraform.d/plugins/	

# Index-Patterns

	resource "elastic_kibana_index_pattern" "test" {
		body = "${file("./out/index-pattern0.json")}"
	}
	
## Inputs	
body : A Kibana exported Index Pattern

	{
	   "attributes" : {
	      "title" : "shakespeare*",
	      "fields" : "[{\"name\":\"_id\",\"type\":	\"string\",\"count\":0,\"scripted\":false,...}]"
		}
	}

## outputs

- index_id : The id of the index pattern

# Saved Searches

	resource "elastic_kibana_saved_search" "test" {
	  title = "Saved Search Title"
	  index_id = "${elastic_kibana_index_pattern.test.id}"
	  body = "${file("./out/search0.json")}"
	}
	
## Inputs

body : A Kibana exported Saved Search

	{
	   "attributes" : {
	      "kibanaSavedObjectMeta" : {
	         "searchSourceJSON" : "{\"index\":\"1f5653a0-46bb-11e8-a995-fb06d59a0729\",...]}"
	      },
	      "sort" : [
	         "_score",
	         "desc"
	      ],
	      "columns" : [
	         "_source"
	      ],
	      "description" : "",
	      "title" : "Henry IV",
	      "hits" : 0,
	      "version" : 1
	   }
	}
	
## outputs

- search_id : The id of the Saved Search

# Todo aka Request for Pull Requests

- [ ] No SSL
- [ ] No Auth
- [ ] No Visualization Support
- [ ] No Dashboard Support
- [ ] No Watcher Support
- [ ] No Graph Support
- [ ] No Canvas Support
- [ ] No TimeLion Support



