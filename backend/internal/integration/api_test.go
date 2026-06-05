package integration_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/carkeeper/backend/config"
	"github.com/carkeeper/backend/database"
	"github.com/carkeeper/backend/internal/app"
	"github.com/carkeeper/backend/internal/handler"
	"github.com/carkeeper/backend/internal/repository"
	"github.com/carkeeper/backend/internal/service"
	"github.com/carkeeper/backend/internal/testsupport"
	"github.com/carkeeper/backend/internal/storage"
	"github.com/google/uuid"
)

var (
	testHandler http.Handler
	testCfg     *config.Config
)

func TestMain(m *testing.M) {
	cfg := config.TestConfig()

	db, err := database.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "skip integration: database unavailable: %v\n", err)
		os.Exit(0)
	}
	defer db.Close()

	if err := db.Ping(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "skip integration: database ping failed: %v\n", err)
		os.Exit(0)
	}

	docRoot := filepath.Join(os.TempDir(), "carkeeper-test-docs")
	_ = os.MkdirAll(docRoot, 0o755)
	cfg.Storage.RootPath = docRoot

	store, err := storage.NewLocal(docRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "skip integration: storage: %v\n", err)
		os.Exit(0)
	}

	repos := repository.New(db)
	service.BootstrapAuthz(context.Background(), repos)
	services := service.New(repos, cfg, store)
	handlers := handler.New(services, cfg)
	testHandler = app.NewRouter(handlers, cfg, db)
	testCfg = cfg

	os.Exit(m.Run())
}

func TestHealth_OK(t *testing.T) {
	rr, _ := testsupport.DoJSON(t, testHandler, http.MethodGet, "/health", nil, "")
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d body %s", rr.Code, rr.Body.String())
	}
}

func TestCatalog_Trims_Public(t *testing.T) {
	rr, resp := testsupport.DoJSON(t, testHandler, http.MethodGet, "/api/catalog/trims", nil, "")
	if rr.Code != http.StatusOK || !resp.Success {
		t.Fatalf("status=%d success=%v err=%s", rr.Code, resp.Success, resp.Error)
	}
}

func TestAuth_RegisterLoginMe(t *testing.T) {
	email := fmt.Sprintf("test_%s@carkeeper.test", uuid.NewString()[:8])
	pass := "TestPass123!"

	rr, resp := testsupport.DoJSON(t, testHandler, http.MethodPost, "/api/auth/register", map[string]any{
		"first_name": "Тест",
		"last_name":  "Пользователь",
		"email":      email,
		"password":   pass,
	}, "")
	if rr.Code != http.StatusOK || !resp.Success {
		t.Fatalf("register status=%d err=%s", rr.Code, resp.Error)
	}

	rr, resp = testsupport.DoJSON(t, testHandler, http.MethodPost, "/api/auth/login", map[string]any{
		"email":    email,
		"password": pass,
	}, "")
	if rr.Code != http.StatusOK || !resp.Success {
		t.Fatalf("login status=%d err=%s", rr.Code, resp.Error)
	}
	token := testsupport.TokenFromLogin(t, rr, resp.Data)

	rr, resp = testsupport.DoJSON(t, testHandler, http.MethodGet, "/api/auth/me", nil, token)
	if rr.Code != http.StatusOK || !resp.Success {
		t.Fatalf("me status=%d err=%s", rr.Code, resp.Error)
	}
	me := testsupport.ParseDataMap(t, resp.Data)
	if me["email"] != email {
		t.Fatalf("email mismatch: %v", me["email"])
	}
	if me["role"] != "customer" {
		t.Fatalf("role=%v", me["role"])
	}
}

func TestAuth_Login_InvalidPassword(t *testing.T) {
	rr, resp := testsupport.DoJSON(t, testHandler, http.MethodPost, "/api/auth/login", map[string]any{
		"email":    "customer@carkeeper.ru",
		"password": "wrong-password-xyz",
	}, "")
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}
	if resp.Success {
		t.Fatal("expected success=false")
	}
}

func TestRBAC_CustomerCannotListAdminOrders(t *testing.T) {
	token := loginSeedUser(t, "customer@carkeeper.ru")
	rr, resp := testsupport.DoJSON(t, testHandler, http.MethodGet, "/api/admin/orders", nil, token)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d err=%s", rr.Code, resp.Error)
	}
}

func TestRBAC_ManagerCanListAdminOrders(t *testing.T) {
	token := loginSeedUser(t, "manager@carkeeper.ru")
	rr, resp := testsupport.DoJSON(t, testHandler, http.MethodGet, "/api/admin/orders", nil, token)
	if rr.Code != http.StatusOK || !resp.Success {
		t.Fatalf("status=%d err=%s", rr.Code, resp.Error)
	}
}

func TestConfigurator_CreateConfiguration(t *testing.T) {
	token := registerFreshCustomer(t)
	trimID := "80000000-0000-0000-0000-000000000001"
	colorID := "90000000-0000-0000-0000-000000000001"

	rr, resp := testsupport.DoJSON(t, testHandler, http.MethodPost, "/api/configurator/configurations", map[string]any{
		"trim_id":    trimID,
		"color_id":   colorID,
		"option_ids": []string{},
	}, token)
	if rr.Code != http.StatusOK || !resp.Success {
		t.Fatalf("create config status=%d err=%s", rr.Code, resp.Error)
	}
	data := testsupport.ParseDataMap(t, resp.Data)
	if data["status"] != "draft" {
		t.Fatalf("status=%v", data["status"])
	}
}

func TestOrder_CreateFromConfiguration(t *testing.T) {
	token := registerFreshCustomer(t)
	trimID := "80000000-0000-0000-0000-000000000001"
	colorID := "90000000-0000-0000-0000-000000000001"

	_, resp := testsupport.DoJSON(t, testHandler, http.MethodPost, "/api/configurator/configurations", map[string]any{
		"trim_id": trimID, "color_id": colorID, "option_ids": []string{},
	}, token)
	cfg := testsupport.ParseDataMap(t, resp.Data)
	cfgID, _ := cfg["configuration_id"].(string)
	if cfgID == "" {
		t.Fatal("missing configuration_id")
	}

	rr, resp := testsupport.DoJSON(t, testHandler, http.MethodPost, "/api/orders", map[string]any{
		"configuration_id": cfgID,
		"final_price":      cfg["total_price"],
	}, token)
	if rr.Code != http.StatusOK || !resp.Success {
		t.Fatalf("order status=%d err=%s", rr.Code, resp.Error)
	}
	order := testsupport.ParseDataMap(t, resp.Data)
	if order["status"] != "pending" {
		t.Fatalf("order status=%v", order["status"])
	}
}

func TestService_AvailabilityRequiresParams(t *testing.T) {
	token := loginSeedUser(t, "customer@carkeeper.ru")
	branchID := "10000000-0000-0000-0000-000000000001"
	rr, _ := testsupport.DoJSON(t, testHandler, http.MethodGet,
		"/api/service/branches/"+branchID+"/availability", nil, token)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

const seedPassword = "password123"

func loginSeedUser(t *testing.T, email string) string {
	t.Helper()
	rr, resp := testsupport.DoJSON(t, testHandler, http.MethodPost, "/api/auth/login", map[string]any{
		"email": email, "password": seedPassword,
	}, "")
	if rr.Code != http.StatusOK || !resp.Success {
		t.Skipf("seed user %s unavailable (run schema.sql + seed.sql; password %q): status=%d err=%s",
			email, seedPassword, rr.Code, resp.Error)
	}
	return testsupport.TokenFromLogin(t, rr, resp.Data)
}

func registerFreshCustomerCredentials(t *testing.T) (email, pass string) {
	t.Helper()
	email = fmt.Sprintf("it_%d@carkeeper.test", time.Now().UnixNano())
	pass = "TestPass123!"
	_, resp := testsupport.DoJSON(t, testHandler, http.MethodPost, "/api/auth/register", map[string]any{
		"first_name": "Ит",
		"last_name":  "Тест",
		"email":      email,
		"password":   pass,
	}, "")
	if !resp.Success {
		t.Fatalf("register failed: %s", resp.Error)
	}
	return email, pass
}

func registerFreshCustomer(t *testing.T) string {
	t.Helper()
	email, pass := registerFreshCustomerCredentials(t)
	rr, resp := testsupport.DoJSON(t, testHandler, http.MethodPost, "/api/auth/login", map[string]any{
		"email": email, "password": pass,
	}, "")
	if rr.Code != http.StatusOK || !resp.Success {
		t.Fatalf("login failed: %d %s", rr.Code, resp.Error)
	}
	return testsupport.TokenFromLogin(t, rr, resp.Data)
}
