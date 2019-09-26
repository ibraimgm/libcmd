package libcmd_test

import (
	"math"
	"strconv"
	"testing"

	"github.com/ibraimgm/libcmd"
)

func TestGet(t *testing.T) {
	tests := []struct {
		cmd      []string
		abool    bool
		aint     int
		auint    uint
		astring  string
		afloat32 float32
		afloat64 float64
		args     []string
	}{
		{cmd: []string{}},
		{cmd: []string{"-b", "-i", "5", "-u", "9", "-s", "foo"}, abool: true, aint: 5, auint: 9, astring: "foo"},
		{cmd: []string{"--abool", "--aint", "5", "--auint", "9", "--astring", "foo"}, abool: true, aint: 5, auint: 9, astring: "foo"},
		{cmd: []string{"--aint=5", "--astring=foo"}, aint: 5, astring: "foo"},
		{cmd: []string{"-b", "--abool=false"}},
		{cmd: []string{"-b", "--no-abool"}},
		{cmd: []string{"-i", "5", "--aint", "6", "-i", "7"}, aint: 7},
		{cmd: []string{"-u", "5", "--auint", "6", "-u", "7"}, auint: 7},
		{cmd: []string{"-b", "-i", "5", "foo", "bar"}, abool: true, aint: 5, args: []string{"foo", "bar"}},
		{cmd: []string{"foo", "bar"}, args: []string{"foo", "bar"}},
		{cmd: []string{"foo", "-i", "5"}, args: []string{"foo", "-i", "5"}},
		{cmd: []string{"-b", "--no-abool=true"}},
		{cmd: []string{"--no-abool=false"}, abool: true},
		{cmd: []string{"-s", "foo", "--astring="}},
		{cmd: []string{"--astring", "--aint", "5"}, astring: "--aint", args: []string{"5"}},
		{cmd: []string{"-i", "5", "-f32", "3.14", "-f64", "3.1415"}, aint: 5, afloat32: float32(3.14), afloat64: float64(3.1415)},
		{cmd: []string{"--afloat32", "3.14", "--afloat64", "3.1415"}, afloat32: float32(3.14), afloat64: float64(3.1415)},
		{cmd: []string{"--afloat32=3.14", "--afloat64=3.1415"}, afloat32: float32(3.14), afloat64: float64(3.1415)},
	}

	for i, test := range tests {
		app := libcmd.NewApp("", "")

		app.Bool("abool", "b", false, "specifies a bool value")
		app.Int("aint", "i", 0, "specifies an int value")
		app.Uint("auint", "u", 0, "specifies an uint value")
		app.String("astring", "s", "", "specifies a string value")
		app.Float32("afloat32", "f32", 0, "specifies a float32 value")
		app.Float64("afloat64", "f64", 0, "specifies a float64 value")

		if err := app.RunArgs(test.cmd); err != nil {
			t.Errorf("Case %d, error parsing args: %v", i, err)
			continue
		}

		abool := app.GetBool("abool")
		aint := app.GetInt("aint")
		auint := app.GetUint("auint")
		astring := app.GetString("astring")
		afloat32 := app.GetFloat32("afloat32")
		afloat64 := app.GetFloat64("afloat64")

		compareValue(t, i, test.abool, *abool)
		compareValue(t, i, test.aint, *aint)
		compareValue(t, i, test.auint, *auint)
		compareValue(t, i, test.astring, *astring)
		compareValue(t, i, test.afloat32, *afloat32)
		compareValue(t, i, test.afloat64, *afloat64)
		compareArgs(t, i, test.args, app.Args())
	}
}

func TestGetIntLimit(t *testing.T) {
	tests := []struct {
		cmd []string
		a   int8
		b   int16
		c   int32
		d   int64
		e   uint8
		f   uint16
		g   uint32
		h   uint64
	}{
		{cmd: []string{"-a", strconv.FormatInt(math.MaxInt8, 10)}, a: math.MaxInt8},
		{cmd: []string{"-b", strconv.FormatInt(math.MaxInt16, 10)}, b: math.MaxInt16},
		{cmd: []string{"-c", strconv.FormatInt(math.MaxInt32, 10)}, c: math.MaxInt32},
		{cmd: []string{"-d", strconv.FormatInt(math.MaxInt64, 10)}, d: math.MaxInt64},
		{cmd: []string{"-e", strconv.FormatUint(math.MaxUint8, 10)}, e: math.MaxUint8},
		{cmd: []string{"-f", strconv.FormatUint(math.MaxUint16, 10)}, f: math.MaxUint16},
		{cmd: []string{"-g", strconv.FormatUint(math.MaxUint32, 10)}, g: math.MaxUint32},
		{cmd: []string{"-h", strconv.FormatUint(math.MaxUint64, 10)}, h: math.MaxUint64},
	}

	for i, test := range tests {
		app := libcmd.NewApp("", "")
		app.Options.SuppressHelpFlag = true

		app.Int8("", "a", 0, "specifies a int8 value")
		app.Int16("", "b", 0, "specifies a int16 value")
		app.Int32("", "c", 0, "specifies a int32 value")
		app.Int64("", "d", 0, "specifies a int64 value")
		app.Uint8("", "e", 0, "specifies a uint8 value")
		app.Uint16("", "f", 0, "specifies a uint16 value")
		app.Uint32("", "g", 0, "specifies a uint32 value")
		app.Uint64("", "h", 0, "specifies a uint64 value")

		if err := app.RunArgs(test.cmd); err != nil {
			t.Errorf("Case %d, error parsing args: %v", i, err)
			continue
		}

		a := app.GetInt8("a")
		b := app.GetInt16("b")
		c := app.GetInt32("c")
		d := app.GetInt64("d")
		e := app.GetUint8("e")
		f := app.GetUint16("f")
		g := app.GetUint32("g")
		h := app.GetUint64("h")

		compareValue(t, i, test.a, *a)
		compareValue(t, i, test.b, *b)
		compareValue(t, i, test.c, *c)
		compareValue(t, i, test.d, *d)
		compareValue(t, i, test.e, *e)
		compareValue(t, i, test.f, *f)
		compareValue(t, i, test.g, *g)
		compareValue(t, i, test.h, *h)
	}
}

func TestGetChoice(t *testing.T) {
	tests := []struct {
		cmd       []string
		expected  string
		expectErr bool
	}{
		{cmd: []string{}},
		{cmd: []string{"-c", "foo"}, expected: "foo"},
		{cmd: []string{"-c", "bar"}, expected: "bar"},
		{cmd: []string{"-c", "baz"}, expected: "baz"},
		{cmd: []string{"-c", "hey"}, expectErr: true},
	}

	for i, test := range tests {
		app := libcmd.NewApp("", "")
		app.Choice([]string{"foo", "bar", "baz"}, "", "c", "", "")

		err := app.RunArgs(test.cmd)
		if !test.expectErr && err != nil {
			t.Errorf("Case %d, error parsing args: %v", i, err)
			continue
		}

		if test.expectErr && !libcmd.IsParserErr(err) {
			t.Errorf("Case %d, expected error but none received", i)
			continue
		}

		s := app.GetChoice("c")
		if *s != test.expected {
			t.Errorf("Case %d, wrong value: expected '%s', received '%s'", i, test.expected, *s)
		}
	}
}
