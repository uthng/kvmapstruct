// Package kvstruct exposes various utility functions to do conversions between:
// Consul KV pairs and native Go Struct or map[string]interface{}.
//
// It also provides several utilities to convert directly:
// nested map to flatten/kv map or Consul kv pairs, flatten/kv map to Go struct,
// Kv map to nested map, etc.
//
// There are some notions that are used in this package.
// Nested map: classic map[string]interface{}.
// Flatten map: map[string]interface{} represents key/value and value can be a normal type including slice or map
// KV map: map[string]interface{} represents key/value but value can not be slice or map.
// A slice will be represented by keys suffixed by 0, 1, 2 etc.
//
// Only the following value types are supported:
// int, bool, string, []int, []bool, []string and map[string]interface{}
package kvstruct

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/fatih/structs"
	consul "github.com/hashicorp/consul/api"
	//"github.com/mitchellh/copystructure"
	//"github.com/mitchellh/mapstructure"
	//"github.com/uthng/common/utils"

	"github.com/spf13/cast"
)

// KVStruct contains consul informations.
type KVStruct struct {
	// Path is consul key parent to store struct's fields
	Path string
	// Client is consul client
	Client *consul.Client
}

// NewKVStruct creates a new *KVStruct.
// Format of URL is ip:port.
func NewKVStruct(url, token, path string) (*KVStruct, error) {
	ks := &KVStruct{}

	ks.Path = path

	// Initialize consul config
	config := consul.DefaultConfig()

	if url != "" {
		config.Address = url
	}

	if token != "" {
		config.Token = token
	}

	// Initialize consul client
	client, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}

	ks.Client = client
	return ks, nil
}

// StructToConsulKV converts and saves the struct to Consul KV store
// Input may be a Go struct.
func (ks *KVStruct) StructToConsulKV(input interface{}) error {
	m := make(map[string]interface{})
	v := reflect.ValueOf(input)
	k := v.Kind()

	if k != reflect.Struct {
		return fmt.Errorf("Error: input is not a Go struct")
	}

	// Convert it to Map
	m = structs.Map(input)

	// Mapping to kvpairs
	pairs, err := ks.MapToKVPairs(m, ks.Path)
	if err != nil {
		return err
	}

	for _, kv := range pairs {
		_, err := ks.Client.KV().Put(kv, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

// MapToConsulKV converts and saves the map to Consul KV store.
// Input may be a map[string]interface{}.
func (ks *KVStruct) MapToConsulKV(input interface{}) error {
	v := reflect.ValueOf(input)
	k := v.Kind()

	if k != reflect.Map {
		return fmt.Errorf("Error: input is not a map[string]interface{}")
	}

	m := input.(map[string]interface{})

	// Mapping to kvpairs
	pairs, err := ks.MapToKVPairs(m, ks.Path)
	if err != nil {
		return err
	}

	for _, kv := range pairs {
		_, err := ks.Client.KV().Put(kv, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

// ConsulKVToStruct gets list of all consul keys from kvstruct path
// and match them to the given struct in argument.
// The Go struct and its substruct must be a pointer and initialized.
func (ks *KVStruct) ConsulKVToStruct(out interface{}) error {
	m := make(map[string]interface{})

	pairs, _, err := ks.Client.KV().List(ks.Path, nil)
	if err != nil {
		return err
	}

	// Build map
	for _, kv := range pairs {
		m[kv.Key] = string(kv.Value)
	}

	err = KVMapToStruct(m, ks.Path, out)

	return err
}

// ConsulKVToMap gets list of all consul keys from kvstruct path
// and match them to a map[string]interface{}.
func (ks *KVStruct) ConsulKVToMap() (map[string]interface{}, error) {
	m := make(map[string]interface{})
	out := make(map[string]interface{})

	pairs, _, err := ks.Client.KV().List(ks.Path, nil)
	if err != nil {
		return nil, err
	}

	// Build map
	for _, kv := range pairs {
		m[kv.Key] = string(kv.Value)
	}

	out, err = KVMapToMap(m, ks.Path)

	return out, err
}

// MapToKVPairs convert a nested map to an array of Consul KV pairs
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

// MapToKVMap convert a nested map to a KV map.
func MapToKVMap(in map[string]interface{}, prefix string) map[string]interface{} {
	out := make(map[string]interface{})
	key := ""

	if prefix != "" {
		key = prefix + "/"
	}

	// Loop map to build
	for k, v := range in {
		kind := reflect.ValueOf(v).Kind()
		if kind == reflect.Map {
			o := MapToKVMap(v.(map[string]interface{}), key+k)
			for k1, v1 := range o {
				out[k1] = v1
			}
		} else if kind == reflect.Slice {
			// TODO: Maybe there is another way to do this more elegant
			switch v.(type) {
			case []int:
				for i, e := range v.([]int) {
					out[key+k+"/"+cast.ToString(i)] = e
				}
			case []string:
				for i, e := range v.([]string) {
					out[key+k+"/"+cast.ToString(i)] = e
				}
			case []bool:
				for i, e := range v.([]bool) {
					out[key+k+"/"+cast.ToString(i)] = e
				}
			}
		} else {
			out[key+k] = v
		}
	}

	return out
}

// MapToFlattenMap converts a nested map to a flatten map.
func MapToFlattenMap(in map[string]interface{}, prefix string) map[string]interface{} {
	out := make(map[string]interface{})
	key := ""

	if prefix != "" {
		key = prefix + "/"
	}

	// Loop map to build
	for k, v := range in {
		kind := reflect.ValueOf(v).Kind()
		if kind == reflect.Map {
			o := MapToFlattenMap(v.(map[string]interface{}), key+k)
			if len(o) <= 0 {
				out[key+k] = o
			} else {
				for k1, v1 := range o {
					out[k1] = v1
				}
			}
		} else if kind == reflect.Slice {
			// TODO: Maybe there is another way to do this more elegant
			switch v.(type) {
			case []int:
				out[key+k] = v.([]int)
			case []string:
				out[key+k] = v.([]string)
			case []bool:
				out[key+k] = v.([]bool)
			}
		} else {
			out[key+k] = v
		}
	}

	return out
}

// KVMapToStruct converts a KV map to a Go struct.
// Go struct and its substruct must be a pointer and initialized.
func KVMapToStruct(in map[string]interface{}, prefix string, out interface{}) error {
	var inVal interface{}

	key := ""

	if prefix != "" {
		key = prefix + "/"
	}

	v := reflect.ValueOf(out)
	// Get value of pointer
	indirect := reflect.Indirect(v)
	k := v.Kind()

	// Check if out is a pointer to a structure
	if k != reflect.Ptr || indirect.Kind() != reflect.Struct {
		return fmt.Errorf("Error of output's type! Only pointer of struct are supported")
	}

	// If struct, convert it to Map
	flattenOut := structs.Map(out)

	for k, v := range flattenOut {
		val := reflect.ValueOf(v)
		kind := val.Kind()

		// Initialize inVal
		inVal = nil

		// If value is not a slice or a map, we assign value directly
		if kind == reflect.Slice {
			i := 0
			arr := []string{}
			// Loop by incremnenting i to get all values of slice
			for {
				v1, ok := in[key+k+"/"+cast.ToString(i)]
				if !ok {
					break
				}
				arr = append(arr, cast.ToString(v1))
				i = i + 1
			}

			inVal = arr
		} else if kind == reflect.Map {
			// Convert kv map to a nested map
			m, err := KVMapToMap(in, key+k)
			if err != nil {
				return err
			}

			inVal = m

		} else {
			// Check in kvmap
			inVal = in[key+k]
		}

		// Assign value following its type
		if inVal != nil {
			switch v.(type) {
			case int:
				flattenOut[k] = cast.ToInt(inVal)
			case bool:
				flattenOut[k] = cast.ToBool(inVal)
			case string:
				flattenOut[k] = cast.ToString(inVal)
			case []int:
				flattenOut[k] = cast.ToIntSlice(inVal)
			case map[string]interface{}:
				flattenOut[k] = inVal
			default:
				return fmt.Errorf("error type at key %s", k)
			}
		}
	}

	// Convert struct's flatten map to struct
	err := FlattenMapToStruct(flattenOut, out)
	//fmt.Println(err, out)
	return err
}

// KVMapToMap converts a KV map to nested map.
func KVMapToMap(in map[string]interface{}, prefix string) (map[string]interface{}, error) {
	var keys []string
	out := make(map[string]interface{})
	key := ""
	parent := ""
	count := 1
	slice := false

	// Create a slice only containing in map's keys to be sorted
	for k := range in {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	// Loop sorted map
	for _, k := range keys {
		key = k

		// If prefix is set, check if key contains prefix and remove it
		if prefix != "" {
			// Only handle key containing prefix, zap others
			re := regexp.MustCompile(prefix + "/.*")
			if re.MatchString(key) {
				key = strings.Replace(key, prefix+"/", "", 1)
			} else {
				continue
			}
		}

		// Assign current out
		outchilds := out

		// Split / to create submap
		childs := strings.Split(key, "/")
		if len(childs) > 0 {
			// Get the last key that will be assigned a value
			key = childs[len(childs)-1]

			// Check if key is an elem of slice(0, 1, 3 etc.)
			pos, err := strconv.Atoi(key)
			if err == nil && pos == count {
				// Get parent of key ==> slice field
				parent = childs[len(childs)-2]
				slice = true
				count = count + 1
			} else {
				// Reinitialize variables for slice case
				slice = false
				count = 0
				parent = ""
			}

			// In case of slice, remove key + its parents (slice itself)
			// Otherwise, remove only the last key
			if slice {
				childs = childs[:len(childs)-2]
			} else {
				childs = childs[:len(childs)-1]
			}

			for _, child := range childs {
				// Check if key exists already. If not, create a new map
				if outchilds[child] == nil {
					outchilds[child] = make(map[string]interface{})
				}

				// Get the child if it exists. If not return an error
				subchild, ok := outchilds[child].(map[string]interface{})
				if !ok {
					return nil, fmt.Errorf("child is both a data item and dir: %s", child)
				}

				// Assign subchild to outchilds to do recursively
				outchilds = subchild
			}

			// Assign value to the last key
			// In case of slice, if 1st elem, check type of slice elem value
			// to initialize slice with the same type. Otherwise add simply elem to slice
			if slice {
				val := in[k]
				switch val.(type) {
				case int:
					if pos == 0 {
						outchilds[parent] = []int{}
					}
					outchilds[parent] = append(outchilds[parent].([]int), val.(int))
				case string:
					if pos == 0 {
						outchilds[parent] = []string{}
					}
					outchilds[parent] = append(outchilds[parent].([]string), val.(string))
				case bool:
					if pos == 0 {
						outchilds[parent] = []bool{}
					}
					outchilds[parent] = append(outchilds[parent].([]bool), val.(bool))
				default:
					return out, fmt.Errorf("Type error! Only int, string, bool are supported")
				}
			} else {
				// Other cases, assign simply value to key
				outchilds[key] = cast.ToString(in[k])
			}
		}
	}

	return out, nil
}

// FlattenMapToStruct converts a flatten map to a Go struct.
// Go struct and its substruct must be a pointer and initialized.
func FlattenMapToStruct(in map[string]interface{}, out interface{}) error {
	if out == nil {
		return fmt.Errorf("go struct is not initialized")
	}

	v := reflect.ValueOf(out)
	// Get value of pointer
	indirect := reflect.Indirect(v)
	k := v.Kind()

	// Check if out is a pointer to a structure
	if k != reflect.Ptr || indirect.Kind() != reflect.Struct {
		return fmt.Errorf("Error of output's type! Only pointer of struct are supported")
	}

	for i := 0; i < indirect.Type().NumField(); i++ {
		field := indirect.Type().Field(i)

		k := field.Type.Kind()
		v := indirect.FieldByName(field.Name)

		if k == reflect.Ptr && reflect.Indirect(v).Kind() == reflect.Struct {
			err := FlattenMapToStruct(cast.ToStringMap(in[field.Name]), v.Interface())
			if err != nil {
				return err
			}
		} else {
			switch t := v.Interface().(type) {
			case string:
				v.SetString(cast.ToString(in[field.Name]))
			case int:
				v.SetInt(cast.ToInt64(in[field.Name]))
			case bool:
				v.SetBool(cast.ToBool(in[field.Name]))
			case []int:
				v.Set(reflect.ValueOf(cast.ToIntSlice(in[field.Name])))
			case []string:
				v.Set(reflect.ValueOf(cast.ToStringSlice(in[field.Name])))
			case []bool:
				v.Set(reflect.ValueOf(cast.ToBoolSlice(in[field.Name])))
			case map[string]interface{}:
				v.Set(reflect.ValueOf(cast.ToStringMap(in[field.Name])))
			default:
				return fmt.Errorf("type error not supported %s", t)
			}
		}
	}

	return nil
}

//////////////////////// PRIVATE FUNCTIONS ///////////////////////
