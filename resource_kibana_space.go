package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	_ "github.com/hashicorp/terraform/helper/validation"
        _ "github.com/hashicorp/terraform/terraform"
	"log"
	"fmt"
	_ "errors"
	"encoding/json"
)

func resourceKibanaSpace() *schema.Resource {
	return &schema.Resource{
		Create: resourceElasticSpaceCreate,
		Read: resourceElasticSpaceRead,
		Update: resourceElasticSpaceUpdate,
		Delete: resourceElasticSpaceDelete,
		Schema: map[string]*schema.Schema{
			"space_id": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"initials": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
                                ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
                                  v := val.(string)
                                  if len(v) > 2 {
                                    errs = append(errs, fmt.Errorf("%k=%v can be a max of 2 characters", key, v))
                                  }
                                  return
                                },
			},
			"color": &schema.Schema{
				Type: schema.TypeString,
				Default: "#aabbcc",
				Optional: true,
			},
		},
	}
}

func spaceFromInput(d *schema.ResourceData) Space {
	log.Printf("SpaceFromInput!")
	space := Space{}
	
	populateStruct(d, "space_id", func(value interface{}) { space.Id = value.(string) })
	populateStruct(d, "name", func(value interface{}) { space.Name = value.(string)} )
	populateStruct(d, "description", func(value interface{}) { space.Description = value.(string)} )
	populateStruct(d, "initials", func(value interface{}) { space.Initials = value.(string)} )
	populateStruct(d, "color", func(value interface{}) { space.Color = value.(string)} )
	return space
}

func resourceElasticSpaceCreate(d *schema.ResourceData, meta interface{}) error {
	url := meta.(*ElasticInfo).kibanaUrl
	username := meta.(*ElasticInfo).kibanaUsername
	password := meta.(*ElasticInfo).kibanaPassword

	space := spaceFromInput(d)
	fmt.Printf("Our Space: %v", space)

	url = fmt.Sprintf("%v/api/spaces/space", url)

	body, err := json.Marshal(&space)
	if err != nil {
		return err
	}
        log.Printf("Create: %v", body)

	_, err = postKibRequest(d, meta, url, username, password, string(body))
	if err != nil {
		return err
	}

	d.SetId(space.Id)

	return err
}

func resourceElasticSpaceRead(d *schema.ResourceData, meta interface{}) error {
	url := meta.(*ElasticInfo).kibanaUrl
	username := meta.(*ElasticInfo).kibanaUsername
	password := meta.(*ElasticInfo).kibanaPassword
	id := d.Id()

	url = fmt.Sprintf("%v/api/spaces/space/%v", url, id)
	respBody, err := getKibRequest(d, meta, url, username, password)
	if err != nil {
       	    return err
	}

	var space Space
	json.Unmarshal(*respBody, &space)
	d.Set("description", space.Description)
	d.Set("name", space.Name)
	d.Set("color", space.Color)
	d.Set("initials", space.Initials)
	return nil
}

func resourceElasticSpaceUpdate(d *schema.ResourceData, meta interface{}) error {

	url := meta.(*ElasticInfo).kibanaUrl
	username := meta.(*ElasticInfo).kibanaUsername
	password := meta.(*ElasticInfo).kibanaPassword
	id := d.Id()

	space := spaceFromInput(d)
	body, err := json.Marshal(&space)
	if err != nil {
		return err
	}

	url = fmt.Sprintf("%v/api/spaces/space/%v", url, id)
	_, err = putKibRequest(d, meta, url, username, password, string(body))
	if err != nil {
		return err
	}
	return nil
}

func resourceElasticSpaceDelete(d *schema.ResourceData, meta interface{}) error {
	url := meta.(*ElasticInfo).kibanaUrl
	username := meta.(*ElasticInfo).kibanaUsername
	password := meta.(*ElasticInfo).kibanaPassword
	id := d.Id()

	url = fmt.Sprintf("%v/api/spaces/space/%v", url, id)
	_, err := deleteKibRequest(d, meta, url, username, password)	
	if err != nil {
       		return err    
	}
	d.SetId("")
	return nil
}


