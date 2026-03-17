package provider

import (
	"testing"

	"github.com/jheddings/ccglow/internal/types"
)

func TestCostProvider(t *testing.T) {
	p := &costProvider{}
	sess := &types.SessionData{
		CWD:  "/tmp",
		Cost: &types.CostInfo{TotalCostUSD: 12.5},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	cost := result.Values["cost"].(map[string]any)
	if cost["usd"] != "$12.50" {
		t.Errorf("expected $12.50, got %s", cost["usd"])
	}
}

func TestCostProviderTotal(t *testing.T) {
	p := &costProvider{}
	sess := &types.SessionData{
		CWD:  "/tmp",
		Cost: &types.CostInfo{TotalCostUSD: 7.89},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	cost := result.Values["cost"].(map[string]any)
	if cost["total"] != 7.89 {
		t.Errorf("expected 7.89, got %v", cost["total"])
	}
	if result.Formats["cost.total"] != "$%.2f" {
		t.Errorf("expected format $%%.2f, got %s", result.Formats["cost.total"])
	}
}

func TestCostProviderNilCost(t *testing.T) {
	p := &costProvider{}
	sess := &types.SessionData{CWD: "/tmp"}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	cost := result.Values["cost"].(map[string]any)
	if cost["total"] != 0.0 {
		t.Errorf("expected 0.0, got %v", cost["total"])
	}
}
