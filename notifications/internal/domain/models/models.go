package models

import (
	"time"
)

type OrderStatusEventSt struct {
	TS      time.Time
	OrderID int64
	Status  string
}

type OrderStatusEventListParsSt struct {
	OrderID *int64
	TsGTE   *time.Time
	TsLTE   *time.Time
}
