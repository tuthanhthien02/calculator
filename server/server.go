package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"time"

	"github.com/tuthanhthien02/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	calculatorpb.UnimplementedCalculatorpbServiceServer
}

// implement API Sum
func (*server) Sum(ctx context.Context, in *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	log.Println("Sum function is called")
	res := &calculatorpb.SumResponse{
		Result: in.GetNum1() + in.GetNum2(),
	}

	return res, nil
}

// implement API SumWithDeadLine
func (*server) SumWithDeadline(ctx context.Context, in *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	log.Println("SumWithDeadline function is called")

	for i := 0; i < 5; i++ {
		if ctx.Err() == context.Canceled {
			log.Println("Context canceled!")
			return nil, status.Errorf(codes.Canceled, "Client canceled request!")
		}
		time.Sleep(1 * time.Second)
	}

	res := &calculatorpb.SumResponse{
		Result: in.GetNum1() + in.GetNum2(),
	}

	return res, nil
}

// implement API PrimeNumberDecomposition
func (*server) PrimeNumberDecomposition(req *calculatorpb.PNDReuqest, stream calculatorpb.CalculatorpbService_PrimeNumberDecompositionServer) error {
	log.Println("PrimeNumberDecomposition function is called")

	k := int32(2)
	n := req.GetNumber()
	for n > 1 {
		if n%k == 0 {
			n = n / k
			// Send to client
			stream.Send(&calculatorpb.PNDResponse{
				Result: k,
			})
			time.Sleep(1000 * time.Millisecond)
		} else {
			k++
			log.Printf("k increase to %v", k)
		}
	}

	return nil
}

// implement API Average
func (*server) Average(stream calculatorpb.CalculatorpbService_AverageServer) error {
	log.Println("Average function is called")
	var total float32
	var count int

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// implement average
			res := &calculatorpb.AverageResponse{
				Result: total / float32(count),
			}

			return stream.SendAndClose(res)
		}

		if err != nil {
			log.Fatalf("err while receiving Average request %v", err)
		}

		log.Printf("receive num %v\n", req.GetNum())
		total += req.GetNum()
		count++
	}

}

// implement API FindMax
func (*server) FindMax(stream calculatorpb.CalculatorpbService_FindMaxServer) error {
	log.Println("FindMax function is called")
	max := int32(0)

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			log.Println("EOF..")
			return nil
		}

		if err != nil {
			log.Fatalf("err while receiving FindMax request %v", err)
			return err
		}

		num := int32(req.GetNum())
		log.Printf("receiving num %v\n", num)
		if num > max {
			max = num
		}

		err = stream.Send(&calculatorpb.FindMaxResponse{
			Max: max,
		})

		if err != nil {
			log.Fatalf("err while sending FindMax response %v", err)
			return err
		}
	}

}

// implement API Square
func (*server) Square(ctx context.Context, in *calculatorpb.SquareRequest) (*calculatorpb.SquareResponse, error) {
	log.Println("Square function is called")

	num := in.GetNum()

	if num < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Expect a positive number")
	}

	return &calculatorpb.SquareResponse{
		SquareRoot: math.Sqrt(float64(num)),
	}, nil
}

const PORT = "0.0.0.0:50068"

func main() {
	lis, err := net.Listen("tcp", PORT)

	// check if error occurred
	if err != nil {
		log.Fatalf("error while listen to port %v", err)
	}

	s := grpc.NewServer()

	calculatorpb.RegisterCalculatorpbServiceServer(s, &server{})

	fmt.Printf("Server is starting at %v\n", PORT)

	err = s.Serve(lis)

	// check if error occurred
	if err != nil {
		log.Fatalf("error while serve %v", err)
	}

}
