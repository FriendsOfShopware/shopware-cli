package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstraints(t *testing.T) {
	MustConstraints(NewConstraint(">=1.0.0"))
	MustConstraints(NewConstraint(">=1.0.0 || <2.0.0"))
	MustConstraints(NewConstraint(">=1.0.0,<2.0.0"))
}

func TestConstraintParsingWhitespaceAnd(t *testing.T) {
	c, err := NewConstraint(">=1.0 <2.0")
	assert.NoError(t, err)

	assert.Equal(t, ">=1.0,<2.0", c.String())
	assert.True(t, c.Check(Must(NewVersion("1.0.0"))))
	assert.False(t, c.Check(Must(NewVersion("2.0.0"))))
}

func TestConstraintParsingWhitespaceAndOr(t *testing.T) {
	c, err := NewConstraint("~6.4 >=6.4.20.0 || ~6.5")
	assert.NoError(t, err)

	assert.Equal(t, "~6.4,>=6.4.20.0||~6.5", c.String())
	assert.True(t, c.Check(Must(NewVersion("6.4.20"))))
	assert.True(t, c.Check(Must(NewVersion("6.4.20.0"))))
	assert.True(t, c.Check(Must(NewVersion("6.5.0"))))
	assert.False(t, c.Check(Must(NewVersion("6.4.0.0"))))
}

func TestConstraintWithoutWhitespace(t *testing.T) {
	c, err := NewConstraint("<6.6.1.0||>=6.3.5.0")
	assert.NoError(t, err)

	assert.Equal(t, "<6.6.1.0||>=6.3.5.0", c.String())
	assert.True(t, c.Check(Must(NewVersion("6.4.0.0"))))
}

func TestConstraintVersionNumber(t *testing.T) {
	c, err := NewConstraint("1.0.0")
	assert.NoError(t, err)

	assert.Equal(t, "1.0.0", c.String())
	assert.True(t, c.Check(Must(NewVersion("1.0.0"))))
	assert.False(t, c.Check(Must(NewVersion("1.0.1"))))
}

func TestTildeConstraintPositive(t *testing.T) {
	assert.True(t, MustConstraints(NewConstraint("~2.4")).Check(Must(NewVersion("2.9.0"))))
	assert.True(t, MustConstraints(NewConstraint("~2.4")).Check(Must(NewVersion("2.4.5"))))

	assert.True(t, MustConstraints(NewConstraint("~1")).Check(Must(NewVersion("1.2.3"))))
	assert.True(t, MustConstraints(NewConstraint("~1.0")).Check(Must(NewVersion("1.4.7"))))

	assert.True(t, MustConstraints(NewConstraint("~1.2.1 >=1.2.3")).Check(Must(NewVersion("1.2.3"))))
	assert.True(t, MustConstraints(NewConstraint("~1.2.1 =1.2.3")).Check(Must(NewVersion("1.2.3"))))
	assert.True(t, MustConstraints(NewConstraint("~1.2.1 1.2.3")).Check(Must(NewVersion("1.2.3"))))
	assert.True(t, MustConstraints(NewConstraint("~1.2.1 >=1.2.3 1.2.3")).Check(Must(NewVersion("1.2.3"))))
	assert.True(t, MustConstraints(NewConstraint("~1.2.1 1.2.3 >=1.2.3")).Check(Must(NewVersion("1.2.3"))))
	assert.True(t, MustConstraints(NewConstraint("~1.2.1 1.2.3")).Check(Must(NewVersion("1.2.3"))))
}

func TestTildeConstraintNegative(t *testing.T) {
	assert.False(t, MustConstraints(NewConstraint("~2.4")).Check(Must(NewVersion("3.0.0"))))
	assert.False(t, MustConstraints(NewConstraint("~2.4")).Check(Must(NewVersion("2.3.9"))))
	assert.False(t, MustConstraints(NewConstraint("~1")).Check(Must(NewVersion("0.2.3"))))
}
