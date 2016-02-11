package sonar

import "bytes"

func ToNmap(results Results) []byte {
	buffer := bytes.NewBuffer([]byte{})

	for _, result := range results {
		for _, addr := range result.Addrs {
			buffer.WriteString(addr + "\n")
		}
	}

	return buffer.Bytes()
}
