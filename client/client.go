package main

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/tuthanhthien02/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	// client connection
	cc, err := grpc.Dial("localhost:50068", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	client := calculatorpb.NewCalculatorpbServiceClient(cc)

	// callSum(client)
	// callPND(client)
	// callAverage(client)
	// callFindMax(client)
	callSquareRoot(client, -4)
}

func callSum(client calculatorpb.CalculatorpbServiceClient) {
	log.Println("calling sum api")

	res, err := client.Sum(context.Background(), &calculatorpb.SumRequest{
		Num1: 5,
		Num2: 5,
	})

	if err != nil {
		log.Fatalf("error while calling Sum RPC: %v", err)
	}

	log.Printf("response from sum api: %v", res.GetResult())
}

func callPND(client calculatorpb.CalculatorpbServiceClient) {
	log.Println("calling PND api")

	stream, err := client.PrimeNumberDecomposition(context.Background(), &calculatorpb.PNDReuqest{
		Number: 120,
	})

	if err != nil {
		log.Fatalf("CallPND error %v", err)
	}

	for {
		res, revErr := stream.Recv()

		if revErr == io.EOF {
			log.Println("Server finish streaming")
			return
		}

		log.Printf("Primary number %v", res.GetResult())
	}
}

func callAverage(client calculatorpb.CalculatorpbServiceClient) {
	log.Println("calling Average api")
	stream, err := client.Average(context.Background())

	if err != nil {
		log.Fatalf("Call Average err %v", err)
	}

	listReq := []calculatorpb.AverageRequest{
		calculatorpb.AverageRequest{
			Num: 5,
		},
		calculatorpb.AverageRequest{
			Num: 10,
		},
		calculatorpb.AverageRequest{
			Num: 12,
		},
		calculatorpb.AverageRequest{
			Num: 3,
		},
		calculatorpb.AverageRequest{
			Num: 4.2,
		},
	}

	for _, req := range listReq {
		err := stream.Send(&req)

		if err != nil {
			log.Fatalf("send average request error: %v", err)
		}
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Receive average response error: %v", err)
	}

	log.Printf("Average response %v", res)
}

func callFindMax(client calculatorpb.CalculatorpbServiceClient) {
	log.Println("calling FindMax api")
	stream, err := client.FindMax(context.Background())

	if err != nil {
		log.Fatalf("call find nax err : %v", err)
	}

	listReq := []calculatorpb.FindMaxRequest{
		calculatorpb.FindMaxRequest{
			Num: 5,
		},
		calculatorpb.FindMaxRequest{
			Num: 10,
		},
		calculatorpb.FindMaxRequest{
			Num: 12,
		},
		calculatorpb.FindMaxRequest{
			Num: 3,
		},
		calculatorpb.FindMaxRequest{
			Num: 4,
		},
	}

	waitChannel := make(chan struct{})

	go func() {
		for _, req := range listReq {
			err := stream.Send(&req)

			if err != nil {
				log.Fatalf("Send find max request err %v", err)
				break
			}
			time.Sleep(1000 * time.Microsecond)
		}

		// stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()

			if err == io.EOF {
				log.Println("ending find max api")
				break
			}

			if err != nil {
				log.Fatalf("receiving find max err %v", err)
				break
			}

			log.Printf("max %v", res.GetMax())
		}
		close(waitChannel)
	}()

	<-waitChannel
}

func callSquareRoot(client calculatorpb.CalculatorpbServiceClient, num int32) {
	log.Println("calling SquareRoot api")

	res, err := client.Square(context.Background(), &calculatorpb.SquareRequest{
		Num: num,
	})

	if err != nil {
		log.Fatalf("call sum api error: %v\n", err)
		if errStatus, ok := status.FromError(err); ok {
			log.Printf("err msg : %v", errStatus.Message())
			log.Printf("err code : %v", errStatus.Code())

			if errStatus.Code() == codes.InvalidArgument {
				log.Printf("Invalid argument num %v", num)
				return
			}
		}
	}

	log.Printf("square root response %v\n", res.GetSquareRoot())
}
