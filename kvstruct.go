package kvstruct

import (
	//"encoding/base64"
	"fmt"
	"reflect"

	"github.com/fatih/structs"
	consul "github.com/hashicorp/consul/api"
	//"github.com/mitchellh/copystructure"
	//"github.com/mitchellh/mapstructure"
	//"github.com/uthng/common/utils"

	"github.com/spf13/cast"
)

// KVStruct contains consul informations
type KVStruct struct {
	ConsulPath string
	//ConsulConfig *consul.Config
}

// NewKVStruct creates a new KVStruct
func NewKVStruct(url, token, path string) *KVStruct {
	ks := &KVStruct{}

	ks.ConsulPath = path

	// Initialize consul config
	return ks
}

// StructToConsulKV converts and saves the struct/map to Consul
// input may be a struct or a map
func (ks *KVStruct) StructToConsulKV(input interface{}) error {
	m := make(map[string]interface{})
	v := reflect.ValueOf(input)
	k := v.Kind()

	if k != reflect.Map && k != reflect.Struct {
		return fmt.Errorf("Error of input's type! Only map or struct are supported")
	}

	// If struct, convert it to Map
	if k == reflect.Struct {
		m = structs.Map(input)
	}

	if k == reflect.Map {
		m = input.(map[string]interface{})
	}

	// Loop map to build
	// Attention: 25 max for consul transaction => split
	ks.MapToKVPairs(m, ks.ConsulPath)
	return nil
}

// MapToKVPairs convert a map to a flatten kv pairs
func (ks *KVStruct) MapToKVPairs(in map[string]interface{}, prefix string) (consul.KVPairs, error) {
	var out consul.KVPairs

	// Convert to flatten map
	m := MapToKVMap(in, prefix)

	for k, v := range m {
		kv := &consul.KVPair{
			Key:   k,
			Value: []byte(cast.ToString(v)),
		}

		out = append(out, kv)
	}

	return out, nil
}

// MapToKVMap convert a map to a flatten kv list
func MapToKVMap(in map[string]interface{}, prefix string) map[string]interface{} {
	out := make(map[string]interface{})
	key := ""

	if prefix != "" {
		key = prefix + "/"
	}

	// Loop map to build
	for k, v := range in {
		kind := reflect.ValueOf(v).Kind()
		if kind != reflect.Map {
			out[key+k] = v
		} else {
			o := kmToKVMap(key+k, v.(map[string]interface{}))
			for k1, v1 := range o {
				out[k1] = v1
			}
		}
	}

	return out
}

//////////////////////// PRIVATE FUNCTIONS ///////////////////////

// kmToKVMap loops & converts recursively nested map to a flatten map
// and preprend "key" to each key
func kmToKVMap(key string, m map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{})

	for k, v := range m {
		newkey := key + "/" + k
		kind := reflect.ValueOf(v).Kind()
		if kind != reflect.Map {
			out[newkey] = v
		} else {
			o := kmToKVMap(newkey, v.(map[string]interface{}))
			for k1, v1 := range o {
				out[k1] = v1
			}
		}
	}

	return out
}
