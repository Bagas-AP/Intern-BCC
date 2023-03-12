package Utils

import "strconv"

func ParseStrToUint(id string) (uint64, error) {
	parsedId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return 0, err
	}
	return parsedId, nil
}
