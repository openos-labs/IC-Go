package agent

import (
	"math/big"
	"reflect"

	"github.com/mix-labs/IC-Go/utils/idl"
)



func Decode(target interface{}, source interface{}) {
	Type := reflect.TypeOf(target).Elem()
	Value := reflect.ValueOf(target)
	_Decode(Value, Type, source)
}

func _Decode(target reflect.Value, targetType reflect.Type,source interface{}) {
	if targetType.Kind() == reflect.Struct {
		if targetType.Name() == "Int" {
			value := source.(*big.Int)
			ptarget := target.Addr().Interface().(*big.Int)
			*ptarget = *value
			return
		}
		//source_ := source.(*idl.FieldValue)
		sourceField := source.(map[string]interface{})

		for k, v := range(sourceField) {
			for i := 0; i < targetType.NumField(); i++ {
				targetFiledType := targetType.Field(i)
				//a := targetFiledType.Tag.Get("ic")
				//fmt.Println(a)
				if idl.Hash(targetFiledType.Tag.Get("ic")).String() == k {
					targetFiledValue := target.Elem().Field(i)		
					_Decode(targetFiledValue.Addr(), targetFiledType.Type, v)
				}
			}
		}
		for k, v := range(sourceField) {
			for i := 0; i < targetType.NumField(); i++ {
				targetFiledType := targetType.Field(i)
				//a := targetFiledType.Tag.Get("ic")
				//fmt.Println(a)
				if targetFiledType.Tag.Get("ic") == k {
					targetFiledValue := target.Elem().Field(i)		
					_Decode(targetFiledValue.Addr(), targetFiledType.Type, v)
				}
			}
		}
	} else if targetType.Kind() == reflect.String {
		sourceFiled := source.(string)
		target.Elem().SetString(sourceFiled)
	} else if targetType.Kind() == reflect.Int || targetType.Kind() == reflect.Int8 || targetType.Kind() == reflect.Int16 || targetType.Kind() == reflect.Int32 || targetType.Kind() == reflect.Int64 {
		source_ := source.(*big.Int)
		sourceFiled := source_.Int64()
		target.Elem().SetInt(sourceFiled)
	} else if targetType.Kind() == reflect.Uint || targetType.Kind() == reflect.Uint8 || targetType.Kind() == reflect.Uint16 || targetType.Kind() == reflect.Uint32 || targetType.Kind() == reflect.Uint64 {
		source_ := source.(*big.Int)
		sourceFiled := source_.Uint64()
		target.Elem().SetUint(sourceFiled)
	} else if targetType.Kind() == reflect.Slice {
		if targetType.Name() == "Principal" {
			sourceFiled := source.([]uint8)
			target.Elem().SetBytes(sourceFiled)
			return
		}
		sourceFiled := source.([]interface{})
		var elem reflect.Value
		Type := targetType.Elem()
		for _, v := range(sourceFiled) {
			elem = reflect.New(Type)
			_Decode(elem, elem.Type().Elem(), v)
			target.Elem().Set(reflect.Append(target.Elem(), elem.Elem()))
		}
	}



}
