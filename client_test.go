package netangelsgo

import (
	"fmt"
	"testing"
)

// Plot in your own api details for testing.
var fixture = CreateNetangelsClient("", "")

type testData struct {
	domain  string
	data    string
	data2   string
	bogusIP string
}

func TestAll(t *testing.T) {
	data := testData{
		domain:  "test2.ddns.example.dk",
		data:    "test_txt_data",
		data2:   "test_txt_data_2",
		bogusIP: "19.19.19.20",
	}
	testAdd(t, data)
	id := testGet(t, data)
	//	testUpdate(t, data, id)
	testRemove(t, data, id)
	//	testddns(t, data)
	//	testddnsWithoutIp(t, data)

}

func testAdd(t *testing.T, data testData) {
	id, err := fixture.AddRecord(data.domain, data.data, "TXT", 300)
	if err != nil {
		t.Fail()
	}
	if id == 0 {
		t.Fail()
	}
	fmt.Println(id)
}

// func testUpdate(t *testing.T, data testData, id int) {
// 	res, err := fixture.UpdateRecord(id, data.domain, data.data2, "TXT")
// 	if err != nil {
// 		t.Fail()
// 	}
// 	if res != true {
// 		t.Fail()
// 	}
// 	fmt.Println(id)
// }

func testRemove(t *testing.T, data testData, id int) {
	res2, _ := fixture.GetRecordID(data.domain, data.data2, "TXT")

	if res2 != id {
		t.Fail()
	}

	res := fixture.RemoveRecord(id)
	if res != nil {
		t.Fail()
	}

}
func testGet(t *testing.T, data testData) int {
	id, _ := fixture.GetRecordID(data.domain, "", "TXT")
	if id == 0 {
		t.Fail()
	}
	return id
}
