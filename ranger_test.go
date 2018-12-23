package ranger

import (
	assert "github.com/stretchr/testify/require"
	"os"
	"testing"
)

type testConfig struct {
	S  string
	SS []string
	I  int
	SI []int
	B  bool
}

func setup(t *testing.T) (*assert.Assertions, *Ranger, *testConfig) {
	t.Parallel()

	r := New()
	r.SetDefault("S", "bar")

	return assert.New(t), r, new(testConfig)
}

func TestUnmarshal(t *testing.T) {
	assert, r, c := setup(t)

	// Test default
	unmarshal(assert, r, c)
	assert.Equal("bar", c.S)

	// Test environment loading
	// - String
	setEnv(assert, r, "S", "Baz")
	unmarshal(assert, r, c)
	assert.Equal("Baz", c.S)

	// - String slice
	setEnv(assert, r, "SS", "Baz,Bin")
	unmarshal(assert, r, c)
	assert.Equal([]string{"Baz", "Bin"}, c.SS)

	// - Int
	setEnv(assert, r, "I", "42")
	unmarshal(assert, r, c)
	assert.Equal(42, c.I)

	// - Int Slice
	setEnv(assert, r, "SI", "7,16")
	unmarshal(assert, r, c)
	assert.Equal([]int{7, 16}, c.SI)

	// - Bool
	setEnv(assert, r, "B", "true")
	unmarshal(assert, r, c)
	assert.True(c.B)
}

func setEnv(assert *assert.Assertions, r *Ranger, key, value string) {
	r.SetRequired(key)
	err := os.Setenv(key, value)
	assert.Nil(err)
	assert.Equal(os.Getenv(key), value)
}

func unmarshal(assert *assert.Assertions, r *Ranger, c *testConfig) {
	err := r.Unmarshal(c)
	assert.Nil(err)
}
