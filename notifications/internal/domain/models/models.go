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
	TsGTE   *time.Time
	TsLTE   *time.Time
	OrderID *int64
}
