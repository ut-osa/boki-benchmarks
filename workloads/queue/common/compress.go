package common

import (
	"bytes"
	"compress/gzip"
	"io"
)

func CompressData(uncompressed []byte) []byte {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	if _, err := zw.Write(uncompressed); err != nil {
		panic(err)
	}
	if err := zw.Close(); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func DecompressFromReader(reader io.Reader) (io.Reader, error) {
	return gzip.NewReader(reader)
}

func DecompressReader(compressed []byte) (io.Reader, error) {
	reader := bytes.NewReader(compressed)
	return DecompressFromReader(reader)
}
