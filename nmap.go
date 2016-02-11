package sonar

import "bytes"

func ToNmap(results Results) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})

	for _, result := range results {
		for _, addr := range result.Addrs {
			_, err := buffer.WriteString(addr + "\n")
			if err != nil {
				return buffer.Bytes(), err
			}
		}
	}

	return buffer.Bytes(), nil
}
