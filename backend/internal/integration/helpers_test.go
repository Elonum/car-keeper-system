package integration_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/carkeeper/backend/internal/testsupport"
)

const (
	seedTrimID   = "80000000-0000-0000-0000-000000000001"
	seedColorID  = "90000000-0000-0000-0000-000000000001"
	seedBranchID = "10000000-0000-0000-0000-000000000001"
	seedServiceTypeID = "b0000000-0000-0000-0000-000000000003"
	seedUserCarID = "d0000000-0000-0000-0000-000000000001"
)

func createDraftConfiguration(t *testing.T, token string) map[string]any {
	t.Helper()
	_, resp := testsupport.DoJSON(t, testHandler, http.MethodPost, "/api/configurator/configurations", map[string]any{
		"trim_id":    seedTrimID,
		"color_id":   seedColorID,
		"option_ids": []string{},
	}, token)
	cfg := testsupport.ParseDataMap(t, resp.Data)
	if cfg["configuration_id"] == nil || cfg["configuration_id"] == "" {
		t.Fatal("missing configuration_id")
	}
	return cfg
}

func createUserCar(t *testing.T, token string) string {
	t.Helper()
	_, resp := testsupport.DoJSON(t, testHandler, http.MethodPost, "/api/profile/cars", map[string]any{
		"trim_id":         seedTrimID,
		"color_id":        seedColorID,
		"vin":             fmt.Sprintf("IT%dABCDEFGHIJK", time.Now().UnixNano()%1e11),
		"year":            2020,
		"current_mileage": 12000,
	}, token)
	car := testsupport.ParseDataMap(t, resp.Data)
	id, _ := car["user_car_id"].(string)
	if id == "" {
		t.Fatal("missing user_car_id")
	}
	return id
}

func nextWeekdayDate(t *testing.T) string {
	t.Helper()
	d := time.Now().UTC().Add(48 * time.Hour)
	for d.Weekday() == time.Sunday {
		d = d.Add(24 * time.Hour)
	}
	return d.Format("2006-01-02")
}

func firstAvailabilitySlot(t *testing.T, token, branchID, date, serviceTypeID string) string {
	t.Helper()
	path := fmt.Sprintf(
		"/api/service/branches/%s/availability?date=%s&service_type_ids=%s",
		branchID, date, serviceTypeID,
	)
	rr, resp := testsupport.DoJSON(t, testHandler, http.MethodGet, path, nil, token)
	if rr.Code != http.StatusOK || !resp.Success {
		t.Fatalf("availability status=%d err=%s", rr.Code, resp.Error)
	}
	data := testsupport.ParseDataMap(t, resp.Data)
	slots, _ := data["slots"].([]any)
	if len(slots) == 0 {
		t.Fatalf("no slots for %s", date)
	}
	slot, _ := slots[0].(string)
	if slot == "" {
		t.Fatal("empty slot")
	}
	return slot
}
