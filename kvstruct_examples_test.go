package kvstruct

import (
	"fmt"
	"reflect"
	"sort"

	consul "github.com/hashicorp/consul/api"
)

func ExampleMapToKVMap() {
	input := map[string]interface{}{
		"key1": "val1",
		"key2": 2,
		"key3": []int{1, 2, 3},
		"key4": map[string]interface{}{
			"key41": "val41",
			"key42": map[string]interface{}{
				"key421": "val421",
				"key422": []string{"one", "two", "three"},
			},
		},
	}

	output := map[string]interface{}{
		"test/key1":                "val1",
		"test/key2":                2,
		"test/key3/0":              1,
		"test/key3/1":              2,
		"test/key3/2":              3,
		"test/key4/key41":          "val41",
		"test/key4/key42/key421":   "val421",
		"test/key4/key42/key422/0": "one",
		"test/key4/key42/key422/1": "two",
		"test/key4/key42/key422/2": "three",
	}

	o := MapToKVMap(input, "test")

	// Compare result with expected output
	fmt.Println(reflect.DeepEqual(o, output))
	// Output:
	// true
}

func ExampleKVStruct_MapToKVPairs() {
	input := map[string]interface{}{
		"key1": "val1",
		"key2": 2,
		"key3": []int{1, 2, 3},
		"key4": map[string]interface{}{
			"key41": "val41",
			"key42": map[string]interface{}{
				"key421": "val421",
				"key422": []string{"one", "two", "three"},
			},
		},
	}

	ks, err := NewKVStruct("", "", "test")
	if err != nil {
		fmt.Println(err)
		return
	}

	o, err := ks.MapToKVPairs(input, ks.Path)
	if err != nil {
		fmt.Println(err)
		return
	}

	keys := []string{}
	out := make(map[string]interface{})

	for _, kv := range o {
		out[kv.Key] = string(kv.Value)
		keys = append(keys, kv.Key)
	}

	sort.Strings(keys)

	// Compare result with expected output
	for _, key := range keys {
		fmt.Println(key, ":", out[key])
	}
	// Output:
	//test/key1 : val1
	//test/key2 : 2
	//test/key3/0 : 1
	//test/key3/1 : 2
	//test/key3/2 : 3
	//test/key4/key41 : val41
	//test/key4/key42/key421 : val421
	//test/key4/key42/key422/0 : one
	//test/key4/key42/key422/1 : two
	//test/key4/key42/key422/2 : three
}

func ExampleKVStruct_MapToConsulKV() {

	input := map[string]interface{}{
		"key1": "val1",
		"key2": 2,
		"key3": []int{1, 2, 3},
		"key4": map[string]interface{}{
			"key41": "val41",
			"key42": map[string]interface{}{
				"key421": "val421",
				"key422": []string{"one", "two", "three"},
			},
		},
	}

	ks, err := NewKVStruct("localhost:8500", "adf4238a-882b-9ddc-4a9d-5b6758e4159e", "test")
	if err != nil {
		return
	}

	out := make(map[string]interface{})
	keys := []string{}

	ks.Path = "nestedmap"
	err = ks.MapToConsulKV(input)
	if err != nil {
		return
	}

	pairs, _, err := ks.Client.KV().List(ks.Path, nil)
	if err != nil {
		return
	}

	for _, kv := range pairs {
		out[kv.Key] = string(kv.Value)
		keys = append(keys, kv.Key)
	}

	sort.Strings(keys)

	// Compare result with expected output
	for _, key := range keys {
		fmt.Println(key, ":", out[key])
	}
	// Output:
	//nestedmap/key1 : val1
	//nestedmap/key2 : 2
	//nestedmap/key3/0 : 1
	//nestedmap/key3/1 : 2
	//nestedmap/key3/2 : 3
	//nestedmap/key4/key41 : val41
	//nestedmap/key4/key42/key421 : val421
	//nestedmap/key4/key42/key422/0 : one
	//nestedmap/key4/key42/key422/1 : two
	//nestedmap/key4/key42/key422/2 : three

	ks.Client.KV().DeleteTree(ks.Path, nil)
}

func ExampleKVStruct_StructToConsulKV() {
	type ExSTChildLevel2 struct {
		Key431 map[string]interface{}
	}

	type ExSTChildLevel1 struct {
		Key41 string
		Key42 map[string]interface{}
		Key43 *ExSTChildLevel2
	}

	type ExST struct {
		Key1 string
		Key2 int
		Key3 []int
		Key4 *ExSTChildLevel1
	}

	input := ExST{
		Key1: "val1",
		Key2: 2,
		Key3: []int{1, 2, 3},
		Key4: &ExSTChildLevel1{
			Key41: "val41",
			Key42: map[string]interface{}{
				"Key421": "val421",
				"Key422": []string{"one", "two", "three"},
			},
			Key43: &ExSTChildLevel2{
				Key431: map[string]interface{}{
					"Key4311": "val4311",
				},
			},
		},
	}

	ks, err := NewKVStruct("localhost:8500", "adf4238a-882b-9ddc-4a9d-5b6758e4159e", "test")
	if err != nil {
		return
	}

	out := make(map[string]interface{})
	keys := []string{}

	ks.Path = "nestedstructmap"

	err = ks.StructToConsulKV(input)
	if err != nil {
		fmt.Println(err)
		return
	}

	pairs, _, err := ks.Client.KV().List(ks.Path, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, kv := range pairs {
		out[kv.Key] = string(kv.Value)
		keys = append(keys, kv.Key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		fmt.Println(key, ":", out[key])
	}
	// Output:
	//nestedstructmap/Key1 : val1
	//nestedstructmap/Key2 : 2
	//nestedstructmap/Key3/0 : 1
	//nestedstructmap/Key3/1 : 2
	//nestedstructmap/Key3/2 : 3
	//nestedstructmap/Key4/Key41 : val41
	//nestedstructmap/Key4/Key42/Key421 : val421
	//nestedstructmap/Key4/Key42/Key422/0 : one
	//nestedstructmap/Key4/Key42/Key422/1 : two
	//nestedstructmap/Key4/Key42/Key422/2 : three
	//nestedstructmap/Key4/Key43/Key431/Key4311 : val4311

	ks.Client.KV().DeleteTree(ks.Path, nil)
}

func ExampleMapToFlattenMap() {
	input := map[string]interface{}{
		"key1": "val1",
		"key2": 2,
		"key3": []int{1, 2, 3},
		"key4": map[string]interface{}{
			"key41": "val41",
			"key42": map[string]interface{}{
				"key421": "val421",
				"key422": []string{"one", "two", "three"},
			},
			"key43": map[string]interface{}{
				"key431": map[string]interface{}{
					"key4311": "val4311",
				},
			},
		},
	}

	o := MapToFlattenMap(input, "test")

	keys := []string{}

	for k := range o {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, key := range keys {
		fmt.Println(key, ":", o[key])
	}
	// Output:
	//test/key1 : val1
	//test/key2 : 2
	//test/key3 : [1 2 3]
	//test/key4/key41 : val41
	//test/key4/key42/key421 : val421
	//test/key4/key42/key422 : [one two three]
	//test/key4/key43/key431/key4311 : val4311
}

func ExampleKVMapToMap() {
	input := map[string]interface{}{
		"test/key1":                "val1",
		"test/key2":                2,
		"test/key3/0":              1,
		"test/key3/1":              2,
		"test/key3/2":              3,
		"test/key4/key41":          "val41",
		"test/key4/key42/key421":   "val421",
		"test/key4/key42/key422/0": "one",
		"test/key4/key42/key422/1": "two",
		"test/key4/key42/key422/2": "three",
	}

	o, err := KVMapToMap(input, "test")
	if err != nil {
		return
	}

	fmt.Println(o)
}

func ExampleFlattenMapToStruct() {
	type ExSTChildLevel2 struct {
		Key431 map[string]interface{}
	}

	type ExSTChildLevel1 struct {
		Key41 string
		Key42 map[string]interface{}
		Key43 *ExSTChildLevel2
	}

	type ExST struct {
		Key1 string
		Key2 int
		Key3 []int
		Key4 *ExSTChildLevel1
	}

	input := map[string]interface{}{
		"Key1": "val1",
		"Key2": 2,
		"Key3": []int{1, 2, 3},
		"Key4": map[string]interface{}{
			"Key41": "val41",
			"Key42": map[string]interface{}{
				"Key421": "val421",
				"Key422": []string{"one", "two", "three"},
			},
			"Key43": map[string]interface{}{
				"Key431": map[string]interface{}{
					"Key4311": "val4311",
				},
			},
		},
	}

	st := &ExST{
		Key4: &ExSTChildLevel1{
			Key43: &ExSTChildLevel2{},
		},
	}

	err := FlattenMapToStruct(input, st)
	if err != nil {
		return
	}

	fmt.Println(st)
}

func ExampleKVMapToStruct() {
	type ExSTChildLevel2 struct {
		Key431 map[string]interface{}
	}

	type ExSTChildLevel1 struct {
		Key41 string
		Key42 map[string]interface{}
		Key43 *ExSTChildLevel2
	}

	type ExST struct {
		Key1 string
		Key2 int
		Key3 []int
		Key4 *ExSTChildLevel1
	}

	input := map[string]interface{}{
		"Key1":                      "val1",
		"Key2":                      2,
		"Key3/0":                    1,
		"Key3/1":                    2,
		"Key3/2":                    3,
		"Key4/Key41":                "val41",
		"Key4/Key42/Key421":         "val421",
		"Key4/Key42/Key422/0":       "one",
		"Key4/Key42/Key422/1":       "two",
		"Key4/Key42/Key422/2":       "three",
		"Key4/Key43/Key431/Key4311": "val4311",
	}

	st := &ExST{
		Key4: &ExSTChildLevel1{
			Key43: &ExSTChildLevel2{},
		},
	}

	err := KVMapToStruct(input, "", st)
	if err != nil {
		return
	}

	fmt.Println(st)
}

func ExampleKVMapToStruct_embeddedStruct() {
	type ExSTChildLevel2 struct {
		Key431 map[string]interface{}
	}

	type ExSTChildLevel1 struct {
		Key41 string
		Key42 map[string]interface{}
		ExSTChildLevel2
	}

	type ExST struct {
		Key1 string
		Key2 int
		Key3 []int
		ExSTChildLevel1
	}

	input := map[string]interface{}{
		"Key1":                                           "val1",
		"Key2":                                           2,
		"Key3/0":                                         1,
		"Key3/1":                                         2,
		"Key3/2":                                         3,
		"ExSTChildLevel1/Key41":                          "val41",
		"ExSTChildLevel1/Key42/Key421":                   "val421",
		"ExSTChildLevel1/Key42/Key422/0":                 "one",
		"ExSTChildLevel1/Key42/Key422/1":                 "two",
		"ExSTChildLevel1/Key42/Key422/2":                 "three",
		"ExSTChildLevel1/ExSTChildLevel2/Key431/Key4311": "val4311",
	}

	st := &ExST{}

	err := KVMapToStruct(input, "", st)
	if err != nil {
		return
	}

	fmt.Printf("%++v\n", st)
	// Output:
	// &{Key1:val1 Key2:2 Key3:[1 2 3] ExSTChildLevel1:{Key41:val41 Key42:map[Key421:val421 Key422:[one two three]] ExSTChildLevel2:{Key431:map[Key4311:val4311]}}}
}

func ExampleKVStruct_ConsulKVToStruct_pointerStruct() {

	type ExSTChildLevel2 struct {
		Key431 map[string]interface{}
	}

	type ExSTChildLevel1 struct {
		Key41 string
		Key42 map[string]interface{}
		Key43 *ExSTChildLevel2
	}

	type ExST struct {
		Key1 string
		Key2 int
		Key3 []int
		Key4 *ExSTChildLevel1
	}

	input := map[string]interface{}{
		"test/Key1":                      "val1",
		"test/Key2":                      "2",
		"test/Key3/0":                    "1",
		"test/Key3/1":                    "2",
		"test/Key3/2":                    "3",
		"test/Key4/Key41":                "val41",
		"test/Key4/Key42/Key421":         "val421",
		"test/Key4/Key42/Key422/0":       "one",
		"test/Key4/Key42/Key422/1":       "two",
		"test/Key4/Key42/Key422/2":       "three",
		"test/Key4/Key43/Key431/Key4311": "val4311",
	}

	ks, err := NewKVStruct("localhost:8500", "adf4238a-882b-9ddc-4a9d-5b6758e4159e", "test")
	if err != nil {
		return
	}

	st := &ExST{
		Key4: &ExSTChildLevel1{
			Key43: &ExSTChildLevel2{},
		},
	}

	ks.Path = "test"

	// Insert data in to consul
	for k, v := range input {
		kv := &consul.KVPair{
			Key:   k,
			Value: []byte(v.(string)),
		}

		_, err := ks.Client.KV().Put(kv, nil)
		if err != nil {
			return
		}
	}

	err = ks.ConsulKVToStruct(st)
	if err != nil {
		return
	}

	fmt.Println(st)
}

func ExampleKVStruct_ConsulKVToStruct_embeddedStruct() {

	type ExSTChildLevel2 struct {
		Key431 map[string]interface{}
	}

	type ExSTChildLevel1 struct {
		Key41 string
		Key42 map[string]interface{}
		ExSTChildLevel2
	}

	type ExST struct {
		Key1 string
		Key2 int
		Key3 []int
		ExSTChildLevel1
	}

	input := map[string]interface{}{
		"test/Key1":                                           "val1",
		"test/Key2":                                           "2",
		"test/Key3/0":                                         "1",
		"test/Key3/1":                                         "2",
		"test/Key3/2":                                         "3",
		"test/ExSTChildLevel1/Key41":                          "val41",
		"test/ExSTChildLevel1/Key42/Key421":                   "val421",
		"test/ExSTChildLevel1/Key42/Key422/0":                 "one",
		"test/ExSTChildLevel1/Key42/Key422/1":                 "two",
		"test/ExSTChildLevel1/Key42/Key422/2":                 "three",
		"test/ExSTChildLevel1/ExSTChildLevel2/Key431/Key4311": "val4311",
	}

	ks, err := NewKVStruct("localhost:8500", "adf4238a-882b-9ddc-4a9d-5b6758e4159e", "test")
	if err != nil {
		return
	}

	st := &ExST{}

	ks.Path = "test"

	// Insert data in to consul
	for k, v := range input {
		kv := &consul.KVPair{
			Key:   k,
			Value: []byte(v.(string)),
		}

		_, err := ks.Client.KV().Put(kv, nil)
		if err != nil {
			return
		}
	}

	err = ks.ConsulKVToStruct(st)
	if err != nil {
		return
	}

	ks.Client.KV().DeleteTree(ks.Path, nil)

	fmt.Printf("%++v\n", st)
	// Output:
	// &{Key1:val1 Key2:2 Key3:[1 2 3] ExSTChildLevel1:{Key41:val41 Key42:map[Key421:val421 Key422:[one two three]] ExSTChildLevel2:{Key431:map[Key4311:val4311]}}}
}

func ExampleKVStruct_ConsulKVToStruct_normalStruct() {

	type ExSTChildLevel2 struct {
		Key431 map[string]interface{}
	}

	type ExSTChildLevel1 struct {
		Key41 string
		Key42 map[string]interface{}
		Key43 ExSTChildLevel2
	}

	type ExST struct {
		Key1 string
		Key2 int
		Key3 []int
		Key4 ExSTChildLevel1
	}

	input := map[string]interface{}{
		"test/Key1":                      "val1",
		"test/Key2":                      "2",
		"test/Key3/0":                    "1",
		"test/Key3/1":                    "2",
		"test/Key3/2":                    "3",
		"test/Key4/Key41":                "val41",
		"test/Key4/Key42/Key421":         "val421",
		"test/Key4/Key42/Key422/0":       "one",
		"test/Key4/Key42/Key422/1":       "two",
		"test/Key4/Key42/Key422/2":       "three",
		"test/Key4/Key43/Key431/Key4311": "val4311",
	}

	ks, err := NewKVStruct("localhost:8500", "adf4238a-882b-9ddc-4a9d-5b6758e4159e", "test")
	if err != nil {
		return
	}

	st := &ExST{}

	ks.Path = "test"

	// Insert data in to consul
	for k, v := range input {
		kv := &consul.KVPair{
			Key:   k,
			Value: []byte(v.(string)),
		}

		_, err := ks.Client.KV().Put(kv, nil)
		if err != nil {
			return
		}
	}

	err = ks.ConsulKVToStruct(st)
	if err != nil {
		return
	}

	ks.Client.KV().DeleteTree(ks.Path, nil)

	fmt.Printf("%++v\n", st)
	// Output:
	// &{Key1:val1 Key2:2 Key3:[1 2 3] Key4:{Key41:val41 Key42:map[Key421:val421 Key422:[one two three]] Key43:{Key431:map[Key4311:val4311]}}}
}

func ExampleKVStruct_ConsulKVToMap() {
	input := map[string]interface{}{
		"test/Key1":                      "val1",
		"test/Key2":                      "2",
		"test/Key3/0":                    "1",
		"test/Key3/1":                    "2",
		"test/Key3/2":                    "3",
		"test/Key4/Key41":                "val41",
		"test/Key4/Key42/Key421":         "val421",
		"test/Key4/Key42/Key422/0":       "one",
		"test/Key4/Key42/Key422/1":       "two",
		"test/Key4/Key42/Key422/2":       "three",
		"test/Key4/Key43/Key431/Key4311": "val4311",
	}

	ks, err := NewKVStruct("localhost:8500", "adf4238a-882b-9ddc-4a9d-5b6758e4159e", "test")
	if err != nil {
		return
	}

	ks.Path = "test"

	// Insert data in to consul
	for k, v := range input {
		kv := &consul.KVPair{
			Key:   k,
			Value: []byte(v.(string)),
		}

		_, err := ks.Client.KV().Put(kv, nil)
		if err != nil {
			return
		}
	}

	out, err := ks.ConsulKVToMap()
	if err != nil {
		return
	}

	fmt.Println(out)
}
