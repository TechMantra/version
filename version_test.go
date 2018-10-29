package version_test

import (
	"fmt"
	"testing"

	"techmantra.io/version"
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
