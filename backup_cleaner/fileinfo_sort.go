package backup_cleaner

import (
	"os"

	"github.com/bborbe/stringutil"
)

type FileInfoByName []os.FileInfo

func (v FileInfoByName) Len() int      { return len(v) }
func (v FileInfoByName) Swap(i, j int) { v[i], v[j] = v[j], v[i] }
func (v FileInfoByName) Less(i, j int) bool {
	return stringutil.StringLess(v[i].Name(), v[j].Name())
}
