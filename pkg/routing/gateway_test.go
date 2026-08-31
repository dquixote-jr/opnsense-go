package routing

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/errs"
)

func newController(t *testing.T) Controller {
	t.Helper()
	return Controller{
		Api: api.NewClient(api.Options{
			Uri:           os.Getenv("OPNSENSE_URI"),
			APIKey:        os.Getenv("OPNSENSE_API_KEY"),
			APISecret:     os.Getenv("OPNSENSE_API_SECRET"),
			AllowInsecure: true,
			MaxBackoff:    30,
			MinBackoff:    1,
			MaxRetries:    4,
		}),
	}
}

// The gateway is created disabled with monitoring off on purpose: an enabled
// gateway starts a dpinger probe and participates in routing on whatever box
// the suite runs against.
func TestGateway(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	gw := &Gateway{
		Disabled:       "1",
		Name:           "gotestgw",
		Description:    "go test gateway",
		Interface:      api.SelectedMap("wan"),
		IPProtocol:     api.SelectedMap("inet"),
		Gateway:        "10.99.99.1",
		FarGateway:     "1",
		MonitorDisable: "1",
		Priority:       "255",
		Weight:         "1",
	}

	id, err := c.AddGateway(ctx, gw)
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	t.Cleanup(func() { _ = c.DeleteGateway(ctx, id) })

	got, err := c.GetGateway(ctx, id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	for _, tc := range []struct{ field, got, want string }{
		{"Name", got.Name, gw.Name},
		{"Gateway", got.Gateway, gw.Gateway},
		{"FarGateway", got.FarGateway, gw.FarGateway},
		{"Interface", got.Interface.String(), gw.Interface.String()},
		{"IPProtocol", got.IPProtocol.String(), gw.IPProtocol.String()},
	} {
		if tc.got != tc.want {
			t.Errorf("round trip mismatch on %s: got %q, want %q", tc.field, tc.got, tc.want)
		}
	}

	got.Description = "go test gateway updated"
	got.Weight = "2"
	if err := c.UpdateGateway(ctx, id, got); err != nil {
		t.Fatalf("update: %v", err)
	}

	back, err := c.GetGateway(ctx, id)
	if err != nil {
		t.Fatalf("re-get: %v", err)
	}
	if back.Description != got.Description || back.Weight != "2" {
		t.Fatalf("update not persisted: descr=%q weight=%q", back.Description, back.Weight)
	}

	if err := c.DeleteGateway(ctx, id); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := c.GetGateway(ctx, id); err == nil {
		t.Fatalf("expected gateway to be gone after delete")
	} else {
		var notFound *errs.NotFoundError
		if !errors.As(err, &notFound) {
			t.Fatalf("expected NotFoundError after delete, got: %v", err)
		}
	}
}
