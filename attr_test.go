package corelog

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnyAttr(t *testing.T) {
	attr := Any("any", "value")
	assert.Equal(t, "any", attr.Key)
	assert.Equal(t, "value", attr.Value.Any())
}

func TestBoolAttr(t *testing.T) {
	attr := Bool("bool", true)
	assert.Equal(t, "bool", attr.Key)
	assert.Equal(t, true, attr.Value.Bool())
}

func TestDurationAttr(t *testing.T) {
	attr := Duration("duration", 5*time.Minute)
	assert.Equal(t, "duration", attr.Key)
	assert.Equal(t, 5*time.Minute, attr.Value.Duration())
}

func TestFloat64Attr(t *testing.T) {
	attr := Float64("float64", float64(1.234))
	assert.Equal(t, "float64", attr.Key)
	assert.Equal(t, float64(1.234), attr.Value.Float64())
}

func TestGroupAttr(t *testing.T) {
	attr := Group("group", "key", "value")
	assert.Equal(t, "group", attr.Key)

	group := attr.Value.Group()
	require.Len(t, group, 1)

	assert.Equal(t, "key", group[0].Key)
	assert.Equal(t, "value", group[0].Value.String())
}

func TestIntAttr(t *testing.T) {
	attr := Int("int", int(10))
	assert.Equal(t, "int", attr.Key)
	assert.Equal(t, int64(10), attr.Value.Int64())
}

func TestInt64Attr(t *testing.T) {
	attr := Int64("int64", int64(-10))
	assert.Equal(t, "int64", attr.Key)
	assert.Equal(t, int64(-10), attr.Value.Int64())
}

func TestStringAttr(t *testing.T) {
	attr := String("string", "value")
	assert.Equal(t, "string", attr.Key)
	assert.Equal(t, "value", attr.Value.String())
}

func TestTimeAttr(t *testing.T) {
	// drop the monotonic portion
	now := time.Unix(0, int64(time.Now().UnixMicro()))

	attr := Time("time", now)
	assert.Equal(t, "time", attr.Key)
	assert.Equal(t, now, attr.Value.Time())
}

func TestUint64Attr(t *testing.T) {
	attr := Uint64("uint64", uint64(10))
	assert.Equal(t, "uint64", attr.Key)
	assert.Equal(t, uint64(10), attr.Value.Uint64())
}
