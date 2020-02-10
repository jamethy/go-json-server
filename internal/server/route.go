package server

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

type (
	Route struct {
		Path     string
		JsonFile string
		Raw      bool
		IdField  string
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
			objectId := object[r.IdField]
			if objectId == id || fmt.Sprint(objectId) == id {
				return object, nil
			}
		}
	}
	return nil, nil
}

func (r Route) saveData(data []interface{}) error {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(r.JsonFile, b, 0)
}

func (r Route) addDataItem(item map[string]interface{}) (map[string]interface{}, error) {
	data, err := r.getData()
	if err != nil {
		return nil, err
	}

	id := r.getId(item)

	if id == nil {
		maxId := -1
		for _, d := range data {
			if i, err := strconv.Atoi(r.getIdStr(d)); err == nil {
				if i > maxId {
					maxId = i
				}
			}
		}
		if maxId != -1 {
			id = maxId + 1
		} else {
			id = base64.StdEncoding.EncodeToString([]byte(time.Now().String()))
		}
		item[r.IdField] = id
	}

	data = append(data, item)
	return item, r.saveData(data)
}

func (r Route) putDataItem(item map[string]interface{}) (map[string]interface{}, error) {
	return r.updateDataItem(item, false)
}

func (r Route) patchDataItem(item map[string]interface{}) (map[string]interface{}, error) {
	return r.updateDataItem(item, true)
}

func (r Route) updateDataItem(item map[string]interface{}, patch bool) (map[string]interface{}, error) {
	data, err := r.getData()
	if err != nil {
		return nil, err
	}

	id := r.getIdStr(item)
	if id == "" {
		return nil, errors.New("no id in item")
	}

	for i, other := range data {
		otherId := r.getIdStr(other)
		if otherId != "" && otherId == id {
			if patch {
				item = patchItem(item, other)
			}
			data[i] = item
			return item, r.saveData(data)
		}
	}

	return nil, errors.New("not found")
}

func patchItem(item map[string]interface{}, original interface{}) map[string]interface{} {
	if obj, isObj := original.(map[string]interface{}); isObj {
		for key, value := range item {
			obj[key] = value
		}
		item = obj
	}
	return item
}

func (r Route) getPage(p pageRequest) (page, error) {
	j, err := r.getData()
	if err != nil {
		return page{}, err
	}
	return paginate(j, p), nil
}

func (r Route) getId(obj interface{}) interface{} {
	idField := "id"
	if r.IdField != "" {
		idField = r.IdField
	}

	if objMap, isMap := obj.(map[string]interface{}); isMap {
		if id, hasId := objMap[idField]; hasId {
			return id
		}
	}
	return nil
}

func (r Route) getIdStr(obj interface{}) string {
	id := r.getId(obj)
	if id == nil {
		return ""
	} else {
		return fmt.Sprint(id)
	}
}
