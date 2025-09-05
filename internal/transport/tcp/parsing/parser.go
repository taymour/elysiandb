package parsing

import (
	"fmt"
)

func FirstWordBytes(b []byte) (cmd, rest []byte) {
	end := len(b)
	for end > 0 && (b[end-1] == '\n' || b[end-1] == '\r') {
		end--
	}
	b = b[:end]

	i := 0
	for i < len(b) && (b[i] == ' ' || b[i] == '\t') {
		i++
	}
	start := i

	for i < len(b) && b[i] != ' ' && b[i] != '\t' {
		i++
	}
	cmd = b[start:i]

	for i < len(b) && (b[i] == ' ' || b[i] == '\t') {
		i++
	}
	rest = b[i:]
	return
}

func EqASCII(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		ca, cb := a[i], b[i]
		if 'a' <= ca && ca <= 'z' {
			ca -= 'a' - 'A'
		}
		if 'a' <= cb && cb <= 'z' {
			cb -= 'a' - 'A'
		}
		if ca != cb {
			return false
		}
	}
	
	return true
}

func ParseDecimalBytes(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, fmt.Errorf("empty input")
	}

	maxInt := int(^uint(0) >> 1)

	n := 0
	read := 0
	for _, c := range b {
		if c < '0' || c > '9' {
			break
		}

		d := int(c - '0')

		if n > maxInt/10 || (n == maxInt/10 && d > maxInt%10) {
			return 0, fmt.Errorf("integer overflow")
		}

		n = n*10 + d
		read++
	}

	if read == 0 {
		return 0, fmt.Errorf("no digits found")
	}

	return n, nil
}
