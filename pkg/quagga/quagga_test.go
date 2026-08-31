package quagga

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

// assertGone confirms a deleted resource is really gone. The API answers a
// missing uuid with an empty body, which the client surfaces as NotFoundError.
func assertGone(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected resource to be gone after delete, but it was still readable")
	}
	var notFound *errs.NotFoundError
	if !errors.As(err, &notFound) {
		t.Fatalf("expected NotFoundError after delete, got: %v", err)
	}
}

var _ = context.Background

func TestBGPNeighbor(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	r := &BGPNeighbor{
		Enabled: "1", Description: "go test bgp nbr", PeerIP: "10.99.99.20",
		RemoteAS: "65010", Weight: "", KeepAlive: "", HoldDown: "", ConnectTimer: "",
	}

	id, err := c.AddBGPNeighbor(ctx, r)
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	t.Cleanup(func() { _ = c.DeleteBGPNeighbor(ctx, id) })

	got, err := c.GetBGPNeighbor(ctx, id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Description != r.Description {
		t.Fatalf("round trip mismatch on Description: got %v, want %v", got.Description, r.Description)
	}

	got.Description = "go test bgp nbr upd"
	if err := c.UpdateBGPNeighbor(ctx, id, got); err != nil {
		t.Fatalf("update: %v", err)
	}

	back, err := c.GetBGPNeighbor(ctx, id)
	if err != nil {
		t.Fatalf("re-get: %v", err)
	}
	if back.Description != got.Description {
		t.Fatalf("update not persisted on Description: got %v, want %v", back.Description, got.Description)
	}

	if err := c.DeleteBGPNeighbor(ctx, id); err != nil {
		t.Fatalf("delete: %v", err)
	}
	_, err = c.GetBGPNeighbor(ctx, id)
	assertGone(t, err)
}

func TestBGPASPath(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	r := &BGPASPath{
		Enabled: "1", Description: "go test aspath", Number: "99",
		Action: api.SelectedMap("permit"), AS: "65010",
	}

	id, err := c.AddBGPASPath(ctx, r)
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	t.Cleanup(func() { _ = c.DeleteBGPASPath(ctx, id) })

	got, err := c.GetBGPASPath(ctx, id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Description != r.Description {
		t.Fatalf("round trip mismatch on Description: got %v, want %v", got.Description, r.Description)
	}

	got.Description = "go test aspath upd"
	if err := c.UpdateBGPASPath(ctx, id, got); err != nil {
		t.Fatalf("update: %v", err)
	}

	back, err := c.GetBGPASPath(ctx, id)
	if err != nil {
		t.Fatalf("re-get: %v", err)
	}
	if back.Description != got.Description {
		t.Fatalf("update not persisted on Description: got %v, want %v", back.Description, got.Description)
	}

	if err := c.DeleteBGPASPath(ctx, id); err != nil {
		t.Fatalf("delete: %v", err)
	}
	_, err = c.GetBGPASPath(ctx, id)
	assertGone(t, err)
}

func TestBGPPrefixList(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	r := &BGPPrefixList{
		Enabled: "1", Description: "go test bgp pl", Name: "gotestbgppl",
		IPVersion: api.SelectedMap("IPv4"), SequenceNumber: "99",
		Action: api.SelectedMap("permit"), Network: "198.51.100.0/24",
	}

	id, err := c.AddBGPPrefixList(ctx, r)
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	t.Cleanup(func() { _ = c.DeleteBGPPrefixList(ctx, id) })

	got, err := c.GetBGPPrefixList(ctx, id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Action != r.Action {
		t.Fatalf("round trip mismatch on Action: got %v, want %v", got.Action, r.Action)
	}

	got.Action = api.SelectedMap("deny")
	if err := c.UpdateBGPPrefixList(ctx, id, got); err != nil {
		t.Fatalf("update: %v", err)
	}

	back, err := c.GetBGPPrefixList(ctx, id)
	if err != nil {
		t.Fatalf("re-get: %v", err)
	}
	if back.Action.String() != got.Action.String() {
		t.Fatalf("update not persisted on Action: got %v, want %v", back.Action.String(), got.Action.String())
	}

	if err := c.DeleteBGPPrefixList(ctx, id); err != nil {
		t.Fatalf("delete: %v", err)
	}
	_, err = c.GetBGPPrefixList(ctx, id)
	assertGone(t, err)
}

func TestBGPCommunityList(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	r := &BGPCommunityList{
		Enabled: "1", Description: "go test cl", Number: "99",
		SequenceNumber: "99", Action: api.SelectedMap("permit"), Community: "65010:1",
	}

	id, err := c.AddBGPCommunityList(ctx, r)
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	t.Cleanup(func() { _ = c.DeleteBGPCommunityList(ctx, id) })

	got, err := c.GetBGPCommunityList(ctx, id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Description != r.Description {
		t.Fatalf("round trip mismatch on Description: got %v, want %v", got.Description, r.Description)
	}

	got.Description = "go test cl upd"
	if err := c.UpdateBGPCommunityList(ctx, id, got); err != nil {
		t.Fatalf("update: %v", err)
	}

	back, err := c.GetBGPCommunityList(ctx, id)
	if err != nil {
		t.Fatalf("re-get: %v", err)
	}
	if back.Description != got.Description {
		t.Fatalf("update not persisted on Description: got %v, want %v", back.Description, got.Description)
	}

	if err := c.DeleteBGPCommunityList(ctx, id); err != nil {
		t.Fatalf("delete: %v", err)
	}
	_, err = c.GetBGPCommunityList(ctx, id)
	assertGone(t, err)
}

func TestBGPRouteMap(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	r := &BGPRouteMap{
		Enabled: "1", Description: "go test bgp rm", Name: "gotestbgprm",
		Action: api.SelectedMap("permit"), RouteMapID: "99",
	}

	id, err := c.AddBGPRouteMap(ctx, r)
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	t.Cleanup(func() { _ = c.DeleteBGPRouteMap(ctx, id) })

	got, err := c.GetBGPRouteMap(ctx, id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Description != r.Description {
		t.Fatalf("round trip mismatch on Description: got %v, want %v", got.Description, r.Description)
	}

	got.Description = "go test bgp rm upd"
	if err := c.UpdateBGPRouteMap(ctx, id, got); err != nil {
		t.Fatalf("update: %v", err)
	}

	back, err := c.GetBGPRouteMap(ctx, id)
	if err != nil {
		t.Fatalf("re-get: %v", err)
	}
	if back.Description != got.Description {
		t.Fatalf("update not persisted on Description: got %v, want %v", back.Description, got.Description)
	}

	if err := c.DeleteBGPRouteMap(ctx, id); err != nil {
		t.Fatalf("delete: %v", err)
	}
	_, err = c.GetBGPRouteMap(ctx, id)
	assertGone(t, err)
}

func TestOSPFInterface(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	r := &OSPFInterface{
		Enabled: "1", Interface: api.SelectedMap("lo0"), Area: "0.0.0.0",
		AuthKeyID: "1", BFD: "0",
	}

	id, err := c.AddOSPFInterface(ctx, r)
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	t.Cleanup(func() { _ = c.DeleteOSPFInterface(ctx, id) })

	got, err := c.GetOSPFInterface(ctx, id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Area != r.Area {
		t.Fatalf("round trip mismatch on Area: got %v, want %v", got.Area, r.Area)
	}

	got.Area = "0.0.0.1"
	if err := c.UpdateOSPFInterface(ctx, id, got); err != nil {
		t.Fatalf("update: %v", err)
	}

	back, err := c.GetOSPFInterface(ctx, id)
	if err != nil {
		t.Fatalf("re-get: %v", err)
	}
	if back.Area != got.Area {
		t.Fatalf("update not persisted on Area: got %v, want %v", back.Area, got.Area)
	}

	if err := c.DeleteOSPFInterface(ctx, id); err != nil {
		t.Fatalf("delete: %v", err)
	}
	_, err = c.GetOSPFInterface(ctx, id)
	assertGone(t, err)
}

func TestOSPFRouteMap(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	r := &OSPFRouteMap{
		Enabled: "1", Name: "gotestospfrm", Action: api.SelectedMap("permit"),
		RouteMapID: "99",
	}

	id, err := c.AddOSPFRouteMap(ctx, r)
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	t.Cleanup(func() { _ = c.DeleteOSPFRouteMap(ctx, id) })

	got, err := c.GetOSPFRouteMap(ctx, id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Name != r.Name {
		t.Fatalf("round trip mismatch on Name: got %v, want %v", got.Name, r.Name)
	}

	got.Name = "gotestospfrm2"
	if err := c.UpdateOSPFRouteMap(ctx, id, got); err != nil {
		t.Fatalf("update: %v", err)
	}

	back, err := c.GetOSPFRouteMap(ctx, id)
	if err != nil {
		t.Fatalf("re-get: %v", err)
	}
	if back.Name != got.Name {
		t.Fatalf("update not persisted on Name: got %v, want %v", back.Name, got.Name)
	}

	if err := c.DeleteOSPFRouteMap(ctx, id); err != nil {
		t.Fatalf("delete: %v", err)
	}
	_, err = c.GetOSPFRouteMap(ctx, id)
	assertGone(t, err)
}

func TestOSPF6Interface(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	r := &OSPF6Interface{
		Enabled: "1", Interface: api.SelectedMap("lo0"), Area: "0.0.0.0",
		Passive: "0", BFD: "0",
	}

	id, err := c.AddOSPF6Interface(ctx, r)
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	t.Cleanup(func() { _ = c.DeleteOSPF6Interface(ctx, id) })

	got, err := c.GetOSPF6Interface(ctx, id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Area != r.Area {
		t.Fatalf("round trip mismatch on Area: got %v, want %v", got.Area, r.Area)
	}

	got.Area = "0.0.0.1"
	if err := c.UpdateOSPF6Interface(ctx, id, got); err != nil {
		t.Fatalf("update: %v", err)
	}

	back, err := c.GetOSPF6Interface(ctx, id)
	if err != nil {
		t.Fatalf("re-get: %v", err)
	}
	if back.Area != got.Area {
		t.Fatalf("update not persisted on Area: got %v, want %v", back.Area, got.Area)
	}

	if err := c.DeleteOSPF6Interface(ctx, id); err != nil {
		t.Fatalf("delete: %v", err)
	}
	_, err = c.GetOSPF6Interface(ctx, id)
	assertGone(t, err)
}

func TestOSPF6RouteMap(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	r := &OSPF6RouteMap{
		Enabled: "1", Name: "gotestospf6rm", Action: api.SelectedMap("permit"),
		RouteMapID: "99",
	}

	id, err := c.AddOSPF6RouteMap(ctx, r)
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	t.Cleanup(func() { _ = c.DeleteOSPF6RouteMap(ctx, id) })

	got, err := c.GetOSPF6RouteMap(ctx, id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Name != r.Name {
		t.Fatalf("round trip mismatch on Name: got %v, want %v", got.Name, r.Name)
	}

	got.Name = "gotestospf6rm2"
	if err := c.UpdateOSPF6RouteMap(ctx, id, got); err != nil {
		t.Fatalf("update: %v", err)
	}

	back, err := c.GetOSPF6RouteMap(ctx, id)
	if err != nil {
		t.Fatalf("re-get: %v", err)
	}
	if back.Name != got.Name {
		t.Fatalf("update not persisted on Name: got %v, want %v", back.Name, got.Name)
	}

	if err := c.DeleteOSPF6RouteMap(ctx, id); err != nil {
		t.Fatalf("delete: %v", err)
	}
	_, err = c.GetOSPF6RouteMap(ctx, id)
	assertGone(t, err)
}

func TestBFDNeighbor(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	r := &BFDNeighbor{
		Enabled: "1", Description: "go test bfd", PeerIP: "10.99.99.9",
		MultiHop: "0", DetectMultiplier: "3", ReceiveInterval: "300", TransmitInterval: "300",
	}

	id, err := c.AddBFDNeighbor(ctx, r)
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	t.Cleanup(func() { _ = c.DeleteBFDNeighbor(ctx, id) })

	got, err := c.GetBFDNeighbor(ctx, id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Description != r.Description {
		t.Fatalf("round trip mismatch on Description: got %v, want %v", got.Description, r.Description)
	}

	got.Description = "go test bfd upd"
	if err := c.UpdateBFDNeighbor(ctx, id, got); err != nil {
		t.Fatalf("update: %v", err)
	}

	back, err := c.GetBFDNeighbor(ctx, id)
	if err != nil {
		t.Fatalf("re-get: %v", err)
	}
	if back.Description != got.Description {
		t.Fatalf("update not persisted on Description: got %v, want %v", back.Description, got.Description)
	}

	if err := c.DeleteBFDNeighbor(ctx, id); err != nil {
		t.Fatalf("delete: %v", err)
	}
	_, err = c.GetBFDNeighbor(ctx, id)
	assertGone(t, err)
}

func TestBGPPeerGroup(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	r := &BGPPeerGroup{
		Enabled: "1", Name: "gotestpg", RemoteASMode: api.SelectedMap("external"),
		Family: api.SelectedMap("ipv4"), NextHopSelf: "0", DefaultRoute: "0",
	}

	id, err := c.AddBGPPeerGroup(ctx, r)
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	t.Cleanup(func() { _ = c.DeleteBGPPeerGroup(ctx, id) })

	got, err := c.GetBGPPeerGroup(ctx, id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.NextHopSelf != r.NextHopSelf {
		t.Fatalf("round trip mismatch on NextHopSelf: got %v, want %v", got.NextHopSelf, r.NextHopSelf)
	}

	got.NextHopSelf = "1"
	if err := c.UpdateBGPPeerGroup(ctx, id, got); err != nil {
		t.Fatalf("update: %v", err)
	}

	back, err := c.GetBGPPeerGroup(ctx, id)
	if err != nil {
		t.Fatalf("re-get: %v", err)
	}
	if back.NextHopSelf != got.NextHopSelf {
		t.Fatalf("update not persisted on NextHopSelf: got %v, want %v", back.NextHopSelf, got.NextHopSelf)
	}

	if err := c.DeleteBGPPeerGroup(ctx, id); err != nil {
		t.Fatalf("delete: %v", err)
	}
	_, err = c.GetBGPPeerGroup(ctx, id)
	assertGone(t, err)
}

func TestBGPRedistribution(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	r := &BGPRedistribution{
		Enabled: "1", Description: "go test bgp redist",
		Redistribute: api.SelectedMap("connected"),
	}

	id, err := c.AddBGPRedistribution(ctx, r)
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	t.Cleanup(func() { _ = c.DeleteBGPRedistribution(ctx, id) })

	got, err := c.GetBGPRedistribution(ctx, id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Description != r.Description {
		t.Fatalf("round trip mismatch on Description: got %v, want %v", got.Description, r.Description)
	}

	got.Description = "go test bgp redist upd"
	if err := c.UpdateBGPRedistribution(ctx, id, got); err != nil {
		t.Fatalf("update: %v", err)
	}

	back, err := c.GetBGPRedistribution(ctx, id)
	if err != nil {
		t.Fatalf("re-get: %v", err)
	}
	if back.Description != got.Description {
		t.Fatalf("update not persisted on Description: got %v, want %v", back.Description, got.Description)
	}

	if err := c.DeleteBGPRedistribution(ctx, id); err != nil {
		t.Fatalf("delete: %v", err)
	}
	_, err = c.GetBGPRedistribution(ctx, id)
	assertGone(t, err)
}

func TestOSPFArea(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	r := &OSPFArea{
		Enabled: "1", AreaID: "0.0.0.99", Type: api.SelectedMap("stub"),
	}

	id, err := c.AddOSPFArea(ctx, r)
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	t.Cleanup(func() { _ = c.DeleteOSPFArea(ctx, id) })

	got, err := c.GetOSPFArea(ctx, id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Type != r.Type {
		t.Fatalf("round trip mismatch on Type: got %v, want %v", got.Type, r.Type)
	}

	got.Type = api.SelectedMap("nssa")
	if err := c.UpdateOSPFArea(ctx, id, got); err != nil {
		t.Fatalf("update: %v", err)
	}

	back, err := c.GetOSPFArea(ctx, id)
	if err != nil {
		t.Fatalf("re-get: %v", err)
	}
	if back.Type.String() != got.Type.String() {
		t.Fatalf("update not persisted on Type: got %v, want %v", back.Type.String(), got.Type.String())
	}

	if err := c.DeleteOSPFArea(ctx, id); err != nil {
		t.Fatalf("delete: %v", err)
	}
	_, err = c.GetOSPFArea(ctx, id)
	assertGone(t, err)
}

func TestOSPFNetwork(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	r := &OSPFNetwork{
		Enabled: "1", IPAddr: "198.51.100.0", Netmask: "24", Area: "0.0.0.0",
	}

	id, err := c.AddOSPFNetwork(ctx, r)
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	t.Cleanup(func() { _ = c.DeleteOSPFNetwork(ctx, id) })

	got, err := c.GetOSPFNetwork(ctx, id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Netmask != r.Netmask {
		t.Fatalf("round trip mismatch on Netmask: got %v, want %v", got.Netmask, r.Netmask)
	}

	got.Netmask = "25"
	if err := c.UpdateOSPFNetwork(ctx, id, got); err != nil {
		t.Fatalf("update: %v", err)
	}

	back, err := c.GetOSPFNetwork(ctx, id)
	if err != nil {
		t.Fatalf("re-get: %v", err)
	}
	if back.Netmask != got.Netmask {
		t.Fatalf("update not persisted on Netmask: got %v, want %v", back.Netmask, got.Netmask)
	}

	if err := c.DeleteOSPFNetwork(ctx, id); err != nil {
		t.Fatalf("delete: %v", err)
	}
	_, err = c.GetOSPFNetwork(ctx, id)
	assertGone(t, err)
}

func TestOSPFNeighbor(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	r := &OSPFNeighbor{
		Enabled: "1", Description: "go test ospf nbr", PeerIP: "198.51.100.9",
	}

	id, err := c.AddOSPFNeighbor(ctx, r)
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	t.Cleanup(func() { _ = c.DeleteOSPFNeighbor(ctx, id) })

	got, err := c.GetOSPFNeighbor(ctx, id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Description != r.Description {
		t.Fatalf("round trip mismatch on Description: got %v, want %v", got.Description, r.Description)
	}

	got.Description = "go test ospf nbr upd"
	if err := c.UpdateOSPFNeighbor(ctx, id, got); err != nil {
		t.Fatalf("update: %v", err)
	}

	back, err := c.GetOSPFNeighbor(ctx, id)
	if err != nil {
		t.Fatalf("re-get: %v", err)
	}
	if back.Description != got.Description {
		t.Fatalf("update not persisted on Description: got %v, want %v", back.Description, got.Description)
	}

	if err := c.DeleteOSPFNeighbor(ctx, id); err != nil {
		t.Fatalf("delete: %v", err)
	}
	_, err = c.GetOSPFNeighbor(ctx, id)
	assertGone(t, err)
}

func TestOSPFPrefixList(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	r := &OSPFPrefixList{
		Enabled: "1", Name: "gotestospfpl", SequenceNumber: "99",
		Action: api.SelectedMap("permit"), Network: "198.51.100.0/24",
	}

	id, err := c.AddOSPFPrefixList(ctx, r)
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	t.Cleanup(func() { _ = c.DeleteOSPFPrefixList(ctx, id) })

	got, err := c.GetOSPFPrefixList(ctx, id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Action != r.Action {
		t.Fatalf("round trip mismatch on Action: got %v, want %v", got.Action, r.Action)
	}

	got.Action = api.SelectedMap("deny")
	if err := c.UpdateOSPFPrefixList(ctx, id, got); err != nil {
		t.Fatalf("update: %v", err)
	}

	back, err := c.GetOSPFPrefixList(ctx, id)
	if err != nil {
		t.Fatalf("re-get: %v", err)
	}
	if back.Action.String() != got.Action.String() {
		t.Fatalf("update not persisted on Action: got %v, want %v", back.Action.String(), got.Action.String())
	}

	if err := c.DeleteOSPFPrefixList(ctx, id); err != nil {
		t.Fatalf("delete: %v", err)
	}
	_, err = c.GetOSPFPrefixList(ctx, id)
	assertGone(t, err)
}

func TestOSPFRedistribution(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	r := &OSPFRedistribution{
		Enabled: "1", Description: "go test ospf redist",
		Redistribute: api.SelectedMap("connected"),
	}

	id, err := c.AddOSPFRedistribution(ctx, r)
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	t.Cleanup(func() { _ = c.DeleteOSPFRedistribution(ctx, id) })

	got, err := c.GetOSPFRedistribution(ctx, id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Description != r.Description {
		t.Fatalf("round trip mismatch on Description: got %v, want %v", got.Description, r.Description)
	}

	got.Description = "go test ospf redist upd"
	if err := c.UpdateOSPFRedistribution(ctx, id, got); err != nil {
		t.Fatalf("update: %v", err)
	}

	back, err := c.GetOSPFRedistribution(ctx, id)
	if err != nil {
		t.Fatalf("re-get: %v", err)
	}
	if back.Description != got.Description {
		t.Fatalf("update not persisted on Description: got %v, want %v", back.Description, got.Description)
	}

	if err := c.DeleteOSPFRedistribution(ctx, id); err != nil {
		t.Fatalf("delete: %v", err)
	}
	_, err = c.GetOSPFRedistribution(ctx, id)
	assertGone(t, err)
}

func TestOSPF6Network(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	r := &OSPF6Network{
		Enabled: "1", IPAddr: "2001:db8::", Netmask: "64", Area: "0.0.0.0",
	}

	id, err := c.AddOSPF6Network(ctx, r)
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	t.Cleanup(func() { _ = c.DeleteOSPF6Network(ctx, id) })

	got, err := c.GetOSPF6Network(ctx, id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Netmask != r.Netmask {
		t.Fatalf("round trip mismatch on Netmask: got %v, want %v", got.Netmask, r.Netmask)
	}

	got.Netmask = "48"
	if err := c.UpdateOSPF6Network(ctx, id, got); err != nil {
		t.Fatalf("update: %v", err)
	}

	back, err := c.GetOSPF6Network(ctx, id)
	if err != nil {
		t.Fatalf("re-get: %v", err)
	}
	if back.Netmask != got.Netmask {
		t.Fatalf("update not persisted on Netmask: got %v, want %v", back.Netmask, got.Netmask)
	}

	if err := c.DeleteOSPF6Network(ctx, id); err != nil {
		t.Fatalf("delete: %v", err)
	}
	_, err = c.GetOSPF6Network(ctx, id)
	assertGone(t, err)
}

func TestOSPF6PrefixList(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	r := &OSPF6PrefixList{
		Enabled: "1", Name: "gotestospf6pl", SequenceNumber: "99",
		Action: api.SelectedMap("permit"), Network: "2001:db8::/64",
	}

	id, err := c.AddOSPF6PrefixList(ctx, r)
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	t.Cleanup(func() { _ = c.DeleteOSPF6PrefixList(ctx, id) })

	got, err := c.GetOSPF6PrefixList(ctx, id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Action != r.Action {
		t.Fatalf("round trip mismatch on Action: got %v, want %v", got.Action, r.Action)
	}

	got.Action = api.SelectedMap("deny")
	if err := c.UpdateOSPF6PrefixList(ctx, id, got); err != nil {
		t.Fatalf("update: %v", err)
	}

	back, err := c.GetOSPF6PrefixList(ctx, id)
	if err != nil {
		t.Fatalf("re-get: %v", err)
	}
	if back.Action.String() != got.Action.String() {
		t.Fatalf("update not persisted on Action: got %v, want %v", back.Action.String(), got.Action.String())
	}

	if err := c.DeleteOSPF6PrefixList(ctx, id); err != nil {
		t.Fatalf("delete: %v", err)
	}
	_, err = c.GetOSPF6PrefixList(ctx, id)
	assertGone(t, err)
}

func TestOSPF6Redistribution(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	r := &OSPF6Redistribution{
		Enabled: "1", Description: "go test ospf6 redist",
		Redistribute: api.SelectedMap("connected"),
	}

	id, err := c.AddOSPF6Redistribution(ctx, r)
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	t.Cleanup(func() { _ = c.DeleteOSPF6Redistribution(ctx, id) })

	got, err := c.GetOSPF6Redistribution(ctx, id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Description != r.Description {
		t.Fatalf("round trip mismatch on Description: got %v, want %v", got.Description, r.Description)
	}

	got.Description = "go test ospf6 redist upd"
	if err := c.UpdateOSPF6Redistribution(ctx, id, got); err != nil {
		t.Fatalf("update: %v", err)
	}

	back, err := c.GetOSPF6Redistribution(ctx, id)
	if err != nil {
		t.Fatalf("re-get: %v", err)
	}
	if back.Description != got.Description {
		t.Fatalf("update not persisted on Description: got %v, want %v", back.Description, got.Description)
	}

	if err := c.DeleteOSPF6Redistribution(ctx, id); err != nil {
		t.Fatalf("delete: %v", err)
	}
	_, err = c.GetOSPF6Redistribution(ctx, id)
	assertGone(t, err)
}

func TestStaticRoute(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	r := &StaticRoute{
		Enabled: "1", Network: "198.51.100.0/24", Gateway: "10.99.99.1", BFD: "0",
	}

	id, err := c.AddStaticRoute(ctx, r)
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	t.Cleanup(func() { _ = c.DeleteStaticRoute(ctx, id) })

	got, err := c.GetStaticRoute(ctx, id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Network != r.Network {
		t.Fatalf("round trip mismatch on Network: got %v, want %v", got.Network, r.Network)
	}

	got.Network = "198.51.100.128/25"
	if err := c.UpdateStaticRoute(ctx, id, got); err != nil {
		t.Fatalf("update: %v", err)
	}

	back, err := c.GetStaticRoute(ctx, id)
	if err != nil {
		t.Fatalf("re-get: %v", err)
	}
	if back.Network != got.Network {
		t.Fatalf("update not persisted on Network: got %v, want %v", back.Network, got.Network)
	}

	if err := c.DeleteStaticRoute(ctx, id); err != nil {
		t.Fatalf("delete: %v", err)
	}
	_, err = c.GetStaticRoute(ctx, id)
	assertGone(t, err)
}

func TestGeneralSettings(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	orig, err := c.GetGeneralSettings(ctx)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	before := orig.EnableSyslog
	t.Cleanup(func() {
		orig.EnableSyslog = before
		_ = c.UpdateGeneralSettings(ctx, orig)
	})

	orig.EnableSyslog = "1"
	if err := c.UpdateGeneralSettings(ctx, orig); err != nil {
		t.Fatalf("update: %v", err)
	}

	back, err := c.GetGeneralSettings(ctx)
	if err != nil {
		t.Fatalf("re-get: %v", err)
	}
	if back.EnableSyslog != "1" {
		t.Fatalf("update not persisted on EnableSyslog: got %v, want %v", back.EnableSyslog, "1")
	}
}

func TestBFDSettings(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	orig, err := c.GetBFDSettings(ctx)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	before := orig.Enabled
	t.Cleanup(func() {
		orig.Enabled = before
		_ = c.UpdateBFDSettings(ctx, orig)
	})

	orig.Enabled = "0"
	if err := c.UpdateBFDSettings(ctx, orig); err != nil {
		t.Fatalf("update: %v", err)
	}

	back, err := c.GetBFDSettings(ctx)
	if err != nil {
		t.Fatalf("re-get: %v", err)
	}
	if back.Enabled != "0" {
		t.Fatalf("update not persisted on Enabled: got %v, want %v", back.Enabled, "0")
	}
}

func TestBGPSettings(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	orig, err := c.GetBGPSettings(ctx)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	before := orig.RouterID
	t.Cleanup(func() {
		orig.RouterID = before
		_ = c.UpdateBGPSettings(ctx, orig)
	})

	orig.RouterID = "10.99.99.99"
	if err := c.UpdateBGPSettings(ctx, orig); err != nil {
		t.Fatalf("update: %v", err)
	}

	back, err := c.GetBGPSettings(ctx)
	if err != nil {
		t.Fatalf("re-get: %v", err)
	}
	if back.RouterID != "10.99.99.99" {
		t.Fatalf("update not persisted on RouterID: got %v, want %v", back.RouterID, "10.99.99.99")
	}
}

func TestOSPFSettings(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	orig, err := c.GetOSPFSettings(ctx)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	before := orig.RouterID
	t.Cleanup(func() {
		orig.RouterID = before
		_ = c.UpdateOSPFSettings(ctx, orig)
	})

	orig.RouterID = "10.99.99.98"
	if err := c.UpdateOSPFSettings(ctx, orig); err != nil {
		t.Fatalf("update: %v", err)
	}

	back, err := c.GetOSPFSettings(ctx)
	if err != nil {
		t.Fatalf("re-get: %v", err)
	}
	if back.RouterID != "10.99.99.98" {
		t.Fatalf("update not persisted on RouterID: got %v, want %v", back.RouterID, "10.99.99.98")
	}
}

func TestOSPF6Settings(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	orig, err := c.GetOSPF6Settings(ctx)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	before := orig.RouterID
	t.Cleanup(func() {
		orig.RouterID = before
		_ = c.UpdateOSPF6Settings(ctx, orig)
	})

	orig.RouterID = "10.99.99.97"
	if err := c.UpdateOSPF6Settings(ctx, orig); err != nil {
		t.Fatalf("update: %v", err)
	}

	back, err := c.GetOSPF6Settings(ctx)
	if err != nil {
		t.Fatalf("re-get: %v", err)
	}
	if back.RouterID != "10.99.99.97" {
		t.Fatalf("update not persisted on RouterID: got %v, want %v", back.RouterID, "10.99.99.97")
	}
}

func TestRIPSettings(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	orig, err := c.GetRIPSettings(ctx)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	before := orig.DefaultMetric
	t.Cleanup(func() {
		orig.DefaultMetric = before
		_ = c.UpdateRIPSettings(ctx, orig)
	})

	orig.DefaultMetric = "3"
	if err := c.UpdateRIPSettings(ctx, orig); err != nil {
		t.Fatalf("update: %v", err)
	}

	back, err := c.GetRIPSettings(ctx)
	if err != nil {
		t.Fatalf("re-get: %v", err)
	}
	if back.DefaultMetric != "3" {
		t.Fatalf("update not persisted on DefaultMetric: got %v, want %v", back.DefaultMetric, "3")
	}
}

func TestStaticSettings(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	orig, err := c.GetStaticSettings(ctx)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	before := orig.Enabled
	t.Cleanup(func() {
		orig.Enabled = before
		_ = c.UpdateStaticSettings(ctx, orig)
	})

	orig.Enabled = "0"
	if err := c.UpdateStaticSettings(ctx, orig); err != nil {
		t.Fatalf("update: %v", err)
	}

	back, err := c.GetStaticSettings(ctx)
	if err != nil {
		t.Fatalf("re-get: %v", err)
	}
	if back.Enabled != "0" {
		t.Fatalf("update not persisted on Enabled: got %v, want %v", back.Enabled, "0")
	}
}
