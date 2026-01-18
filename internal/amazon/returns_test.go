package amazon

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidateReturnReason(t *testing.T) {
	tests := []struct {
		name    string
		reason  string
		wantErr bool
	}{
		{
			name:    "valid defective reason",
			reason:  "defective",
			wantErr: false,
		},
		{
			name:    "valid wrong_item reason",
			reason:  "wrong_item",
			wantErr: false,
		},
		{
			name:    "valid not_as_described reason",
			reason:  "not_as_described",
			wantErr: false,
		},
		{
			name:    "valid no_longer_needed reason",
			reason:  "no_longer_needed",
			wantErr: false,
		},
		{
			name:    "valid better_price reason",
			reason:  "better_price",
			wantErr: false,
		},
		{
			name:    "valid other reason",
			reason:  "other",
			wantErr: false,
		},
		{
			name:    "invalid reason",
			reason:  "invalid_reason",
			wantErr: true,
		},
		{
			name:    "empty reason",
			reason:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateReturnReason(tt.reason)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateReturnReason() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_GetReturnableItems(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	tests := []struct {
		name    string
		client  *Client
		wantErr bool
	}{
		{
			name:    "nil client returns error",
			client:  nil,
			wantErr: true,
		},
		{
			name: "valid client returns items",
			client: &Client{
				httpClient: &http.Client{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items, err := tt.client.GetReturnableItems()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetReturnableItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && items == nil {
				t.Errorf("GetReturnableItems() returned nil items")
			}
		})
	}
}

func TestClient_GetReturnOptions(t *testing.T) {
	tests := []struct {
		name    string
		client  *Client
		orderID string
		itemID  string
		wantErr bool
	}{
		{
			name:    "nil client returns error",
			client:  nil,
			orderID: "123-4567890-1234567",
			itemID:  "ITEM123",
			wantErr: true,
		},
		{
			name: "empty orderID returns error",
			client: &Client{
				httpClient: &http.Client{},
			},
			orderID: "",
			itemID:  "ITEM123",
			wantErr: true,
		},
		{
			name: "empty itemID returns error",
			client: &Client{
				httpClient: &http.Client{},
			},
			orderID: "123-4567890-1234567",
			itemID:  "",
			wantErr: true,
		},
		{
			name: "valid parameters",
			client: &Client{
				httpClient: &http.Client{},
			},
			orderID: "123-4567890-1234567",
			itemID:  "ITEM123",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options, err := tt.client.GetReturnOptions(tt.orderID, tt.itemID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetReturnOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && options == nil {
				t.Errorf("GetReturnOptions() returned nil options")
			}
		})
	}
}

func TestClient_CreateReturn(t *testing.T) {
	tests := []struct {
		name    string
		client  *Client
		orderID string
		itemID  string
		reason  ReturnReason
		wantErr bool
	}{
		{
			name:    "nil client returns error",
			client:  nil,
			orderID: "123-4567890-1234567",
			itemID:  "ITEM123",
			reason:  ReasonDefective,
			wantErr: true,
		},
		{
			name: "empty orderID returns error",
			client: &Client{
				httpClient: &http.Client{},
			},
			orderID: "",
			itemID:  "ITEM123",
			reason:  ReasonDefective,
			wantErr: true,
		},
		{
			name: "empty itemID returns error",
			client: &Client{
				httpClient: &http.Client{},
			},
			orderID: "123-4567890-1234567",
			itemID:  "",
			reason:  ReasonDefective,
			wantErr: true,
		},
		{
			name: "invalid reason returns error",
			client: &Client{
				httpClient: &http.Client{},
			},
			orderID: "123-4567890-1234567",
			itemID:  "ITEM123",
			reason:  ReturnReason("invalid"),
			wantErr: true,
		},
		{
			name: "valid parameters with defective reason",
			client: &Client{
				httpClient: &http.Client{},
			},
			orderID: "123-4567890-1234567",
			itemID:  "ITEM123",
			reason:  ReasonDefective,
			wantErr: false,
		},
		{
			name: "valid parameters with wrong_item reason",
			client: &Client{
				httpClient: &http.Client{},
			},
			orderID: "123-4567890-1234567",
			itemID:  "ITEM123",
			reason:  ReasonWrongItem,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ret, err := tt.client.CreateReturn(tt.orderID, tt.itemID, tt.reason)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateReturn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if ret == nil {
					t.Errorf("CreateReturn() returned nil return")
				} else {
					if ret.OrderID != tt.orderID {
						t.Errorf("CreateReturn() orderID = %v, want %v", ret.OrderID, tt.orderID)
					}
					if ret.ItemID != tt.itemID {
						t.Errorf("CreateReturn() itemID = %v, want %v", ret.ItemID, tt.itemID)
					}
					if ret.Reason != string(tt.reason) {
						t.Errorf("CreateReturn() reason = %v, want %v", ret.Reason, string(tt.reason))
					}
					if ret.Status != "initiated" {
						t.Errorf("CreateReturn() status = %v, want initiated", ret.Status)
					}
				}
			}
		})
	}
}

func TestClient_GetReturnLabel(t *testing.T) {
	tests := []struct {
		name     string
		client   *Client
		returnID string
		wantErr  bool
	}{
		{
			name:     "nil client returns error",
			client:   nil,
			returnID: "RET123",
			wantErr:  true,
		},
		{
			name: "empty returnID returns error",
			client: &Client{
				httpClient: &http.Client{},
			},
			returnID: "",
			wantErr:  true,
		},
		{
			name: "valid returnID",
			client: &Client{
				httpClient: &http.Client{},
			},
			returnID: "RET123",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			label, err := tt.client.GetReturnLabel(tt.returnID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetReturnLabel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && label == nil {
				t.Errorf("GetReturnLabel() returned nil label")
			}
		})
	}
}

func TestClient_GetReturnStatus(t *testing.T) {
	tests := []struct {
		name     string
		client   *Client
		returnID string
		wantErr  bool
	}{
		{
			name:     "nil client returns error",
			client:   nil,
			returnID: "RET123",
			wantErr:  true,
		},
		{
			name: "empty returnID returns error",
			client: &Client{
				httpClient: &http.Client{},
			},
			returnID: "",
			wantErr:  true,
		},
		{
			name: "valid returnID",
			client: &Client{
				httpClient: &http.Client{},
			},
			returnID: "RET123",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, err := tt.client.GetReturnStatus(tt.returnID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetReturnStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && status == nil {
				t.Errorf("GetReturnStatus() returned nil status")
			}
		})
	}
}

func TestValidReturnReasons(t *testing.T) {
	expectedReasons := []ReturnReason{
		ReasonDefective,
		ReasonWrongItem,
		ReasonNotAsDescribed,
		ReasonNoLongerNeeded,
		ReasonBetterPrice,
		ReasonOther,
	}

	for _, reason := range expectedReasons {
		if _, exists := ValidReturnReasons[reason]; !exists {
			t.Errorf("ValidReturnReasons missing expected reason: %s", reason)
		}
	}

	if len(ValidReturnReasons) != len(expectedReasons) {
		t.Errorf("ValidReturnReasons length = %d, want %d", len(ValidReturnReasons), len(expectedReasons))
	}
}
