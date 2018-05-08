package fsutils

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bogem/gel/pools"
)

func CreateMissingDir(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			if err = os.MkdirAll(path, 0755); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func ChangeExt(path, newExt string) string {
	return DeleteExt(path) + newExt
}

func DeleteExt(path string) string {
	return strings.TrimSuffix(path, filepath.Ext(path))
}

// Copied from https://github.com/otiai10/copy.
func Copy(dst, src string) error {
	fi, err := os.Stat(src)
	if err != nil {
		return err
	}
	return cp(dst, src, fi)
}

func cp(dst, src string, srcFI os.FileInfo) error {
	if srcFI.IsDir() {
		return cpDir(dst, src, srcFI)
	}
	return cpFile(dst, src, srcFI)
}

func cpFile(dst, src string, srcFI os.FileInfo) error {
	dstFile, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, srcFI.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	buf := pools.GetByteSlice(128 * 1024)
	_, err = io.CopyBuffer(dstFile, srcFile, buf)
	pools.PutByteSlice(buf)

	return err
}

func cpDir(dst, src string, srcFI os.FileInfo) error {
	if err := os.MkdirAll(dst, srcFI.Mode()); err != nil {
		return err
	}

	fis, err := readDir(src)
	if err != nil {
		return err
	}

	for _, fi := range fis {
		err := cp(filepath.Join(dst, fi.Name()), filepath.Join(src, fi.Name()), fi)
		if err != nil {
			return err
		}
	}

	return nil
}

// Copy of ioutil.ReadDir but without sort.
func readDir(dirname string) ([]os.FileInfo, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return f.Readdir(-1)
}
