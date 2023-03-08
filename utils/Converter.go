package Utils

import "strconv"

func UintToString(id string) uint64 {
	parsedId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return 0
	}
	return parsedId
}
