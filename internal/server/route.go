package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type (
	Route struct {
		Path     string
		JsonFile string
		Plain    bool
	}
)

func (r Route) getRawData() ([]byte, error) {
	jsonFile, err := os.Open(r.JsonFile)
	if err != nil {
		return nil, err
	}
	defer logClose(jsonFile)
	return ioutil.ReadAll(jsonFile)
}

func (r Route) getData() ([]interface{}, error) {
	data, err := r.getRawData()
	if err != nil {
		return nil, err
	}
	var jsonData []interface{}
	err = json.Unmarshal(data, &jsonData)
	return jsonData, err
}

func (r Route) getDataItem(id string) (interface{}, error) {
	data, err := r.getData()
	if err != nil {
		return nil, err
	}

	for _, d := range data {
		if object, ok := d.(map[string]interface{}); ok {
			objectId := object["id"]
			if objectId == id || fmt.Sprint(objectId) == id {
				return object, nil
			}
		}
	}
	return nil, nil
}

func (r Route) getPage(p pageRequest) (page, error) {
	j, err := r.getData()
	if err != nil {
		return page{}, err
	}
	return paginate(j, p), nil
}
