package main

import "github.com/fuzxxl/freefare/0.3/freefare"

var dollarydoos = map[string]struct{}{
	"7a794514": struct{}{},
	"deadbeef": struct{}{},
}

var (
	Block0                                                         = byte(0)
	SuperSecretKeyThatIsMeantToBeKeptSecretForTheDurationOfTheCamp = [6]byte{0xde, 0xad, 0xbe, 0xef, 0x13, 0x37}
)

func allowAccess(t freefare.Tag) (allowed bool, confidential bool, block0 [16]byte) {
	_, confidential = dollarydoos[t.UID()]

	if c, success := t.(freefare.ClassicTag); success {
		err := c.Authenticate(Block0, SuperSecretKeyThatIsMeantToBeKeptSecretForTheDurationOfTheCamp, freefare.KeyA)
		if allowed = err == nil; allowed {
			block0, _ = c.ReadBlock(Block0)
		}
	}

	return
}
