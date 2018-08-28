package kvstruct

import (
	//"fmt"
	"reflect"
	"testing"

	consul "github.com/hashicorp/consul/api"
)

func TestMapToKVMap(t *testing.T) {
	testCases := []struct {
		name   string
		prefix string
		input  map[string]interface{}
		output map[string]interface{}
	}{
		{
			"NestedMapWithoutPrefix",
			"",
			map[string]interface{}{
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
			},
			map[string]interface{}{
				"key1":                "val1",
				"key2":                2,
				"key3/0":              1,
				"key3/1":              2,
				"key3/2":              3,
				"key4/key41":          "val41",
				"key4/key42/key421":   "val421",
				"key4/key42/key422/0": "one",
				"key4/key42/key422/1": "two",
				"key4/key42/key422/2": "three",
			},
		},
		{
			"NestedMapWithPrefix",
			"test",
			map[string]interface{}{
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
			},
			map[string]interface{}{
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
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			o := MapToKVMap(tc.input, tc.prefix)

			if !reflect.DeepEqual(o, tc.output) {
				t.Errorf("\nwant:\n%s\nhave:\n%s", tc.output, o)
			}
		})
	}

}

func TestMapToKVPairs(t *testing.T) {
	testCases := []struct {
		name   string
		prefix string
		input  map[string]interface{}
		output map[string]interface{}
	}{
		{
			"NestedMapWithoutPrefix",
			"",
			map[string]interface{}{
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
			},
			map[string]interface{}{
				"key1":                "val1",
				"key2":                "2",
				"key3/0":              "1",
				"key3/1":              "2",
				"key3/2":              "3",
				"key4/key41":          "val41",
				"key4/key42/key421":   "val421",
				"key4/key42/key422/0": "one",
				"key4/key42/key422/1": "two",
				"key4/key42/key422/2": "three",
			},
		},
		{
			"NestedMapWithPrefix",
			"test",
			map[string]interface{}{
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
			},
			map[string]interface{}{
				"test/key1":                "val1",
				"test/key2":                "2",
				"test/key3/0":              "1",
				"test/key3/1":              "2",
				"test/key3/2":              "3",
				"test/key4/key41":          "val41",
				"test/key4/key42/key421":   "val421",
				"test/key4/key42/key422/0": "one",
				"test/key4/key42/key422/1": "two",
				"test/key4/key42/key422/2": "three",
			},
		},
	}

	ks, err := NewKVStruct("", "", "test")
	if err != nil {
		t.Errorf("%s", err.Error())
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out := make(map[string]interface{})

			o, err := ks.MapToKVPairs(tc.input, tc.prefix)
			if err != nil {
				t.Errorf("%s", err.Error())
			}

			for _, kv := range o {
				out[kv.Key] = string(kv.Value)
			}

			if !reflect.DeepEqual(out, tc.output) {
				t.Errorf("\nwant:\n%s\nhave:\n%s", tc.output, out)
			}
		})
	}

}

func TestMapToConsulKV(t *testing.T) {
	testCases := []struct {
		name   string
		prefix string
		input  interface{}
		output map[string]interface{}
	}{
		{
			"NestedMap",
			"nestedmap",
			map[string]interface{}{
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
			},
			map[string]interface{}{
				"nestedmap/key1":                "val1",
				"nestedmap/key2":                "2",
				"nestedmap/key3/0":              "1",
				"nestedmap/key3/1":              "2",
				"nestedmap/key3/2":              "3",
				"nestedmap/key4/key41":          "val41",
				"nestedmap/key4/key42/key421":   "val421",
				"nestedmap/key4/key42/key422/0": "one",
				"nestedmap/key4/key42/key422/1": "two",
				"nestedmap/key4/key42/key422/2": "three",
			},
		},
	}

	ks, err := NewKVStruct("localhost:8500", "adf4238a-882b-9ddc-4a9d-5b6758e4159e", "test")
	if err != nil {
		t.Errorf("%s", err.Error())
	}

	config := consul.DefaultConfig()
	config.Address = "localhost:8500"
	config.Token = "adf4238a-882b-9ddc-4a9d-5b6758e4159e"

	client, err := consul.NewClient(config)
	if err != nil {
		t.Errorf("%s", err.Error())
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out := make(map[string]interface{})

			ks.Path = tc.prefix
			err := ks.StructToConsulKV(tc.input)
			if err != nil {
				t.Errorf("%s", err.Error())
			}

			pairs, _, err := client.KV().List(tc.prefix, nil)
			if err != nil {
				t.Errorf("%s", err.Error())
			}

			for _, kv := range pairs {
				out[kv.Key] = string(kv.Value)
			}

			if !reflect.DeepEqual(out, tc.output) {
				t.Errorf("\nwant:\n%s\nhave:\n%s", tc.output, out)
			}

			client.KV().DeleteTree(tc.prefix, nil)
		})
	}

}

type ST struct {
	Key1 string
	Key2 int
	Key3 []int
	Key4 *STChildLevel1
}

type STChildLevel1 struct {
	Key41 string
	Key42 map[string]interface{}
	Key43 *STChildLevel2
}

type STChildLevel2 struct {
	Key431 map[string]interface{}
}

func TestStructToConsulKV(t *testing.T) {
	testCases := []struct {
		name   string
		prefix string
		input  ST
		output map[string]interface{}
	}{
		{
			"Struct",
			"nestedstructmap",
			ST{
				Key1: "val1",
				Key2: 2,
				Key3: []int{1, 2, 3},
				Key4: &STChildLevel1{
					Key41: "val41",
					Key42: map[string]interface{}{
						"Key421": "val421",
						"Key422": []string{"one", "two", "three"},
					},
					Key43: &STChildLevel2{
						Key431: map[string]interface{}{
							"Key4311": "val4311",
						},
					},
				},
			},
			map[string]interface{}{
				"nestedstructmap/Key1":                      "val1",
				"nestedstructmap/Key2":                      "2",
				"nestedstructmap/Key3/0":                    "1",
				"nestedstructmap/Key3/1":                    "2",
				"nestedstructmap/Key3/2":                    "3",
				"nestedstructmap/Key4/Key41":                "val41",
				"nestedstructmap/Key4/Key42/Key421":         "val421",
				"nestedstructmap/Key4/Key42/Key422/0":       "one",
				"nestedstructmap/Key4/Key42/Key422/1":       "two",
				"nestedstructmap/Key4/Key42/Key422/2":       "three",
				"nestedstructmap/Key4/Key43/Key431/Key4311": "val4311",
			},
		},
	}

	ks, err := NewKVStruct("localhost:8500", "adf4238a-882b-9ddc-4a9d-5b6758e4159e", "test")
	if err != nil {
		t.Errorf("%s", err.Error())
	}

	config := consul.DefaultConfig()
	config.Address = "localhost:8500"
	config.Token = "adf4238a-882b-9ddc-4a9d-5b6758e4159e"

	client, err := consul.NewClient(config)
	if err != nil {
		t.Errorf("%s", err.Error())
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out := make(map[string]interface{})

			ks.Path = tc.prefix

			err := ks.StructToConsulKV(tc.input)
			if err != nil {
				t.Errorf("%s", err.Error())
			}

			pairs, _, err := client.KV().List(tc.prefix, nil)
			if err != nil {
				t.Errorf("%s", err.Error())
			}

			for _, kv := range pairs {
				out[kv.Key] = string(kv.Value)
			}

			if !reflect.DeepEqual(out, tc.output) {
				t.Errorf("\nwant:\n%s\nhave:\n%s", tc.output, out)
			}

			client.KV().DeleteTree(tc.prefix, nil)
		})
	}

}

func TestMapToFlattenMap(t *testing.T) {
	testCases := []struct {
		name   string
		prefix string
		input  map[string]interface{}
		output map[string]interface{}
	}{
		{
			"NestedMapWithPrefix",
			"test",
			map[string]interface{}{
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
			},
			map[string]interface{}{
				"test/key1":                      "val1",
				"test/key2":                      2,
				"test/key3":                      []int{1, 2, 3},
				"test/key4/key41":                "val41",
				"test/key4/key42/key421":         "val421",
				"test/key4/key42/key422":         []string{"one", "two", "three"},
				"test/key4/key43/key431/key4311": "val4311",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			o := MapToFlattenMap(tc.input, tc.prefix)

			if !reflect.DeepEqual(o, tc.output) {
				t.Errorf("\nwant:\n%s\nhave:\n%s", tc.output, o)
			}
		})
	}

}

func TestKVMapToMap(t *testing.T) {
	testCases := []struct {
		name   string
		prefix string
		input  map[string]interface{}
		output map[string]interface{}
	}{
		{
			"NestedMapWithPrefix",
			"test",
			map[string]interface{}{
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
			},
			map[string]interface{}{
				"key1": "val1",
				"key2": "2",
				"key3": []int{1, 2, 3},
				"key4": map[string]interface{}{
					"key41": "val41",
					"key42": map[string]interface{}{
						"key421": "val421",
						"key422": []string{"one", "two", "three"},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			o, err := KVMapToMap(tc.input, tc.prefix)
			if err != nil {
				t.Errorf("%s", err)
			}

			if !reflect.DeepEqual(o, tc.output) {
				t.Errorf("\nwant:\n%s\nhave:\n%s", tc.output, o)
			}
		})
	}

}

func TestFlattenMapToStruct(t *testing.T) {
	testCases := []struct {
		name   string
		prefix string
		input  map[string]interface{}
		output *ST
	}{
		{
			"NestedMapWithPrefix",
			"",
			map[string]interface{}{
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
			},
			&ST{
				Key1: "val1",
				Key2: 2,
				Key3: []int{1, 2, 3},
				Key4: &STChildLevel1{
					Key41: "val41",
					Key42: map[string]interface{}{
						"Key421": "val421",
						"Key422": []string{"one", "two", "three"},
					},
					Key43: &STChildLevel2{
						Key431: map[string]interface{}{
							"Key4311": "val4311",
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			st := &ST{
				Key4: &STChildLevel1{
					Key43: &STChildLevel2{},
				},
			}

			err := FlattenMapToStruct(tc.input, tc.prefix, st)
			if err != nil {
				t.Errorf("%s", err)
			}

			if !reflect.DeepEqual(st, tc.output) {
				t.Errorf("\nwant:\n%v\nhave:\n%v", tc.output, st)
			}
		})
	}

}

func TestKVMapToStruct(t *testing.T) {
	testCases := []struct {
		name   string
		prefix string
		input  map[string]interface{}
		output *ST
	}{
		{
			"NestedMapWithPrefix",
			"",
			map[string]interface{}{
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
			},
			&ST{
				Key1: "val1",
				Key2: 2,
				Key3: []int{1, 2, 3},
				Key4: &STChildLevel1{
					Key41: "val41",
					Key42: map[string]interface{}{
						"Key421": "val421",
						"Key422": []string{"one", "two", "three"},
					},
					Key43: &STChildLevel2{
						Key431: map[string]interface{}{
							"Key4311": "val4311",
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			st := &ST{
				Key4: &STChildLevel1{
					Key43: &STChildLevel2{},
				},
			}

			err := KVMapToStruct(tc.input, tc.prefix, st)
			if err != nil {
				t.Errorf("%s", err)
			}

			if !reflect.DeepEqual(st, tc.output) {
				t.Errorf("\nwant:\n%v\nhave:\n%v", tc.output, st)
			}
		})
	}

}
