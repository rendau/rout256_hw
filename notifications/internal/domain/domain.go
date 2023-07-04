package domain

import (
	"fmt"
	"log"
)

type Domain struct {
	messanger                      IMessanger
	orderStatusChangeEventTemplate string
}

func New(messanger IMessanger, orderStatusChangeEventTemplate string) *Domain {
	return &Domain{
		messanger:                      messanger,
		orderStatusChangeEventTemplate: orderStatusChangeEventTemplate,
	}
}

type OrderStatusChangeEventSt struct {
	OrderID int64
	Status  string
}

func (d *Domain) HandleOrderStatusChangeEvent(obj OrderStatusChangeEventSt) error {
	err := d.messanger.Send(fmt.Sprintf(d.orderStatusChangeEventTemplate, obj.OrderID, obj.Status))
	if err != nil {
		log.Println("messanger.Send: ", err)
	}

	return nil
}
