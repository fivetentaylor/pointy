package dag

import "reflect"

func GetNodeType(node Node) string {
	t := reflect.TypeOf(node)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t.Name()
}
