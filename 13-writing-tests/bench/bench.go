package bench

import (
	"io"
	"os"
)

func FileLen(fName string, bufSize int) (int, error) {
	file, err := os.Open(fName)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	count := 0
	for {
		buf := make([]byte, bufSize)
		num, err := file.Read(buf)
		count += num
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}
	}
	return count, nil
}

func FileLen2(fName string, bufSize int) (int, error) {
	file, err := os.Open(fName)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	count := 0
	// Moving the buffer creation out of the loop prevents a lot of memory allocations
	buf := make([]byte, bufSize)
	for {
		num, err := file.Read(buf)
		count += num
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}
	}
	return count, nil
}
