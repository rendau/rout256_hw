package domain

import (
	"encoding/json"
	"fmt"
	"route256/loms/internal/domain/models"
	"strconv"
)

func (d *Domain) NotificationSendOrderStatusChange(obj models.NotificationOrderStatusChangeSt) error {
	if d.eventProducerOrderStatusChange == nil {
		fmt.Println("Order status changed", obj)
		return nil
	}

	rawData, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	err = d.eventProducerOrderStatusChange.SendMessage(strconv.FormatInt(obj.OrderID, 10), rawData)
	if err != nil {
		return fmt.Errorf("eventProducerOrderStatusChange.SendMessage: %w", err)
	}

	return nil
}
