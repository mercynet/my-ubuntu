package catalog

import "testing"

func TestVendorFlatpakReferenceIdentity(t *testing.T) {
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	item := *c.Item("navicat-premium-lite")
	if _, err = FlatpakReference(&item); err != nil {
		t.Fatal(err)
	}
	item.Flatpak = "another.application"
	if _, err = FlatpakReference(&item); err == nil {
		t.Fatal("mismatched app accepted")
	}
	item.FlatpakRef = "../escape"
	if _, err = FlatpakReference(&item); err == nil {
		t.Fatal("path escape accepted")
	}
}
