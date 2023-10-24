package mongoutils

import (
	"crypto/md5"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/gomig/utils"
)

type Checksum struct {
	data map[string]any
}

func NewChecksum(data map[string]any) Checksum {
	return Checksum{data: data}
}

func (recv Checksum) MD5() string {
	md5 := md5.Sum([]byte(recv.Normalize()))
	return fmt.Sprintf("%x", md5)
}

func (recv Checksum) Normalize() string {
	val := reflect.ValueOf(recv.data)
	if recv.isNil(val) || recv.isEmptyString(recv.data) {
		return ""
	}

	// convert object to map
	result := make(map[string]string)
	if (len(recv.data)) > 0 {
		recv.flattern(recv.data, "", result)
	}

	// sort
	keys := make([]string, 0, len(result))
	for k := range result {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// concatenate
	resStr := ""
	for _, k := range keys {
		resStr += fmt.Sprintf("%s:%s|", k, result[k])
	}
	resStr = strings.TrimSuffix(resStr, "|")
	return resStr
}

func (recv Checksum) flattern(v any, key string, out map[string]string) {
	val := reflect.ValueOf(v)
	if recv.isNil(val) || recv.isEmptyString(v) {
		out[key] = ""
	}

	kind := val.Kind()
	if kind == reflect.Ptr {
		if val.IsNil() {
			out[key] = ""
		} else {
			recv.flattern(val.Elem().Interface(), key, out)
		}
	} else if recv.isSimple(kind) {
		if val.CanInt() {
			out[key] = strconv.FormatInt(val.Int(), 10)
		} else if val.CanUint() {
			out[key] = strconv.FormatUint(val.Uint(), 10)
		} else if val.CanFloat() {
			out[key] = fmt.Sprintf("%.0f", val.Float())
		} else {
			out[key] = fmt.Sprint(v)
		}
	} else if kind == reflect.Array || kind == reflect.Slice {
		for i := 0; i < val.Len(); i++ {
			recv.flattern(val.Index(i).Interface(), recv.keyOf(key, fmt.Sprintf("E%d", i)), out)
		}
	} else if kind == reflect.Map {
		for _, k := range val.MapKeys() {
			name := fmt.Sprint(k.Interface())
			value := val.MapIndex(k)
			recv.flattern(value.Interface(), recv.keyOf(key, name), out)
		}
	}
}

func (Checksum) keyOf(root, key string) string {
	return utils.If(root == "", key, root+"."+key)
}

func (Checksum) isSimple(kind reflect.Kind) bool {
	return !utils.Contains([]reflect.Kind{
		reflect.Array,
		reflect.Slice,
		reflect.Map,
		reflect.Struct,
	}, kind)
}

func (Checksum) isNil(val reflect.Value) bool {
	return utils.Contains([]reflect.Kind{
		reflect.Interface,
		reflect.Map,
		reflect.Pointer,
		reflect.Slice,
	}, val.Kind()) && val.IsNil()
}

func (Checksum) isEmptyString(v any) bool {
	return fmt.Sprint(v) == ""
}
