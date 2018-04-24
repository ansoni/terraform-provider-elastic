package main

import (
	"github.com/hashicorp/terraform/helper/schema"
        _ "github.com/hashicorp/terraform/terraform"
	"io/ioutil"
	"net/http"
	"bytes"
	"log"
	"fmt"
	_ "errors"
	"encoding/json"
)

/**
 * Example JSON Document 
{
   "attributes" : {
      "title" : "Henry IV",
      "sort" : [
         "_score",
         "desc"
      ],
      "columns" : [
         "_source"
      ],
      "description" : "",
      "hits" : 0,
      "version" : 1,
      "kibanaSavedObjectMeta" : {
         "searchSourceJSON" : "{\"index\":\"1f5653a0-46bb-11e8-a995-fb06d59a0729\",\"highlightAll\":true,\"version\":true,\"query\":{\"query\":\"\",\"language\":\"lucene\"},\"filter\":[{\"meta\":{\"negate\":false,\"index\":\"1f5653a0-46bb-11e8-a995-fb06d59a0729\",\"type\":\"phrase\",\"key\":\"play_name\",\"value\":\"Henry IV\",\"params\":{\"query\":\"Henry IV\",\"type\":\"phrase\"},\"disabled\":false,\"alias\":null},\"query\":{\"match\":{\"play_name\":{\"query\":\"Henry IV\",\"type\":\"phrase\"}}},\"$state\":{\"store\":\"appState\"}}]}"
      }
   }
}
*/


// SavedSearch header
type SavedSearchHeader struct {
	Id string `json:"id,omitempty"`
        ObjectType string `json:"type,omitempty"`
        UpdatedAt string `json:"updated_at,omitempty"`
        Version int `json:"version,omitempty"`
	Attributes SavedSearch `json:"attributes"`
}

type SavedSearch struct {
	Title string `json:"title"`
	Description string `json:"description"`
	Sort []string `json:"sort"`
	Columns []string `json:"columns"`
	UpdatedAt string `json:"updated_at,omitempty"`
	Hits int `json:"hits"`
	Version int `json:"version"`
	KibanaSavedObjectMeta map[string]string `json:"kibanaSavedObjectMeta"`
}

/**
 * Example
{
  "filter": [
    {
      "query": {
        "match": {
          "play_name": {
            "query": "Henry IV",
            "type": "phrase"
          }
        }
      },
      "meta": {
        "index": "1f5653a0-46bb-11e8-a995-fb06d59a0729",
        "value": "Henry IV",
        "disabled": false,
        "alias": null,
        "params": {
          "query": "Henry IV",
          "type": "phrase"
        },
        "key": "play_name",
        "negate": false,
        "type": "phrase"
      },
      "$state": {
        "store": "appState"
      }
    }
  ],
  "index": "1f5653a0-46bb-11e8-a995-fb06d59a0729",
  "version": true,
  "highlightAll": true,
  "query": {
    "query": "",
    "language": "lucene"
  }
}
*/

type SearchSourceJSON struct {
	Index string `json:"index"`
	HighlightAll bool `json:"highlightAll"`
	Version bool `json:"version"`
	Query map[string]string `json:"query"`
	Filter []map[string]map[string]interface{} `json:"filter"`
}




func resourceKibanaSavedSearch() *schema.Resource {
	return &schema.Resource{
		Create: resourceKibanaSavedSearchCreate,
		Read:   resourceKibanaSavedSearchRead,
		Update: resourceKibanaSavedSearchUpdate,
		Delete: resourceKibanaSavedSearchDelete,
		Schema: map[string]*schema.Schema{
			"search_id": &schema.Schema{
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

func generateSavedSearch(d *schema.ResourceData) ([]byte, error) {
	body := d.Get("body").(string)
	// enrich
	var savedSearchHeader SavedSearchHeader
	json.Unmarshal([]byte(body), &savedSearchHeader)
	title, found := d.GetOk("title")
	if found {
		savedSearchHeader.Attributes.Title = title.(string)
	}
	description, found := d.GetOk("description")
	if found {
		savedSearchHeader.Attributes.Description = description.(string)
	}
	index_id, found := d.GetOk("index_id")
	if found {
		searchSourceJSONText := savedSearchHeader.Attributes.KibanaSavedObjectMeta["searchSourceJSON"]
		var searchSourceJSON SearchSourceJSON
		json.Unmarshal([]byte(searchSourceJSONText), &searchSourceJSON)
		searchSourceJSON.Index = index_id.(string)
		searchSourceJSON.Filter[0]["meta"]["index"] = index_id.(string)	

		searchSourceJSONBytes, _ := json.Marshal(&searchSourceJSON)
		savedSearchHeader.Attributes.KibanaSavedObjectMeta["searchSourceJSON"] = string(searchSourceJSONBytes)
	}

	return json.Marshal(&savedSearchHeader)

}

func resourceKibanaSavedSearchCreate(d *schema.ResourceData, meta interface{}) error {
	url := meta.(*ElasticInfo).kibanaUrl

	bodyJson, err := generateSavedSearch(d)
	if err != nil {
		return err
	}

	url = fmt.Sprintf("%v/api/saved_objects/search", url)

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJson))
	if err != nil {
		return err
	}
	req.Header.Add("kbn-xsrf", "true")
	req.Header.Add("content-type", "application/json")

	log.Printf("POST new search to %v", url)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var savedSearchResponse SavedSearchHeader
	json.Unmarshal(respBody, &savedSearchResponse)	

	log.Printf("Raw Body: %s", respBody)
	log.Printf("ID: %s", savedSearchResponse.Id)
	log.Printf("ObjectType: %s", savedSearchResponse.ObjectType)
	log.Printf("UpdatedAt: %s", savedSearchResponse.UpdatedAt)
	log.Printf("Version: %v", savedSearchResponse.Version)

	d.SetId(savedSearchResponse.Id)
	d.Set("search_id", savedSearchResponse.Id)
	d.Set("version", savedSearchResponse.Version)

	return err
}

func resourceKibanaSavedSearchRead(d *schema.ResourceData, meta interface{}) error {
	url := meta.(*ElasticInfo).kibanaUrl
	objectId := d.Get("object_id")
	//version := d.Get("version")

	url = fmt.Sprintf("%v/api/saved_objects/search/%v", url, objectId)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("kbn-xsrf", "true")
	req.Header.Add("content-type", "application/json")
	log.Printf("Get existing saved-search with %v", url)
	resp, err := client.Do(req)
	if err != nil {
       	    return err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
       	    return err
	}
	log.Printf("Read %v => %v", url, respBody)
	var savedSearch SavedSearch
	json.Unmarshal(respBody, &savedSearch)	
/*
	if savedSearch.Version != version {
		return errors.New(fmt.Sprintf("Index Pattern has been modified.  Found version %v, Expected version %v", err, savedSearch.Version, version))
	}
*/

	return nil
}

func resourceKibanaSavedSearchUpdate(d *schema.ResourceData, meta interface{}) error {
	url := meta.(*ElasticInfo).kibanaUrl
	objectId := d.Get("object_id")
	bodyJson, err := generateSavedSearch(d)
	if err != nil {
		return err
	}

	url = fmt.Sprintf("%v/api/saved_objects/search/%v", url, objectId)
	client := &http.Client{}
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(bodyJson))
	if err != nil {
		return err
	}
	req.Header.Add("kbn-xsrf", "true")
	req.Header.Add("content-type", "application/json")
	log.Printf("Update saved-search to %v", url)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var savedSearchHeader SavedSearchHeader
	json.Unmarshal(respBody, &savedSearchHeader)	
	log.Printf("Raw Body: %s", respBody)
	log.Printf("ID: %s", savedSearchHeader.Id)
	log.Printf("ObjectType: %s", savedSearchHeader.ObjectType)
	log.Printf("UpdatedAt: %s", savedSearchHeader.UpdatedAt)
	log.Printf("Version: %v", savedSearchHeader.Version)
	d.Set("object_id", savedSearchHeader.Id)
	d.Set("version", savedSearchHeader.Version)
	return nil
}

func resourceKibanaSavedSearchDelete(d *schema.ResourceData, meta interface{}) error {
	url := meta.(*ElasticInfo).kibanaUrl
	objectId := d.Get("object_id")

	d.Set("object_id", nil)
	d.Set("version", nil)

	url = fmt.Sprintf("%v/api/saved_objects/search/%v", url, objectId)
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("kbn-xsrf", "true")
	req.Header.Add("content-type", "application/json")
	log.Printf("DELETE saved-search to %v", url)
	_, err = client.Do(req)
	if err != nil {
       		return err    
	}
	return nil
}


