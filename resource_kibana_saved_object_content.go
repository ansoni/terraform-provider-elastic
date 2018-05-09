package main

import (
	"github.com/hashicorp/terraform/helper/schema"
        _ "github.com/hashicorp/terraform/terraform"
	"log"
	"fmt"
	"errors"
	"encoding/json"
)

func resourceKibanaSavedObjectContent() *schema.Resource {
	return &schema.Resource{
		Create: resourceElasticSavedObjectContentCreate,
		Read: resourceElasticSavedObjectContentRead,
		Update: resourceElasticSavedObjectContentUpdate,
		Delete: resourceElasticSavedObjectContentDelete,
		Schema: map[string]*schema.Schema{
			"saved_object_id": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"saved_object_type": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"attributes": &schema.Schema{
				Type:             schema.TypeString,
				Required:         true,
			},
		},
	}
}

func resourceElasticSavedObjectContentCreate(d *schema.ResourceData, meta interface{}) error {
	url := meta.(*ElasticInfo).kibanaUrl
	attributes := d.Get("attributes").(string)
	object_id := d.Get("saved_object_id").(string)
	saved_object_type := d.Get("saved_object_type").(string)

	url = fmt.Sprintf("%v/api/saved_objects/%v/%v", url, saved_object_type, object_id)
	existingSavedObjectBytes, err := getKibRequest(d, meta, url)
	if err != nil {
		return err
	}

	var existingSavedObject SavedObjectHeader
	json.Unmarshal(*existingSavedObjectBytes, &existingSavedObject)
	existingSavedObject.UpdatedAt=""
	existingSavedObject.Id=""
	existingSavedObject.ObjectType=""

	if len(existingSavedObject.Attributes) != 0 {
		errors.New(fmt.Sprintf("Existing Object: %s, already has content: %v", object_id, existingSavedObject.Attributes))
	}

	json.Unmarshal([]byte(attributes), &existingSavedObject.Attributes)

	body, err := json.Marshal(&existingSavedObject)
	if err != nil {
		return err
	}

	respBody, err := putKibRequest(d, meta, url, string(body))
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

func resourceElasticSavedObjectContentRead(d *schema.ResourceData, meta interface{}) error {
	url := meta.(*ElasticInfo).kibanaUrl
	id := d.Id()
	saved_object_type := d.Get("saved_object_type").(string)

	url = fmt.Sprintf("%v/api/saved_objects/%v/%v", url, saved_object_type, id)
	respBody, err := getKibRequest(d, meta, url)
	if err != nil {
       	    return err
	}

	var savedObject SavedObjectHeader
	json.Unmarshal(*respBody, &savedObject)
	return nil
}

func resourceElasticSavedObjectContentUpdate(d *schema.ResourceData, meta interface{}) error {
	url := meta.(*ElasticInfo).kibanaUrl
	id := d.Id()
	saved_object_type := d.Get("saved_object_type").(string)

	attributes := d.Get("attributes").(string)
	var savedObjectHeader SavedObjectHeader
	json.Unmarshal([]byte(attributes), &savedObjectHeader.Attributes)
	body, err := json.Marshal(&savedObjectHeader)
	if err != nil {
		return err
	}

	savedObjectHeader.UpdatedAt=""
	savedObjectHeader.Id=""
	savedObjectHeader.ObjectType=""

	url = fmt.Sprintf("%v/api/saved_objects/%v/%v", url, saved_object_type, id)
	respBody, err := putKibRequest(d, meta, url, string(body))
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

func resourceElasticSavedObjectContentDelete(d *schema.ResourceData, meta interface{}) error {
	url := meta.(*ElasticInfo).kibanaUrl
	id := d.Id()
	saved_object_type := d.Get("saved_object_type").(string)

	url = fmt.Sprintf("%v/api/saved_objects/%v/%v", url, saved_object_type, id)
	_, err := deleteKibRequest(d, meta, url)	
	if err != nil {
       		return err    
	}
	d.Set("version", nil)
	d.SetId("")
	return nil
}


