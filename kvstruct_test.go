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

func TestStructToConsulKV(t *testing.T) {
	testCases := []struct {
		name   string
		prefix string
		input  map[string]interface{}
		output map[string]interface{}
	}{
		{
			"NestedMapWithoutPrefix",
			"test1",
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
				"test1/key1":                "val1",
				"test1/key2":                "2",
				"test1/key3/0":              "1",
				"test1/key3/1":              "2",
				"test1/key3/2":              "3",
				"test1/key4/key41":          "val41",
				"test1/key4/key42/key421":   "val421",
				"test1/key4/key42/key422/0": "one",
				"test1/key4/key42/key422/1": "two",
				"test1/key4/key42/key422/2": "three",
			},
		},
		{
			"NestedMapWithPrefix",
			"test2",
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
				"test2/key1":                "val1",
				"test2/key2":                "2",
				"test2/key3/0":              "1",
				"test2/key3/1":              "2",
				"test2/key3/2":              "3",
				"test2/key4/key41":          "val41",
				"test2/key4/key42/key421":   "val421",
				"test2/key4/key42/key422/0": "one",
				"test2/key4/key42/key422/1": "two",
				"test2/key4/key42/key422/2": "three",
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
		})
	}

}
