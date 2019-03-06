# terraform-provider-elastic

A Terraform Provider for Elastic.co products (Elasticsearch, Kibana, etc).  Please test and provide feedback and examples of things that don't work.  Most of these work along the lines of taking in a raw body output of a Kibana SavedObject.  From their, we modify it with supplied variables which are typically the primary keys of index-patterns, visualizations, etc.  Our goal is to one day be able to create objects without requiring a body at all.

# Build

	go build

# Install (globally on machine)

Mac/Linux instructions

	mkdir -p ~/.terraform.d/plugins
	cp terraform-provider-elastic ~/.terraform.d/plugins/	

Windows

	write me! (Fork, Edit, Pull Request)

# Run Example

After installing the provider...

	cd example
	docker-compose up
	cd basic
	terraform init && terraform apply
	cd ../xpack # need to accept X-Pack Trial in Kibana
	terraform init && terraform apply

You should have a few visualizations built off of the Elastic Shakespeare data.  I tried to have an example of every resource in the example folder. 

# Kibana Space Provisioning

Provision a Kibana Space

	resource "elastic_kibana_space" "engineering" {
	  space_id = "engineering" 
	  name = "Engineering"
	  description = "Where Engineering Keeps Its Gold!"
	  initials = "pe"
	}

# Canvas

Canvas API is undocumented, but this code does work against 6.6.1

	resource "elastic_kibana_canvas" "engineering" {
	  canvas_id = "rer"                     # Your unique name for this canvas
	  name = "mah"                          # Display name for this canvas
	  space_id = "engineering"              # Optional Space to use
	  contents = "${file("./canvas.json")}" # document exported from Kibana Web Console/API
	}

# Generic Saved-Object

Supported Types:  timelion-sheet, index-pattern, search, visualization, dashboard

Use a template to modify a SavedSearch Attributes item.  The attributes field should refence the attributes field inside of a saved_objects api Get request. 

	data "template_file" "timelion" {
		template = <<EOF
	{
	   "timelion_interval" : "auto",
	   "timelion_columns" : 2,
	   "timelion_chart_height" : 275,
	   "hits" : 0,
	   "description" : "",
	   "timelion_sheet" : [
	      ".es(*).abs()"
	   ],
	   "version" : 1,
	   "timelion_rows" : 2,
	   "title" : "$${title}"
	}
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

outputs

id
saved_object_type
title 

# Generic Saved-Object-Content

More advanced Dashboards in Kibana have deep-link to dashboards within visualizations within the Dashboard.  This creates a "chicken and the egg"/circular dependency where we need to know the dashboard id before it has been created.

In the previous example, you would remove the attributes field and could deploy it later in your terraform with this snippet:

	resource "elastic_kibana_saved_object_content" "test" {
                attributes = "${data.template_file.timelion.rendered}"
                saved_object_type = "${elastic_kibana_saved_object.test.saved_object_type}"
                saved_object_id = "${elastic_kibana_saved_object.test.id}"
        }

# Todo aka Request for Pull Requests

- [ ] No Graph Support
- [ ] No ML Support
