package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
        _ "github.com/hashicorp/terraform/terraform"
	"log"
	"fmt"
	_ "errors"
	"encoding/json"
)

func resourceKibanaSavedObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceElasticSavedObjectCreate,
		Read: resourceElasticSavedObjectRead,
		Update: resourceElasticSavedObjectUpdate,
		Delete: resourceElasticSavedObjectDelete,
		Schema: map[string]*schema.Schema{
			"saved_object_type": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{ "index-pattern", "visualization", "search", "timelion-sheet","dashboard",}, false),
			},
			"object_id": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"attributes": &schema.Schema{
				Type:             schema.TypeString,
				Default: "",
				Optional:         true,
			},
		},
	}
}

func resourceElasticSavedObjectCreate(d *schema.ResourceData, meta interface{}) error {
	url := meta.(*ElasticInfo).kibanaUrl
	username := meta.(*ElasticInfo).kibanaUsername
	password := meta.(*ElasticInfo).kibanaPassword
	attributes := d.Get("attributes").(string)

	saved_object_type := d.Get("saved_object_type").(string)

	url = fmt.Sprintf("%v/api/saved_objects/%v", url, saved_object_type)

	var savedObjectHeader SavedObjectHeader
	json.Unmarshal([]byte(attributes), &savedObjectHeader.Attributes)
	if savedObjectHeader.Attributes == nil {
		savedObjectHeader.Attributes=make(map[string]interface{})		
	}
	name, found := d.GetOk("name")
	if found {
		savedObjectHeader.Attributes["title"]=name.(string)
	}

	body, err := json.Marshal(&savedObjectHeader)
	if err != nil {
		return err
	}
        log.Printf("Create: %v", body)

	respBody, err := postKibRequest(d, meta, url, username, password, string(body))
	if err != nil {
		return err
	}

	var savedObject SavedObjectHeader
	json.Unmarshal(*respBody, &savedObject)	
	log.Printf("Raw Body: %s", respBody)
	log.Printf("ID: %s", savedObject.Id)
	log.Printf("UpdatedAt: %s", savedObject.UpdatedAt)
	log.Printf("Version: %v", savedObject.Version)

	d.SetId(savedObject.Id)
	d.Set("version", savedObject.Version)

	return err
}

func resourceElasticSavedObjectRead(d *schema.ResourceData, meta interface{}) error {
	url := meta.(*ElasticInfo).kibanaUrl
	username := meta.(*ElasticInfo).kibanaUsername
	password := meta.(*ElasticInfo).kibanaPassword
	id := d.Id()
	saved_object_type := d.Get("saved_object_type").(string)

	url = fmt.Sprintf("%v/api/saved_objects/%v/%v", url, saved_object_type, id)
	respBody, err := getKibRequest(d, meta, url, username, password)
	if err != nil {
       	    return err
	}

	var savedObject SavedObjectHeader
	json.Unmarshal(*respBody, &savedObject)
	return nil
}

func resourceElasticSavedObjectUpdate(d *schema.ResourceData, meta interface{}) error {

	url := meta.(*ElasticInfo).kibanaUrl
	username := meta.(*ElasticInfo).kibanaUsername
	password := meta.(*ElasticInfo).kibanaPassword
	id := d.Id()
	saved_object_type := d.Get("saved_object_type").(string)

	attributes := d.Get("attributes").(string)
	var savedObjectHeader SavedObjectHeader
	json.Unmarshal([]byte(attributes), &savedObjectHeader.Attributes)
	name, found := d.GetOk("name")
	if found {
		savedObjectHeader.Attributes["title"]=name.(string)
	}
	body, err := json.Marshal(&savedObjectHeader)
	if err != nil {
		return err
	}

	url = fmt.Sprintf("%v/api/saved_objects/%v/%v", url, saved_object_type, id)
	respBody, err := putKibRequest(d, meta, url, username, password, string(body))
	if err != nil {
		return err
	}

	var savedObject SavedObjectHeader
	json.Unmarshal(*respBody, &savedObject)	
	log.Printf("Raw Body: %s", respBody)
	log.Printf("ID: %s", savedObject.Id)
	log.Printf("UpdatedAt: %s", savedObject.UpdatedAt)
	log.Printf("Version: %v", savedObject.Version)
	d.Set("version", savedObject.Version)
	return nil
}

func resourceElasticSavedObjectDelete(d *schema.ResourceData, meta interface{}) error {
	url := meta.(*ElasticInfo).kibanaUrl
	username := meta.(*ElasticInfo).kibanaUsername
	password := meta.(*ElasticInfo).kibanaPassword
	id := d.Id()
	saved_object_type := d.Get("saved_object_type").(string)

	url = fmt.Sprintf("%v/api/saved_objects/%v/%v", url, saved_object_type, id)
	_, err := deleteKibRequest(d, meta, url, username, password)	
	if err != nil {
       		return err    
	}
	d.Set("version", nil)
	d.SetId("")
	return nil
}


