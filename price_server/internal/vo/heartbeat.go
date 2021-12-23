package vo

type HEARTBEAT_INFO struct {
	ExpectResources int   `json:"expect_resources"`
	ActualResources int   `json:"actual_resources"`
	LatestTimestamp int64 `json:"latest_timestamp"`
	Interval        int64 `json:"interval"`
}
