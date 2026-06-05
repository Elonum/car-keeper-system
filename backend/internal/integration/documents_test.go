package integration_test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/carkeeper/backend/internal/testsupport"
)

func TestDocuments_ManagerUploadAppearsInAccessibleLists(t *testing.T) {
	customerToken := registerFreshCustomer(t)
	cfg := createDraftConfiguration(t, customerToken)
	cfgID, _ := cfg["configuration_id"].(string)

	rr, resp := testsupport.DoJSON(t, testHandler, http.MethodPost, "/api/orders", map[string]any{
		"configuration_id": cfgID,
	}, customerToken)
	if rr.Code != http.StatusOK || !resp.Success {
		t.Fatalf("create order status=%d err=%s", rr.Code, resp.Error)
	}
	order := testsupport.ParseDataMap(t, resp.Data)
	orderID, _ := order["order_id"].(string)
	if orderID == "" {
		t.Fatal("missing order_id")
	}

	managerToken := loginSeedUser(t, "manager@carkeeper.ru")
	rr, resp = uploadDocument(t, managerToken, map[string]string{
		"document_type": "order_contract",
		"order_id":      orderID,
	}, "contract.pdf", []byte("%PDF-1.4\n%test\n"))
	if rr.Code != http.StatusOK || !resp.Success {
		t.Fatalf("upload status=%d err=%s body=%s", rr.Code, resp.Error, rr.Body.String())
	}
	created := testsupport.ParseDataMap(t, resp.Data)
	documentID, _ := created["document_id"].(string)
	if documentID == "" {
		t.Fatal("missing document_id")
	}
	if created["file_available"] != true {
		t.Fatalf("created file_available=%v, want true", created["file_available"])
	}
	if created["owner_name"] == "" || created["owner_email"] == "" {
		t.Fatalf("created document missing owner context: %#v", created)
	}
	if created["attachment_kind"] != "order" || created["attachment_label"] == "" {
		t.Fatalf("created document missing attachment context: %#v", created)
	}

	assertDocumentInList(t, customerToken, documentID)
	assertDocumentInList(t, managerToken, documentID)
}

func uploadDocument(t *testing.T, token string, fields map[string]string, fileName string, content []byte) (*httptest.ResponseRecorder, testsupport.APIResponse) {
	t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			t.Fatal(err)
		}
	}
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/documents/", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rr := httptest.NewRecorder()
	testHandler.ServeHTTP(rr, req)

	var resp testsupport.APIResponse
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	return rr, resp
}

func assertDocumentInList(t *testing.T, token, documentID string) {
	t.Helper()
	rr, resp := testsupport.DoJSON(t, testHandler, http.MethodGet, "/api/documents", nil, token)
	if rr.Code != http.StatusOK || !resp.Success {
		t.Fatalf("list documents status=%d err=%s", rr.Code, resp.Error)
	}
	list := testsupport.ParseDataArray(t, resp.Data)
	for _, doc := range list {
		if doc["document_id"] == documentID {
			if doc["file_available"] != true {
				t.Fatalf("listed file_available=%v, want true", doc["file_available"])
			}
			if doc["owner_name"] == "" || doc["attachment_label"] == "" {
				t.Fatalf("listed document missing context: %#v", doc)
			}
			return
		}
	}
	t.Fatalf("document %s not found in list (%d docs)", documentID, len(list))
}
