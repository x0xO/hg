package hg

import (
	"os"
	"path/filepath"
)

// NewHDir returns a new HDir instance with the given path.
func NewHDir(path HString) *HDir { return &HDir{path: path} }

// CopyDir copies the contents of the current directory to the destination directory.
//
// Parameters:
//
// - dest (HString): The destination directory where the contents of the current
// directory should be copied.
//
// Returns:
//
// - *HDir: A pointer to a new HDir instance representing the destination directory.
//
// Example usage:
//
//	sourceDir := hg.NewHDir("path/to/source")
//	destinationDir := sourceDir.CopyDir("path/to/destination")
func (hd *HDir) CopyDir(dest HString) *HDir {
	if err := filepath.Walk(hd.path.String(), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(hd.path.String(), path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(dest.String(), relPath)

		if info.IsDir() {
			dPath := NewHDir(HString(destPath)).MkdirAll(info.Mode())
			return dPath.err
		}

		hPath := NewHFile(HString(path)).Copy(HString(destPath), info.Mode())
		return hPath.err
	}); err != nil {
		hd.err = err
		return hd
	}

	return NewHDir(dest)
}

// Mkdir creates a new directory with the specified mode (optional).
//
// Parameters:
//
// - mode (os.FileMode, optional): The file mode for the new directory.
// If not provided, it defaults to DirDefault (0755).
//
// Returns:
//
// - *HDir: A pointer to the HDir instance on which the method was called.
//
// Example usage:
//
//	dir := hg.NewHDir("path/to/directory")
//	createdDir := dir.Mkdir(0755) // Optional mode argument
func (hd *HDir) Mkdir(mode ...os.FileMode) *HDir {
	dmode := DirDefault
	if len(mode) != 0 {
		dmode = mode[0]
	}

	hd.err = os.Mkdir(hd.path.String(), dmode)

	return hd
}

// Join joins the current directory path with the given path elements, returning the joined path.
//
// Parameters:
//
// - elem (...HString): One or more HString values representing path elements to
// be joined with the current directory path.
//
// Returns:
//
// - HString: The resulting joined path as an HString.
//
// Example usage:
//
//	dir := hg.NewHDir("path/to/directory")
//	joinedPath := dir.Join("subdir", "file.txt")
func (hd *HDir) Join(elem ...HString) HString {
	paths := HSliceOf(elem...).Insert(0, hd.Path()).ToStringSlice()
	return HString(filepath.Join(paths...))
}

// SetPath sets the path of the current directory.
//
// Parameters:
//
// - path (HString): The new path to be set for the current directory.
//
// Returns:
//
// - *HDir: A pointer to the updated HDir instance with the new path.
//
// Example usage:
//
//	dir := hg.NewHDir("path/to/directory")
//	dir.SetPath("new/path/to/directory")
func (hd *HDir) SetPath(path HString) *HDir {
	hd.path = path
	return hd
}

// MkdirAll creates all directories along the given path, with the specified mode (optional).
//
// Parameters:
//
// - mode ...os.FileMode (optional): The file mode to be used when creating the directories.
// If not provided, it defaults to the value of DirDefault constant (0755).
//
// Returns:
//
// - *HDir: A pointer to the HDir instance representing the created directories.
//
// Example usage:
//
//	dir := hg.NewHDir("path/to/directory")
//	dir.MkdirAll()
//	dir.MkdirAll(0755)
func (hd *HDir) MkdirAll(mode ...os.FileMode) *HDir {
	if hd.Exist() {
		return hd
	}

	dmode := DirDefault
	if len(mode) != 0 {
		dmode = mode[0]
	}

	hd.err = os.MkdirAll(hd.Path().String(), dmode)

	return hd
}

// Rename renames the current directory to the new path.
//
// Parameters:
//
// - newpath HString: The new path for the directory.
//
// Returns:
//
// - *HDir: A pointer to the HDir instance representing the renamed directory.
// If an error occurs, the original HDir instance is returned with the error stored in hd.err,
// which can be checked using the Error() method.
//
// Example usage:
//
//	dir := hg.NewHDir("path/to/directory")
//	dir.Rename("path/to/new_directory")
func (hd *HDir) Rename(newpath HString) *HDir {
	err := os.Rename(hd.path.String(), newpath.String())
	if err != nil {
		hd.err = err
		return hd
	}

	return NewHDir(newpath)
}

// Path returns the absolute path of the current directory.
//
// Returns:
//
// - HString: The absolute path of the current directory as an HString.
// If an error occurs while converting the path to an absolute path,
// the error is stored in hd.err, which can be checked using the Error() method.
//
// Example usage:
//
//	dir := hg.NewHDir("path/to/directory")
//	absPath := dir.Path()
func (hd *HDir) Path() HString {
	path, err := filepath.Abs(hd.path.String())
	hd.err = err

	return HString(path)
}

// Exist checks if the current directory exists.
//
// Returns:
//
// - bool: true if the current directory exists, false otherwise.
//
// Example usage:
//
//	dir := hg.NewHDir("path/to/directory")
//	exists := dir.Exist()
func (hd HDir) Exist() bool {
	_, err := os.Stat(hd.Path().String())
	return !os.IsNotExist(err)
}

// ReadDir reads the content of the current directory and returns a slice of HFile instances.
//
// Returns:
//
// - []*HFile: A slice of HFile instances representing the files and directories
// in the current directory.
//
// Example usage:
//
//	dir := hg.NewHDir("path/to/directory")
//	files := dir.ReadDir()
func (hd *HDir) ReadDir() []*HFile {
	dirs, err := os.ReadDir(hd.path.String())
	if err != nil {
		hd.err = err
		return nil
	}

	hfiles := make([]*HFile, 0, len(dirs))

	for _, dir := range dirs {
		path, _ := filepath.Abs(filepath.Join(hd.path.String(), dir.Name()))
		hfiles = append(hfiles, NewHFile(HString(path)))
	}

	return hfiles
}

// Glob matches files in the current directory using the path pattern and
// returns a slice of HFile instances.
//
// Returns:
//
// - []*HFile: A slice of HFile instances representing the files that match the
// provided pattern in the current directory.
//
// Example usage:
//
//	dir := hg.NewHDir("path/to/directory")
//	files := dir.Glob("*.txt")
func (hd *HDir) Glob() []*HFile {
	files, err := filepath.Glob(hd.path.String())
	if err != nil {
		hd.err = err
		return nil
	}

	hfiles := make([]*HFile, 0, len(files))

	for _, file := range files {
		path, _ := filepath.Abs(filepath.Join(hd.path.String(), file))
		hfiles = append(hfiles, NewHFile(HString(path)))
	}

	return hfiles
}

// HString returns the HString representation of the current directory's path.
func (hd HDir) HString() HString { return hd.path }

// Error returns the latest error that occurred during an operation.
func (hd HDir) Error() error { return hd.err }
