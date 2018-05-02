package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"net/http"
	"bytes"
	"io/ioutil"
)

type SavedObjectHeader struct {
        Id string `json:"id,omitempty"`
        ObjectType string `json:"type,omitempty"`
        UpdatedAt string `json:"updated_at,omitempty"`
        Version int `json:"version,omitempty"`
	Attributes interface{} `json:"attributes"`
}




func getKibClient(d *schema.ResourceData, meta interface{}) *http.Client {
	return &http.Client{}
}

func genericKibRequest(requestType string, d *schema.ResourceData, meta interface{}, url string, buffer string) (*[]byte, error ) {

	client := getKibClient(d, meta)

	var req *http.Request
	var err error
	if buffer != "" {
		req, err = http.NewRequest(requestType, url, bytes.NewBufferString(buffer))
	} else {
		req, err = http.NewRequest(requestType, url, nil)
	}
	if err != nil {
		return nil, err
	}

	req.Header.Add("kbn-xsrf", "true")
        req.Header.Add("content-type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	return &respBody, nil
}

func postKibRequest(d *schema.ResourceData, meta interface{}, url string, buffer string) (*[]byte, error) {
	return genericKibRequest("POST", d, meta, url, buffer)
}

func getKibRequest(d *schema.ResourceData, meta interface{}, url string) (*[]byte, error) {
	return genericKibRequest("GET", d, meta, url, "")
}

func putKibRequest(d *schema.ResourceData, meta interface{}, url string, buffer string) (*[]byte, error) {
	return genericKibRequest("PUT", d, meta, url, buffer)
}

func deleteKibRequest(d *schema.ResourceData, meta interface{}, url string) (*[]byte, error) {
	return genericKibRequest("DELETE", d, meta, url, "")
}
