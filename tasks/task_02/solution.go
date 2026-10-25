package main

func rotateRunes(s string, shift int) string {
	runes := []rune(s)
	n := len(runes)
	if n == 0 {
		return ""
	}

	k := shift % n
	if k < 0 {
		k += n
	}
	return string(runes[k:]) + string(runes[:k])
}
