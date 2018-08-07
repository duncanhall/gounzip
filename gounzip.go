package gounzip

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Options passed to the Unzip method
type Options struct {
	// Destination directory of the unzipped contents.
	// If the directory does not exist it will be created.
	// Defaults to match path and name of the source zip, without the file extension
	Destination string

	// Number of leading components to strip from file paths (same as tar --strip-components)
	// Defaults to 0
	StripComponents uint

	// Whether the original zip should be removed after unzipping
	// Defaults to false
	DeleteSource bool
}

// Option is a setter for assigning values to Options
type Option func(*Options)

// Destination creates a setter for Options.Destination values
func Destination(dest string) Option {
	return func(args *Options) {
		args.Destination = dest
	}
}

// StripComponents creates a setter for Options.StripComponents values
func StripComponents(n uint) Option {
	return func(args *Options) {
		args.StripComponents = n
	}
}

// DeleteSource creates a setter for Options.DeleteSource values
func DeleteSource(d bool) Option {
	return func(args *Options) {
		args.DeleteSource = d
	}
}

// Unzip accepts a path to a compressed .zip file and outputs the contents
// Optional arguments are passed in as Option setters
func Unzip(src string, options ...Option) error {
	b := filepath.Base(src)

	// Set deafult option values
	args := &Options{
		// The default output destination will be a folder adjacent to the target file
		// with the same name as the target, without the file extension.
		// Eg, a src '~/work/project/file.zip' would have a default destination of '~/work/project/file/'
		Destination: filepath.Join(filepath.Dir(src), strings.TrimSuffix(b, filepath.Ext(b))),

		// The default StripComponent value is left at 0, which will not strip anything
		StripComponents: 0,

		// Don't delete source files by default
		DeleteSource: false,
	}

	// Apply any user set options to override defaults
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

	// Each file in the .zip is enumerated regardless of nested directories
	// meaning range here will iterate over every file and give us its
	// fully nested path
	for _, f := range r.File {
		err := extractFile(f, args.Destination, args.StripComponents)
		if err != nil {
			return err
		}
	}

	if args.DeleteSource {
		err := os.Remove(src)
		if err != nil {
			return err
		}
	}

	return nil
}

// extractFile outputs the contents of a zip file to a given destination
func extractFile(f *zip.File, dest string, c uint) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer func() {
		if err := rc.Close(); err != nil {
			panic(err)
		}
	}()

	// Split ihe zip file path into its components
	s := strings.Split(f.Name, string(os.PathSeparator))
	// Strip the given number of components
	fc := s[c:]
	if len(fc) == 0 {
		// If there are no components left, skip the file for output
		return nil
	}

	// Recreate the file path with stripped components removed
	fp := strings.Join(fc, string(os.PathSeparator))
	path := filepath.Join(dest, fp)
	ogMode := f.Mode()

	if f.FileInfo().IsDir() {
		os.MkdirAll(path, ogMode)
	} else {
		os.MkdirAll(filepath.Dir(path), ogMode)
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, ogMode)
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
		if err = os.Chmod(f.Name, ogMode); err != nil {
			return err
		}
	}
	return nil
}
