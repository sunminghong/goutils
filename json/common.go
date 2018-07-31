// Copyright 2017 Bo-Yi Wu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.


package json

import (
  "io"
  "bytes"
)

const EnableDecoderUseNumber bool = true

func Decode(body []byte, obj interface{}) error {
	return DecodeFromReader(bytes.NewReader(body), obj)
}

func DecodeFromReader(r io.Reader, obj interface{}) error {
	decoder := json.NewDecoder(r)
	if EnableDecoderUseNumber {
		decoder.UseNumber()
	}
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return nil
}

func Encode(obj interface{}) (string,error) {
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(jsonBytes),nil
}

func EncodeBytes(obj interface{}) ([]byte,error) {
	jsonBytes, err := json.MarshalIndent(obj,"","  ")
	if err != nil {
		return nil, err
	}
	return jsonBytes,nil
}
