package details

import "unicode"

func filter(str string) string {
	var newStr string
	for _, s := range str {
		if unicode.IsLetter(s) || unicode.IsNumber(s) || s == ' ' {
			newStr += string(s)
		}
	}
	return newStr
}

func saveDetails() {

}
