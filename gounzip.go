package gounzip

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Options are all optional values that can be passed to Unzip
type Options struct {
	// Destination directory of the unzipped contents, If the directory does not exist it will be created.
	// Defaults to match path and name of the source zip, without the file extension
	Destination string

	// Number of leading components to strip from file paths (same as tar --strip-components)
	// Defaults to 0
	StripComponents uint
}

// Option is a utility type for setting options
type Option func(*Options)

func Destination(dest string) Option {
	return func(args *Options) {
		args.Destination = dest
	}
}

func StripComponents(n uint) Option {
	return func(args *Options) {
		args.StripComponents = n
	}
}

func Unzip(src string, options ...Option) error {
	b := filepath.Base(src)

	// Defaults
	args := &Options{

		// The default output destination will be a folder adjacent to the target file
		// with the same name as the target, without the file extension.
		// Eg, a src '~/work/project/file.zip' would have a default destination of '~/work/project/file/'
		Destination: filepath.Join(filepath.Dir(src), strings.TrimSuffix(b, filepath.Ext(b))),

		//
		StripComponents: 1,
	}

	for _, option := range options {
		option(args)
	}

	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(args.Destination, 0755)

	for _, f := range r.File {
		err := extractAndWriteFile(f, args.Destination, args.StripComponents)
		if err != nil {
			return err
		}
	}

	return nil
}

// Closure to address file descriptors issue with all the deferred .Close() methods
func extractAndWriteFile(f *zip.File, dest string, c uint) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer func() {
		if err := rc.Close(); err != nil {
			panic(err)
		}
	}()

	s := strings.Split(f.Name, string(os.PathSeparator))
	fc := s[c:]
	if len(fc) == 0 {
		return nil
	}

	fp := strings.Join(fc, string(os.PathSeparator))
	path := filepath.Join(dest, fp)

	if f.FileInfo().IsDir() {
		os.MkdirAll(path, f.Mode())
	} else {
		os.MkdirAll(filepath.Dir(path), f.Mode())
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer func() {
			if err := f.Close(); err != nil {
				panic(err)
			}
		}()

		_, err = io.Copy(f, rc)
		if err != nil {
			return err
		}
	}
	return nil
}
