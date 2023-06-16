package domain

import (
	"testing"

	"route256/loms/internal/domain/models"
)

func TestDistributeItemsInStock(t *testing.T) {
	dm := &Domain{}

	// create sample input data
	items := []*models.OrderItemSt{
		{Sku: 1, Count: 5},
		{Sku: 2, Count: 8},
		{Sku: 3, Count: 15},
	}
	stock := []*models.StockSt{
		{WarehouseID: 1, Sku: 1, Count: 2},
		{WarehouseID: 2, Sku: 1, Count: 5},
		{WarehouseID: 1, Sku: 2, Count: 8},
		{WarehouseID: 2, Sku: 3, Count: 20},
	}
	expected := []*models.StockReserveSt{
		{WarehouseID: 1, Sku: 1, Count: 2},
		{WarehouseID: 2, Sku: 1, Count: 3},
		{WarehouseID: 1, Sku: 2, Count: 8},
		{WarehouseID: 2, Sku: 3, Count: 15},
	}
	expectedErr := error(nil)
	// call the function and get the actual result
	result, resultErr := dm.distributeItemsInStock(items, stock)
	if resultErr != expectedErr {
		t.Errorf("Expected error %v, but got %v", expectedErr, resultErr)
	}
	// check if the actual result matches the expected result
	if len(result) != len(expected) {
		t.Errorf("Expected %d results, but got %d", len(expected), len(result))
	}
	for i := range result {
		if result[i].WarehouseID != expected[i].WarehouseID {
			t.Errorf("Result[%d].WarehouseID = %d, expected %d", i, result[i].WarehouseID, expected[i].WarehouseID)
		}
		if result[i].Sku != expected[i].Sku {
			t.Errorf("Result[%d].Sku = %d, expected %d", i, result[i].Sku, expected[i].Sku)
		}
		if result[i].Count != expected[i].Count {
			t.Errorf("Result[%d].Count = %d, expected %d", i, result[i].Count, expected[i].Count)
		}
	}
	// test with empty items list
	items = []*models.OrderItemSt{}
	stock = []*models.StockSt{
		{WarehouseID: 1, Sku: 1, Count: 2},
		{WarehouseID: 2, Sku: 2, Count: 5},
	}
	expected = []*models.StockReserveSt{}
	expectedErr = nil
	result, resultErr = dm.distributeItemsInStock(items, stock)
	if resultErr != expectedErr {
		t.Errorf("Expected error %v, but got %v", expectedErr, resultErr)
	}
	if len(result) != len(expected) {
		t.Errorf("Expected %d results, but got %d", len(expected), len(result))
	}
	// test with empty stock list
	items = []*models.OrderItemSt{
		{Sku: 1, Count: 5},
		{Sku: 2, Count: 10},
	}
	stock = []*models.StockSt{}
	expected = []*models.StockReserveSt{}
	expectedErr = ErrStockInsufficient
	result, resultErr = dm.distributeItemsInStock(items, stock)
	if resultErr != expectedErr {
		t.Errorf("Expected error %v, but got %v", expectedErr, resultErr)
	}
	// test with items not matching any stock
	items = []*models.OrderItemSt{
		{Sku: 11, Count: 5},
		{Sku: 12, Count: 10},
	}
	stock = []*models.StockSt{
		{WarehouseID: 1, Sku: 1, Count: 2},
		{WarehouseID: 2, Sku: 2, Count: 5},
	}
	expected = []*models.StockReserveSt{}
	expectedErr = ErrStockInsufficient
	result, resultErr = dm.distributeItemsInStock(items, stock)
	if resultErr != expectedErr {
		t.Errorf("Expected error %v, but got %v", expectedErr, resultErr)
	}
	// test with stock count less than item count
	items = []*models.OrderItemSt{
		{Sku: 1, Count: 5},
		{Sku: 2, Count: 10},
	}
	stock = []*models.StockSt{
		{WarehouseID: 1, Sku: 1, Count: 2},
		{WarehouseID: 2, Sku: 1, Count: 2},
		{WarehouseID: 1, Sku: 2, Count: 8},
	}
	expected = []*models.StockReserveSt{
		{WarehouseID: 1, Sku: 1, Count: 2},
		{WarehouseID: 2, Sku: 1, Count: 2},
		{WarehouseID: 1, Sku: 2, Count: 8},
	}
	expectedErr = ErrStockInsufficient
	result, resultErr = dm.distributeItemsInStock(items, stock)
	if resultErr != expectedErr {
		t.Errorf("Expected error %v, but got %v", expectedErr, resultErr)
	}
	// test with stock count equal to item count
	items = []*models.OrderItemSt{
		{Sku: 1, Count: 5},
		{Sku: 2, Count: 10},
	}
	stock = []*models.StockSt{
		{WarehouseID: 1, Sku: 1, Count: 2},
		{WarehouseID: 2, Sku: 1, Count: 3},
		{WarehouseID: 1, Sku: 2, Count: 8},
		{WarehouseID: 2, Sku: 2, Count: 2},
	}
	expected = []*models.StockReserveSt{
		{WarehouseID: 1, Sku: 1, Count: 2},
		{WarehouseID: 2, Sku: 1, Count: 3},
		{WarehouseID: 1, Sku: 2, Count: 8},
		{WarehouseID: 2, Sku: 2, Count: 2},
	}
	expectedErr = nil
	result, resultErr = dm.distributeItemsInStock(items, stock)
	if resultErr != expectedErr {
		t.Errorf("Expected error %v, but got %v", expectedErr, resultErr)
	}
	if len(result) != len(expected) {
		t.Errorf("Expected %d results, but got %d", len(expected), len(result))
	}
	for i := range result {
		if result[i].WarehouseID != expected[i].WarehouseID {
			t.Errorf("Result[%d].WarehouseID = %d, expected %d", i, result[i].WarehouseID, expected[i].WarehouseID)
		}
		if result[i].Sku != expected[i].Sku {
			t.Errorf("Result[%d].Sku = %d, expected %d", i, result[i].Sku, expected[i].Sku)
		}
		if result[i].Count != expected[i].Count {
			t.Errorf("Result[%d].Count = %d, expected %d", i, result[i].Count, expected[i].Count)
		}
	}
	// test with stock count greater than item count
	items = []*models.OrderItemSt{
		{Sku: 1, Count: 5},
		{Sku: 2, Count: 10},
	}
	stock = []*models.StockSt{
		{WarehouseID: 1, Sku: 1, Count: 2},
		{WarehouseID: 2, Sku: 1, Count: 5},
		{WarehouseID: 1, Sku: 2, Count: 8},
		{WarehouseID: 2, Sku: 2, Count: 15},
	}
	expected = []*models.StockReserveSt{
		{WarehouseID: 1, Sku: 1, Count: 2},
		{WarehouseID: 2, Sku: 1, Count: 3},
		{WarehouseID: 1, Sku: 2, Count: 8},
		{WarehouseID: 2, Sku: 2, Count: 2},
	}
	expectedErr = nil
	result, resultErr = dm.distributeItemsInStock(items, stock)
	if resultErr != expectedErr {
		t.Errorf("Expected error %v, but got %v", expectedErr, resultErr)
	}
	if len(result) != len(expected) {
		t.Errorf("Expected %d results, but got %d", len(expected), len(result))
	}
	for i := range result {
		if result[i].WarehouseID != expected[i].WarehouseID {
			t.Errorf("Result[%d].WarehouseID = %d, expected %d", i, result[i].WarehouseID, expected[i].WarehouseID)
		}
		if result[i].Sku != expected[i].Sku {
			t.Errorf("Result[%d].Sku = %d, expected %d", i, result[i].Sku, expected[i].Sku)
		}
		if result[i].Count != expected[i].Count {
			t.Errorf("Result[%d].Count = %d, expected %d", i, result[i].Count, expected[i].Count)
		}
	}
}
