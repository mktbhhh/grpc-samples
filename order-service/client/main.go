package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "ordermgt/client/ecommerce"

	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	/// setting up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewOrderManagementClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Add Order
	order1 := pb.Order{
		Id:          "101",
		Items:       []string{"iPhone XS", "Mac Book Pro"},
		Destination: "San Jose, CA",
		Price:       2300.00,
	}
	order2 := pb.Order{
		Id:          "102",
		Items:       []string{"Google Pixel 3A", "Mac Book Pro"},
		Destination: "Mountain View, CA",
		Price:       1800.00,
	}
	order3 := pb.Order{
		Id:          "103",
		Items:       []string{"Apple Watch S4"},
		Destination: "San Jose, CA",
		Price:       400.00,
	}
	order4 := pb.Order{
		Id:          "104",
		Items:       []string{"Google Home Mini", "Google Nest Hub"},
		Destination: "Mountain View, CA",
		Price:       400.00,
	}
	order5 := pb.Order{
		Id:          "105",
		Items:       []string{"Amazon Echo"},
		Destination: "San Jose, CA",
		Price:       30.00,
	}
	order6 := pb.Order{
		Id:          "106",
		Items:       []string{"Amazon Echo", "Apple iPhone XS"},
		Destination: "Mountain View, CA",
		Price:       300.00,
	}

	orders := []*pb.Order{}
	orders = append(orders, &order1)
	orders = append(orders, &order2)
	orders = append(orders, &order3)
	orders = append(orders, &order4)
	orders = append(orders, &order5)
	orders = append(orders, &order6)

	for _, order := range orders {
		res, _ := client.AddOrder(ctx, order)
		if res != nil {
			log.Print("AddOrder Response -> ", res.Value)
		}
	}

	// Get Order
	retrievedOrder, err := client.GetOrder(ctx, &wrappers.StringValue{Value: "105"})
	log.Print("GetOrder Response -> : ", retrievedOrder)
	if err != nil {
		log.Fatalf("Could not get product: %v", err)
	}

	// Search Order : Server streaming scenario 设想
	searchStream, _ := client.SearchOrders(ctx, &wrappers.StringValue{Value: "Google"})
	for {
		searchOrder, err := searchStream.Recv()
		if err == io.EOF {
			log.Print("EOF")
			break
		}

		if err == nil {
			log.Print("Search Result : ", searchOrder)
		}
	}

	// Update Orders : Client streaming scenario
	updOrder1 := pb.Order{
		Id:          "102",
		Items:       []string{"Google Pixel 3A", "Google Pixel Book"},
		Destination: "Mountain View, CA",
		Price:       1100.00,
	}
	updOrder2 := pb.Order{
		Id:          "103",
		Items:       []string{"Apple Watch S4", "Mac Book Pro", "iPad Pro"},
		Destination: "San Jose, CA",
		Price:       2800.00,
	}
	updOrder3 := pb.Order{
		Id:          "104",
		Items:       []string{"Google Home Mini", "Google Nest Hub", "iPad Mini"},
		Destination: "Mountain View, CA",
		Price:       2200.00,
	}

	updateStream, err := client.UpdateOrders(ctx)

	if err != nil {
		log.Fatalf("%v.UpdateOrders(_) = _, %v", client, err)
	}

	// Updating order 1
	if err := updateStream.Send(&updOrder1); err != nil {
		log.Fatalf("%v.Send(%v) = %v", updateStream, updOrder1, err)
	}

	// Updating order 2
	if err := updateStream.Send(&updOrder2); err != nil {
		log.Fatalf("%v.Send(%v) = %v", updateStream, updOrder2, err)
	}

	// Updating order 1
	if err := updateStream.Send(&updOrder3); err != nil {
		log.Fatalf("%v.Send(%v) = %v", updateStream, updOrder3, err)
	}

	updateRes, err := updateStream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", updateStream, err, nil)
	}
	log.Printf("Update Orders Res : %s", updateRes)
}
