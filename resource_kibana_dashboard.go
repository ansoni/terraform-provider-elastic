package main

import (
	"github.com/hashicorp/terraform/helper/schema"
        _ "github.com/hashicorp/terraform/terraform"
	"log"
	"fmt"
	_ "errors"
	"encoding/json"
)

type DashboardHeader struct {
	Id string `json:"id,omitempty"`
        ObjectType string `json:"type,omitempty"`
        UpdatedAt string `json:"updated_at,omitempty"`
        Version int `json:"version,omitempty"`
	Attributes Dashboard `json:"attributes"`
}

type Dashboard struct {
	Title string `json:"title"`
	Description string `json:"description"`
	Hits int `json:"hits"`
	TimeRestore bool `json:"timeRestore"`
	Version int `json:"version"`
	KibanaSavedObjectMeta map[string]string `json:"kibanaSavedObjectMeta"`
	OptionsJSON string `json:"optionsJSON"`
	PanelsJSON string `json:"panelsJSON"`
}

type Panels struct {
	PanelIndex string `json:"panelIndex"`
	GridData GridData `json:"gridData"`
	Type string `json:"type"`
	Version string `json:"version"`
	Id string `json:"id"`
}

type GridData struct {
	X int `json:"x"`  
	Y int `json:"y"`  
	W int `json:"w"`  
	H int `json:"h"`  
	I string `json:"i"`  
}

func resourceKibanaDashboard() *schema.Resource {
	return &schema.Resource{
		Create: resourceKibanaDashboardCreate,
		Read: resourceKibanaDashboardRead,
		Update: resourceKibanaDashboardUpdate,
		Delete: resourceKibanaDashboardDelete,
		Schema: map[string]*schema.Schema{
			"object_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"title": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"panels": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"panel_index": {
							Type:     schema.TypeString,
							Required: true,
						},
						"visualization_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"body": &schema.Schema{
				Type:             schema.TypeString,
				Required:         true,
                                ForceNew: true,
			},
		},
	}
}

func generateDashboard(d *schema.ResourceData) ([]byte, error) {
	body := d.Get("body").(string)
	// enrich
	var dashboardHeader DashboardHeader
	json.Unmarshal([]byte(body), &dashboardHeader)
	title, found := d.GetOk("title")
	if found {
		dashboardHeader.Attributes.Title = title.(string)
	}
	description, found := d.GetOk("description")
	if found {
		dashboardHeader.Attributes.Description = description.(string)
	}
	panelsConfigTmp, found := d.GetOk("panels")
	
	if found {
		panelsConfig := panelsConfigTmp.([]interface{})
		log.Printf("Panels Config %v ", panelsConfig)

		panelsJSONText := dashboardHeader.Attributes.PanelsJSON
                var panels []*Panels
                json.Unmarshal([]byte(panelsJSONText), &panels)

		for _, panel := range panels {
			for _, panelConfigTmp := range panelsConfig {
				panelConfig := panelConfigTmp.(map[string]interface{})
				if panel.PanelIndex == panelConfig["panel_index"].(string) {
					panel.Id = panelConfig["visualization_id"].(string)
				}
			}		
    		}		

		log.Printf("Panels %v ", panels)


                panelsJSONBytes, _ := json.Marshal(&panels)
                dashboardHeader.Attributes.PanelsJSON = string(panelsJSONBytes)

	}
	return json.Marshal(&dashboardHeader)
}

func resourceKibanaDashboardCreate(d *schema.ResourceData, meta interface{}) error {
	url := meta.(*ElasticInfo).kibanaUrl

	bodyJson, err := generateDashboard(d)
	if err != nil {
		return err
	}

	url = fmt.Sprintf("%v/api/saved_objects/dashboard", url)

	log.Printf("Create new Resource using %v with data %s", url, bodyJson)
	respBody, err := postRequest(d, meta, url, string(bodyJson))
	if err != nil {
		return err
	}

	var newObjectResponse SavedObjectHeader
	json.Unmarshal(*respBody, &newObjectResponse)	

	log.Printf("Raw Body: %s", respBody)
	log.Printf("ID: %s", newObjectResponse.Id)
	log.Printf("ObjectType: %s", newObjectResponse.ObjectType)
	log.Printf("UpdatedAt: %s", newObjectResponse.UpdatedAt)
	log.Printf("Version: %v", newObjectResponse.Version)

	d.SetId(newObjectResponse.Id)
	d.Set("object_id", newObjectResponse.Id)
	d.Set("version", newObjectResponse.Version)

	return err
}

func resourceKibanaDashboardRead(d *schema.ResourceData, meta interface{}) error {
	url := meta.(*ElasticInfo).kibanaUrl
	objectId := d.Get("object_id")
	//version := d.Get("version")

	url = fmt.Sprintf("%v/api/saved_objects/dashboard/%v", url, objectId)
	respBody, err := getRequest(d, meta, url)
	if err != nil {
       	    return err
	}
	log.Printf("Read %v => %v", url, respBody)
	var objectResponse SavedObjectHeader
	json.Unmarshal(*respBody, &objectResponse)	

	return nil
}

func resourceKibanaDashboardUpdate(d *schema.ResourceData, meta interface{}) error {
	url := meta.(*ElasticInfo).kibanaUrl
	objectId := d.Get("object_id")
	bodyJson, err := generateSavedSearch(d)
	if err != nil {
		return err
	}

	url = fmt.Sprintf("%v/api/saved_objects/dashboard/%v", url, objectId)
	respBody, err := putRequest(d, meta, url, string(bodyJson))
	if err != nil {
		return err
	}

	var newObjectResponse SavedObjectHeader
	json.Unmarshal(*respBody, newObjectResponse)	

	log.Printf("Raw Body: %s", respBody)
	log.Printf("ID: %s", newObjectResponse.Id)
	log.Printf("ObjectType: %s", newObjectResponse.ObjectType)
	log.Printf("UpdatedAt: %s", newObjectResponse.UpdatedAt)
	log.Printf("Version: %v", newObjectResponse.Version)

	d.SetId(newObjectResponse.Id)
	d.Set("object_id", newObjectResponse.Id)
	d.Set("version", newObjectResponse.Version)

	return err
}

func resourceKibanaDashboardDelete(d *schema.ResourceData, meta interface{}) error {
	url := meta.(*ElasticInfo).kibanaUrl
	objectId := d.Get("object_id")

	d.Set("object_id", nil)
	d.Set("version", nil)

	url = fmt.Sprintf("%v/api/saved_objects/dashboard/%v", url, objectId)
	_, err := deleteRequest(d, meta, url)
	if err != nil {
       		return err    
	}
	return nil
}


