package runtime

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/jandedobbeleer/oh-my-posh/src/log"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileContentLogging(t *testing.T) {
	if os.Getenv("OMP_FILE_CONTENT_LOGGING_TEST") != "1" {
		executable, err := os.Executable()
		require.NoError(t, err)
		// The logger has no reset API; isolate its global state from other tests.
		cmd := exec.Command(executable, "-test.run=^TestFileContentLogging$")
		cmd.Env = append(os.Environ(), "OMP_FILE_CONTENT_LOGGING_TEST=1")
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "%s", output)
		return
	}

	dir, err := os.MkdirTemp(".", "filecontent-test-")
	require.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, os.RemoveAll(dir)) })
	dir, err = filepath.Abs(dir)
	require.NoError(t, err)
	log.Enable(true)
	term := &Terminal{}

	cases := []struct {
		File      string
		Sensitive bool
	}{
		{File: ".cfconfig", Sensitive: true},
		{File: "credentials.json", Sensitive: true},
		{File: filepath.Join("unrelated", "credentials.json"), Sensitive: true},
		{File: "config.json"},
		{File: ".cfconfig.backup"},
		{File: "not-credentials.json"},
	}

	for _, tc := range cases {
		t.Run(tc.File, func(t *testing.T) {
			file := filepath.Join(dir, tc.File)
			require.NoError(t, os.MkdirAll(filepath.Dir(file), 0700))
			const content = "synthetic-file-content-sentinel\nunchanged second line"
			require.NoError(t, os.WriteFile(file, []byte(content), 0600))

			start := len(log.String())
			assert.Equal(t, content, term.FileContent(file))
			output := log.String()[start:]
			assert.Contains(t, output, file)
			if tc.Sensitive {
				assert.NotContains(t, output, content)
				assert.NotContains(t, output, "synthetic-file-content-sentinel")
				assert.NotContains(t, output, "unchanged second line")
				assert.Contains(t, output, "[REDACTED]")
				return
			}
			assert.Contains(t, output, "synthetic-file-content-sentinel")
			assert.Contains(t, output, "unchanged second line")
			assert.NotContains(t, output, "[REDACTED]")
		})
	}
}

func TestNormalHostName(t *testing.T) {
	hostName := "hello"
	assert.Equal(t, hostName, cleanHostName(hostName))
}

func TestHostNameWithLocal(t *testing.T) {
	hostName := "hello.local"
	assert.Equal(t, "hello", cleanHostName(hostName))
}

func TestHostNameWithLan(t *testing.T) {
	hostName := "hello.lan"
	cleanHostName := cleanHostName(hostName)
	assert.Equal(t, "hello", cleanHostName)
}

func TestDirMatchesOneOf(t *testing.T) {
	cases := []struct {
		GOOS     string
		HomeDir  string
		Dir      string
		Pattern  string
		Expected bool
	}{
		{GOOS: LINUX, HomeDir: "/home/bill", Dir: "/home/bill", Pattern: "/home/bill", Expected: true},
		{GOOS: LINUX, HomeDir: "/home/bill", Dir: "/home/bill/foo", Pattern: "~/foo", Expected: true},
		{GOOS: LINUX, HomeDir: "/home/bill", Dir: "/home/bill/foo", Pattern: "~/Foo", Expected: false},
		{GOOS: LINUX, HomeDir: "/home/bill", Dir: "/home/bill/foo", Pattern: "~\\\\foo", Expected: true},
		{GOOS: LINUX, HomeDir: "/home/bill", Dir: "/home/bill/foo/bar", Pattern: "~/fo.*", Expected: true},
		{GOOS: LINUX, HomeDir: "/home/bill", Dir: "/home/bill/foo", Pattern: "~/fo\\w", Expected: true},

		{GOOS: WINDOWS, HomeDir: "C:\\Users\\Bill", Dir: "C:\\Users\\Bill", Pattern: "C:\\\\Users\\\\Bill", Expected: true},
		{GOOS: WINDOWS, HomeDir: "C:\\Users\\Bill", Dir: "C:\\Users\\Bill", Pattern: "C:/Users/Bill", Expected: true},
		{GOOS: WINDOWS, HomeDir: "C:\\Users\\Bill", Dir: "C:\\Users\\Bill", Pattern: "c:/users/bill", Expected: true},
		{GOOS: WINDOWS, HomeDir: "C:\\Users\\Bill", Dir: "C:\\Users\\Bill", Pattern: "~", Expected: true},
		{GOOS: WINDOWS, HomeDir: "C:\\Users\\Bill", Dir: "C:\\Users\\Bill\\Foo", Pattern: "~/Foo", Expected: true},
		{GOOS: WINDOWS, HomeDir: "C:\\Users\\Bill", Dir: "C:\\Users\\Bill\\Foo", Pattern: "~/foo", Expected: true},
		{GOOS: WINDOWS, HomeDir: "C:\\Users\\Bill", Dir: "C:\\Users\\Bill\\Foo\\Bar", Pattern: "~/fo.*", Expected: true},
		{GOOS: WINDOWS, HomeDir: "C:\\Users\\Bill", Dir: "C:\\Users\\Bill\\Foo", Pattern: "~/fo\\w", Expected: true},
	}

	for _, tc := range cases {
		got := dirMatchesOneOf(tc.Dir, tc.HomeDir, tc.GOOS, []string{tc.Pattern})
		assert.Equal(t, tc.Expected, got)
	}
}

func TestDirMatchesOneOfRegexInverted(t *testing.T) {
	// detect panic(thrown by MustCompile)
	defer func() {
		if err := recover(); err != nil {
			// display a message explaining omp failed(with the err)
			assert.Equal(t, "regexp: Compile(`^(?!Projects[\\/]).*$`): error parsing regexp: invalid or unsupported Perl syntax: `(?!`", err)
		}
	}()
	_ = dirMatchesOneOf("Projects/oh-my-posh", "", LINUX, []string{"(?!Projects[\\/]).*"})
}

func TestDirMatchesOneOfRegexInvertedNonEscaped(t *testing.T) {
	// detect panic(thrown by MustCompile)
	defer func() {
		if err := recover(); err != nil {
			// display a message explaining omp failed(with the err)
			assert.Equal(t, "regexp: Compile(`^(?!Projects/).*$`): error parsing regexp: invalid or unsupported Perl syntax: `(?!`", err)
		}
	}()
	_ = dirMatchesOneOf("Projects/oh-my-posh", "", LINUX, []string{"(?!Projects/).*"})
}
