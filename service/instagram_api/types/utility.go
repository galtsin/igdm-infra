package types

import "strconv"

func ValueToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int64:
		return strconv.FormatInt(v, 10)
	case int:
		return strconv.Itoa(v)
	default:
		return ""
	}
}

func ValueToInt(value interface{}) int {
	switch v := value.(type) {
	case int64:
		return int(v)
	case int:
		return v
	case string:
		i, err := strconv.Atoi(v)
		if err != nil {
			return 0
		}

		return i
	default:
		return 0
	}
}
