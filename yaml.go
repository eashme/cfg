package coord_cfg

import (
	"fmt"
	"github.com/alediaferia/stackgo"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"reflect"
	"strconv"
	"sync"
)

var (
	yamlOnce  sync.Once
	yamlCache map[string]interface{}
)

// 获取 K V yaml文件配置值
func getFromYaml(k string) string {
	yamlOnce.Do(func() {
		// 获取yaml 配置文件路径
		fn := getFromEnv(YamlFileEnvKey)
		if fn == "" {
			fn = YamlDefaultPath
		}
		// 读取配置文件,绑定到cfg
		b, err := ioutil.ReadFile(fn)
		if err != nil {
			log.Print("failed load cfg ", err)
			return
		}
		yamlCache = Yaml2Map(b)
		log.Printf("[yaml cfg] %s", Json(yamlCache))
	})
	v, _ := yamlCache[StandCode(k)]
	if v == nil{
		return ""
	}
	return fmt.Sprintf("%v", v)
}

type I interface{}

var (
	mKeys      = make(map[string]interface{})
	currentKey string
	stKeys     = stackgo.NewStack()
)

func Yaml2Map(data3 []byte) map[string]interface{} {

	mKeys0 := make(map[string]interface{})
	stKeys0 := stackgo.NewStack()
	mKeys = mKeys0
	stKeys = stKeys0

	m := make(map[string]interface{})
	checkError(yaml.Unmarshal(data3, &m))
	for k, v := range m {
		currentKey = k
		stKeys.Push(currentKey)
		extract(v)
	}
	return mKeys
}

func extract(obj interface{}) interface{} {
	// Wrap the original in a reflect.Value
	original := reflect.ValueOf(obj)

	cp := reflect.New(original.Type()).Elem()
	extractRecursive(cp, original)

	// Remove the reflection wrapper
	return cp.Interface()
}

func extractRecursive(copy, original reflect.Value) {
	var existingValue string
	var NewValue string
	switch original.Kind() {
	// The first cases handle nested structures and extract them recursively

	// If it is a pointer we need to unwrap and call once again
	case reflect.Ptr:
		// To get the actual value of the original we have to call Elem()
		// At the same time this unwraps the pointer so we don't end up in
		// an infinite recursion
		originalValue := original.Elem()
		// Check if the pointer is nil
		if !originalValue.IsValid() {
			return
		}
		// Allocate a new object and set the pointer to it
		copy.Set(reflect.New(originalValue.Type()))
		// Unwrap the newly created pointer
		extractRecursive(copy.Elem(), originalValue)

	// If it is an interface (which is very similar to a pointer), do basically the
	// same as for the pointer. Though a pointer is not the same as an interface so
	// note that we have to call Elem() after creating a new object because otherwise
	// we would end up with an actual pointer
	case reflect.Interface:
		// Get rid of the wrapping interface
		originalValue := original.Elem()
		// Create a new object. Now new gives us a pointer, but we want the value it
		// points to, so we have to call Elem() to unwrap it
		copyValue := reflect.New(originalValue.Type()).Elem()
		extractRecursive(copyValue, originalValue)
		copy.Set(copyValue)

	// If it is a struct we extract each field
	case reflect.Struct:
		for i := 0; i < original.NumField(); i += 1 {
			extractRecursive(copy.Field(i), original.Field(i))
		}

	// If it is a slice we create a new slice and extract each element
	case reflect.Slice:
		copy.Set(reflect.MakeSlice(original.Type(), original.Len(), original.Cap()))
		for i := 0; i < original.Len(); i += 1 {
			extractRecursive(copy.Index(i), original.Index(i))
		}

	// If it is a map we create a new map and extract each value
	case reflect.Map:
		copy.Set(reflect.MakeMap(original.Type()))
		stKeys.Push(currentKey)
		for _, key := range original.MapKeys() {
			originalValue := original.MapIndex(key)

			stKeys.Push(currentKey)
			currentKey = StandCode(fmt.Sprintf("%s_%s", currentKey, key))

			// New gives us a pointer, but again we want the value
			copyValue := reflect.New(originalValue.Type()).Elem()
			extractRecursive(copyValue, originalValue)
			copy.SetMapIndex(key, copyValue)
			currentKey = stKeys.Pop().(string)
		}
		currentKey = stKeys.Pop().(string)

	// Otherwise we cannot traverse anywhere so this finishes the the recursion

	// If it is a string extract it (yay finally we're doing what we came for)
	case reflect.String:
		extractString := original.Interface().(string)
		copy.SetString(extractString)
		if val, ok := mKeys[currentKey]; ok {
			existingValue = val.(string)
			NewValue = existingValue + "," + extractString
			mKeys[currentKey] = NewValue
		} else {
			mKeys[currentKey] = extractString
		}

		// A bool type will always be a value, convert it to string before saving
	case reflect.Bool:
		var tf = original.Bool()
		extractString := strconv.FormatBool(tf)
		mKeys[currentKey] = extractString
		copy.Set(original)

	// And everything else will simply be taken from the original
	default:
		copy.Set(original)
	}
}
