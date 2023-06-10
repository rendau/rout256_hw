package test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"route256/test/pkg/proto/checkout_v1"
	"route256/test/pkg/proto/loms_v1"
)

var (
	checkoutClient checkout_v1.CheckoutClient
	lomsClient     loms_v1.LomsClient
)

func TestMain(m *testing.M) {
	viper.AutomaticEnv()

	checkoutUrl := viper.GetString("SERVICES_CHECKOUT_URL")
	lomsUrl := viper.GetString("SERVICES_LOMS_URL")

	conn, err := grpc.Dial(checkoutUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	checkoutClient = checkout_v1.NewCheckoutClient(conn)

	conn, err = grpc.Dial(lomsUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	lomsClient = loms_v1.NewLomsClient(conn)

	os.Exit(m.Run())
}

func TestGrpc(t *testing.T) {
	ctx := context.Background()

	checkoutAddToCartResponse, err := checkoutClient.AddToCart(ctx, &checkout_v1.AddToCartRequest{
		User:  7,
		Sku:   1,
		Count: 1,
	})
	require.Nil(t, err)
	require.NotNil(t, checkoutAddToCartResponse)

	checkoutDeleteFromCartResponse, err := checkoutClient.DeleteFromCart(ctx, &checkout_v1.DeleteFromCartRequest{
		User:  7,
		Sku:   1,
		Count: 1,
	})
	require.Nil(t, err)
	require.NotNil(t, checkoutDeleteFromCartResponse)

	checkoutListCartResponse, err := checkoutClient.ListCart(ctx, &checkout_v1.ListCartRequest{
		User: 7,
	})
	require.Nil(t, err)
	require.NotNil(t, checkoutListCartResponse)

	fmt.Printf("%+v\n", checkoutListCartResponse)

	checkoutPurchaseResponse, err := checkoutClient.Purchase(ctx, &checkout_v1.PurchaseRequest{
		User: 7,
	})
	require.Nil(t, err)
	require.NotNil(t, checkoutPurchaseResponse)
}
