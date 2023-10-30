package APIHandler

import (
	"bytes"
	"encoding/json"
	"testing"
)

func Test_JSONRPCRequest(t *testing.T) {
	request1, _ := NewJSONRPCRequest(1, "test", map[string]interface{}{
		"data1": 1,
		"data2": "2",
		"data3": 3.0,
	})

	if jsonBytes, marshalErr := json.Marshal(request1); marshalErr != nil {
		t.Fatalf("cannot encode to JSON: %f", marshalErr)
	} else {
		if request2, parseErr := ParseJSONRequest(bytes.NewBuffer(jsonBytes)); parseErr == nil {
			if request1.JSONRPC.Version != request2.JSONRPC.Version {
				t.Fatalf("not matched version: %s vs %s", request1.JSONRPC.Version, request2.JSONRPC.Version)
			}

			if request1.JSONRPC.IsIDNum() != request2.JSONRPC.IsIDNum() {
				t.Fatalf("not matched id format(number): %v vs %v", request1.JSONRPC.id, request2.JSONRPC.id)
			}
			if request1.JSONRPC.IsIDString() != request2.JSONRPC.IsIDString() {
				t.Fatalf("not matched id format(text): %v vs %v", request1.JSONRPC.id, request2.JSONRPC.id)
			}
			if request1.JSONRPC.IsIDNum() {
				id1, _ := request1.JSONRPC.IDNum()
				id2, _ := request2.JSONRPC.IDNum()

				if id1 != id2 {
					t.Fatalf("not matched id value(number): %f vs %f", id1, id2)
				}
			} else if request1.JSONRPC.IsIDString() {
				id1, _ := request1.JSONRPC.IDString()
				id2, _ := request2.JSONRPC.IDString()

				if id1 != id2 {
					t.Fatalf("not matched id value(text): %s vs %s", id1, id2)
				}
			}
		} else {
			t.Fatalf("cannot decode to JSON: %f", parseErr)
		}
	}
}
