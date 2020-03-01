package cmn

// GetFailMsg - dummy error message
func GetFailMsg() string {
	var items = []string{
		"Семь раз за Пупу и один раз за Лупу",
		"Missing Data...",
		"Optimizing scan results...",
		"FATAL ERROR",
		"CKEY is full",
		"шото пошло не так",
		"не в этот раз",
		"logic_error",
		"no",
		"нет"}

	return GetOneMsgFromMany(items...)
}

// GetOneMsgFromMany - wrapper for taking one msg from a set of
func GetOneMsgFromMany(items ...string) string {
	if len(items) == 0 {
		return ""
	}
	itemNo := Rnd.Intn(len(items))
	return items[itemNo]
}
