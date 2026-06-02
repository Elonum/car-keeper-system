package integration_test

import (
	"net/http"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/carkeeper/backend/internal/testsupport"
)

func TestBooking_ParallelSameSlot_RespectsConcurrentBays(t *testing.T) {
	tokenA := registerFreshCustomer(t)
	tokenB := registerFreshCustomer(t)

	carA := createUserCar(t, tokenA)
	carB := createUserCar(t, tokenB)

	date := nextWeekdayDate(t)
	slot := firstAvailabilitySlot(t, tokenA, seedBranchID, date, seedServiceTypeID)

	payloadA := map[string]any{
		"user_car_id":       carA,
		"branch_id":         seedBranchID,
		"appointment_date":  slot,
		"service_type_ids":  []string{seedServiceTypeID},
		"description":       "parallel test A",
	}
	payloadB := map[string]any{
		"user_car_id":       carB,
		"branch_id":         seedBranchID,
		"appointment_date":  slot,
		"service_type_ids":  []string{seedServiceTypeID},
		"description":       "parallel test B",
	}

	const attempts = 3
	codes := make(chan int, attempts)
	var wg sync.WaitGroup
	wg.Add(attempts)

	for i := 0; i < attempts; i++ {
		go func(n int) {
			defer wg.Done()
			tok := tokenA
			body := payloadA
			if n == 1 {
				tok = tokenB
				body = payloadB
			}
			rr, _ := testsupport.DoJSON(t, testHandler, http.MethodPost, "/api/service/appointments", body, tok)
			codes <- rr.Code
		}(i)
	}
	wg.Wait()
	close(codes)

	var okCount, conflictCount atomic.Int32
	for code := range codes {
		switch code {
		case http.StatusOK:
			okCount.Add(1)
		case http.StatusConflict:
			conflictCount.Add(1)
		default:
			t.Errorf("unexpected status %d", code)
		}
	}

	if okCount.Load() != 2 {
		t.Fatalf("expected 2 successful bookings (concurrent_bays=2), got %d", okCount.Load())
	}
	if conflictCount.Load() != 1 {
		t.Fatalf("expected 1 conflict, got %d", conflictCount.Load())
	}
}

func TestBooking_SeedUserCar_Appointment(t *testing.T) {
	token := loginSeedUser(t, "customer@carkeeper.ru")
	date := nextWeekdayDate(t)
	slot := firstAvailabilitySlot(t, token, seedBranchID, date, seedServiceTypeID)

	rr, resp := testsupport.DoJSON(t, testHandler, http.MethodPost, "/api/service/appointments", map[string]any{
		"user_car_id":       seedUserCarID,
		"branch_id":         seedBranchID,
		"appointment_date":  slot,
		"service_type_ids":  []string{seedServiceTypeID},
	}, token)
	if rr.Code != http.StatusOK || !resp.Success {
		t.Fatalf("appointment status=%d err=%s", rr.Code, resp.Error)
	}
	appt := testsupport.ParseDataMap(t, resp.Data)
	if appt["status"] != "scheduled" {
		t.Fatalf("status=%v", appt["status"])
	}
}
