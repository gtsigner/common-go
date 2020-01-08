package convert

import "strconv"

func Int2String(val int) string {
	return strconv.Itoa(val)
}

func Int642String(val int64) string {
	return strconv.FormatInt(val, 10)
}
