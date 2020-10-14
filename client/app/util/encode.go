package util

import "encoding/json"

func Encode(s interface{}) ([]byte, error) {
	spec, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return spec, nil
}

func PrettyEncode(s interface{}) ([]byte, error) {
	spec, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		return nil, err
	}
	return spec, nil
}
