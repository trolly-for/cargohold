package policy_test

import (
	"errors"
	"testing"

	"cargohold/internal/policy"
)

func TestProductionDeniesWrite(t *testing.T) {
	p := policy.New()
	err := p.CheckWrite("production", false)
	if err == nil {
		t.Fatal("expected error for production write, got nil")
	}
	if !errors.Is(err, policy.ErrDenied) {
		t.Fatalf("expected ErrDenied, got %v", err)
	}
}

func TestProductionDeniesWriteEvenWithConfirmation(t *testing.T) {
	p := policy.New()
	err := p.CheckWrite("production", true)
	if err == nil {
		t.Fatal("expected error for production write even with confirmation")
	}
}

func TestStagingRequiresConfirmation(t *testing.T) {
	p := policy.New()
	if err := p.CheckWrite("staging", false); err == nil {
		t.Fatal("expected error for staging write without confirmation")
	}
	if err := p.CheckWrite("staging", true); err != nil {
		t.Fatalf("expected no error for staging write with confirmation, got %v", err)
	}
}

func TestDevelopmentAllowsWriteWithoutConfirmation(t *testing.T) {
	p := policy.New()
	if err := p.CheckWrite("development", false); err != nil {
		t.Fatalf("unexpected error for dev write: %v", err)
	}
}

func TestUnknownEnvironmentAllowed(t *testing.T) {
	p := policy.New()
	if err := p.CheckWrite("custom-env", false); err != nil {
		t.Fatalf("expected unknown env to be allowed, got %v", err)
	}
}

func TestSetOverridesRule(t *testing.T) {
	p := policy.New()
	p.Set("development", policy.Rule{AllowWrite: false, RequireConfirmation: false})
	err := p.CheckWrite("development", false)
	if err == nil {
		t.Fatal("expected write to be denied after rule override")
	}
}

func TestNormalizationApplied(t *testing.T) {
	p := policy.New()
	// "PRODUCTION" should normalise to "production" and be denied.
	if p.IsWriteAllowed("PRODUCTION", true) {
		t.Fatal("expected production (uppercased) to be denied")
	}
}

func TestIsWriteAllowedHelper(t *testing.T) {
	p := policy.New()
	if p.IsWriteAllowed("prod", false) {
		t.Fatal("expected IsWriteAllowed to return false for prod")
	}
	if !p.IsWriteAllowed("dev", false) {
		t.Fatal("expected IsWriteAllowed to return true for dev")
	}
}
