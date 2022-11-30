package workceptor

import (
	"testing"
)

func TestIsCompatible(t *testing.T) {
	type versiontest struct {
		version      string
		isCompatible bool
		isOCP        bool
	}

	kw := &kubeUnit{}

	scenarios := []versiontest{
		// OCP compatible versions
		{version: "v4.12.0", isCompatible: true, isOCP: true},
		{version: "v4.11.16", isCompatible: true, isOCP: true},
		{version: "v4.10.44", isCompatible: true, isOCP: true},

		// OCP Z stream >
		{version: "v4.12.99", isCompatible: true, isOCP: true},
		{version: "v4.11.99", isCompatible: true, isOCP: true},
		{version: "v4.10.99", isCompatible: true, isOCP: true},

		// OCP Z stream <
		{version: "v4.11.15", isCompatible: false, isOCP: true},
		{version: "v4.10.43", isCompatible: false, isOCP: true},

		// OCP X stream >
		{version: "v5.12.0", isCompatible: true, isOCP: true},
		{version: "v5.11.16", isCompatible: true, isOCP: true},
		{version: "v5.10.44", isCompatible: true, isOCP: true},

		// OCP X stream <
		{version: "v3.12.0", isCompatible: false, isOCP: true},
		{version: "v3.11.16", isCompatible: false, isOCP: true},
		{version: "v3.10.44", isCompatible: false, isOCP: true},

		// K8S compatible versions
		{version: "v1.24.8", isCompatible: true, isOCP: false},
		{version: "v1.25.4", isCompatible: true, isOCP: false},
		{version: "v1.23.14", isCompatible: true, isOCP: false},

		// K8S Z stream >
		{version: "v1.24.99", isCompatible: true, isOCP: false},
		{version: "v1.25.99", isCompatible: true, isOCP: false},
		{version: "v1.23.99", isCompatible: true, isOCP: false},

		// K8S Z stream <
		{version: "v1.24.7", isCompatible: false, isOCP: false},
		{version: "v1.25.3", isCompatible: false, isOCP: false},
		{version: "v1.23.13", isCompatible: false, isOCP: false},

		// K8S X stream >
		{version: "v2.24.8", isCompatible: true, isOCP: false},
		{version: "v2.25.4", isCompatible: true, isOCP: false},
		{version: "v2.23.14", isCompatible: true, isOCP: false},

		// K8S X stream <
		{version: "v0.24.8", isCompatible: false, isOCP: false},
		{version: "v0.25.4", isCompatible: false, isOCP: false},
		{version: "v0.23.14", isCompatible: false, isOCP: false},
	}
	var comp bool
	for _, s := range scenarios {
		t.Logf("version: %s, isCompatible: %t, isOCP: %t", s.version, s.isCompatible, s.isOCP)
		if s.isOCP {
			comp = isCompatibleOCP(kw, s.version)
		} else {
			comp = isCompatibleK8S(kw, s.version)
		}
		if comp != s.isCompatible {
			t.Fatal()
		}
	}
}
