// Package version handles common SemVer based version operations
package version

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var (
	// ErrInvalid is return in case of invalid version input
	ErrInvalid = errors.New("Invalid version")
)

// Version is a SemVer based version
type Version struct {
	major      int
	minor      int
	patch      int
	preRelease string
	metadata   string
}

// New returns a new version
func New(v string) (*Version, error) {
	if v == "" {
		return &Version{}, nil
	}
	s := strings.Split(v, ".")
	if len(s) < 3 {
		return nil, ErrInvalid
	}
	Major, err := strconv.Atoi(s[0])
	if err != nil {
		return nil, errors.Wrap(err, "Can't parse Major version number")
	}
	Minor, err := strconv.Atoi(s[1])
	if err != nil {
		return nil, errors.Wrap(err, "Can't parse Minor version number")
	}
	Patch, err := strconv.Atoi(s[2])
	if err != nil {
		return nil, errors.Wrap(err, "Can't parse Patch version number")
	}
	return &Version{
		Major,
		Minor,
		Patch,
		"",
		"",
	}, nil
}

// Bump bumps the Patch version number
func (v *Version) Bump(breakingChanges, changes bool) {
	if breakingChanges {
		v.major++
		v.minor = 0
		v.patch = 0
		return
	}
	if changes {
		v.minor++
		v.patch = 0
		return
	}
	v.patch++
}

// Alpha marks a version a alpha pre-release
func (v *Version) Alpha() {
	v.preRelease = "alpha"
}

// Beta marks a version a alpha pre-release
func (v *Version) Beta() {
	v.preRelease = "beta"
}

// AddMetadata adds metadata to a given version
// This can often be build date or commit information
func (v *Version) AddMetadata(s string) {
	v.metadata = s
}

func (v *Version) String() string {
	switch {
	case v.preRelease != "" && v.metadata == "":
		return fmt.Sprintf("%d.%d.%d-%s", v.major, v.minor, v.patch, v.preRelease)
	case v.preRelease == "" && v.metadata != "":
		return fmt.Sprintf("%d.%d.%d+%s", v.major, v.minor, v.patch, v.metadata)
	case v.preRelease != "" && v.metadata != "":
		return fmt.Sprintf("%d.%d.%d-%s+%s", v.major, v.minor, v.patch, v.preRelease, v.metadata)
	}
	return fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)
}

// After returns whether or not a version is after another
func (v *Version) After(b *Version) bool {
	if v.major > b.major {
		return true
	}
	if v.minor > b.minor {
		return true
	}
	if v.patch > b.patch {
		return true
	}
	return false
}

// Before returns whether or not a version is before an other
func (v *Version) Before(b *Version) bool {
	return !v.After(b)
}

// Equal returns whether or not a version is equal to another
func (v *Version) Equal(b *Version) bool {
	if v.major != b.major {
		return false
	}
	if v.minor != b.minor {
		return false
	}
	if v.patch != b.patch {
		return false
	}
	return true
}
