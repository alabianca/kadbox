package cli

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
)

func writeTarball(writer io.Writer, dir string) (int64, error) {
	tw := tar.NewWriter(writer)

	defer tw.Close()

	// walk path
	var bytesWritten int64
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}

		header.Name = path
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		f, err := os.Open(path)
		defer f.Close()
		if err != nil {
			return err
		}

		var n int64
		if n, err = io.Copy(tw, f); err != nil {
			return err
		}
		bytesWritten+=n

		return nil
	})

	return bytesWritten, err
}

// ReadTarball reads from reader and creates the resulting directory at target
func readTarball(reader io.Reader, target string) error {

	tr := tar.NewReader(reader)

	for {
		header, err := tr.Next()

		switch {
		case err == io.EOF:
			return nil

		case err != nil:
			return err

		case header == nil:
			continue

		}

		target := filepath.Join(target, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			//it is a directory. create it if it does not exist
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		case tar.TypeReg:
			//regular file. create it
			f, err := os.Create(header.Name)
			//f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			f.Close()
		}
	}
}

