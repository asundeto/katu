package handlers

import (
	"crypto/sha256"
	"fmt"
)

func Hash256(s string) string {
	s += "H^3D-(3@17" // SALT
	h := sha256.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x\n", bs)
}
