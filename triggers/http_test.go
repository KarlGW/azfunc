package triggers

import (
	"bytes"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/KarlGW/azfunc/data"
	"github.com/google/go-cmp/cmp"
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
							UtcNow:     _testTime1,
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

var _testTime1, _ = time.Parse("2006-01-02T15:04:05.999999Z", "2023-10-12T20:13:49.640002Z")
