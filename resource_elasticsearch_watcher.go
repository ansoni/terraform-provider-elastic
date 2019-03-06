package main

import (
	"github.com/hashicorp/terraform/helper/schema"
        _ "github.com/hashicorp/terraform/terraform"
	"log"
	"fmt"
	_ "errors"
	"encoding/json"
)

func resourceElasticsearchWatcher() *schema.Resource {
	return &schema.Resource{
		Create: resourceElasticWatcherCreate,
		Read: resourceElasticWatcherRead,
		Update: resourceElasticWatcherUpdate,
		Delete: resourceElasticWatcherDelete,
		Schema: map[string]*schema.Schema{
			"active": &schema.Schema{
				Type:     schema.TypeBool,
				Default: true,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"trigger": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"input": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"condition": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"actions": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default: "{}",
			},
			"metadata": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default: "{}",
			},
			"throttle_period": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func genericMap(data string) map[string]interface{} {
	contentMap := make(map[string]interface{})
        err := json.Unmarshal([]byte(data), &contentMap)
	if err != nil {
		log.Printf("Error %v", err)
		return nil
	}
	return contentMap
}

func watcherFromInput(d *schema.ResourceData) Watcher {
        log.Printf("WatcherFromInput!")
        obj := Watcher{}

        populateStruct(d, "trigger", func(value interface{}) { obj.Trigger = genericMap(value.(string)) })
        populateStruct(d, "actions", func(value interface{}) { obj.Actions = genericMap(value.(string)) } )
        populateStruct(d, "metadata", func(value interface{}) { obj.Metadata = genericMap(value.(string)) } )
        populateStruct(d, "condition", func(value interface{}) { obj.Condition = genericMap(value.(string)) } )
        populateStruct(d, "input", func(value interface{}) { obj.Input = genericMap(value.(string)) } )
        populateStruct(d, "throttle_period", func(value interface{}) { obj.ThrottlePeriod = value.(string)} )

        return obj
}

func resourceElasticWatcherCreate(d *schema.ResourceData, meta interface{}) error {
	url := meta.(*ElasticInfo).elasticsearchUrl
        username := meta.(*ElasticInfo).elasticsearchUsername
        password := meta.(*ElasticInfo).elasticsearchPassword

	name := d.Get("name").(string)
	active := d.Get("active").(bool)

	watcher := watcherFromInput(d)

	url = fmt.Sprintf("%v/_xpack/watcher/watch/%v?active=%v", url, name, active)

	body, err := json.Marshal(&watcher)
	if err != nil {
		return err
	}	
        log.Printf("Create: %v", body)
	respBody, err := postKibRequest(d, meta, url, username, password, string(body))
	if err != nil {
		return err
	}
	log.Printf("Response: %v", respBody) 

	d.SetId(name)

	return err
}

func resourceElasticWatcherRead(d *schema.ResourceData, meta interface{}) error {
	url := meta.(*ElasticInfo).elasticsearchUrl
        username := meta.(*ElasticInfo).elasticsearchUsername
        password := meta.(*ElasticInfo).elasticsearchPassword
	id := d.Id()

	url = fmt.Sprintf("%v/_xpack/watcher/watch/%v", url, id)
	respBody, err := getKibRequest(d, meta, url, username, password)
	if err != nil {
       	    return err
	}
	log.Printf("Response: %v", respBody) 

	return nil
}

func resourceElasticWatcherUpdate(d *schema.ResourceData, meta interface{}) error {
	url := meta.(*ElasticInfo).elasticsearchUrl
        username := meta.(*ElasticInfo).elasticsearchUsername
        password := meta.(*ElasticInfo).elasticsearchPassword

	id := d.Id()
	active := d.Get("active").(bool)

	watcher := watcherFromInput(d)

	url = fmt.Sprintf("%v/_xpack/watcher/watch/%v?active=%v", url, id, active)

	body, err := json.Marshal(&watcher)
	if err != nil {
		return err
	}	
        log.Printf("Create: %v", body)
	respBody, err := postKibRequest(d, meta, url, username, password, string(body))
	if err != nil {
		return err
	}
	log.Printf("Response: %v", respBody) 

	return nil
}

func resourceElasticWatcherDelete(d *schema.ResourceData, meta interface{}) error {
	url := meta.(*ElasticInfo).elasticsearchUrl
        username := meta.(*ElasticInfo).elasticsearchUsername
        password := meta.(*ElasticInfo).elasticsearchPassword
	id := d.Id()

	url = fmt.Sprintf("%v/_xpack/watcher/watch/%v", url, id)
	_, err := deleteKibRequest(d, meta, url, username, password)	
	if err != nil {
       		return err    
	}
	d.SetId("")
	return nil
}


