package goutils

import "github.com/mitchellh/mapstructure"

// Decode mapstructure.Decode
func Decode(in, out interface{}) error {
	var mapstructureConfig = &mapstructure.DecoderConfig{
		TagName: "json",
		Result:  out,
	}

	decoder, err := mapstructure.NewDecoder(mapstructureConfig)
	if err != nil {
		return err
	}
	err = decoder.Decode(in)
	return err

}
