package azfunc

import (
	"net/http"
	"testing"

	"github.com/KarlGW/azfunc/bindings"
	"github.com/KarlGW/azfunc/data"
	"github.com/google/go-cmp/cmp"
)

func TestContext_Output_HTTP_Write(t *testing.T) {
	ctx := Context{
		Output: bindings.Output{},
	}
	ctx.Output.AddBindings(bindings.NewHTTP())
	ctx.Output.HTTP().Write([]byte(`{"message":"hello"}`))

	want := &bindings.HTTP{Body: data.Raw(`{"message":"hello"}`)}
	got := ctx.Output.HTTP()

	if diff := cmp.Diff(want.Body, got.Body, cmp.AllowUnexported(bindings.HTTP{})); diff != "" {
		t.Errorf("Write() = unexpected result (-want +got)\n%s\n", diff)
	}
}

func TestContext_Output_HTTP_WriteHeader(t *testing.T) {
	ctx := Context{
		Output: bindings.Output{},
	}
	ctx.Output.AddBindings(bindings.NewHTTP())
	ctx.Output.HTTP().WriteHeader(http.StatusNotFound)

	want := &bindings.HTTP{}
	want.WriteHeader(http.StatusNotFound)
	got := ctx.Output.HTTP()

	if diff := cmp.Diff(want.StatusCode, got.StatusCode, cmp.AllowUnexported(bindings.HTTP{})); diff != "" {
		t.Errorf("Write() = unexpected result (-want +got)\n%s\n", diff)
	}
}

func TestContext_Output_HTTP_Header_Add(t *testing.T) {
	ctx := Context{
		Output: bindings.Output{},
	}
	ctx.Output.AddBindings(bindings.NewHTTP())
	ctx.Output.HTTP().Header().Add("Content-Type", "application/json")

	want := &bindings.HTTP{}
	want.Header().Add("Content-Type", "application/json")
	got := ctx.Output.HTTP()

	if diff := cmp.Diff(want.Header(), got.Header(), cmp.AllowUnexported(bindings.HTTP{})); diff != "" {
		t.Errorf("Write() = unexpected result (-want +got)\n%s\n", diff)
	}
}

func TestContext_Output_Trigger_Write(t *testing.T) {
	ctx := Context{
		Output: bindings.Output{},
	}
	ctx.Output.AddBindings(bindings.NewBase("queue"))
	ctx.Output.Binding("queue").Write([]byte(`{"message":"hello"}`))
	got := ctx.Output.Binding("queue")

	want := bindings.NewBase("queue")
	want.Write([]byte(`{"message":"hello"}`))

	if diff := cmp.Diff(want, got, cmp.AllowUnexported(bindings.Base{})); diff != "" {
		t.Errorf("Write() = unexpected result (-want +got)\n%s\n", diff)
	}
}
