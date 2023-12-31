// Copyright (c) 2023, The Goki Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tomls

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"goki.dev/glop/dirs"
	"goki.dev/grows"
)

// Decoder is needed to return a standard Decode function signature for toml.
// Just wraps [toml.Decoder] to satisfy our [grows.Decoder] interface.
// Should not need to use in your code.
type Decoder struct {
	*toml.Decoder
}

// Decode implements the standard [grows.Decoder] signature
func (d *Decoder) Decode(v any) error {
	_, err := d.Decoder.Decode(v)
	return err
}

// NewDecoder returns a new [Decoder]
func NewDecoder(r io.Reader) grows.Decoder { return &Decoder{toml.NewDecoder(r)} }

// Open reads the given object from the given filename using TOML encoding
func Open(v any, filename string) error {
	return grows.Open(v, filename, NewDecoder)
}

// OpenFiles reads the given object from the given filenames using TOML encoding
func OpenFiles(v any, filenames []string) error {
	return grows.OpenFiles(v, filenames, NewDecoder)
}

// OpenFS reads the given object from the given filename using TOML encoding,
// using the given [fs.FS] filesystem (e.g., for embed files)
func OpenFS(v any, fsys fs.FS, filename string) error {
	return grows.OpenFS(v, fsys, filename, NewDecoder)
}

// OpenFilesFS reads the given object from the given filenames using TOML encoding,
// using the given [fs.FS] filesystem (e.g., for embed files)
func OpenFilesFS(v any, fsys fs.FS, filenames []string) error {
	return grows.OpenFilesFS(v, fsys, filenames, NewDecoder)
}

// Read reads the given object from the given reader,
// using TOML encoding
func Read(v any, reader io.Reader) error {
	return grows.Read(v, reader, NewDecoder)
}

// ReadBytes reads the given object from the given bytes,
// using TOML encoding
func ReadBytes(v any, data []byte) error {
	return grows.ReadBytes(v, data, NewDecoder)
}

// Save writes the given object to the given filename using TOML encoding
func Save(v any, filename string) error {
	return grows.Save(v, filename, grows.NewEncoderFunc(toml.NewEncoder))
}

// Write writes the given object using TOML encoding
func Write(v any, writer io.Writer) error {
	return grows.Write(v, writer, grows.NewEncoderFunc(toml.NewEncoder))
}

// WriteBytes writes the given object, returning bytes of the encoding,
// using TOML encoding
func WriteBytes(v any) ([]byte, error) {
	return grows.WriteBytes(v, grows.NewEncoderFunc(toml.NewEncoder))
}

// OpenFromPaths reads object from given TOML file,
// looking on paths for the file.
func OpenFromPaths(obj any, file string, paths []string) error {
	filenames := dirs.FindFilesOnPaths(paths, file)
	if len(filenames) == 0 {
		err := fmt.Errorf("OpenFromPaths: No files found")
		log.Println(err)
		return err
	}
	// _, err = toml.DecodeFile(fp, obj)
	fp, err := os.Open(filenames[0])
	defer fp.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	return Read(obj, bufio.NewReader(fp))
}
