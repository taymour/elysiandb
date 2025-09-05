package parser

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
