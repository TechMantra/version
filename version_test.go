package version_test

import (
	"fmt"
	"sort"
	"testing"

	"pkg.glorieux.io/version"
)

func TestNewVersion(t *testing.T) {
	t.Run("Empty version", func(t *testing.T) {
		v, err := version.New("")
		if err != nil {
			t.Error(err)
		}
		if v == nil {
			t.Error("Should return a Version")
		}
	})

	t.Run("Parse basic 1.2.3 version", func(t *testing.T) {
		v, err := version.New("1.2.3")
		if err != nil {
			t.Error(err)
		}
		sv := fmt.Sprint(v)
		if sv != "1.2.3" {
			t.Errorf("Should return version 1.2.3 got %s", sv)
		}
	})

	t.Run("Handles invalid version input", func(t *testing.T) {
		_, err := version.New("42")
		if err == nil {
			t.Error("Should be an invalid version number")
		}

		_, err = version.New("a.b.c")
		if err == nil {
			t.Error("Should be an invalid version number")
		}

		_, err = version.New("1.b.c")
		if err == nil {
			t.Error("Should be an invalid version number")
		}

		_, err = version.New("1.2.c")
		if err == nil {
			t.Error("Should be an invalid version number")
		}
	})
}

func TestBump(t *testing.T) {
	v, err := version.New("")
	if err != nil {
		t.Error(err)
	}
	v.Bump(false, false)
	if fmt.Sprint(v) != "0.0.1" {
		t.Error("Should have bumped the patch version number")
	}
}

func TestBumpWithChanges(t *testing.T) {
	v, err := version.New("")
	if err != nil {
		t.Error(err)
	}
	v.Bump(false, true)
	if fmt.Sprint(v) != "0.1.0" {
		t.Error("Should have bumped the minor version number")
	}
}

func TestBumpWithBreakingChanges(t *testing.T) {
	v, err := version.New("")
	if err != nil {
		t.Error(err)
	}
	v.Bump(true, true)
	if fmt.Sprint(v) != "1.0.0" {
		t.Error("Should have bumped the major version number")
	}
}

func TestAddMetadata(t *testing.T) {
	v, err := version.New("")
	if err != nil {
		t.Error(err)
	}
	v.AddMetadata("test")
	sv := fmt.Sprint(v)
	if sv != "0.0.0+test" {
		t.Errorf("Expected 0.0.0+test got %s", sv)
	}
}

func TestAlpha(t *testing.T) {
	v, err := version.New("")
	if err != nil {
		t.Error(err)
	}
	v.Alpha()
	sv := fmt.Sprint(v)
	if sv != "0.0.0-alpha" {
		t.Errorf("Expected 0.0.0-alpha got %s", sv)
	}
}

func TestBeta(t *testing.T) {
	v, err := version.New("")
	if err != nil {
		t.Error(err)
	}
	v.Beta()
	sv := fmt.Sprint(v)
	if sv != "0.0.0-beta" {
		t.Errorf("Expected 0.0.0-beta got %s", sv)
	}
}

func TestString(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		v, err := version.New("")
		if err != nil {
			t.Error(err)
		}
		if fmt.Sprint(v) != "0.0.0" {
			t.Error("Should return a string formatted new version")
		}
	})
	t.Run("Pre-release and metadata", func(t *testing.T) {
		v, err := version.New("")
		if err != nil {
			t.Error(err)
		}
		v.Beta()
		v.AddMetadata("test")
		sv := fmt.Sprint(v)
		if sv != "0.0.0-beta+test" {
			t.Errorf("Should return 0.0.0-beta+test got %s", sv)
		}
	})
}

func TestBeforeAfter(t *testing.T) {
	type test struct {
		versionA string
		versionB string
		expected bool
	}

	cases := []test{
		{"0.1.2", "0.1.3", false},
		{"0.1.2", "1.2.3", false},
		{"0.2.0", "1.1.0", false},
		{"5.4.3", "1.2.4", true},
		{"3.1.2", "1.2.3", true},
		{"1.1.2", "1.2.3", false},
	}

	for _, c := range cases {
		versionA, _ := version.New(c.versionA)
		versionB, _ := version.New(c.versionB)

		if versionA.After(versionB) != c.expected {
			t.Errorf("Expected %s to be after %s", versionA, versionB)
		}
		if versionA.Before(versionB) != !c.expected {
			t.Errorf("Expected %s to be before %s", versionA, versionB)
		}
	}
}

func TestSort(t *testing.T) {
	versions, err := version.Versions("5.4.3", "1.2.4", "1.3.3", "0.1.2")
	if err != nil {
		t.Error(err)
	}
	t.Run("Ascending", func(t *testing.T) {
		sort.Sort(version.Ascending(versions))
		expected, err := version.Versions("0.1.2", "1.2.4", "1.3.3", "5.4.3")
		if err != nil {
			t.Error(err)
		}

		for i, _ := range expected {
			if !versions[i].Equal(expected[i]) {
				t.Errorf("Expected %s got %s", expected[i].String(), versions[i].String())
			}
		}
	})

	t.Run("Descending", func(t *testing.T) {
		sort.Sort(version.Descending(versions))
		expected, err := version.Versions("5.4.3", "1.3.3", "1.2.4", "0.1.2")
		if err != nil {
			t.Error(err)
		}

		for i, _ := range expected {
			if !versions[i].Equal(expected[i]) {
				t.Errorf("Expected %s got %s", expected[i].String(), versions[i].String())
			}
		}
	})

}
