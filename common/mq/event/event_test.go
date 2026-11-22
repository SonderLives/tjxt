package event

import (
	"encoding/json"
	"testing"
	"time"
)

// 事件 json 字段名是跨服务契约：pay/trade/learning/search 之间靠它互通，
// 任何一侧改动字段名都会静默破坏链路，这里锁定字段名防止无意破坏。

func TestPaySuccessEventJSONContract(t *testing.T) {
	data, err := json.Marshal(PaySuccessEvent{
		BizOrderNo: 1,
		PayOrderNo: 2,
		Amount:     3,
		PayChannel: "mock",
		PayTime:    time.Unix(1700000000, 0).UTC(),
	})
	if err != nil {
		t.Fatal(err)
	}
	want := `{"bizOrderNo":1,"payOrderNo":2,"amount":3,"payChannel":"mock","payTime":"2023-11-14T22:13:20Z"}`
	if string(data) != want {
		t.Fatalf("PaySuccessEvent JSON = %s, want %s", data, want)
	}
}

func TestOrderPayEventJSONContract(t *testing.T) {
	data, err := json.Marshal(OrderPayEvent{OrderBasic: OrderBasic{
		OrderID:    10,
		UserID:     20,
		CourseIDs:  []int64{30, 40},
		FinishTime: time.Unix(1700000000, 0).UTC(),
	}})
	if err != nil {
		t.Fatal(err)
	}
	want := `{"orderId":10,"userId":20,"courseIds":[30,40],"finishTime":"2023-11-14T22:13:20Z"}`
	if string(data) != want {
		t.Fatalf("OrderPayEvent JSON = %s, want %s", data, want)
	}

	// 反序列化侧同样锁定（消费方按字段名解析）
	var back OrderRefundEvent
	if err := json.Unmarshal(data, &back); err != nil {
		t.Fatal(err)
	}
	if back.OrderID != 10 || back.UserID != 20 || len(back.CourseIDs) != 2 {
		t.Fatalf("roundtrip mismatch: %+v", back)
	}
}

func TestCourseEventJSONContract(t *testing.T) {
	data, err := json.Marshal(CourseEvent{CourseID: 123})
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"courseId":123}` {
		t.Fatalf("CourseEvent JSON = %s", data)
	}
}
