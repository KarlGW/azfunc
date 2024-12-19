package azfunc

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestSetupLogger(t *testing.T) {
	var tests = []struct {
		name  string
		input map[string]string
		want  Logger
	}{
		{
			name:  "defaults",
			input: map[string]string{},
			want:  NewLogger(),
		},
		{
			name: "disabled with environment variable",
			input: map[string]string{
				functionsDisableLogging: "true",
			},
			want: noOpLogger{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Cleanup(func() {
				for k := range test.input {
					os.Unsetenv(k)
				}
			})

			for k, v := range test.input {
				t.Setenv(k, v)
			}

			got := setupLogger()
			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(logger{}, noOpLogger{}), cmpopts.IgnoreFields(logger{}, "stdout", "stderr")); diff != "" {
				t.Errorf("setupLogger() = unexpected result (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestParseBool(t *testing.T) {
	var tests = []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "input: true",
			input: "true",
			want:  true,
		},
		{
			name:  "input: 1",
			input: "1",
			want:  true,
		},
		{
			name:  "input: false",
			input: "false",
			want:  false,
		},
		{
			name:  "input: 0",
			input: "0",
			want:  false,
		},
		{
			name:  "input: empty",
			input: "",
			want:  false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := parseBool(test.input)
			if test.want != got {
				t.Errorf("parseBool() = unexpected result, want: %v, got: %v\n", test.want, got)
			}
		})
	}
}

func TestWithDisableLogging(t *testing.T) {
	want := noOpLogger{}

	f := functionApp{
		log: NewLogger(),
	}
	WithDisableLogging()(&f)

	got := f.log
	if diff := cmp.Diff(want, got, cmp.AllowUnexported(noOpLogger{})); diff != "" {
		t.Errorf("WithDisableLogging() = unexpected result (-want +got)\n%s\n", diff)
	}
}
