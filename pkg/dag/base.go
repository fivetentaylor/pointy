package dag

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

type Base struct {
	InputAliases map[string]string
}

func (b *Base) Run(ctx context.Context) (Node, error) {
	return nil, errors.New("not implemented")
}

func (b *Base) Alias(from, to string) {
	b.InputAliases[from] = to
}

func (b *Base) Hydrate(ctx context.Context, input any) error {
	state := GetState(ctx)

	vo := reflect.ValueOf(input)
	if vo.Kind() != reflect.Ptr {
		return fmt.Errorf("expected pointer, got %T", input)
	}
	v := vo.Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.CanSet() {
			tag := t.Field(i).Tag.Get("key")
			if alias, ok := b.InputAliases[tag]; ok {
				tag = alias
			}
			if tag != "" {
				val := state.Get(tag)
				if val != nil {
					field.Set(reflect.ValueOf(val).Convert(field.Type()))
				}
			}
		}
	}

	return nil

}
