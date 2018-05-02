package main

import (
	"github.com/hashicorp/terraform/helper/schema"
        _ "github.com/hashicorp/terraform/terraform"
	"log"
	"fmt"
	_ "errors"
	"encoding/json"
)

type IndexPattern struct {
	Id string `json:"id"`
	ObjectType string `json:"type"`
	UpdatedAt string `json:"updated_at"`
	Version int `json:"version"`
}

func resourceKibanaIndexPattern() *schema.Resource {
	return &schema.Resource{
		Create: resourceElasticIndexPatternCreate,
		Read: resourceElasticIndexPatternRead,
		Update: resourceElasticIndexPatternUpdate,
		Delete: resourceElasticIndexPatternDelete,
		Schema: map[string]*schema.Schema{
			"object_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"body": &schema.Schema{
				Type:             schema.TypeString,
				Required:         true,
			},
		},
	}
}

func resourceElasticIndexPatternCreate(d *schema.ResourceData, meta interface{}) error {
	url := meta.(*ElasticInfo).kibanaUrl
	_ = d.Get("name").(string)

	body := d.Get("body").(string)

	url = fmt.Sprintf("%v/api/saved_objects/index-pattern", url)

	respBody, err := postKibRequest(d, meta, url, body)
	if err != nil {
		return err
	}

	var indexPattern IndexPattern
	json.Unmarshal(*respBody, &indexPattern)	
	log.Printf("Raw Body: %s", respBody)
	log.Printf("ID: %s", indexPattern.Id)
	log.Printf("ObjectType: %s", indexPattern.ObjectType)
	log.Printf("UpdatedAt: %s", indexPattern.UpdatedAt)
	log.Printf("Version: %v", indexPattern.Version)

	d.SetId(indexPattern.Id)
	d.Set("object_id", indexPattern.Id)
	d.Set("version", indexPattern.Version)

	return err
}

func resourceElasticIndexPatternRead(d *schema.ResourceData, meta interface{}) error {
	url := meta.(*ElasticInfo).kibanaUrl
	objectId := d.Get("object_id")
	//version := d.Get("version")

	url = fmt.Sprintf("%v/api/saved_objects/index-pattern/%v", url, objectId)
	respBody, err := getKibRequest(d, meta, url)
	if err != nil {
       	    return err
	}

	var indexPattern IndexPattern
	json.Unmarshal(*respBody, &indexPattern)
/*
	if indexPattern.Version != version {
		return errors.New(fmt.Sprintf("Index Pattern has been modified.  Found version %v, Expected version %v", err, indexPattern.Version, version))
	}
*/

	return nil
}

func resourceElasticIndexPatternUpdate(d *schema.ResourceData, meta interface{}) error {
	url := meta.(*ElasticInfo).kibanaUrl
	objectId := d.Get("object_id")

	body := d.Get("body").(string)

	url = fmt.Sprintf("%v/api/saved_objects/index-pattern/%v", url, objectId)
	respBody, err := putKibRequest(d, meta, url, body)
	if err != nil {
		return err
	}

	var indexPattern IndexPattern
	json.Unmarshal(*respBody, &indexPattern)	
	log.Printf("Raw Body: %s", respBody)
	log.Printf("ID: %s", indexPattern.Id)
	log.Printf("ObjectType: %s", indexPattern.ObjectType)
	log.Printf("UpdatedAt: %s", indexPattern.UpdatedAt)
	log.Printf("Version: %v", indexPattern.Version)
	d.Set("object_id", indexPattern.Id)
	d.Set("version", indexPattern.Version)
	return nil
}

func resourceElasticIndexPatternDelete(d *schema.ResourceData, meta interface{}) error {
	url := meta.(*ElasticInfo).kibanaUrl
	objectId := d.Get("object_id")

	d.Set("object_id", nil)
	d.Set("version", nil)

	url = fmt.Sprintf("%v/api/saved_objects/index-pattern/%v", url, objectId)
	_, err := deleteKibRequest(d, meta, url)	
	if err != nil {
       		return err    
	}
	return nil
}


