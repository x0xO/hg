package hg

import (
	"bufio"
	"bytes"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

// NewHFile returns a new HFile instance with the given name.
func NewHFile(name HString) *HFile { return &HFile{name: name} }

// SeekToLine moves the file pointer to the specified line number and reads the
// specified number of lines from that position.
// The function returns the new position and a concatenation of the lines read as an HString.
//
// Parameters:
//
// - position int64: The starting position in the file to read from
//
// - linesRead int: The number of lines to read.
//
// Returns:
//
// - int64: The new file position after reading the specified number of lines
//
// - HString: A concatenation of the lines read as an HString.
//
// Example usage:
//
//	hf := hg.NewHFile("file.txt")
//	position, content := hf.SeekToLine(0, 5) // Read 5 lines from the beginning of the file
func (hf *HFile) SeekToLine(position int64, linesRead int) (int64, HString) {
	iterator := hf.Iterator()

	if _, err := hf.file.Seek(position, 0); err != nil {
		hf.err = err
		return 0, ""
	}

	var (
		content     strings.Builder
		linesReaded int
	)

	for line := iterator.Lines(); line.Next(); linesReaded++ {
		if linesReaded == linesRead {
			hf.Close()
			break
		}

		content.WriteString(line.HString().Add("\n").String())
		position += int64(line.HBytes().Len() + 1) // Add 1 for the newline character
	}

	return position, HString(content.String())
}

// TempFile creates a new temporary file in the specified directory with the
// specified name pattern, and returns a pointer to the HFile.
// If no directory is specified, the default directory for temporary files is used.
// If no name pattern is specified, the default pattern "*" is used.
//
// Parameters:
//
// - args ...HString: A variadic parameter specifying the directory and/or name
// pattern for the temporary file.
//
// Returns:
//
// - *HFile: A pointer to the HFile representing the temporary file.
//
// Example usage:
//
//	hf := hg.NewHFile("")
//	tmpfile := hf.TempFile()                     // Creates a temporary file with default settings
//	tmpfileWithDir := hf.TempFile("mydir")       // Creates a temporary file in "mydir" directory
//	tmpfileWithPattern := hf.TempFile("", "tmp") // Creates a temporary file with "tmp" pattern
func (hf *HFile) TempFile(args ...HString) *HFile {
	dir := ""
	pattern := "*"

	if len(args) != 0 {
		if len(args) > 1 {
			pattern = args[1].String()
		}

		dir = args[0].String()
	}

	tmpfile, err := os.CreateTemp(dir, pattern)
	if err != nil {
		hf.err = err
		return hf
	}

	htmpfile := NewHFile(HString(tmpfile.Name()))
	htmpfile.file = tmpfile

	defer htmpfile.Close()

	*hf = *htmpfile

	return hf
}

// Append appends the given content to the file, with the specified mode (optional).
// If no FileMode is provided, the default FileMode (0644) is used.
func (hf *HFile) Append(content HString, mode ...os.FileMode) *HFile {
	if hf.file == nil {
		hf.mkdirAll()

		fmode := FileDefault
		if len(mode) != 0 {
			fmode = mode[0]
		}

		file, err := os.OpenFile(hf.name.String(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, fmode)
		if err != nil {
			hf.err = err
			return hf
		}

		hf.file = file
	}

	_, hf.err = hf.file.WriteString(content.String())

	return hf
}

// Write writes the given content to the file, with the specified mode (optional).
// If no FileMode is provided, the default FileMode (0644) is used.
func (hf *HFile) Write(content HString, mode ...os.FileMode) *HFile {
	if hf.file == nil {
		hf.mkdirAll()
	} else {
		hf.Close()
	}

	fmode := FileDefault
	if len(mode) != 0 {
		fmode = mode[0]
	}

	hf.err = os.WriteFile(hf.name.String(), content.Bytes(), fmode)

	return hf
}

// MimeType returns the MIME type of the file as an HString.
func (hf *HFile) MimeType() HString {
	hf.open()
	defer hf.Close()

	const bufferSize = 512

	buff := make([]byte, bufferSize)

	bytesRead, err := hf.file.ReadAt(buff, 0)
	if err != nil && err != io.EOF {
		hf.err = err
		return ""
	}

	buff = buff[:bytesRead]

	return HString(http.DetectContentType(buff))
}

// Read reads the content of the file and returns it as an HString.
func (hf *HFile) Read() HString {
	var content []byte
	content, hf.err = os.ReadFile(hf.name.String())

	return HString(content)
}

// Copy copies the file to the specified destination, with the specified mode (optional).
// If no mode is provided, the default FileMode (0644) is used.
func (hf *HFile) Copy(dest HString, mode ...os.FileMode) *HFile {
	hf.open()
	defer hf.Close()

	return NewHFile(dest).WriteFromReader(hf.file, mode...)
}

// WriteFromReader takes an io.Reader (scr) as input and writes the data from the reader into the file.
// If no FileMode is provided, the default FileMode (0644) is used.
func (hf *HFile) WriteFromReader(scr io.Reader, mode ...os.FileMode) *HFile {
	if hf.file == nil {
		hf.mkdirAll()
	} else {
		hf.Close()
	}

	fmode := FileDefault
	if len(mode) != 0 {
		fmode = mode[0]
	}

	file, err := os.OpenFile(hf.filePath().String(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, fmode)
	if err != nil {
		hf.err = err
		return hf
	}

	defer file.Close()

	_, err = io.Copy(file, scr)
	if err != nil {
		hf.err = err
		return hf
	}

	err = file.Sync()
	if err != nil {
		hf.err = err
		return hf
	}

	return hf
}

// Iterator returns a new hfiter instance that can be used to read the file
// line by line, word by word, rune by rune, or byte by byte.
//
// Returns:
//
// - *hfiter: A pointer to the new hfiter instance.
//
// Example usage:
//
//	hf := hg.NewHFile("file.txt")
//	iterator := hf.Iterator() // Returns a new hfiter instance for the file
func (hf *HFile) Iterator() *hfiter {
	hf.open().hfiter = &hfiter{
		scanner: bufio.NewScanner(hf.file),
		hfile:   hf,
	}

	return hf.hfiter
}

func (hf *HFile) mkdirAll() *HFile {
	if hf.dirPath() != "" && hf.err == nil && !hf.Exist() {
		hf.err = os.MkdirAll(hf.dirPath().String(), DirDefault)
	}

	return hf
}

// Rename renames the file to the specified new path.
func (hf *HFile) Rename(newpath HString) *HFile {
	if !hf.Exist() {
		return hf
	}

	nhf := NewHFile(newpath).mkdirAll()

	err := os.Rename(hf.name.String(), newpath.String())
	if err != nil {
		hf.err = err
		return hf
	}

	*hf = *nhf

	return hf
}

// Exist checks if the file exists.
func (hf HFile) Exist() bool {
	if hf.dirPath() != "" && hf.err == nil {
		_, err := os.Stat(hf.filePath().String())
		return !os.IsNotExist(err)
	}

	return false
}

// HDir returns the directory the file is in as an HDir instance.
func (hf HFile) HDir() *HDir {
	if hf.Exist() {
		return NewHDir(hf.dirPath())
	}

	return NewHDir("")
}

// Split splits the file path into its directory and file components.
func (hf HFile) Split() (*HDir, *HFile) {
	dir, file := filepath.Split(hf.Path().String())
	return NewHDir(HString(dir)), NewHFile(HString(file))
}

// Path returns the absolute path of the file.
func (hf HFile) Path() HString {
	if hf.Exist() {
		return hf.filePath()
	}

	return ""
}

// filePath returns the full file path, including the directory and file name.
func (hf HFile) filePath() HString {
	return HString(filepath.Join(hf.dirPath().String(), filepath.Base(hf.name.String())))
}

// dirPath returns the absolute path of the directory containing the file.
func (hf *HFile) dirPath() HString {
	var path string
	path, hf.err = filepath.Abs(filepath.Dir(hf.name.String()))

	return HString(path)
}

// stat returns the fs.FileInfo of the file.
// It calls the file's Stat method if the file is open, or os.Stat otherwise.
func (hf *HFile) stat() fs.FileInfo {
	var stats fs.FileInfo

	if hf.file != nil {
		stats, hf.err = hf.file.Stat()
	} else {
		stats, hf.err = os.Stat(hf.name.String())
	}

	return stats
}

// Chmod changes the mode of the file.
func (hf *HFile) Chmod(mode os.FileMode) *HFile {
	if hf.file != nil {
		hf.err = hf.file.Chmod(mode)
	} else {
		hf.err = os.Chmod(hf.name.String(), mode)
	}

	return hf
}

// Chown changes the owner of the file.
func (hf *HFile) Chown(uid, gid int) *HFile {
	if hf.file != nil {
		hf.err = hf.file.Chown(uid, gid)
	} else {
		hf.err = os.Chown(hf.name.String(), uid, gid)
	}

	return hf
}

// open opens the HFile's underlying file for reading.
func (hf *HFile) open() *HFile {
	hf.file, hf.err = os.Open(hf.name.String())
	return hf
}

// Close closes the HFile's underlying file, if it is not already closed.
func (hf *HFile) Close() {
	if hf.file != nil {
		hf.err = hf.file.Close()
	}
}

// // panicOnError panics if there is an error stored in the HFile instance.
// func (hf *HFile) panicOnError() *HFile {
// 	if hf.err != nil {
// 		panic(hf.err)
// 	}

// 	return hf
// }

// Remove removes the file.
func (hf *HFile) Remove() { hf.err = os.Remove(hf.name.String()) }

// Error returns the latest error that occurred during an operation.
func (hf HFile) Error() error { return hf.err }

// Ext returns the file extension.
func (hf HFile) Ext() HString { return HString(filepath.Ext(hf.name.String())) }

// File returns the underlying *os.File instance.
func (hf HFile) File() *os.File { return hf.file }

// IsDir checks if the file is a directory.
func (hf HFile) IsDir() bool { return hf.stat().IsDir() }

// ModTime returns the modification time of the file.
func (hf HFile) ModTime() time.Time { return hf.stat().ModTime() }

// Mode returns the file mode.
func (hf HFile) Mode() fs.FileMode { return hf.stat().Mode() }

// Name returns the name of the file.
func (hf HFile) Name() string { return hf.file.Name() }

// ReadLines reads the file and returns its content as a slice of lines.
func (hf HFile) ReadLines() HSlice[HString] { return hf.Read().Split("\n") }

// Size returns the size of the file.
func (hf HFile) Size() int64 { return hf.stat().Size() }

// Bytes sets the iterator to read the file byte by byte.
//
// Returns:
//
// - An hfiter instance with the scanner configured to read the file byte by byte.
//
// Example usage:
//
//	myHFile := hg.NewHFile("path/to/myfile.txt")
//	defer myHFile.Close()
//
//	iterator := myHFile.Iterator().Bytes()
//
//	for iterator.Next() {
//	    fmt.Println(iterator.HString())
//	}
func (hfit hfiter) Bytes() hfiter {
	hfit.By(func(data []byte, atEOF bool) (int, []byte, error) {
		if atEOF && len(data) == 0 {
			hfit.hfile.Close()
			return 0, nil, nil
		}

		return 1, data[0:1], nil
	})

	return hfit
}

// Lines sets the iterator to read the file line by line.
//
// Returns:
//
// - An hfiter instance with the scanner configured to read the file line by line.
//
// Example usage:
//
//	myHFile := hg.NewHFile("path/to/myfile.txt")
//	defer myHFile.Close()
//
//	iterator := myHFile.Iterator().Lines()
//
//	for iterator.Next() {
//	    fmt.Printf("%c", iterator.HBytes())
//	}
func (hfit hfiter) Lines() hfiter {
	dropCR := func(data []byte) []byte {
		if len(data) > 0 && data[len(data)-1] == '\r' {
			return data[0 : len(data)-1]
		}

		return data
	}

	hfit.By(func(data []byte, atEOF bool) (int, []byte, error) {
		if atEOF && len(data) == 0 {
			hfit.hfile.Close()
			return 0, nil, nil
		}

		if i := bytes.IndexByte(data, '\n'); i >= 0 {
			return i + 1, dropCR(data[0:i]), nil
		}

		if atEOF {
			return len(data), dropCR(data), nil
		}

		return 0, nil, nil
	})

	return hfit
}

// Words sets the iterator to read the file word by word.
//
// Returns:
//
// - An hfiter instance with the scanner configured to read the file word by word.
//
// Example usage:
//
//	myHFile := hg.NewHFile("path/to/myfile.txt")
//	defer myHFile.Close()
//
//	iterator := myHFile.Iterator().Words()
//
//	for iterator.Next() {
//	    fmt.Println(iterator.HString())
//	}
func (hfit hfiter) Words() hfiter {
	hfit.By(func(data []byte, atEOF bool) (int, []byte, error) {
		if atEOF && len(data) == 0 {
			hfit.hfile.Close()
			return 0, nil, nil
		}

		start := 0
		for width := 0; start < len(data); start += width {
			var r rune
			r, width = utf8.DecodeRune(data[start:])
			if !unicode.IsSpace(r) {
				break
			}
		}

		for width, i := 0, start; i < len(data); i += width {
			var r rune
			r, width = utf8.DecodeRune(data[i:])
			if unicode.IsSpace(r) {
				return i + width, data[start:i], nil
			}
		}

		if atEOF && len(data) > start {
			return len(data), data[start:], nil
		}

		return start, nil, nil
	})

	return hfit
}

// Runes sets the iterator to read the file rune by rune.
//
// Returns:
//
// - An hfiter instance with the scanner configured to read the file rune by rune.
//
// Example usage:
//
//	myHFile := hg.NewHFile("path/to/myfile.txt")
//	defer myHFile.Close()
//
//	iterator := myHFile.Iterator().Runes()
//
//	for iterator.Next() {
//	    fmt.Printf("%c", iterator.HString())
//	}
func (hfit hfiter) Runes() hfiter {
	hfit.By(func(data []byte, atEOF bool) (int, []byte, error) {
		if atEOF && len(data) == 0 {
			hfit.hfile.Close()
			return 0, nil, nil
		}

		if data[0] < utf8.RuneSelf {
			return 1, data[0:1], nil
		}

		_, width := utf8.DecodeRune(data)
		if width > 1 {
			return width, data[0:width], nil
		}

		if !atEOF && !utf8.FullRune(data) {
			return 0, nil, nil
		}

		return 1, []byte(string(utf8.RuneError)), nil
	})

	return hfit
}

// By configures the hfiter instance's scanner to use a custom split function.
//
// The custom split function should take a byte slice and a boolean indicating whether this is the
// end of file.
// It should return the advance count, the token, and any encountered error.
//
// Parameters:
//
// - f: A split function of the form func(data []byte, atEOF bool) (advance int, token []byte, err
// error).
//
// Returns:
//
// - An hfiter instance with the scanner configured to use the provided custom split function.
//
// Example usage:
//
//	myHFile := hg.NewHFile("path/to/myfile.txt")
//
//	customSplitFunc := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
//	    // Custom implementation here
//	}
//
//	iterator := myHFile.Iterator().By(customSplitFunc)
//
//	for iterator.Next() {
//	    fmt.Printf("%s", iterator.HString())
//	}
func (hfit hfiter) By(
	f func(data []byte, atEOF bool) (advance int, token []byte, err error),
) hfiter {
	hfit.scanner.Split(bufio.SplitFunc(f))
	return hfit
}

// Buffer sets the initial buffer to use when iterating and the maximum size of the buffer.
//
// By default, Iterator uses an internal buffer and grows it as large as necessary.
// This method allows you to use a custom buffer and limit its size.
//
// Parameters:
//
// - buf: A byte slice that will be used as a buffer.
//
// - max: The maximum size of the buffer.
//
// Example usage:
//
//	myHFile := hg.NewHFile("path/to/myfile.txt")
//	defer myHFile.Close()
//
//	iterator := myHFile.Iterator().Runes()
//
//	customBuffer := make([]byte, 1024)
//	iterator.Buffer(customBuffer, 4096)
//
//	for iterator.Next() {
//	    fmt.Printf("%c", iterator.HString())
//	}
func (hfit hfiter) Buffer(buf []byte, max int) { hfit.scanner.Buffer(buf, max) }

// Next advances the iterator to the next item (byte, line, word, or rune) and
// returns true if successful or false if there are no more items to read.
//
// Returns:
//
// - A boolean value indicating whether the iterator successfully advanced to the next item.
// Returns false if there are no more items to read.
//
// Example usage:
//
//	myHFile := hg.NewHFile("path/to/myfile.txt")
//	defer myHFile.Close()
//
//	iterator := myHFile.Iterator()
//	lines := iterator.Lines() // or iterator.Words() or iterator.Runes()
//
//	for lines.Next() {
//	    fmt.Println(iterator.HString())
//	}
func (hfit hfiter) Next() bool { return hfit.scanner.Scan() }

// HBytes returns the current item as an HBytes instance.
//
// Returns:
//
// - An HBytes instance containing the current item in the iterator.
//
// Example usage:
//
//	myHFile := hg.NewHFile("path/to/myfile.txt")
//	defer myHFile.Close()
//
//	iterator := myHFile.Iterator().Bytes() // Sets the iterator to read the file byte by byte
//
//	for iterator.Next() {
//	    fmt.Println(iterator.HBytes())
//	}
func (hfit hfiter) HBytes() HBytes { return HBytes(hfit.scanner.Bytes()) }

// HString returns the current item as an HString instance.
//
// Returns:
//
// - An HString instance containing the current item in the iterator.
//
// Example usage:
//
//	myHFile := hg.NewHFile("path/to/myfile.txt")
//	defer myHFile.Close()
//
//	iterator := myHFile.Iterator().Lines() // Sets the iterator to read the file line by line
//
//	for iterator.Next() {
//	    fmt.Println(iterator.HString())
//	}
func (hfit hfiter) HString() HString { return HString(hfit.scanner.Text()) }

// Error returns the first non-EOF error encountered by the Iterator.
//
// Call this method after an iteration loop has finished to check if any errors occurred
// during the iteration process.
//
// Returns:
//
// - An error encountered by the iterator, or nil if no errors occurred.
//
// Example usage:
//
//	myHFile := hg.NewHFile("path/to/myfile.txt")
//	defer myHFile.Close()
//
//	iterator := myHFile.Iterator().Runes()
//
//	for iterator.Next() {
//	    fmt.Printf("%c", iterator.HString())
//	}
//
//	if err := iterator.Error(); err != nil {
//	    log.Printf("Error while iterating: %v", err)
//	}
func (hfit hfiter) Error() error { return hfit.scanner.Err() }
