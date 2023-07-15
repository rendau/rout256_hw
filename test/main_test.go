package test

import (
	"context"
	"os"
	"route256/libs/logger"
	"route256/test/pkg/proto/checkout_v1"
	"route256/test/pkg/proto/loms_v1"
	"route256/test/pkg/proto/notifications_v1"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	checkoutClient      checkout_v1.CheckoutClient
	lomsClient          loms_v1.LomsClient
	notificationsClient notifications_v1.NotificationsClient
)

func TestMain(m *testing.M) {
	logger.Init("info", true)

	//time.Sleep(3 * time.Second)

	viper.AutomaticEnv()

	checkoutUrl := viper.GetString("SERVICES_CHECKOUT_URL")
	lomsUrl := viper.GetString("SERVICES_LOMS_URL")
	notificationsUrl := viper.GetString("SERVICES_NOTIFICATIONS_URL")

	// checkout-client
	conn, err := grpc.Dial(checkoutUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalw(nil, err, "grpc.Dial")
	}
	checkoutClient = checkout_v1.NewCheckoutClient(conn)

	// loms-client
	conn, err = grpc.Dial(lomsUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalw(nil, err, "grpc.Dial")
	}
	lomsClient = loms_v1.NewLomsClient(conn)

	// notifications-client
	conn, err = grpc.Dial(notificationsUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalw(nil, err, "grpc.Dial")
	}
	notificationsClient = notifications_v1.NewNotificationsClient(conn)

	os.Exit(m.Run())
}

func TestGrpc(t *testing.T) {
	var err error

	ctx := context.Background()

	const usrId int64 = 7

	const sku1 uint32 = 4678816
	const sku2 uint32 = 4288068
	const sku3 uint32 = 4487693

	stock := []*loms_v1.StockAddRequest{
		{WarehouseId: 1, Sku: sku1, Count: 10},
		{WarehouseId: 1, Sku: sku2, Count: 10},
		{WarehouseId: 1, Sku: sku3, Count: 10},
	}

	for _, s := range stock {
		_, err = lomsClient.StockRemove(ctx, &loms_v1.StockRemoveRequest{
			WarehouseId: s.WarehouseId,
			Sku:         s.Sku,
		})

		_, err = lomsClient.StockAdd(ctx, s)
		require.Nil(t, err, errMsg(err))
	}

	_, err = checkoutClient.DeleteFromCart(ctx, &checkout_v1.DeleteFromCartRequest{
		User:  usrId,
		Sku:   sku1,
		Count: 1000000,
	})
	require.Nil(t, err, errMsg(err))

	checkoutAddToCartResponse, err := checkoutClient.AddToCart(ctx, &checkout_v1.AddToCartRequest{
		User:  usrId,
		Sku:   sku1,
		Count: 7,
	})
	require.Nil(t, err, errMsg(err))
	require.NotNil(t, checkoutAddToCartResponse)

	checkoutDeleteFromCartResponse, err := checkoutClient.DeleteFromCart(ctx, &checkout_v1.DeleteFromCartRequest{
		User:  usrId,
		Sku:   sku1,
		Count: 2,
	})
	require.Nil(t, err, errMsg(err))
	require.NotNil(t, checkoutDeleteFromCartResponse)

	checkoutListCartResponse, err := checkoutClient.ListCart(ctx, &checkout_v1.ListCartRequest{
		User: usrId,
	})
	require.Nil(t, err, errMsg(err))
	require.NotNil(t, checkoutListCartResponse)
	require.Len(t, checkoutListCartResponse.Items, 1)
	cartItem := checkoutListCartResponse.Items[0]
	require.Equal(t, sku1, cartItem.Sku)
	require.Equal(t, uint32(5), cartItem.Count)

	skuStock, err := lomsClient.Stocks(ctx, &loms_v1.StocksRequest{
		Sku: sku1,
	})
	require.Nil(t, err, errMsg(err))
	require.NotNil(t, skuStock)
	require.Len(t, skuStock.Stocks, 1)
	require.Equal(t, uint64(10), skuStock.Stocks[0].Count)

	checkoutPurchaseResponse, err := checkoutClient.Purchase(ctx, &checkout_v1.PurchaseRequest{
		User: usrId,
	})
	require.Nil(t, err, errMsg(err))
	require.NotNil(t, checkoutPurchaseResponse)
	ordId := checkoutPurchaseResponse.OrderID
	require.Greater(t, ordId, int64(0))

	skuStock, err = lomsClient.Stocks(ctx, &loms_v1.StocksRequest{
		Sku: sku1,
	})
	require.Nil(t, err, errMsg(err))
	require.NotNil(t, skuStock)
	require.Len(t, skuStock.Stocks, 1)
	require.Equal(t, uint64(5), skuStock.Stocks[0].Count)

	_, err = lomsClient.CancelOrder(ctx, &loms_v1.CancelOrderRequest{
		OrderID: ordId,
	})
	require.Nil(t, err, errMsg(err))

	skuStock, err = lomsClient.Stocks(ctx, &loms_v1.StocksRequest{
		Sku: sku1,
	})
	require.Nil(t, err, errMsg(err))
	require.NotNil(t, skuStock)
	require.Len(t, skuStock.Stocks, 1)
	require.Equal(t, uint64(10), skuStock.Stocks[0].Count)

	time.Sleep(3 * time.Second) // wait for notifications

	notificationsListResponse, err := notificationsClient.ListOrderStatusEvent(ctx, &notifications_v1.ListOrderStatusEventRequest{
		TsGTE: timestamppb.New(time.Now().Add(-5 * time.Second)),
		TsLTE: timestamppb.New(time.Now().Add(5 * time.Second)),
	})
	require.Nil(t, err, errMsg(err))
	require.NotNil(t, notificationsListResponse)
	require.Greater(t, len(notificationsListResponse.Items), 0)
}

//func TestKafka(t *testing.T) {
//	brokers := viper.GetStringSlice("KAFKA_BROKERS")
//	groupId := viper.GetString("KAFKA_GROUP")
//	topic := viper.GetString("KAFKA_TOPIC")
//
//	var err error
//
//	ctx := context.Background()
//
//	// producer
//	producer, err := kafka_producer.NewKafkaProducer(kafka_producer.KafkaProducerConfig{
//		Brokers:        brokers,
//		Topic:          topic,
//		Compress:       false,
//		AssuranceLevel: kafka_producer.AssuranceLevelExactlyOnce,
//	})
//	require.Nil(t, err)
//	defer func() {
//		err = producer.Close()
//		require.Nil(t, err)
//	}()
//
//	srcMsg, err := uuid.GenerateUUID()
//	require.Nil(t, err)
//
//	msgHitChan := make(chan string, 0)
//
//	// consumer
//	consumer, err := kafka_consumer.NewKafkaConsumer(kafka_consumer.KafkaConsumerConfig{
//		Context: ctx,
//		Brokers: brokers,
//		GroupId: groupId,
//		Topic:   topic,
//		Handler: func(ctx context.Context, topic string, msg []byte) bool {
//			if string(msg) == srcMsg {
//				close(msgHitChan)
//			}
//			return true
//		},
//		RetryInterval: 3 * time.Second,
//	})
//	require.Nil(t, nil)
//	defer consumer.Stop()
//
//	// start consumer
//	consumer.Start()
//
//	// produce message
//	err = producer.SendMessage("", []byte(srcMsg))
//	require.Nil(t, err)
//
//	timer := time.NewTimer(10 * time.Second)
//	select {
//	case <-timer.C:
//		t.Fatal("timeout")
//	case <-msgHitChan:
//		if !timer.Stop() {
//			<-timer.C
//		}
//	}
//
//	fmt.Println("DONE")
//}

func errMsg(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
