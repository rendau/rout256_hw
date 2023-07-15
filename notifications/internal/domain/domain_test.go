package domain

import (
	"errors"
	"route256/notifications/internal/domain/mocks"
	"route256/notifications/internal/domain/models"
	repoMocks "route256/notifications/internal/repo/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestListOrderStatusEvent_noCacheHasItems_success(t *testing.T) {
	repo := &repoMocks.Repo{}
	repo.On("OrderStatusEventList", mock.Anything, mock.Anything).Return([]*models.OrderStatusEventSt{
		{
			TS:      time.Time{},
			OrderID: 7,
			Status:  "created",
		},
	}, nil)

	cache := mocks.ICache{}
	cache.On("GetJsonObj", mock.Anything, mock.Anything, mock.Anything).Return(false, nil)
	cache.On("SetJsonObj", mock.Anything, mock.Anything, []*models.OrderStatusEventSt{
		{
			TS:      time.Time{},
			OrderID: 7,
			Status:  "created",
		},
	}, mock.Anything).Return(nil)

	domain := New(repo, &cache, nil, "")

	result, err := domain.ListOrderStatusEvent(nil, &models.OrderStatusEventListParsSt{})
	require.Nil(t, err)
	require.Equal(t, []*models.OrderStatusEventSt{
		{
			TS:      time.Time{},
			OrderID: 7,
			Status:  "created",
		},
	}, result)

	cache.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestListOrderStatusEvent_hasCache_success(t *testing.T) {
	repo := &repoMocks.Repo{}

	cache := mocks.ICache{}
	cache.On("GetJsonObj", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)

	domain := New(repo, &cache, nil, "")

	result, err := domain.ListOrderStatusEvent(nil, &models.OrderStatusEventListParsSt{})
	require.Nil(t, err)
	require.Equal(t, []*models.OrderStatusEventSt{}, result)

	cache.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestListOrderStatusEvent_repoFails_error(t *testing.T) {
	repoErr := errors.New("repo fails")

	repo := &repoMocks.Repo{}
	repo.On("OrderStatusEventList", mock.Anything, mock.Anything).Return(nil, repoErr)

	cache := mocks.ICache{}
	cache.On("GetJsonObj", mock.Anything, mock.Anything, mock.Anything).Return(false, nil)

	domain := New(repo, &cache, nil, "")

	_, err := domain.ListOrderStatusEvent(nil, &models.OrderStatusEventListParsSt{})
	require.NotNil(t, err)
	require.True(t, errors.Is(err, repoErr))

	cache.AssertExpectations(t)
	repo.AssertExpectations(t)
}
