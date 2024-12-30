package menu

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"syscall"
)

func NewUpdater() (*updater, error) {
	content, err := readFileFromServer("manifest.json")
	if err != nil {
		return nil, err
	}

	manifest := struct{ Files []string }{}
	if err := json.Unmarshal(content, &manifest); err != nil {
		return nil, err
	}

	expectedPrefix := fmt.Sprintf("lunar-fever-%s-%s", runtime.GOOS, runtime.GOARCH)

	var matches []string
	for _, file := range manifest.Files {
		if strings.HasPrefix(file, expectedPrefix) {
			matches = append(matches, file)
		}
	}

	if len(matches) != 2 {
		return &updater{}, nil
	}

	sort.Strings(matches)

	checksum, err := readFileFromServer(manifest.Files[1])
	if err != nil {
		return nil, err
	}

	return &updater{
		binaryBasename: manifest.Files[0],
		updateChecksum: strings.TrimSpace(string(checksum)),
	}, nil
}

type updater struct {
	binaryBasename string
	updateChecksum string
}

func (u *updater) HasUpdate() (bool, error) {
	if u.updateChecksum == "" {
		return false, nil
	}

	binaryChecksum, err := checksumFromPath(os.Args[0])
	if err != nil {
		return false, err
	}

	return binaryChecksum != u.updateChecksum, nil
}

func (u *updater) Update() error {
	contents, err := u.getBinary()
	if err != nil {
		return err
	}

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

func (u *updater) getBinary() ([]byte, error) {
	contents, err := readFileFromServer(u.binaryBasename)
	if err != nil {
		return nil, err
	}

	s, err := checksumFromReader(bytes.NewReader(contents))
	if err != nil {
		return nil, err
	}
	if s != u.updateChecksum {
		return nil, fmt.Errorf("mismatched binary checksum")
	}

	return contents, nil
}

const baseURL = "http://lunarfever.com/assets/"

func readFileFromServer(filename string) ([]byte, error) {
	resp, err := http.Get(baseURL + filename)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func checksumFromPath(path string) (string, error) {
	r, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer r.Close()

	return checksumFromReader(r)
}

func checksumFromReader(r io.Reader) (string, error) {
	sum := sha256.New()
	if _, err := io.Copy(sum, r); err != nil {
		return "", err
	}

	return hex.EncodeToString(sum.Sum(nil)), nil
}
