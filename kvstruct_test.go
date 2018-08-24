package kvstruct

import (
	//"encoding/base64"
	"fmt"
	"reflect"
	"testing"
	//consul "github.com/hashicorp/consul/api"
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
				"key1":              "val1",
				"key2":              2,
				"key3":              []int{1, 2, 3},
				"key4/key41":        "val41",
				"key4/key42/key421": "val421",
				"key4/key42/key422": []string{"one", "two", "three"},
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
				"test/key1":              "val1",
				"test/key2":              2,
				"test/key3":              []int{1, 2, 3},
				"test/key4/key41":        "val41",
				"test/key4/key42/key421": "val421",
				"test/key4/key42/key422": []string{"one", "two", "three"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			o := MapToKVMap(tc.input, tc.prefix)
			//fmt.Println(o)

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
				"key1":              "val1",
				"key2":              2,
				"key3":              []int{1, 2, 3},
				"key4/key41":        "val41",
				"key4/key42/key421": "val421",
				"key4/key42/key422": []string{"one", "two", "three"},
			},
		},
		//{
		//"NestedMapWithPrefix",
		//"test",
		//map[string]interface{}{
		//"key1": "val1",
		//"key2": 2,
		//"key3": []int{1, 2, 3},
		//"key4": map[string]interface{}{
		//"key41": "val41",
		//"key42": map[string]interface{}{
		//"key421": "val421",
		//"key422": []string{"one", "two", "three"},
		//},
		//},
		//},
		//map[string]interface{}{
		//"test/key1":              "val1",
		//"test/key2":              2,
		//"test/key3":              []int{1, 2, 3},
		//"test/key4/key41":        "val41",
		//"test/key4/key42/key421": "val421",
		//"test/key4/key42/key422": []string{"one", "two", "three"},
		//},
		//},
	}

	ks := NewKVStruct("", "", "test")

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			o, err := ks.MapToKVPairs(tc.input, tc.prefix)
			if err != nil {
				fmt.Println(err)
			}
			for _, kv := range o {
				fmt.Println(kv.Key, string(kv.Value))
			}

			//if !reflect.DeepEqual(o, tc.output) {
			//t.Errorf("\nwant:\n%s\nhave:\n%s", tc.output, o)
			//}
		})
	}

}
