provider "elastic" {
  kibana_url = "http://127.0.0.1:5601"
  elasticsearch_url = "http://127.0.0.1:9200"
}

# Import some data!
resource "elastic_elasticsearch_index_data" "shakespeare" {
  index_name = "shakespeare",
  file_url = "https://download.elastic.co/demos/kibana/gettingstarted/shakespeare_6.0.json"
}

# Create a Kibana Space
resource "elastic_kibana_space" "engineering" {
  space_id = "engineering"
  name = "Engineering-"
  description = "Where- Engineering Keeps Its Gold!"
  initials = "aa"
}

# Index Pattern into workspace
resource "elastic_kibana_saved_object" "index_pattern_space" {
	space_id = "${elastic_kibana_space.engineering.id}"
	depends_on = [ "elastic_kibana_saved_object.index_pattern" ]
	saved_object_type = "index-pattern"
	name = "shakespeare"
	description = "Shakespeare Index Pattern"
}

# Index Pattern into default workspace
resource "elastic_kibana_saved_object" "index_pattern" {
	depends_on = [ "elastic_kibana_saved_object.index_pattern" ]
	saved_object_type = "index-pattern"
	name = "shakespeare"
	description = "Shakespeare Index Pattern"
}

# Index Pattern (actual content/attributes)
resource "elastic_kibana_saved_object_content" "index_pattern_content" {
  saved_object_id = "${elastic_kibana_saved_object.index_pattern.id}"
  saved_object_type = "${elastic_kibana_saved_object.index_pattern.saved_object_type}"
  attributes = "${file("./objects/index-pattern0.json")}"
}

# Use a template so we can replace the index
data "template_file" "search_Henry_IV" {
  template = "${file("./objects/search0.json")}"
  vars {
    index-pattern = "${elastic_kibana_saved_object.index_pattern.id}"

  }
}

resource "elastic_kibana_saved_object" "search" {
	saved_object_type = "search"
	name = "Henry IV"
	description = "Henry IV saved search"
        attributes = "${data.template_file.search_Henry_IV.rendered}"
}

data "template_file" "visualization" {
  template = "${file("./objects/visualization0.json")}"
  vars {
    index-pattern = "${elastic_kibana_saved_object.index_pattern.id}"
  }
}

resource "elastic_kibana_saved_object" "visualization" {
        count = 2 
	saved_object_type = "visualization"
	name = "Shakespeare Awesome - ${count.index}"
	description = "Shakespeare Visualization - ${count.index}"
        attributes = "${data.template_file.visualization.rendered}"
}

output "workspace" {
  value = "${elastic_kibana_space.engineering.id}"
}
