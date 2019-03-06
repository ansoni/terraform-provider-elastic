package main

import (
	"github.com/hashicorp/terraform/helper/schema"
        _ "github.com/hashicorp/terraform/terraform"
	"log"
	"fmt"
	_ "errors"
	"encoding/json"
)

func resourceKibanaCanvas() *schema.Resource {
	return &schema.Resource{
		Create: resourceElasticCanvasCreate,
		Read: resourceElasticCanvasRead,
		Update: resourceElasticCanvasUpdate,
		Delete: resourceElasticCanvasDelete,
		Schema: map[string]*schema.Schema{
			"canvas_id": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"space_id": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"contents": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceElasticCanvasCreate(d *schema.ResourceData, meta interface{}) error {
	url := kibanaUrl(d, meta)
	username := meta.(*ElasticInfo).kibanaUsername
	password := meta.(*ElasticInfo).kibanaPassword

	canvasId := d.Get("canvas_id").(string)
	name := d.Get("name").(string)
	contents := d.Get("contents").(string)

	url = fmt.Sprintf("%v/api/canvas/workpad", url)

	contentMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(contents), &contentMap)
	if err != nil {
		return err
	}	

	contentMap["id"]=canvasId
	contentMap["name"]=name

	body, err := json.Marshal(&contentMap)
	if err != nil {
		return err
	}	
        log.Printf("Create: %v", body)
	respBody, err := postKibRequest(d, meta, url, username, password, string(body))
	if err != nil {
		return err
	}
	log.Printf("Response: %v", respBody) 

	d.SetId(canvasId)

	return err
}

func resourceElasticCanvasRead(d *schema.ResourceData, meta interface{}) error {
	url := kibanaUrl(d, meta)
	username := meta.(*ElasticInfo).kibanaUsername
	password := meta.(*ElasticInfo).kibanaPassword
	id := d.Id()

	url = fmt.Sprintf("%v/api/canvas/workpad/%v", url, id)
	respBody, err := getKibRequest(d, meta, url, username, password)
	if err != nil {
       	    return err
	}
	log.Printf("Response: %v", respBody) 

	//var savedObject CanvasHeader
	//json.Unmarshal(*respBody, &savedObject)
	return nil
}

func resourceElasticCanvasUpdate(d *schema.ResourceData, meta interface{}) error {

	url := kibanaUrl(d, meta)
	username := meta.(*ElasticInfo).kibanaUsername
	password := meta.(*ElasticInfo).kibanaPassword
	id := d.Id()
	canvasId := d.Get("canvas_id").(string)
	name := d.Get("name").(string)
	contents := d.Get("contents").(string)

	url = fmt.Sprintf("%v/api/canvas/workpad", url)

	contentMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(contents), &contentMap)
	if err != nil {
		return err
	}	

	contentMap["id"]=canvasId
	contentMap["name"]=name

	body, err := json.Marshal(&contentMap)
	if err != nil {
		return err
	}	

	url = fmt.Sprintf("%v/api/canvas/workpad/%v", url, id)
	respBody, err := putKibRequest(d, meta, url, username, password, string(body))
	if err != nil {
		return err
	}
	log.Printf("Response: %v", respBody)

	d.SetId(canvasId)
	return nil
}

func resourceElasticCanvasDelete(d *schema.ResourceData, meta interface{}) error {
	url := kibanaUrl(d, meta)
	username := meta.(*ElasticInfo).kibanaUsername
	password := meta.(*ElasticInfo).kibanaPassword
	id := d.Id()

	url = fmt.Sprintf("%v/api/canvas/workpad/%v", url, id)
	_, err := deleteKibRequest(d, meta, url, username, password)	
	if err != nil {
       		return err    
	}
	d.SetId("")
	return nil
}


