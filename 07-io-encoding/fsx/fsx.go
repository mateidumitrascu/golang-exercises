// Package fsx is about io/fs - the interface that made Go's filesystem code
// testable - and about writing files without corrupting them.
//
// Functions that take an fs.FS can be tested with fstest.MapFS: no temp
// directories, no cleanup, no disk. Write your code against fs.FS by default
// and take an os path only where you genuinely must (writing).
package fsx

import (
	"io/fs"
	"time"
)

// FindFiles returns the paths of all regular files under root whose name ends
// in ext (".go", say), sorted lexically. An empty ext matches everything.
// Directories are not returned. Errors from walking are returned, except that
// a permission error on a single directory should be skipped rather than
// aborting the whole walk (fs.SkipDir is not the right tool - think about what
// the WalkDirFunc's err parameter is for).
func FindFiles(fsys fs.FS, root, ext string) ([]string, error) { panic("TODO: implement FindFiles") }

// Match is one search hit.
type Match struct {
	Path string
	Line int    // 1-based
	Text string // the whole line, without its newline
}

// Grep searches every file for a substring and returns the matches in path
// order, then line order. Files that cannot be read are skipped silently;
// binary files (containing a NUL byte in the first 512 bytes) are skipped too.
func Grep(fsys fs.FS, substr string) ([]Match, error) { panic("TODO: implement Grep") }

// Stats summarises a tree.
type Stats struct {
	Files   int
	Dirs    int
	Bytes   int64
	Largest string // path of the largest file, "" if there are none; ties go to the first in walk order
	ByExt   map[string]int
}

// Walk collects statistics in a single pass. ByExt is keyed by the extension
// including the dot; files with no extension are counted under "".
func Walk(fsys fs.FS) (Stats, error) { panic("TODO: implement Walk") }

// CopyTree copies every regular file from src into the real directory dstDir,
// creating directories as needed, preserving the relative layout. It must not
// follow anything outside dstDir - reject paths containing ".." even though
// fs.FS forbids them, because defence in depth is cheap.
func CopyTree(src fs.FS, dstDir string) error { panic("TODO: implement CopyTree") }

// WriteFileAtomic writes data to path so that a reader either sees the old
// content or the complete new content, never a half-written file, even if the
// process is killed mid-write.
//
// The recipe:
//
//	create a temp file in the SAME directory (a rename across filesystems is
//	    not atomic and may not even work)
//	write the data
//	Sync() - otherwise the rename can land before the data does
//	Close
//	Rename over the target
//	clean up the temp file if anything failed
//
// Set the file mode to perm.
func WriteFileAtomic(path string, data []byte, perm fs.FileMode) error {
	panic("TODO: implement WriteFileAtomic")
}

// TailLines returns the last n lines of a file, reading only the end of it -
// no matter how large the file is. Seek to the end, read a block backwards at a
// time, count newlines, stop as soon as you have enough.
//
// Returns them oldest-first. A trailing newline at the end of the file does not
// count as an empty final line.
func TailLines(path string, n int) ([]string, error) { panic("TODO: implement TailLines") }

// Newest returns the path of the most recently modified file in the tree, and
// its modification time.
func Newest(fsys fs.FS) (string, time.Time, error) { panic("TODO: implement Newest") }
