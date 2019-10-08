// Package version handles common SemVer based version operations
package version

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
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
		return nil, fmt.Errorf("Can't parse Major version number. %w", err)
	}
	Minor, err := strconv.Atoi(s[1])
	if err != nil {
		return nil, fmt.Errorf("Can't parse Minor version number. %w", err)
	}
	Patch, err := strconv.Atoi(s[2])
	if err != nil {
		return nil, fmt.Errorf("Can't parse Patch version number. %w", err)
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
	if v.Equal(b) {
		return false
	}
	switch {
	case v.major < b.major:
		return false
	case v.major > b.major:
		return true
	case v.minor < b.minor:
		return false
	case v.minor > b.minor:
		return true
	case v.patch < b.patch:
		return false
	}
	return true
}

// Before returns whether or not a version is before an other
func (v *Version) Before(b *Version) bool {
	if v.Equal(b) {
		return false
	}
	switch {
	case v.major > b.major:
		return false
	case v.major < b.major:
		return true
	case v.minor > b.minor:
		return false
	case v.minor < b.minor:
		return true
	case v.patch > b.patch:
		return false
	}
	return true
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

// Versions parses multiple versions at once
func Versions(stringVersions ...string) ([]*Version, error) {
	versions := []*Version{}
	for _, v := range stringVersions {
		version, err := New(v)
		if err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}
	return versions, nil
}

// Ascending is used to sort version in ascending order
// by calling sort.Sort(Ascending(versions))
type Ascending []*Version

func (a Ascending) Len() int           { return len(a) }
func (a Ascending) Less(i, j int) bool { return a[i].Before(a[j]) }
func (a Ascending) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// Descending is used to sort version in descending order
// by calling sort.Sort(Descending(versions))
type Descending []*Version

func (d Descending) Len() int           { return len(d) }
func (d Descending) Less(i, j int) bool { return d[i].After(d[j]) }
func (d Descending) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
