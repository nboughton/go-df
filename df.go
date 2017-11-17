package df

import (
	"math"
	"syscall"
)

// DiskFree packages the same types of information one gets from the GNU Linux "df" cmd
type DiskFree struct {
	total       float64
	used        float64
	avail       float64
	percentUsed int
}

// size provides a specific type for data size category conversion
type size int64

// Exported constants for file size calculations
const (
	KB size = 1024    // Kilobyte
	MB      = KB * KB // Megabyte
	GB      = MB * KB // Gigabyte
	TB      = GB * KB // Terabyte
	PB      = TB * KB // Petabyte
)

// NewDf creates a DiskFree struct for the specified mount point. One can then use the
// struct methods to get the appropriate values in whatever size increment they wish.
// Values are stored as float64s in order to allow for a reasonable degree of accuracy when
// dealing with large value numbers i.e 1.7GB, 3.4TB etc...
func NewDf(mountPoint string) (f DiskFree, err error) {
	s := syscall.Statfs_t{}
	if err = syscall.Statfs(mountPoint, &s); err != nil {
		return f, err
	}

	f = DiskFree{
		total: float64(s.Blocks) * float64(s.Bsize),
		used:  float64(s.Blocks-s.Bfree) * float64(s.Bsize),
	}

	f.avail = f.total - f.used
	f.percentUsed = int(math.Ceil(f.used / f.total * 100))

	return f, err
}

// Total takes a size (df.MB, df.GB etc) and returns the amount of total space
//
// Example:
//    d, _ := df.NewDf("/mnt/fs")
//    fmt.Println(df.Total(df.GB))
func (df DiskFree) Total(s size) float64 {
	return df.total / float64(s)
}

// Used takes a size (df.MB, df.GB etc) and returns the amount of space used.
//
// Example:
//    d, _ := df.NewDf("/mnt/fs")
//    fmt.Println(df.Used(df.GB))
func (df DiskFree) Used(s size) float64 {
	return df.used / float64(s)
}

// Avail takes a size (df.MB, df.GB etc) and returns the amount of space available.
//
// Example:
//    d, _ := df.NewDf("/mnt/fs")
//    fmt.Println(df.Avail(df.GB))
func (df DiskFree) Avail(s size) float64 {
	return df.avail / float64(s)
}

// PercentUsed returns the percentage of space used as an int
func (df DiskFree) PercentUsed() int {
	return df.percentUsed
}
