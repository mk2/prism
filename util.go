package prism

import "strconv"

func s2b(str string) []byte {
	return []byte(str)
}

func b2s(buf []byte) string {

	if buf == nil {
		return ""
	}

	n := len(buf)
	return string(buf[:n])
}

func b2i(buf []byte) int {

	s := b2s(buf)
	i, _ := strconv.Atoi(s)

	return i
}

func i2b(i int) []byte {

	s := strconv.Itoa(i)
	buf := s2b(s)

	return buf
}

func i2s(i int) string {
	return strconv.Itoa(i)
}

func s2i(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
