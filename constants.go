package hg

import (
	"os"
)

const (
	ASCII_LETTERS   = ASCII_LOWERCASE + ASCII_UPPERCASE
	ASCII_LOWERCASE = "abcdefghijklmnopqrstuvwxyz"
	ASCII_UPPERCASE = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	DIGITS          = "0123456789"
	HEXDIGITS       = "0123456789abcdefABCDEF"
	OCTDIGITS       = "01234567"
	PUNCTUATION     = `!"#$%&'()*+,-./:;<=>?@[\]^{|}~` + "`"

	FileDefault os.FileMode = 0o644
	DirDefault  os.FileMode = 0o755
	FullAccess  os.FileMode = 0o777
)
