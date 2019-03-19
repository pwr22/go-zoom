package main

import (
	"io/ioutil"
	"os"
	"reflect"
	"runtime"
	"testing"
)

// if this doesn't work then lots of tests here are going to blow up
func TestParseArgsCanBeCalledMoreThanOnce(t *testing.T) {
	parseArgs()
	parseArgs()
}

func TestParrellelismDefault(t *testing.T) {
	p, ep := *parallelism, runtime.NumCPU()
	if p != ep {
		t.Fatalf("default -j is %v but should be %v", p, ep)
	}
}

func TestKeepOrderDefault(t *testing.T) {
	k, ek := *keepOrder, false
	if k != ek {
		t.Fatalf("default -k is %v but should be %v", k, ek)
	}
}

func TestParseArgsWithNoArgs(t *testing.T) {
	os.Args = []string{"zoom"}
	parseArgs()
}

func TestParseArgsWithVersion(t *testing.T) {
	os.Args = []string{"zoom", "--version"}

	exitEarly, err := parseArgs()
	defer func() { *printVersion = false }() // reset version for later tests

	if err != nil {
		t.Fatal(err)
	} else if !exitEarly {
		t.Fatalf("exitEarly should be true")
	}
}

func TestParseArgsWithDryRun(t *testing.T) {
	os.Args = []string{"zoom", "--dry-run", ":::", "cmd"}

	exitEarly, err := parseArgs()
	defer func() { *dryRun = false }()

	if err != nil {
		t.Fatal(err)
	} else if !exitEarly {
		t.Fatalf("exitEarly should be true")
	}
}

func TestParseArgsWithDryRunAndFailedCmdsBuild(t *testing.T) {
	os.Args = []string{"zoom", "--dry-run", "::::", "non-existent-file"}

	_, err := parseArgs()
	defer func() { *dryRun = false }()

	if err == nil {
		t.Fatalf("should have errored")
	}
}

func TestGetArgSets(t *testing.T) {
	// setup file contents for file based tests
	tmpfile, err := ioutil.TempFile("", "zoomtest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	expFileArgs := []string{"a", "b", "c"}
	for _, a := range expFileArgs {
		if _, err := tmpfile.WriteString(a + "\n"); err != nil {
			t.Fatal(err)
		}
	}

	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	expArgList := []string{"1", "2", "3"}

	testCases := []struct {
		desc string
		args []string
		exp  [][]string
	}{
		{
			desc: "no args",
			args: []string{"zoom", "echo", "-n"},
		},
		{
			desc: "arg list",
			args: append([]string{"zoom", "echo", ":::"}, expArgList...),
			exp:  [][]string{expArgList},
		},
		{
			desc: "two arg lists",
			args: append(append(append([]string{"zoom", "echo", ":::"}, expArgList...), ":::"), expArgList...),
			exp:  [][]string{expArgList, expArgList},
		},
		{
			desc: "file list",
			args: []string{"zoom", "echo", "::::", tmpfile.Name()},
			exp:  [][]string{expFileArgs},
		},
		{
			desc: "two file lists",
			args: []string{"zoom", "echo", "::::", tmpfile.Name(), "::::", tmpfile.Name()},
			exp:  [][]string{expFileArgs, expFileArgs},
		},
		{
			desc: "arglist then file list",
			args: append(append([]string{"zoom", "echo", ":::"}, expArgList...), "::::", tmpfile.Name()),
			exp:  [][]string{expArgList, expFileArgs},
		},
		{
			desc: "file list then arg list",
			args: append([]string{"zoom", "echo", "::::", tmpfile.Name(), ":::"}, expArgList...),
			exp:  [][]string{expFileArgs, expArgList},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			// flags.Args() is used so gotta call this here
			// TODO refactor so this isn't called directly and we pass arguments instead
			os.Args = tC.args
			parseArgs()

			as, err := getArgSets()
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(as, tC.exp) {
				t.Fatalf("expected %v but got %v", tC.exp, as)
			}
		})
	}
}

func TestGetArgSetsErrs(t *testing.T) {
	// setup file contents for file based tests
	tmpfile, err := ioutil.TempFile("", "zoomtest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	expFileArgs := []string{"a", "b", "c"}
	for _, a := range expFileArgs {
		if _, err := tmpfile.WriteString(a + "\n"); err != nil {
			t.Fatal(err)
		}
	}

	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		desc string
		args []string
	}{
		{
			desc: "trailing empty :::",
			args: []string{"zoom", "echo", ":::"},
		},
		{
			desc: "trailing empty ::::",
			args: []string{"zoom", "echo", "::::"},
		},
		{
			desc: "beginning empty ::: followed by :::",
			args: []string{"zoom", "echo", ":::", ":::", "arg"},
		},
		{
			desc: "beginning empty :::: followed by :::",
			args: []string{"zoom", "echo", "::::", ":::", "arg"},
		},
		{
			desc: "beginning empty ::: followed by ::::",
			args: []string{"zoom", "echo", ":::", "::::", tmpfile.Name()},
		},
		{
			desc: "beginning empty :::: followed by :::",
			args: []string{"zoom", "echo", "::::", "::::", tmpfile.Name()},
		},
		{
			desc: "trailing double standard input",
			args: []string{"zoom", "echo", "::::", "-", "-"},
		},
		{
			desc: "beginning double standard input followed by :::",
			args: []string{"zoom", "echo", "::::", "-", "-", ":::", "arg"},
		},
		{
			desc: "beginning double standard input followed by ::::",
			args: []string{"zoom", "echo", "::::", "-", "-", "::::", tmpfile.Name()},
		},
		{
			desc: "beginning non-existent file followed by :::",
			args: []string{"zoom", "echo", "::::", "non-existent-file", ":::", "arg"},
		},
		{
			desc: "beginning non-existent file followed by ::::",
			args: []string{"zoom", "echo", "::::", "non-existent-file", "::::", tmpfile.Name()},
		},
		{
			desc: "trailing non-existent file",
			args: []string{"zoom", "echo", "::::", "non-existent-file"},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			// flags.Args() is used so gotta call this here
			// TODO refactor so this isn't called directly and we pass arguments instead
			os.Args = tC.args
			parseArgs()

			_, err := getArgSets()
			if err == nil {
				t.Fatalf("should have errored")
			}
		})
	}
}

func TestGetCmdPrefix(t *testing.T) {
	testCases := []struct {
		desc string
		args []string
		exp  string
	}{
		{
			desc: "just a command",
			args: []string{"zoom", "echo", "-n"},
			exp:  "echo -n",
		},
		{
			desc: "command followed by arg list",
			args: []string{"zoom", "echo", ":::", "1", "2", "3"},
			exp:  "echo",
		},
		{
			desc: "command followed by file list",
			args: []string{"zoom", "echo", "::::", "1.txt", "2.txt", "3.txt"},
			exp:  "echo",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			// flags.Args() is used so gotta call this here
			// TODO refactor so this isn't called directly and we pass arguments instead
			os.Args = tC.args
			parseArgs()
			if p := getCmdPrefix(); p != tC.exp {
				t.Fatalf("expected %v but got %v", tC.exp, p)
			}
		})
	}
}

func TestPermuteArgSets(t *testing.T) {
	testCases := []struct {
		desc    string
		argSets [][]string
		exp     []string
	}{
		{
			desc:    "no argsets",
			argSets: [][]string{},
			exp:     []string{},
		},
		{
			desc:    "single argset",
			argSets: [][]string{[]string{"a", "b", "c"}},
			exp:     []string{"a", "b", "c"},
		},
		{
			desc:    "two argsets",
			argSets: [][]string{[]string{"a", "b", "c"}, []string{"1", "2"}},
			exp:     []string{"a 1", "a 2", "b 1", "b 2", "c 1", "c 2"},
		},
		{
			desc:    "three argsets",
			argSets: [][]string{[]string{"a", "b", "c"}, []string{"1", "2"}, []string{"x", "y", "z"}},
			exp:     []string{"a 1 x", "a 1 y", "a 1 z", "a 2 x", "a 2 y", "a 2 z", "b 1 x", "b 1 y", "b 1 z", "b 2 x", "b 2 y", "b 2 z", "c 1 x", "c 1 y", "c 1 z", "c 2 x", "c 2 y", "c 2 z"},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if p := permuteArgSets(tC.argSets); !reflect.DeepEqual(p, tC.exp) {
				t.Fatalf("expected %v but got %v", tC.exp, p)
			}
		})
	}
}

func TestPlaceHolder(t *testing.T) {
	if placeHolder != "{}" {
		t.Fatalf("expected %v got %v", "{}", placeHolder)
	}
}

func TestReadCmdsFromFile(t *testing.T) {
	exp := []string{"1", "2", "3"}
	tmpfile, err := ioutil.TempFile("", "zoomtest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	for _, arg := range exp {
		if _, err := tmpfile.WriteString(arg + "\n"); err != nil {
			t.Fatal(err)
		}
	}

	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	cmds, err := readCmdsFromFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(cmds, exp) {
		t.Fatalf("expected %v but got %v", exp, cmds)
	}

}

func TestReadCmdsFromStdin(t *testing.T) {
	var exp []string

	cmds, err := readCmdsFromFile("-")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(cmds, exp) {
		t.Fatalf("expected %v but got %v", exp, cmds)
	}

}

// TODO somehow test reading from stdin
func TestGetCmdStrings(t *testing.T) {
	testCases := []struct {
		desc string
		args []string
		exp  []string
	}{
		{
			desc: "arglist with no placeholder ",
			args: []string{"zoom", "echo", ":::", "1", "2"},
			exp:  []string{"echo 1", "echo 2"},
		},
		{
			desc: "arglist with placeholder",
			args: []string{"zoom", "echo", "{}", "-n", ":::", "1", "2"},
			exp:  []string{"echo 1 -n", "echo 2 -n"},
		},
		{
			desc: "args from stdin",
			args: []string{"zoom", "echo"},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			// flags.Args() is used so gotta call this here
			// TODO refactor so this isn't called directly and we pass arguments instead
			os.Args = tC.args
			parseArgs()

			c, err := getCmdStrings()
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(c, tC.exp) {
				t.Fatalf("expected %v but got %v", tC.exp, c)
			}
		})
	}
}

// benchmarks

// a single arglist

func benchmarkArglist(n int, b *testing.B) {
	os.Args = make([]string, 0, n+3)
	os.Args = append(os.Args, "zoom", "echo", ":::")
	for i := 0; i < n; i++ {
		os.Args = append(os.Args, "foo")
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := getCmdStrings(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkArglist1(b *testing.B) {
	benchmarkArglist(1, b)
}

func BenchmarkArglist10(b *testing.B) {
	benchmarkArglist(10, b)
}

func BenchmarkArglist100(b *testing.B) {
	benchmarkArglist(100, b)
}

func BenchmarkArglist1000(b *testing.B) {
	benchmarkArglist(1000, b)
}

// two arglists

func benchmarkPermutedArglists(n int, b *testing.B) {
	os.Args = make([]string, 0, 2*n+4)
	os.Args = append(os.Args, "zoom", "echo", ":::")
	for i := 0; i < n; i++ {
		os.Args = append(os.Args, "foo")
	}
	os.Args = append(os.Args, ":::")
	for i := 0; i < n; i++ {
		os.Args = append(os.Args, "bar")
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := getCmdStrings(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPermutedArglists1(b *testing.B) {
	benchmarkPermutedArglists(1, b)
}

func BenchmarkPermutedArglists10(b *testing.B) {
	benchmarkPermutedArglists(10, b)
}

func BenchmarkPermutedArglists100(b *testing.B) {
	benchmarkPermutedArglists(100, b)
}

func BenchmarkPermutedArglists1000(b *testing.B) {
	benchmarkPermutedArglists(1000, b)
}

// a single arg file

func benchmarkArgfile(n int, b *testing.B) {
	// setup file contents for file based tests
	tmpfile, err := ioutil.TempFile("", "zoomtest")
	if err != nil {
		b.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	for i := 0; i < n; i++ {
		if _, err := tmpfile.WriteString("foo\n"); err != nil {
			b.Fatal(err)
		}
	}

	if err := tmpfile.Close(); err != nil {
		b.Fatal(err)
	}

	os.Args = []string{"zoom", "echo", "::::", tmpfile.Name()}
	parseArgs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := getCmdStrings(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkArgfile1(b *testing.B) {
	benchmarkArgfile(1, b)
}

func BenchmarkArgfile10(b *testing.B) {
	benchmarkArgfile(10, b)
}

func BenchmarkArgfile100(b *testing.B) {
	benchmarkArgfile(100, b)
}

func BenchmarkArgfile1000(b *testing.B) {
	benchmarkArgfile(1000, b)
}

// two arg file

func benchmarkPermutedArgfiles(n int, b *testing.B) {
	// setup file contents for file based tests
	tmpfile, err := ioutil.TempFile("", "zoomtest")
	if err != nil {
		b.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	for i := 0; i < n; i++ {
		if _, err := tmpfile.WriteString("foo\n"); err != nil {
			b.Fatal(err)
		}
	}

	if err := tmpfile.Close(); err != nil {
		b.Fatal(err)
	}

	os.Args = []string{"zoom", "echo", "::::", tmpfile.Name(), "::::", tmpfile.Name()}
	parseArgs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := getCmdStrings(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPermutedArgfiles1(b *testing.B) {
	benchmarkPermutedArgfiles(1, b)
}

func BenchmarkPermutedArgfiles10(b *testing.B) {
	benchmarkPermutedArgfiles(10, b)
}

func BenchmarkPermutedArgfiles100(b *testing.B) {
	benchmarkPermutedArgfiles(100, b)
}

func BenchmarkPermutedArgfiles1000(b *testing.B) {
	benchmarkPermutedArgfiles(1000, b)
}
