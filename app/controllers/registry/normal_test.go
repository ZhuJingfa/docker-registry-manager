package registry

import (
	"testing"
)

func TestShort(t *testing.T) {
	s := "registry.alishui.com:443"
	short := "registry.alishui.com"

	t.Logf("short: %v", RegistryNameShort(s))
	if RegistryNameShort(s) != short {
		t.Fatalf("RegistryNameShort err", )
	}

	t.Logf("short: %v", FormatRegistryName(s))
	if FormatRegistryName(s) != s {
		t.Fatalf("FormatRegistryName err", )
	}

}

func TestLong(t *testing.T) {
	s := "aliyun.alishui.com:7788"
	short := "aliyun.alishui.com:7788"

	t.Logf("short: %v", RegistryNameShort(s))
	if RegistryNameShort(s) != short {
		t.Fatalf("RegistryNameShort err", )
	}

	t.Logf("short: %v", FormatRegistryName(s))
	if FormatRegistryName(s) != s {
		t.Fatalf("FormatRegistryName err", )
	}

}
