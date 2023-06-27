package version

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchingRCWithTilde(t *testing.T) {
	vs := []*Version{
		Must(NewVersion("6.4.4.0")),
		Must(NewVersion("6.5.0.0-rc1")),
	}

	constraint, _ := NewConstraint("~6.5.0")

	match := ""

	for _, v := range vs {
		if constraint.Check(v) {
			match = v.String()
			break
		}
	}

	assert.Equal(t, "6.5.0.0-rc1", match)
}

func TestMatchingRCWithCaret(t *testing.T) {
	vs := []*Version{
		Must(NewVersion("6.4.4.0")),
		Must(NewVersion("6.5.0.0-rc1")),
	}

	constraint, _ := NewConstraint("^6.5")

	match := ""

	for _, v := range vs {
		if constraint.Check(v) {
			match = v.String()
			break
		}
	}

	assert.Equal(t, "6.5.0.0-rc1", match)
}

func TestMatchingRCWithCaretThreeNumbers(t *testing.T) {
	vs := []*Version{
		Must(NewVersion("6.4.4.0")),
		Must(NewVersion("6.5.0.0-rc1")),
	}

	constraint, _ := NewConstraint("^6.5.0")

	match := ""

	for _, v := range vs {
		if constraint.Check(v) {
			match = v.String()
			break
		}
	}

	assert.Equal(t, "6.5.0.0-rc1", match)
}

func TestMatchingRCWithGreaterThanEqual(t *testing.T) {
	vs := []*Version{
		Must(NewVersion("6.4.4.0")),
		Must(NewVersion("6.5.0.0-rc1")),
	}

	constraint, _ := NewConstraint(">=6.5")

	match := ""

	for _, v := range vs {
		if constraint.Check(v) {
			match = v.String()
			break
		}
	}

	assert.Equal(t, "6.5.0.0-rc1", match)
}

func TestCaretConstraint(t *testing.T) {
	constraint, _ := NewConstraint("^6.4.0")

	assert.Equal(t, false, constraint.Check(Must(NewVersion("6.3.0.0"))))
	assert.Equal(t, true, constraint.Check(Must(NewVersion("6.4.0.0"))))
	assert.Equal(t, true, constraint.Check(Must(NewVersion("6.4.0.1"))))
	assert.Equal(t, true, constraint.Check(Must(NewVersion("6.4.1.0"))))
	assert.Equal(t, true, constraint.Check(Must(NewVersion("6.4.5.0"))))
	assert.Equal(t, true, constraint.Check(Must(NewVersion("6.5.5.5"))))
	assert.Equal(t, true, constraint.Check(Must(NewVersion("6.9.9.9"))))
	assert.Equal(t, false, constraint.Check(Must(NewVersion("7.0.0"))))
}

func TestVersionWithoutOperator(t *testing.T) {
	constraint, err := NewConstraint("6.4.0.0")

	assert.NoError(t, err)

	assert.Equal(t, false, constraint.Check(Must(NewVersion("6.3.0.0"))))
	assert.Equal(t, true, constraint.Check(Must(NewVersion("6.4.0.0"))))
	assert.Equal(t, false, constraint.Check(Must(NewVersion("6.5.0.0"))))
}

func TestSortingVersions(t *testing.T) {
	vs := []*Version{
		Must(NewVersion("6.5.0.0-rc2")),
		Must(NewVersion("6.3.1.0")),
		Must(NewVersion("6.5.0.0-rc1")),
		Must(NewVersion("6.2.0")),
		Must(NewVersion("6.4.8.0")),
		Must(NewVersion("6.5.0.0")),
	}

	sort.Sort(Collection(vs))

	assert.Equal(t, "6.2.0", vs[0].String())
	assert.Equal(t, "6.3.1.0", vs[1].String())
	assert.Equal(t, "6.4.8.0", vs[2].String())
	assert.Equal(t, "6.5.0.0-rc1", vs[3].String())
	assert.Equal(t, "6.5.0.0-rc2", vs[4].String())
	assert.Equal(t, "6.5.0.0", vs[5].String())
}

func TestVersionIncrease(t *testing.T) {
	version := Must(NewVersion("1.2.3"))
	version.Increase()
	assert.Equal(t, "1.2.4", version.String())
}

func TestVersionString(t *testing.T) {
	cases := [][]string{
		{"1.2.3", "1.2.3"},
		{"1.2-beta", "1.2.0-beta"},
		{"1.2.0-x.Y.0", "1.2.0-x.Y.0"},
		{"1.2.0-x.Y.0+metadata", "1.2.0-x.Y.0+metadata"},
		{"1.2.0-metadata-1.2.0+metadata~dist", "1.2.0-metadata-1.2.0+metadata~dist"},
		{"17.03.0-ce", "17.3.0-ce"}, // zero-padded fields
	}

	for _, tc := range cases {
		v, err := NewVersion(tc[0])
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		actual := v.String()
		expected := tc[1]
		if actual != expected {
			t.Fatalf("expected: %s\nactual: %s", expected, actual)
		}
		if actual := v.Original(); actual != tc[0] {
			t.Fatalf("expected original: %q\nactual: %q", tc[0], actual)
		}
	}
}

func TestEqual(t *testing.T) {
	cases := []struct {
		v1       string
		v2       string
		expected bool
	}{
		{"1.2.3", "1.4.5", false},
		{"1.2-beta", "1.2-beta", true},
		{"1.2", "1.1.4", false},
		{"1.2", "1.2-beta", false},
		{"1.2+foo", "1.2+beta", true},
		{"v1.2", "v1.2-beta", false},
		{"v1.2+foo", "v1.2+beta", true},
		{"v1.2.3.4", "v1.2.3.4", true},
		{"v1.2.0.0", "v1.2", true},
		{"v1.2.0.0.1", "v1.2", false},
		{"v1.2", "v1.2.0.0", true},
		{"v1.2", "v1.2.0.0.1", false},
		{"v1.2.0.0", "v1.2.0.0.1", false},
		{"v1.2.3.0", "v1.2.3.4", false},
		{"1.7rc2", "1.7rc1", false},
		{"1.7rc2", "1.7", false},
		{"1.2.0", "1.2.0-X-1.2.0+metadata~dist", false},
	}

	for _, tc := range cases {
		v1, err := NewVersion(tc.v1)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		v2, err := NewVersion(tc.v2)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		actual := v1.Equal(v2)
		expected := tc.expected
		if actual != expected {
			t.Fatalf(
				"%s <=> %s\nexpected: %t\nactual: %t",
				tc.v1, tc.v2,
				expected, actual)
		}
	}
}

func TestGreaterThan(t *testing.T) {
	cases := []struct {
		v1       string
		v2       string
		expected bool
	}{
		{"1.2.3", "1.4.5", false},
		{"1.2-beta", "1.2-beta", false},
		{"1.2", "1.1.4", true},
		{"1.2", "1.2-beta", true},
		{"1.2+foo", "1.2+beta", false},
		{"v1.2", "v1.2-beta", true},
		{"v1.2+foo", "v1.2+beta", false},
		{"v1.2.3.4", "v1.2.3.4", false},
		{"v1.2.0.0", "v1.2", false},
		{"v1.2.0.0.1", "v1.2", true},
		{"v1.2", "v1.2.0.0", false},
		{"v1.2", "v1.2.0.0.1", false},
		{"v1.2.0.0", "v1.2.0.0.1", false},
		{"v1.2.3.0", "v1.2.3.4", false},
		{"1.7rc2", "1.7rc1", true},
		{"1.7rc2", "1.7", false},
		{"1.2.0", "1.2.0-X-1.2.0+metadata~dist", true},
	}

	for _, tc := range cases {
		v1, err := NewVersion(tc.v1)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		v2, err := NewVersion(tc.v2)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		actual := v1.GreaterThan(v2)
		expected := tc.expected
		if actual != expected {
			t.Fatalf(
				"%s > %s\nexpected: %t\nactual: %t",
				tc.v1, tc.v2,
				expected, actual)
		}
	}
}

func TestLessThan(t *testing.T) {
	cases := []struct {
		v1       string
		v2       string
		expected bool
	}{
		{"1.2.3", "1.4.5", true},
		{"1.2-beta", "1.2-beta", false},
		{"1.2", "1.1.4", false},
		{"1.2", "1.2-beta", false},
		{"1.2+foo", "1.2+beta", false},
		{"v1.2", "v1.2-beta", false},
		{"v1.2+foo", "v1.2+beta", false},
		{"v1.2.3.4", "v1.2.3.4", false},
		{"v1.2.0.0", "v1.2", false},
		{"v1.2.0.0.1", "v1.2", false},
		{"v1.2", "v1.2.0.0", false},
		{"v1.2", "v1.2.0.0.1", true},
		{"v1.2.0.0", "v1.2.0.0.1", true},
		{"v1.2.3.0", "v1.2.3.4", true},
		{"1.7rc2", "1.7rc1", false},
		{"1.7rc2", "1.7", true},
		{"1.2.0", "1.2.0-X-1.2.0+metadata~dist", false},
	}

	for _, tc := range cases {
		v1, err := NewVersion(tc.v1)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		v2, err := NewVersion(tc.v2)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		actual := v1.LessThan(v2)
		expected := tc.expected
		if actual != expected {
			t.Fatalf(
				"%s < %s\nexpected: %t\nactual: %t",
				tc.v1, tc.v2,
				expected, actual)
		}
	}
}

func TestGreaterThanOrEqual(t *testing.T) {
	cases := []struct {
		v1       string
		v2       string
		expected bool
	}{
		{"1.2.3", "1.4.5", false},
		{"1.2-beta", "1.2-beta", true},
		{"1.2", "1.1.4", true},
		{"1.2", "1.2-beta", true},
		{"1.2+foo", "1.2+beta", true},
		{"v1.2", "v1.2-beta", true},
		{"v1.2+foo", "v1.2+beta", true},
		{"v1.2.3.4", "v1.2.3.4", true},
		{"v1.2.0.0", "v1.2", true},
		{"v1.2.0.0.1", "v1.2", true},
		{"v1.2", "v1.2.0.0", true},
		{"v1.2", "v1.2.0.0.1", false},
		{"v1.2.0.0", "v1.2.0.0.1", false},
		{"v1.2.3.0", "v1.2.3.4", false},
		{"1.7rc2", "1.7rc1", true},
		{"1.7rc2", "1.7", false},
		{"1.2.0", "1.2.0-X-1.2.0+metadata~dist", true},
	}

	for _, tc := range cases {
		v1, err := NewVersion(tc.v1)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		v2, err := NewVersion(tc.v2)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		actual := v1.GreaterThanOrEqual(v2)
		expected := tc.expected
		if actual != expected {
			t.Fatalf(
				"%s >= %s\nexpected: %t\nactual: %t",
				tc.v1, tc.v2,
				expected, actual)
		}
	}
}

func TestLessThanOrEqual(t *testing.T) {
	cases := []struct {
		v1       string
		v2       string
		expected bool
	}{
		{"1.2.3", "1.4.5", true},
		{"1.2-beta", "1.2-beta", true},
		{"1.2", "1.1.4", false},
		{"1.2", "1.2-beta", false},
		{"1.2+foo", "1.2+beta", true},
		{"v1.2", "v1.2-beta", false},
		{"v1.2+foo", "v1.2+beta", true},
		{"v1.2.3.4", "v1.2.3.4", true},
		{"v1.2.0.0", "v1.2", true},
		{"v1.2.0.0.1", "v1.2", false},
		{"v1.2", "v1.2.0.0", true},
		{"v1.2", "v1.2.0.0.1", true},
		{"v1.2.0.0", "v1.2.0.0.1", true},
		{"v1.2.3.0", "v1.2.3.4", true},
		{"1.7rc2", "1.7rc1", false},
		{"1.7rc2", "1.7", true},
		{"1.2.0", "1.2.0-X-1.2.0+metadata~dist", false},
	}

	for _, tc := range cases {
		v1, err := NewVersion(tc.v1)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		v2, err := NewVersion(tc.v2)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		actual := v1.LessThanOrEqual(v2)
		expected := tc.expected
		if actual != expected {
			t.Fatalf(
				"%s <= %s\nexpected: %t\nactual: %t",
				tc.v1, tc.v2,
				expected, actual)
		}
	}
}

func TestConstraintPrerelease(t *testing.T) {
	cases := []struct {
		constraint string
		prerelease bool
	}{
		{"= 1.0", false},
		{"= 1.0-beta", true},
		{"~> 2.1.0", false},
		{"~> 2.1.0-dev", true},
		{"> 2.0", false},
		{">= 2.1.0-a", true},
	}

	for _, tc := range cases {
		c, err := parseSingle(tc.constraint)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		actual := c.Prerelease()
		expected := tc.prerelease
		if actual != expected {
			t.Fatalf("Constraint: %s\nExpected: %#v",
				tc.constraint, expected)
		}
	}
}
