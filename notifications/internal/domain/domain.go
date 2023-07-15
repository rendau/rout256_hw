package domain

import (
	"context"
	"fmt"
	"log"
	"route256/notifications/internal/domain/models"
	"route256/notifications/internal/repo"
	"time"
)

type Domain struct {
	repo                           repo.Repo
	messanger                      IMessanger
	orderStatusChangeEventTemplate string
}

func New(repo repo.Repo, messanger IMessanger, orderStatusChangeEventTemplate string) *Domain {
	return &Domain{
		repo:                           repo,
		messanger:                      messanger,
		orderStatusChangeEventTemplate: orderStatusChangeEventTemplate,
	}
}

func (d *Domain) HandleOrderStatusEvent(obj *models.OrderStatusEventSt) error {
	ctx := context.Background()

	obj.TS = time.Now()

	err := d.repo.OrderStatusEventCreate(ctx, obj)
	if err != nil {
		return fmt.Errorf("repo.OrderStatusEventCreate: %w", err)
	}

	err = d.messanger.Send(fmt.Sprintf(d.orderStatusChangeEventTemplate, obj.OrderID, obj.Status))
	if err != nil {
		log.Println("messanger.Send: ", err)
	}

	return nil
}

func (d *Domain) ListOrderStatusEvent(ctx context.Context, pars *models.OrderStatusEventListParsSt) ([]*models.OrderStatusEventSt, error) {
	result, err := d.repo.OrderStatusEventList(ctx, pars)
	if err != nil {
		return nil, fmt.Errorf("repo.OrderStatusEventList: %w", err)
	}

	return result, nil
}
