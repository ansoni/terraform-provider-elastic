package main

import (
        "github.com/hashicorp/terraform/helper/schema"
        "github.com/hashicorp/terraform/terraform"
	_ "net/http"
)

func Provider() terraform.ResourceProvider {
        return &schema.Provider{
                Schema: map[string]*schema.Schema{
			"elasticsearch_url": &schema.Schema{
                                Type:        schema.TypeString,
                                Optional:    true,
                                DefaultFunc: schema.EnvDefaultFunc("ELASTICSEARCH_URL", nil),
                                Description: "Elasticsearch URL",
                        },
			"elasticsearch_username": &schema.Schema{
                                Type:        schema.TypeString,
                                Optional:    true,
                                DefaultFunc: schema.EnvDefaultFunc("ELASTICSEARCH_USERNAME", nil),
                                Description: "Elasticsearch Username",
                        },
			"elasticsearch_password": &schema.Schema{
                                Type:        schema.TypeString,
                                Optional:    true,
                                DefaultFunc: schema.EnvDefaultFunc("ELASTICSEARCH_PASSWORD", nil),
                                Description: "Elasticsearch Password",
                        },
			"kibana_url": &schema.Schema{
                                Type:        schema.TypeString,
                                Optional:    true,
                                DefaultFunc: schema.EnvDefaultFunc("KIBANA_URL", nil),
                                Description: "Elasticsearch URL",
                        },
			"kibana_username": &schema.Schema{
                                Type:        schema.TypeString,
                                Optional:    true,
                                DefaultFunc: schema.EnvDefaultFunc("KIBANA_USERNAME", nil),
                                Description: "Elasticsearch Username",
                        },
			"kibana_password": &schema.Schema{
                                Type:        schema.TypeString,
                                Optional:    true,
                                DefaultFunc: schema.EnvDefaultFunc("KIBANA_PASSWORD", nil),
                                Description: "Elasticsearch Password",
                        },
		},
                ResourcesMap: map[string]*schema.Resource{
			"elastic_kibana_saved_object":   resourceKibanaSavedObject(),
			"elastic_kibana_saved_object_content":   resourceKibanaSavedObjectContent(),
		},
		ConfigureFunc: providerConfigure,
	}
}

type ElasticInfo struct {
	kibanaUrl string
	kibanaUsername string
	kibanaPassword string
	elasticsearchUrl string

}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	url := d.Get("kibana_url").(string)
	kibanaUsername := d.Get("kibana_username").(string)
	kibanaPassword := d.Get("kibana_password").(string)
	esUrl := d.Get("elasticsearch_url").(string)
	elasticInfo := &ElasticInfo{kibanaUrl: url, elasticsearchUrl: esUrl, kibanaUsername:kibanaUsername, kibanaPassword:kibanaPassword }	
	return elasticInfo, nil
}
