package utils

type FilterObj struct {
	FieldName string        `json:"field_name"`
	Value     string        `json:"value"`
	Values    []interface{} `json:"values"`
	Operation string        `json:"operation"`
	FieldType string        `json:"field_type"`
}

type RequestObj struct {
	PageNo        int            `json:"page_no"`
	PerPage       int            `json:"per_page"`
	Sort          string         `json:"sort"`
	Filter        []FilterObj    `json:"filter"`
	DefaultParam  []DefaultParam `json:"default_param"`
	GlobOperation string         `json:"glob_operation"`
}

type DefaultParam struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type AnalyticFilterByLog struct {
	PageNo  int    `json:"page_no"`
	PerPage int    `json:"per_page"`
	Sort    string `json:"sort"`
	SartDay int    `json:"start_day"`
	EndDay  int    `json:"end_day"`
	AppID   string `json:"app_id"`
}
type AnalyticFilterByGraph struct {
	SartDay   int      `json:"start_day"`
	EndDay    int      `json:"end_day"`
	AppID     string   `json:"app_id"`
	ServiceID []string `json:"service_id"`
	EndPontID []string `json:"endpoint_id"`
}
