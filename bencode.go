package bt

import (
	"bufio"
	"io"
	"strconv"
)

func ReadString(reader *bufio.Reader) (string, error) {
	sLength, err := reader.ReadString(':')
	if err != nil {
		return "", err
	}
	sLength = sLength[0 : len(sLength)-1]

	length, err := strconv.ParseInt(sLength, 10, 64)
	if err != nil {
		return "", err
	}

	stringBuf := make([]byte, length)
	_, err = io.ReadFull(reader, stringBuf)
	if err != nil {
		return "", err
	}

	return string(stringBuf), nil
}

func ReadInt64(reader *bufio.Reader) (int64, error) {
	sNum, err := reader.ReadString('e')
	if err != nil {
		return 0, err
	}
	sNum = sNum[1 : len(sNum)-1]

	return strconv.ParseInt(sNum, 10, 64)
}
