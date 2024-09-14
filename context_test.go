package azfunc

import (
	"net/http"
	"testing"

	"github.com/KarlGW/azfunc/data"
	"github.com/KarlGW/azfunc/output"
	"github.com/google/go-cmp/cmp"
)

func TestContext_Output_HTTP_Write(t *testing.T) {
	ctx := Context{
		Outputs: &outputs{},
	}
	ctx.Outputs.Add(output.NewHTTP())
	ctx.Outputs.HTTP().Write([]byte(`{"message":"hello"}`))

	want := output.NewHTTP(func(o *output.HTTPOptions) {
		o.Body = data.Raw(`{"message":"hello"}`)
	})
	got := ctx.Outputs.HTTP()

	if diff := cmp.Diff(want.Data(), got.Data(), cmp.AllowUnexported(output.HTTP{})); diff != "" {
		t.Errorf("Write() = unexpected result (-want +got)\n%s\n", diff)
	}
}

func TestContext_Output_HTTP_WriteHeader(t *testing.T) {
	ctx := Context{
		Outputs: &outputs{},
	}
	ctx.Outputs.Add(output.NewHTTP())
	ctx.Outputs.HTTP().WriteHeader(http.StatusNotFound)

	want := output.NewHTTP()
	want.WriteHeader(http.StatusNotFound)
	got := ctx.Outputs.HTTP()

	if diff := cmp.Diff(want, got, cmp.AllowUnexported(output.HTTP{})); diff != "" {
		t.Errorf("Write() = unexpected result (-want +got)\n%s\n", diff)
	}
}

func TestContext_Output_HTTP_Header_Add(t *testing.T) {
	ctx := Context{
		Outputs: &outputs{},
	}
	ctx.Outputs.Add(output.NewHTTP())
	ctx.Outputs.HTTP().Header().Add("Content-Type", "application/json")

	want := output.NewHTTP()
	want.Header().Add("Content-Type", "application/json")
	got := ctx.Outputs.HTTP()

	if diff := cmp.Diff(want.Header(), got.Header(), cmp.AllowUnexported(output.HTTP{})); diff != "" {
		t.Errorf("Write() = unexpected result (-want +got)\n%s\n", diff)
	}
}

func TestContext_Output_Trigger_Write(t *testing.T) {
	ctx := Context{
		Outputs: &outputs{},
	}
	ctx.Outputs.Add(output.NewGeneric("queue"))
	ctx.Outputs.Binding("queue").Write([]byte(`{"message":"hello"}`))
	got := ctx.Outputs.Binding("queue")

	want := output.NewGeneric("queue")
	want.Write([]byte(`{"message":"hello"}`))

	if diff := cmp.Diff(want, got, cmp.AllowUnexported(output.Generic{})); diff != "" {
		t.Errorf("Write() = unexpected result (-want +got)\n%s\n", diff)
	}
}
