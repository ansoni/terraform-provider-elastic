package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"net/http"
	"bytes"
	"io/ioutil"
)

type SavedObjectHeader struct {
        Id string `json:"id"`
        ObjectType string `json:"type"`
        UpdatedAt string `json:"updated_at"`
        Version int `json:"version"`
	Attributes interface{} `json:"attributes"`
}




func getClient(d *schema.ResourceData, meta interface{}) *http.Client {
	return &http.Client{}
}

func genericRequest(requestType string, d *schema.ResourceData, meta interface{}, url string, buffer string) (*[]byte, error ) {

	client := getClient(d, meta)

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

	defer resp.Body.Close()
	
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	return &respBody, nil
}

func postRequest(d *schema.ResourceData, meta interface{}, url string, buffer string) (*[]byte, error) {
	return genericRequest("POST", d, meta, url, buffer)
}

func getRequest(d *schema.ResourceData, meta interface{}, url string) (*[]byte, error) {
	return genericRequest("GET", d, meta, url, "")
}

func putRequest(d *schema.ResourceData, meta interface{}, url string, buffer string) (*[]byte, error) {
	return genericRequest("PUT", d, meta, url, buffer)
}

func deleteRequest(d *schema.ResourceData, meta interface{}, url string) (*[]byte, error) {
	return genericRequest("DELETE", d, meta, url, "")
}
