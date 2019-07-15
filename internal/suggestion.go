package Internal

import (
	"strings"
)

//Suggestion func for matching the text with the availables
func Suggestion(text1 string, text2 string) bool {
	array1 := strings.Fields(text1)
	array2 := strings.Fields(text2)
	count := len(array2)
	count = count / 2
	var total int
	for _, a := range array1 {
		for _, b := range array2 {
			if a == b {
				total = total + 1
			}
		}
	}
	if total > count {
		return true
	}
	return false
}
