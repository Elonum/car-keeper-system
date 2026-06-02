package integration_test

import (
	"net/http"
	"testing"

	"github.com/carkeeper/backend/internal/testsupport"
)

func TestOrder_FullFlow_ConfigurationBecomesOrdered(t *testing.T) {
	token := registerFreshCustomer(t)
	cfg := createDraftConfiguration(t, token)
	cfgID, _ := cfg["configuration_id"].(string)

	rr, resp := testsupport.DoJSON(t, testHandler, http.MethodPost, "/api/orders", map[string]any{
		"configuration_id": cfgID,
	}, token)
	if rr.Code != http.StatusOK || !resp.Success {
		t.Fatalf("order status=%d err=%s", rr.Code, resp.Error)
	}
	order := testsupport.ParseDataMap(t, resp.Data)
	if order["status"] != "pending" {
		t.Fatalf("order status=%v", order["status"])
	}

	rr, resp = testsupport.DoJSON(t, testHandler, http.MethodGet, "/api/configurator/configurations/"+cfgID, nil, token)
	if rr.Code != http.StatusOK || !resp.Success {
		t.Fatalf("get config status=%d err=%s", rr.Code, resp.Error)
	}
	updated := testsupport.ParseDataMap(t, resp.Data)
	if updated["status"] != "ordered" {
		t.Fatalf("configuration status=%v want ordered", updated["status"])
	}

	rr, resp = testsupport.DoJSON(t, testHandler, http.MethodPost, "/api/orders", map[string]any{
		"configuration_id": cfgID,
	}, token)
	if rr.Code != http.StatusBadRequest && rr.Code != http.StatusConflict {
		t.Fatalf("duplicate order expected 400/409, got %d err=%s", rr.Code, resp.Error)
	}
}

func TestOrder_ListMyOrdersAfterCreate(t *testing.T) {
	token := registerFreshCustomer(t)
	cfg := createDraftConfiguration(t, token)
	cfgID, _ := cfg["configuration_id"].(string)

	_, resp := testsupport.DoJSON(t, testHandler, http.MethodPost, "/api/orders", map[string]any{
		"configuration_id": cfgID,
	}, token)
	if !resp.Success {
		t.Fatalf("create order failed: %s", resp.Error)
	}

	rr, resp := testsupport.DoJSON(t, testHandler, http.MethodGet, "/api/orders", nil, token)
	if rr.Code != http.StatusOK || !resp.Success {
		t.Fatalf("list orders status=%d err=%s", rr.Code, resp.Error)
	}
	orders := testsupport.ParseDataArray(t, resp.Data)
	if len(orders) < 1 {
		t.Fatal("expected at least one order")
	}
}
