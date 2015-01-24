package rainmaker

import "strconv"

func ParseInt(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}

	parsedValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		panic(err)
	}

	return int(parsedValue)
}
