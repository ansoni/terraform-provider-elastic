package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"net/http"
	"bytes"
	"log"
	"io/ioutil"
        "io"
        "os"
        "fmt"
        _ "reflect"
)

type SavedObjectHeader struct {
        Id string `json:"id,omitempty"`
        ObjectType string `json:"type,omitempty"`
        UpdatedAt string `json:"updated_at,omitempty"`
        Version int `json:"version,omitempty"`
	Attributes map[string]interface{} `json:"attributes"`
}

type Space struct {
        Id string `json:"id,omitempty"`
        Name string `json:"name,omitempty"`
        Description string `json:"description,omitempty"`
        Color string `json:"color,omitempty"`
        Initials string `json:"initials,omitempty"`
}

func kibanaUrl(d *schema.ResourceData, meta interface{}) string {
	url := meta.(*ElasticInfo).kibanaUrl
	value, found := d.GetOk("space_id")
	if ! found {
		return url
	}
	return fmt.Sprintf("%v/s/%v", url, value)
}

func populateStruct(d *schema.ResourceData, key string, callback func(interface{})) {
	value, found := d.GetOk(key)
        if ! found {
		log.Printf("Not Found: %v in %v", key, d)
		return
	} 
	callback(value)
}

func getFileOrUrlReader(d *schema.ResourceData) (io.Reader, error) {
	var filePath string
	if v, ok := d.GetOk("file_path"); ok {
		filePath = v.(string)
	}
	if v, ok := d.GetOk("file_url"); ok {
		fileUrl := v.(string)
		resp, err := http.Get(fileUrl)
		if err != nil {
			return nil, err
                }
		//defer resp.Body.Close()
		return resp.Body, nil
        }
	
	if filePath == "" {
		return nil, fmt.Errorf("Invalid File Path")
        }
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("Error opening file_path: %s", err)
	}
	//defer file.Close()
	return file, nil
}


func getKibClient(d *schema.ResourceData, meta interface{}) *http.Client {
	return &http.Client{}
}

func getElasticsearchClient(d *schema.ResourceData, meta interface{}) *http.Client {
	return &http.Client{}
}

func streamingRequest(requestType string, d *schema.ResourceData, meta interface{}, url string, contentType string, username string, password string, buffer io.Reader) (io.Reader, error) {
	client := getElasticsearchClient(d, meta)
	var req *http.Request
	var err error
	req, err = http.NewRequest(requestType, url, buffer)
	if err != nil {
		return nil, err
	}

	if username != "" {
		req.SetBasicAuth(username, password)
	}

	req.Header.Add("kbn-xsrf", "true")
        req.Header.Add("content-type", contentType)
	log.Printf("Reauest %v", req)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	//defer resp.Body.Close() 	
	return resp.Body, nil
}

func genericKibRequest(requestType string, d *schema.ResourceData, meta interface{}, url string, username string, password string, buffer string) (*[]byte, error ) {

	client := getKibClient(d, meta)

	var req *http.Request
	var err error
	if buffer != "" {
		log.Printf("Reauest %s, %s => %s", requestType, url, buffer)
		req, err = http.NewRequest(requestType, url, bytes.NewBufferString(buffer))
	} else {
		log.Printf("Reauest %s, %s", requestType, url)
		req, err = http.NewRequest(requestType, url, bytes.NewBufferString(buffer))
		req, err = http.NewRequest(requestType, url, nil)
	}
	if err != nil {
		return nil, err
	}

	if username != "" {
		req.SetBasicAuth(username, password)
	}

	req.Header.Add("kbn-xsrf", "true")
        req.Header.Add("content-type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	
	respBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("Response: %s", respBody)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		return nil, fmt.Errorf("Request %v failed with error %v, Body: %s", url, resp.StatusCode, respBody)
	}
	
	return &respBody, nil
}

func postKibRequest(d *schema.ResourceData, meta interface{}, url string, username string, password string, buffer string) (*[]byte, error) {
	return genericKibRequest("POST", d, meta, url, username, password, buffer)
}

func getKibRequest(d *schema.ResourceData, meta interface{}, url string, username string, password string) (*[]byte, error) {
	return genericKibRequest("GET", d, meta, url, username, password, "")
}

func putKibRequest(d *schema.ResourceData, meta interface{}, url string, username string, password string, buffer string) (*[]byte, error) {
	return genericKibRequest("PUT", d, meta, url, username, password, buffer)
}

func deleteKibRequest(d *schema.ResourceData, meta interface{}, url string, username string, password string) (*[]byte, error) {
	return genericKibRequest("DELETE", d, meta, url, username, password, "")
}
