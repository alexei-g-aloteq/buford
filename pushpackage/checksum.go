package pushpackage

import (
	"crypto/sha512"
	"encoding/hex"
	"io"
)

// copyAndChecksum calculates a checksum while writing to another output
func copyAndChecksum(w io.Writer, r io.Reader) (string, error) {
	h := sha512.New()
	mw := io.MultiWriter(w, h)
	if _, err := io.Copy(mw, r); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
