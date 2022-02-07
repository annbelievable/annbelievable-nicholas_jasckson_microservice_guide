package main

import (
	"fmt"
	"testing"

	"github.com/annbelievable/nicholas_jackson_microservice_guide/product-api/sdk/client"
	"github.com/annbelievable/nicholas_jackson_microservice_guide/product-api/sdk/client/products"
)

func TestOurClient(t *testing.T) {
	cfg := client.DefaultTransportConfig().WithHost("localhost:9090")
	c := client.NewHTTPClientWithConfig(nil, cfg)

	params := products.NewListProductsParams()
	prod, err := c.Products.ListProducts(params)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%#v", prod.GetPayLoad()[0])
	t.Fail()
}
