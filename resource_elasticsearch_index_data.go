package main

import (
	"github.com/hashicorp/terraform/helper/schema"
        _ "github.com/hashicorp/terraform/terraform"
	"log"
	"fmt"
        "io/ioutil"
	_ "errors"
)

func resourceElasticsearchIndexData() *schema.Resource {
	return &schema.Resource{
		Create: resourceElasticsearchIndexDataCreate,
		Read: resourceElasticsearchIndexDataRead,
		Update: resourceElasticsearchIndexDataUpdate,
		Delete: resourceElasticsearchIndexDataDelete,
		Schema: map[string]*schema.Schema{
			"index_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"file_path": &schema.Schema{
				Type:     schema.TypeString,
				Optional:  true,
			},
			"file_url": &schema.Schema{
				Type:     schema.TypeString,
                                ConflictsWith: []string{"file_path"},
				Optional:  true,
			},
			"content_type": &schema.Schema{
				Type:     schema.TypeString,
                                Default: "application/x-ndjson",
				Optional:  true,
			},
		},
	}
}

func resourceElasticsearchIndexDataCreate(d *schema.ResourceData, meta interface{}) error {
	url := meta.(*ElasticInfo).elasticsearchUrl
	username := meta.(*ElasticInfo).elasticsearchUsername
	password := meta.(*ElasticInfo).elasticsearchPassword
	index := d.Get("index_name").(string)
	contentType := d.Get("content_type").(string)
        reader, err := getFileOrUrlReader(d)
        if err != nil {
		return err
        }

	url = fmt.Sprintf("%v/%v/doc/_bulk?pretty", url, index)

	respBody, err := streamingRequest("POST", d, meta, url, contentType, username, password, reader)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(respBody)
      	log.Printf("Respose: %s", body)
        if err != nil {
                return err
        }

	d.SetId("a")

	return err
}

func resourceElasticsearchIndexDataRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceElasticsearchIndexDataUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceElasticsearchIndexDataDelete(d *schema.ResourceData, meta interface{}) error {
	url := meta.(*ElasticInfo).elasticsearchUrl
	username := meta.(*ElasticInfo).elasticsearchUsername
	password := meta.(*ElasticInfo).elasticsearchPassword
	index := d.Get("index_name").(string)

	url = fmt.Sprintf("%v/%v", url, index)
	_, err := deleteKibRequest(d, meta, url, username, password)	
	if err != nil {
       		return err    
	}
	d.SetId("")
	return nil
}


