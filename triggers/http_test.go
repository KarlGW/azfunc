package triggers

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"testing"
	"time"

	"github.com/KarlGW/azfunc/data"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewHTTP(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			req     *http.Request
			options []Option
		}
		want    *HTTP
		wantErr error
	}{
		{
			name: "NewHTTP",
			input: struct {
				req     *http.Request
				options []Option
			}{
				req: &http.Request{
					Body: io.NopCloser(bytes.NewBuffer(httpRequest1)),
				},
			},
			want: &HTTP{
				URL:    "http://localhost:7071/api/endpoint",
				Method: http.MethodPost,
				Body:   data.Raw(`{"message":"hello","number":2}`),
				Headers: http.Header{
					"Content-Type": {"application/json"},
				},
				Params: map[string]string{},
				Query:  map[string]string{},
				Identities: []HTTPIdentity{
					{
						AuthenticationType: "WebJobsAuthLevel",
						IsAuthenticated:    true,
						Actor:              nil,
						BootstrapContext:   nil,
						Claims: []HTTPIdentityClaims{
							{
								Issuer:         "LOCAL AUTHORITY",
								OriginalIssuer: "LOCAL AUTHORITY",
								Properties:     map[string]string{},
								Type:           "http://schemas.microsoft.com/2017/07/functions/claims/authlevel",
								Value:          "Admin",
								ValueType:      "http://www.w3.org/2001/XMLSchema#string",
							},
						},
						Label:         nil,
						Name:          nil,
						NameClaimType: "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/name",
						RoleClaimType: "http://schemas.microsoft.com/ws/2008/06/identity/claims/role",
					},
				},
				Metadata: HTTPMetadata{
					Params: map[string]string{},
					Query:  map[string]string{},
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
					Metadata: Metadata{
						Sys: MetadataSys{
							MethodName: "helloHTTP",
							UtcNow:     _testHTTPTime1,
							RandGuid:   "4e773554-f6b7-4ea2-b07d-4c5fd5aba741",
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := NewHTTP(test.input.req)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("NewHTTP() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr); diff != "" {
				t.Errorf("NewHTTP() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestHTTP_Form(t *testing.T) {
	var tests = []struct {
		name    string
		input   *HTTP
		want    url.Values
		wantErr error
	}{
		{
			name: "Parse Form",
			input: &HTTP{
				Headers: http.Header{
					"Content-Type": {"application/x-www-form-urlencoded"},
				},
				Body: []byte(`field1=value1&field2=value2`),
			},
			want: url.Values{
				"field1": {"value1"},
				"field2": {"value2"},
			},
			wantErr: nil,
		},
		{
			name:    "Parse Form - invalid content type",
			input:   &HTTP{},
			want:    nil,
			wantErr: ErrHTTPInvalidContentType,
		},
		{
			name: "Parse Form - invalid body",
			input: &HTTP{
				Headers: http.Header{
					"Content-Type": {"application/x-www-form-urlencoded"},
				},
				Body: []byte(`{"message":"hello"}`),
			},
			want:    nil,
			wantErr: ErrHTTPInvalidBody,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := test.input.Form()

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Form() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Form() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestHTTP_MultipartForm(t *testing.T) {
	var tests = []struct {
		name    string
		input   *HTTP
		want    *multipart.Form
		wantErr error
	}{
		{
			name: "Parse Multipart Form",
			input: &HTTP{
				Headers: http.Header{
					"Content-Type": {
						"multipart/form-data; boundary=------------------------458d15332083a867",
					},
				},
				Body: []byte(mpfd),
			},
			want: &multipart.Form{
				Value: map[string][]string{},
				File: map[string][]*multipart.FileHeader{
					"file": {
						{
							Filename: "test.txt",
							Header: textproto.MIMEHeader{
								"Content-Disposition": {`form-data; name="file"; filename="test.txt"`},
								"Content-Type":        {"text/plain"},
							},
							Size: 7,
						},
					},
				},
			},
		},
		{
			name: "Parse Multipart form - invalid content type",
			input: &HTTP{
				Body: []byte(mpfd),
			},
			want:    nil,
			wantErr: ErrHTTPInvalidContentType,
		},
		{
			name: "Parse Multipart form - invalid body",
			input: &HTTP{
				Headers: http.Header{
					"Content-Type": {
						"multipart/form-data; boundary=------------------------458d15332083a867",
					},
				},
				Body: []byte(`test`),
			},
			want:    nil,
			wantErr: ErrHTTPInvalidBody,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := test.input.MultipartForm(32 << 20)

			if diff := cmp.Diff(test.want, got, cmpopts.IgnoreUnexported(multipart.FileHeader{})); diff != "" {
				t.Errorf("MultipartForm() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("MultipartForm() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

var httpRequest1 = []byte(`{
	"Data": {
	  "req": {
		"Url": "http://localhost:7071/api/endpoint",
		"Method": "POST",
		"Body": "{\"message\":\"hello\",\"number\":2}",
		"Params": {},
		"Query": {},
		"Headers": {
		  "Content-Type": [
			"application/json"
		  ]
		},
		"Identities": [
		  {
			"AuthenticationType": "WebJobsAuthLevel",
			"IsAuthenticated": true,
			"Actor": null,
			"BootstrapContext": null,
			"Claims": [
			  {
				"Issuer": "LOCAL AUTHORITY",
				"OriginalIssuer": "LOCAL AUTHORITY",
				"Properties": {},
				"Type": "http://schemas.microsoft.com/2017/07/functions/claims/authlevel",
				"Value": "Admin",
				"ValueType": "http://www.w3.org/2001/XMLSchema#string"
			  }
			],
			"Label": null,
			"Name": null,
			"NameClaimType": "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/name",
			"RoleClaimType": "http://schemas.microsoft.com/ws/2008/06/identity/claims/role"
		  }
		]
	  }
	},
	"Metadata": {
	  "Params": {},
	  "Query": {},
	  "Headers": {
		"Content-Type": "application/json"
	  },
	  "sys": {
		"MethodName": "helloHTTP",
		"UtcNow": "2023-10-12T20:13:49.640002Z",
		"RandGuid": "4e773554-f6b7-4ea2-b07d-4c5fd5aba741"
	  }
	}
  }
`)

var mpfd = `--------------------------458d15332083a867
Content-Disposition: form-data; name="file"; filename="test.txt"
Content-Type: text/plain

a file

--------------------------458d15332083a867--
`

var _testHTTPTime1, _ = time.Parse("2006-01-02T15:04:05.999999Z", "2023-10-12T20:13:49.640002Z")
