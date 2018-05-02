package main

import (
	"github.com/hashicorp/terraform/helper/schema"
        _ "github.com/hashicorp/terraform/terraform"
	"log"
	"fmt"
	_ "errors"
	"encoding/json"
)

/**
 * Example JSON Document 
{
   "attributes" : {
      "description" : "",
      "title" : "test",
      "uiStateJSON" : "{}",
      "version" : 1,
      "visState" : "{\"title\":\"test\",\"type\":\"pie\",\"params\":{\"type\":\"pie\",\"addTooltip\":true,\"addLegend\":true,\"legendPosition\":\"right\",\"isDonut\":true,\"labels\":{\"show\":false,\"values\":true,\"last_level\":true,\"truncate\":100}},\"aggs\":[{\"id\":\"1\",\"enabled\":true,\"type\":\"count\",\"schema\":\"metric\",\"params\":{\"customLabel\":\"lines\"}},{\"id\":\"2\",\"enabled\":true,\"type\":\"terms\",\"schema\":\"segment\",\"params\":{\"field\":\"play_name.keyword\",\"otherBucket\":false,\"otherBucketLabel\":\"Other\",\"missingBucket\":false,\"missingBucketLabel\":\"Missing\",\"size\":30,\"order\":\"desc\",\"orderBy\":\"1\",\"customLabel\":\"Play\"}},{\"id\":\"3\",\"enabled\":true,\"type\":\"terms\",\"schema\":\"segment\",\"params\":{\"field\":\"speaker.keyword\",\"otherBucket\":true,\"otherBucketLabel\":\"Other\",\"missingBucket\":false,\"missingBucketLabel\":\"Missing\",\"size\":15,\"order\":\"desc\",\"orderBy\":\"1\",\"customLabel\":\"Speaker\"}}]}",
      "kibanaSavedObjectMeta" : {
         "searchSourceJSON" : "{\"index\":\"1f5653a0-46bb-11e8-a995-fb06d59a0729\",\"filter\":[],\"query\":{\"query\":\"\",\"language\":\"lucene\"}}"
      }
   }
}
*/


type VisualizationHeader struct {
	Id string `json:"id,omitempty"`
        ObjectType string `json:"type,omitempty"`
        UpdatedAt string `json:"updated_at,omitempty"`
        Version int `json:"version,omitempty"`
	Attributes Visualization `json:"attributes"`
}

type Visualization struct {
	Title string `json:"title"`
	Description string `json:"description"`
	Version int `json:"version"`
	VisState string `json:"visState"`
	KibanaSavedObjectMeta map[string]string `json:"kibanaSavedObjectMeta"`
}

func resourceKibanaVisualization() *schema.Resource {
	return &schema.Resource{
		Create: resourceKibanaVisualizationCreate,
		Read: resourceKibanaVisualizationRead,
		Update: resourceKibanaVisualizationUpdate,
		Delete: resourceKibanaVisualizationDelete,
		Schema: map[string]*schema.Schema{
			"object_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"index_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
                                ForceNew: true,
			},
			"title": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"body": &schema.Schema{
				Type:             schema.TypeString,
				Required:         true,
                                ForceNew: true,
			},
		},
	}
}

func generateVisualization(d *schema.ResourceData) ([]byte, error) {
	body := d.Get("body").(string)
	// enrich
	var visualizationHeader VisualizationHeader
	json.Unmarshal([]byte(body), &visualizationHeader)
	title, found := d.GetOk("title")
	if found {
		visualizationHeader.Attributes.Title = title.(string)
	}
	description, found := d.GetOk("description")
	if found {
		visualizationHeader.Attributes.Description = description.(string)
	}
	index_id, found := d.GetOk("index_id")
	if found {
		searchSourceJSONText := visualizationHeader.Attributes.KibanaSavedObjectMeta["searchSourceJSON"]
                var searchSourceJSON SearchSourceJSON
                json.Unmarshal([]byte(searchSourceJSONText), &searchSourceJSON)
                searchSourceJSON.Index = index_id.(string)
		if len(searchSourceJSON.Filter) > 0 {
                	searchSourceJSON.Filter[0]["meta"]["index"] = index_id.(string)
		}

                searchSourceJSONBytes, _ := json.Marshal(&searchSourceJSON)
                visualizationHeader.Attributes.KibanaSavedObjectMeta["searchSourceJSON"] = string(searchSourceJSONBytes)

	}
	return json.Marshal(&visualizationHeader)
}

func resourceKibanaVisualizationCreate(d *schema.ResourceData, meta interface{}) error {
	url := meta.(*ElasticInfo).kibanaUrl

	bodyJson, err := generateVisualization(d)
	if err != nil {
		return err
	}

	url = fmt.Sprintf("%v/api/saved_objects/visualization", url)

	log.Printf("Create new Resource using %v with data %s", url, bodyJson)
	respBody, err := postKibRequest(d, meta, url, string(bodyJson))
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

func resourceKibanaVisualizationRead(d *schema.ResourceData, meta interface{}) error {
	url := meta.(*ElasticInfo).kibanaUrl
	objectId := d.Get("object_id")
	//version := d.Get("version")

	url = fmt.Sprintf("%v/api/saved_objects/visualization/%v", url, objectId)
	respBody, err := getKibRequest(d, meta, url)
	if err != nil {
       	    return err
	}
	log.Printf("Read %v => %v", url, respBody)
	var objectResponse SavedObjectHeader
	json.Unmarshal(*respBody, &objectResponse)	
/*
	if savedSearch.Version != version {
		return errors.New(fmt.Sprintf("Index Pattern has been modified.  Found version %v, Expected version %v", err, savedSearch.Version, version))
	}
*/

	return nil
}

func resourceKibanaVisualizationUpdate(d *schema.ResourceData, meta interface{}) error {
	url := meta.(*ElasticInfo).kibanaUrl
	objectId := d.Get("object_id")
	bodyJson, err := generateSavedSearch(d)
	if err != nil {
		return err
	}

	url = fmt.Sprintf("%v/api/saved_objects/visualization/%v", url, objectId)
	respBody, err := putKibRequest(d, meta, url, string(bodyJson))
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

func resourceKibanaVisualizationDelete(d *schema.ResourceData, meta interface{}) error {
	url := meta.(*ElasticInfo).kibanaUrl
	objectId := d.Get("object_id")

	d.Set("object_id", nil)
	d.Set("version", nil)

	url = fmt.Sprintf("%v/api/saved_objects/visualization/%v", url, objectId)
	_, err := deleteKibRequest(d, meta, url)
	if err != nil {
       		return err    
	}
	return nil
}


