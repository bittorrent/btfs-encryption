package utils

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func TarGiz(source string) (out []byte, err error) {
	base := filepath.Base(source)

	var buf bytes.Buffer
	defer func() {
		if err != nil {
			return
		}
		out = buf.Bytes()
	}()

	zr := gzip.NewWriter(&buf)
	defer func() {
		if err != nil {
			return
		}
		err = zr.Close()
	}()

	tw := tar.NewWriter(zr)
	defer func() {
		if err != nil {
			return
		}
		err = tw.Close()
	}()

	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Mode()&os.ModeSymlink == os.ModeSymlink {
			return nil
		}

		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}

		header.Name = filepath.Join(base, strings.TrimPrefix(path, source))
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		defer file.Close()

		_, err = io.Copy(tw, file)
		return err
	})

	return
}

func UnTarGiz(content []byte, target string) error {
	gr, err := gzip.NewReader(bytes.NewBuffer(content))
	if err != nil {
		return err
	}
	defer gr.Close()

	tr := tar.NewReader(gr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		path := filepath.Join(target, header.Name)
		info := header.FileInfo()

		if info.IsDir() {
			err = os.MkdirAll(path, info.Mode())
			if err != nil {
				return err
			}
			continue
		}

		err = func() (err error) {
			file, err := os.OpenFile(
				path,
				os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
				info.Mode(),
			)
			if err != nil {
				return
			}
			defer file.Close()
			_, err = io.Copy(file, tr)
			return
		}()
		if err != nil {
			return err
		}

	}

	return nil
}
