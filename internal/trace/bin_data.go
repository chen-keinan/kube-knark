package trace

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
)

var assets map[string][]byte

func asset(key string) ([]byte, error) {
	if assets == nil {
		assets = map[string][]byte{}

		var value []byte
		value, _ = base64.StdEncoding.DecodeString("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")
		value, _ = gzipDecode(value)
		assets["exec.o"] = value
	}

	if value, found := assets[key]; found {
		return value, nil
	}
	return nil, fmt.Errorf("asset not found for key=%v", key)
}

func gzipDecode(data []byte) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	out := new(bytes.Buffer)
	if _, err = io.Copy(out, gz); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}
