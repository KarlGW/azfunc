package data

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestRaw_MarshalJSON(t *testing.T) {
	var tests = []struct {
		name    string
		input   []byte
		want    []byte
		wantErr error
	}{
		{
			name:    "string",
			input:   []byte(`hello`),
			want:    []byte(`{"body":"hello"}`),
			wantErr: nil,
		},
		{
			name:  "JSON",
			input: []byte(`{"message":"hello"}`),
			want:  []byte(`{"body":"{\"message\":\"hello\"}"}`),
		},
		{
			name:  "escaped JSON",
			input: []byte(`{\"message\":\"hello\"}`),
			want:  []byte(`{"body":"{\\\"message\\\":\\\"hello\\\"}"}`),
		},
		{
			name:  "raw bytes",
			input: []byte{118, 134, 6, 38, 145, 183, 207, 177},
			want:  []byte(`{"body":"doYGJpG3z7E="}`),

			wantErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			type testRawData struct {
				Body Raw `json:"body"`
			}

			r := testRawData{
				Body: test.input,
			}
			got, gotErr := json.Marshal(r)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("MarshalJSON() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("MarshalJSON() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestRaw_UnmarshalJSON(t *testing.T) {
	var tests = []struct {
		name    string
		input   []byte
		want    Raw
		wantErr error
	}{
		{
			name:    "string",
			input:   []byte(`{"body":"hello"}`),
			want:    Raw(`hello`),
			wantErr: nil,
		},
		{
			name:    "form data string",
			input:   []byte(`{"body":"field1=value1&field2=value2"}`),
			want:    Raw(`field1=value1&field2=value2`),
			wantErr: nil,
		},
		{
			name:    "JSON",
			input:   []byte(`{"body":"{\"message\":\"hello\"}"}`),
			want:    []byte(`{"message":"hello"}`),
			wantErr: nil,
		},
		{
			name:    "escaped JSON",
			input:   []byte(`{"body":"{\\\"message\\\":\\\"hello\\\"}"}`),
			want:    []byte(`{\"message\":\"hello\"}`),
			wantErr: nil,
		},
		{
			name:  "raw bytes",
			input: []byte(`{"body":"doYGJpG3z7E="}`),
			want:  []byte{118, 134, 6, 38, 145, 183, 207, 177},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			type testRawData struct {
				Body Raw `json:"body"`
			}

			var p testRawData
			gotErr := json.Unmarshal(test.input, &p)
			got := p.Body

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("MarshalJSON() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("MarshalJSON() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}
