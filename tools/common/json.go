package common

import jsoniter "github.com/json-iterator/go"

var myjson = jsoniter.ConfigCompatibleWithStandardLibrary

func Marshal(v interface{}) ([]byte, error) {

	return myjson.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {

	return myjson.Unmarshal(data, v)
}
