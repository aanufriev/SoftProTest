package grpcserver

import (
	context "context"
	"io"
	"net"
	"reflect"
	"testing"

	grpc "google.golang.org/grpc"
)

type fakeStorage struct{}

func (fs fakeStorage) Ping() error {
	return nil
}

func (fs fakeStorage) Save(sport string, line string) error {
	return nil
}

func (fs fakeStorage) GetLastLine(sport string) (float32, error) {
	switch sport {
	case "baseball":
		return 1.5, nil
	case "football":
		return 4.2, nil
	case "soccer":
		return 0.1, nil
	}

	return 0, nil
}

func TestGRPC(t *testing.T) {
	go func() {
		fs := fakeStorage{}
		lis, err := net.Listen("tcp", ":10000")
		if err != nil {
			t.Errorf("Test failed. Error: %v", err)
		}

		server := grpc.NewServer()

		RegisterSubscribeOnSportsLinesServer(server, newServer(fs))

		err = server.Serve(lis)
		if err != nil {
			t.Errorf("Test failed. Error: %v", err)
		}
	}()

	grcpConn, err := grpc.Dial(
		"localhost:10000",
		grpc.WithInsecure(),
	)
	if err != nil {
		t.Errorf("Test failed. Error: %v", err)
	}
	defer grcpConn.Close()

	sub := NewSubscribeOnSportsLinesClient(grcpConn)

	ctx := context.Background()
	client, err := sub.Subscribe(ctx)
	if err != nil {
		t.Errorf("Test failed. Error: %v", err)
	}

	sports := []string{"baseball", "football", "soccer"}

	err = client.Send(&Request{
		Sports:   sports,
		Interval: 1,
	})
	if err != nil {
		t.Errorf("Test failed. Error: %v", err)
	}

	expected := map[string]float32{"baseball": 1.5, "football": 4.2, "soccer": 0.1}
	out, err := client.Recv()
	if err == io.EOF {
		t.Errorf("Test failed. Error: %v", err)
	} else if err != nil {
		t.Errorf("Test failed. Error: %v", err)
	}

	if !reflect.DeepEqual(expected, out.Lines) {
		t.Errorf("Test failed. Results not match\nGot:\n%v\nExpected:\n%v", out.Lines, expected)
	}
}
