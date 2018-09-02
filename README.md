kvmapstruct
========
Package `kvstruct` exposes various utility functions to do conversions between: Consul KV pairs and native Go Struct or map[string]interface{}.

It also provides several utilities to convert directly:
- Nested map to flatten/kv map or Consul kv pairs
- Flatten/kv map to Go struct
- Kv map to nested map, etc.

There are some notions that are used in this package.
- Nested map: classic nested map[string]interface{}.
- Flatten map: map[string]interface{} represents key/value. It means that no nested map will be value. Value can be a normal type including slice.
- KV map: map[string]interface{} represents key/value but value can not be slice or map. A slice will be represented by keys suffixed by 0, 1, 2 etc.

This package only supports the following value types:
int, bool, string, []int, []bool, []string and map[string]interface{}

Documentation
-----------
See the [Godoc](https://godoc.org/github.com/uthng/kvmapstruct)