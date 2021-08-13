package main

import (
	"encoding/json"
	"errors"
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
)

// Js2kv convert json to key value
// data: json string
// sep: separation
// return: key-value string
func Js2kv(data, sep string)(string, error){
	if !gjson.Valid(data){
		return "", errors.New("invalid json")
	}
	decoder := json.NewDecoder(strings.NewReader(data))
	decoder.UseNumber()
	jsBox := make(map[string]interface{})
	err := decoder.Decode(&jsBox)
	if err!=nil{
		return "", err
	}
	var kv = make(map[string]interface{})
	parse(jsBox, "", sep, kv)
	kvb, err := json.Marshal(kv)
	if err!=nil{
		return "", err
	}
	return string(kvb), nil
}

// recursive parse json
func parse(m interface{}, prefix, sep string, kv map[string]interface{}) {
	switch m.(type) {
	case map[string]interface{}:
		for k, v := range m.(map[string]interface{}) {
			key := prefix + k
			switch vv := v.(type) {
			case map[string]interface{}:
				parse(vv, key+sep, sep, kv)
			case []interface{}:
				for i, v := range vv {
					parse(v, key+"["+strconv.Itoa(i+1)+"]"+sep, sep, kv)
				}
			default:
				kv[key] = vv
			}
		}
	case []interface{}:
		for i, v := range m.([]interface{}) {
			parse(v, prefix+"["+strconv.Itoa(i+1)+"]"+sep, sep, kv)
		}
	default:
		if strings.HasSuffix(prefix, sep) {
			kv[prefix[:len(prefix)-len(sep)]] = m
		} else {
			kv[prefix] = m
		}
	}
}