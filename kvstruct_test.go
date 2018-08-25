package kvstruct

import (
	//"fmt"
	"reflect"
	"testing"
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

	ks := NewKVStruct("", "", "test")

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
