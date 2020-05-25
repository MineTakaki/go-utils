package utils

import "reflect"

//IsNil インターフェイスがnilかどうかを判定します
func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	return reflect.ValueOf(i).IsNil()
}
