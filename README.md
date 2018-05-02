# terraform-provider-elastic

A Terraform Provider for Elastic.co products (Elasticsearch, Kibana, etc).  Please test and provide feedback and examples of things that don't work.  Most of these work along the lines of taking in a raw body output of a Kibana SavedObject.  From their, we modify it with supplied variables which are typically the primary keys of index-patterns, visualizations, etc.  Our goal is to one day be able to create objects without requiring a body at all.

# Build

	go get
	go build

# Install

	cp terraform-provider-elastic ~/.terraform.d/plugins/	

# Run Example

	cd example
	docker-compose up
	bash load_data.sh
	terraform init
	terraform apply

You should have a whole bunch of visualizations built off of the Shakespeare example

# Generic Saved-Object

Supported Types:  timelion-sheet, index-pattern, search, visualization, dashboard

Use a template to modify a SavedSearch Attributes item

	data "template_file" "timelion" {
		template = <<EOF
	{"hits": 0, "description": "", "timelion_rows": 2, "title": "$${title}", "version": 1, "timelion_chart_height": 275, "timelion_sheet": [".es(*).abs()"], "timelion_columns": 2, "timelion_interval": "auto"}
	EOF
		vars {
			title = "rar"
		}
	}
	
        resource "elastic_kibana_saved_object" "test" {
		attributes = "${data.template_file.timelion.rendered}"
		saved_object_type = "timelion-sheet"
		name = "awesome"
		description = "w00t"
	}

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
- [ ] No Watcher Support
- [ ] No Graph Support
- [ ] No Canvas Support
- [ ] No TimeLion Support



