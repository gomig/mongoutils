package mongoutils

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/mongo"
)

type countResult struct {
	Count int64 `bson:"count" json:"count"`
}

type MongoOption struct {
	DebugPipe   bool
	DebugResult bool
	Pipeline    string
	Params      []any
	Database    *mongo.Database
}

// optionOf get option of dynamic params or return empty option
func optionOf(opts ...MongoOption) MongoOption {
	opt := MongoOption{}
	if len(opts) > 0 {
		opt = opts[0]
	}
	if opt.Pipeline == "" {
		opt.Pipeline = "Pipeline"
	}
	return opt
}

// parsePipeline get pipeline from CallMethod result or return nil
func parsePipeline(res []reflect.Value) MongoPipeline {
	if len(res) > 0 {
		if v, ok := res[0].Interface().(MongoPipeline); ok {
			return v
		}
	}
	return nil
}

// callMethod call object method dynamically
func callMethod(obj any, method string, params ...any) ([]reflect.Value, error) {
	_type := reflect.TypeOf(obj)
	for i := 0; i < _type.NumMethod(); i++ {
		_method := _type.Method(i)
		if method == _method.Name {
			vals := make([]reflect.Value, 0)
			vals = append(vals, reflect.ValueOf(obj))
			for _, p := range params {
				vals = append(vals, reflect.ValueOf(p))
			}
			return _method.Func.Call(vals), nil
		}
	}
	return nil, errors.New("method " + method + " not defined!")
}

// prettyLog log data to output using json indent format
func prettyLog(data any) {
	_bytes, _ := json.MarshalIndent(data, "", "    ")
	fmt.Println(string(_bytes))
}

func modelChecksum(v any) (string, Backup) {
	if model, ok := parseAsInterface[Backup](v); ok && model != nil {
		return NewChecksum(model.ToMap()).MD5(), model
	}
	return "", nil
}

// modelSafe convert v to github.com/gomig/mongoutils.Model or panic
func modelSafe[T any](v T) Model {
	if _v, ok := any(v).(Model); !ok {
		panic("T must implements github.com/gomig/mongoutils.Model")
	} else {
		return _v
	}
}

// Get new instance of github.com/gomig/mongoutils.Model or panic if T not implement model
func typeModelSafe[T any]() Model {
	return modelSafe(new(T))
}

func parseAsInterface[T any](v any) (T, bool) {
	i, ok := v.(T)
	return i, ok
}
