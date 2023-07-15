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
	cache                          ICache
	messanger                      IMessanger
	orderStatusChangeEventTemplate string
}

func New(repo repo.Repo, cache ICache, messanger IMessanger, orderStatusChangeEventTemplate string) *Domain {
	return &Domain{
		repo:                           repo,
		cache:                          cache,
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
	cacheKey := d.listOrderStatusEventGenerateCacheKey(pars)

	result := make([]*models.OrderStatusEventSt, 0)

	ok, err := d.cache.GetJsonObj(ctx, cacheKey, &result)
	if err == nil && ok {
		return result, nil
	}

	result, err = d.repo.OrderStatusEventList(ctx, pars)
	if err != nil {
		return nil, fmt.Errorf("repo.OrderStatusEventList: %w", err)
	}

	_ = d.cache.SetJsonObj(ctx, cacheKey, result, 0)

	return result, nil
}

func (d *Domain) listOrderStatusEventGenerateCacheKey(pars *models.OrderStatusEventListParsSt) string {
	return fmt.Sprintf("order_status_event_list:%v", pars)
}
