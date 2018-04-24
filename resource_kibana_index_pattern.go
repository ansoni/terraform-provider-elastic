package main

import (
	"github.com/hashicorp/terraform/helper/schema"
        _ "github.com/hashicorp/terraform/terraform"
	"io/ioutil"
	"net/http"
	"bytes"
	"log"
	"fmt"
	"errors"
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

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(body))
	if err != nil {
		return err
	}
	req.Header.Add("kbn-xsrf", "true")
	req.Header.Add("content-type", "application/json")

	log.Printf("POST new index-pattern to %v", url)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var indexPattern IndexPattern
	json.Unmarshal(respBody, &indexPattern)	
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
	version := d.Get("version")

	url = fmt.Sprintf("%v/api/saved_objects/index-pattern/%v", url, objectId)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("kbn-xsrf", "true")
	req.Header.Add("content-type", "application/json")
	log.Printf("Get existing index-pattern with %v", url)
	resp, err := client.Do(req)
	if err != nil {
       	    return err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
       	    return err
	}

	var indexPattern IndexPattern
	json.Unmarshal(respBody, &indexPattern)	
	if indexPattern.Version != version {
		return errors.New(fmt.Sprintf("Index Pattern has been modified.  Found version %v, Expected version %v", err, indexPattern.Version, version))
	}

	return nil
}

func resourceElasticIndexPatternUpdate(d *schema.ResourceData, meta interface{}) error {
	url := meta.(*ElasticInfo).kibanaUrl
	objectId := d.Get("object_id")

	body := d.Get("body").(string)

	url = fmt.Sprintf("%v/api/saved_objects/index-pattern/%v", url, objectId)
	client := &http.Client{}
	req, err := http.NewRequest("PUT", url, bytes.NewBufferString(body))
	if err != nil {
		return err
	}
	req.Header.Add("kbn-xsrf", "true")
	req.Header.Add("content-type", "application/json")
	log.Printf("Update index-pattern to %v", url)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var indexPattern IndexPattern
	json.Unmarshal(respBody, &indexPattern)	
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
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("kbn-xsrf", "true")
	req.Header.Add("content-type", "application/json")
	log.Printf("DELETE index-pattern to %v", url)
	_, err = client.Do(req)
	if err != nil {
       		return err    
	}
	return nil
}


