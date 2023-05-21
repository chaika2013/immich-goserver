package helper

import "strconv"

func StringID(id uint) string {
	return strconv.FormatUint(uint64(id), 10)
}
