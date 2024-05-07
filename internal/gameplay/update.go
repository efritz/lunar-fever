package gameplay

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

// TODO
const targetBinaryPath = "/Users/efritz/dev/efritz/lunar-fever/lunar-fever-macos-arm64"
const targetChecksumPath = "/Users/efritz/dev/efritz/lunar-fever/lunar-fever-macos-arm64.sum256"

func HasUpdate() (bool, error) {
	contents, err := os.ReadFile(targetChecksumPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}
	expectedSum := strings.TrimSpace(string(contents))

	s, err := ChecksumPath(os.Args[0])
	if err != nil {
		return false, err
	}

	return expectedSum != s, nil
}

func ChecksumPath(path string) (string, error) {
	r, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer r.Close()

	return ChecksumReader(r)
}

func ChecksumReader(r io.Reader) (string, error) {
	sum := sha256.New()
	if _, err := io.Copy(sum, r); err != nil {
		return "", err
	}

	return hex.EncodeToString(sum.Sum(nil)), nil
}

func Update() error {
	contents, err := os.ReadFile(targetBinaryPath)
	if err != nil {
		return err
	}

	checksum, err := ChecksumReader(bytes.NewReader(contents))
	if err != nil {
		return err
	}
	_ = checksum // TODO - verify

	targetPath, err := filepath.Abs(os.Args[0])
	if err != nil {
		return err
	}

	updateDir := filepath.Dir(targetPath)
	filename := filepath.Base(targetPath)
	newPath := filepath.Join(updateDir, fmt.Sprintf(".%s.new", filename))
	oldPath := filepath.Join(updateDir, fmt.Sprintf(".%s.old", filename))

	fmt.Printf("writing new binary %s\n", newPath)
	fp, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer fp.Close()

	if _, err := io.Copy(fp, bytes.NewReader(contents)); err != nil {
		return err
	}

	// remove rename target of next rename
	fmt.Printf("removing %s\n", oldPath)
	_ = os.Remove(oldPath)

	// Back up current binary
	fmt.Printf("renaming %s -> %s\n", targetPath, oldPath)
	if err := os.Rename(targetPath, oldPath); err != nil {
		return err
	}

	fmt.Printf("renaming %s -> %s\n", newPath, targetPath)
	if err := os.Rename(newPath, targetPath); err != nil {
		// TODO - try to rollback
		return err
	}

	// TODO - remove old binary

	return syscall.Exec(targetPath, append([]string{targetPath}, os.Args[1:]...), os.Environ())
}
