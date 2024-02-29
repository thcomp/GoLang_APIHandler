package jsonrpc

import (
	"bytes"
	"testing"

	ThcompUtility "github.com/thcomp/GoLang_Utility"
)

func Test_JSONRPCRequest1(t *testing.T) {
	request1, _ := NewJSONRPCRequest(1, "test", map[string]interface{}{
		"data1": 1,
		"data2": "2",
		"data3": 3.0,
	})

	if jsonBytes, encodeErr := request1.EncodeByJSON(); encodeErr != nil {
		t.Fatalf("cannot encode to JSON: %f", encodeErr)
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

			req1ParamsMap, assertionOK1 := request1.Params.(map[string]interface{})
			req2ParamsMap, assertionOK2 := request2.Params.(map[string]interface{})

			if assertionOK1 && assertionOK2 {
				data1OfReq1ValueInf, existData1OfReq1 := req1ParamsMap["data1"]
				data2OfReq1ValueInf, existData2OfReq1 := req1ParamsMap["data2"]
				data3OfReq1ValueInf, existData3OfReq1 := req1ParamsMap["data3"]
				data1OfReq2ValueInf, existData1OfReq2 := req2ParamsMap["data1"]
				data2OfReq2ValueInf, existData2OfReq2 := req2ParamsMap["data2"]
				data3OfReq2ValueInf, existData3OfReq2 := req2ParamsMap["data3"]
				dataList := [][]interface{}{
					[]interface{}{data1OfReq1ValueInf, data2OfReq1ValueInf, data3OfReq1ValueInf},
					[]interface{}{data1OfReq2ValueInf, data2OfReq2ValueInf, data3OfReq2ValueInf},
				}
				existList := [][]bool{
					[]bool{existData1OfReq1, existData2OfReq1, existData3OfReq1},
					[]bool{existData1OfReq2, existData2OfReq2, existData3OfReq2},
				}

				for i := 0; i < len(dataList[0]); i++ {
					if existList[0][i] && existList[1][i] {
						helperOfReq1 := ThcompUtility.NewInterfaceHelper(dataList[0][i])
						helperOfReq2 := ThcompUtility.NewInterfaceHelper(dataList[1][i])

						if helperOfReq1.IsNumber() != helperOfReq2.IsNumber() {
							t.Fatalf("unknown value type(data%d): %d vs %d", i+1, helperOfReq1.GetKind(), helperOfReq2.GetKind())
						}
					} else {
						t.Fatalf("not exist value(data%d): %t vs %t", i+1, existList[0][i], existList[1][i])
					}
				}
			} else {
				t.Fatalf("not match params type: %v vs %v", request1.Params, request2.Params)
			}
		} else {
			t.Fatalf("cannot decode to JSON: %f", parseErr)
		}
	}
}

func Test_JSONRPCRequest2(t *testing.T) {
	request1, _ := NewJSONRPCRequest(1, "test", map[string]interface{}{
		"data1": 1,
		"data2": "2",
		"data3": 3.0,
	})
	type Params struct {
		Data1 int     `json:"data1"`
		Data2 string  `json:"data2"`
		Data3 float64 `json:"data3"`
	}

	if request1Bytes, encodeErr := request1.EncodeByJSON(); encodeErr == nil {
		request1ByteBUffer := bytes.NewBuffer(request1Bytes)
		if request2, parseErr := ParseJSONRequest(request1ByteBUffer); parseErr == nil {
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

			req1Params, req2Params := Params{}, Params{}
			parseErr1 := request1.ParseParams(&req1Params)
			parseErr2 := request1.ParseParams(&req2Params)

			if parseErr1 == nil && parseErr2 == nil {
				if req1Params.Data1 != req1Params.Data1 {
					t.Fatalf("not matched data1: %v, %v", req1Params.Data1, req2Params.Data1)
				}
				if req1Params.Data2 != req1Params.Data2 {
					t.Fatalf("not matched data2: %v, %v", req1Params.Data2, req2Params.Data2)
				}
				if req1Params.Data3 != req1Params.Data3 {
					t.Fatalf("not matched data3: %v, %v", req1Params.Data3, req2Params.Data3)
				}
			} else {
				t.Fatalf("fail to parse: %v, %v", parseErr1, parseErr2)
			}
		}
	} else {
		t.Fatalf("cannot encode to JSON: %v", encodeErr)
	}

}

func Test_JSONRPCResponse1(t *testing.T) {
	response1 := NewJSONRPCResponse(1)
	response1.Result = map[string]interface{}{
		"data1": 1,
		"data2": "2",
		"data3": 3.0,
	}

	if !response1.IsIDNum() {
		t.Fatalf("id type is not number: %v", response1.id)
	}

	tempResult1 := map[string]interface{}{}
	if parseErr := response1.ParseResult(&tempResult1); parseErr == nil {
		if dataInf, exist := tempResult1["data1"]; !exist {
			t.Fatalf("data1 is not exist: %v", tempResult1)
		} else {
			dataHelper := ThcompUtility.NewInterfaceHelper(dataInf)
			if data, ret := dataHelper.GetNumber(); !ret {
				t.Fatalf("data1 is not number")
			} else if data != 1 {
				t.Fatalf("data1 is not number 1: %v", data)
			}
		}
		if dataInf, exist := tempResult1["data2"]; !exist {
			t.Fatalf("data2 is not exist: %v", tempResult1)
		} else {
			dataHelper := ThcompUtility.NewInterfaceHelper(dataInf)
			if data, ret := dataHelper.GetString(); !ret {
				t.Fatalf("data2 is not text")
			} else if data != "2" {
				t.Fatalf("data2 is not text 2: %v", data)
			}
		}
		if dataInf, exist := tempResult1["data3"]; !exist {
			t.Fatalf("data3 is not exist: %v", tempResult1)
		} else {
			dataHelper := ThcompUtility.NewInterfaceHelper(dataInf)
			if data, ret := dataHelper.GetNumber(); !ret {
				t.Fatalf("data3 is not number")
			} else if data != 3.0 {
				t.Fatalf("data3 is not number 3.0: %v", data)
			}
		}
	} else {
		t.Fatalf("fail to parse result: %v", parseErr)
	}

	type Result struct {
		Data1 int     `json:"data1"`
		Data2 string  `json:"data2"`
		Data3 float64 `json:"data3"`
	}
	tempResult2 := Result{}
	if parseErr := response1.ParseResult(&tempResult2); parseErr == nil {
		if tempResult2.Data1 != 1 {
			t.Fatalf("data1 is not number 1: %v", tempResult2.Data1)
		}
		if tempResult2.Data2 != "2" {
			t.Fatalf("data2 is not text 2: %s", tempResult2.Data2)
		}
		if tempResult2.Data3 != 3.0 {
			t.Fatalf("data3 is not number 3.0: %f", tempResult2.Data3)
		}
	} else {
		t.Fatalf("fail to parse result: %v", parseErr)
	}
}

func Test_JSONRPCResponse2(t *testing.T) {
	request1, _ := NewJSONRPCRequest(1, "test", map[string]interface{}{
		"data1": 1,
		"data2": "2",
		"data3": 3.0,
	})
	response1 := NewJSONRPCResponseFromRequest(request1)
	response1.Error = NewJSONRPCError(JSONRPCInvalidRequest, "unknown method", request1.Params)

	if !response1.IsIDNum() {
		t.Fatalf("id type is not number: %v", response1.id)
	}

	tempResult1 := map[string]interface{}{}
	if parseErr := response1.ParseResult(&tempResult1); parseErr == nil {
		if dataInf, exist := tempResult1["data1"]; !exist {
			t.Fatalf("data1 is not exist: %v", tempResult1)
		} else {
			dataHelper := ThcompUtility.NewInterfaceHelper(dataInf)
			if data, ret := dataHelper.GetNumber(); !ret {
				t.Fatalf("data1 is not number")
			} else if data != 1 {
				t.Fatalf("data1 is not number 1: %v", data)
			}
		}
		if dataInf, exist := tempResult1["data2"]; !exist {
			t.Fatalf("data2 is not exist: %v", tempResult1)
		} else {
			dataHelper := ThcompUtility.NewInterfaceHelper(dataInf)
			if data, ret := dataHelper.GetString(); !ret {
				t.Fatalf("data2 is not text")
			} else if data != "2" {
				t.Fatalf("data2 is not text 2: %v", data)
			}
		}
		if dataInf, exist := tempResult1["data3"]; !exist {
			t.Fatalf("data3 is not exist: %v", tempResult1)
		} else {
			dataHelper := ThcompUtility.NewInterfaceHelper(dataInf)
			if data, ret := dataHelper.GetNumber(); !ret {
				t.Fatalf("data3 is not number")
			} else if data != 3.0 {
				t.Fatalf("data3 is not number 3.0: %v", data)
			}
		}
	} else {
		t.Fatalf("fail to parse result: %v", parseErr)
	}

	type Result struct {
		Data1 int     `json:"data1"`
		Data2 string  `json:"data2"`
		Data3 float64 `json:"data3"`
	}
	tempResult2 := Result{}
	if parseErr := response1.ParseResult(&tempResult2); parseErr == nil {
		if tempResult2.Data1 != 1 {
			t.Fatalf("data1 is not number 1: %v", tempResult2.Data1)
		}
		if tempResult2.Data2 != "2" {
			t.Fatalf("data2 is not text 2: %s", tempResult2.Data2)
		}
		if tempResult2.Data3 != 3.0 {
			t.Fatalf("data3 is not number 3.0: %f", tempResult2.Data3)
		}
	} else {
		t.Fatalf("fail to parse result: %v", parseErr)
	}
}
