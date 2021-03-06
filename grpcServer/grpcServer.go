package grpcserver

import (
	"io"
	"log"
	"net"
	reflect "reflect"
	"time"

	"github.com/aanufriev/SoftProTest/storage"
	"github.com/sirupsen/logrus"
	grpc "google.golang.org/grpc"
)

// StartSubServer creates a grps server to handle
// the clients subscription to the sports line
func StartSubServer(port string, storage storage.DatabaseInterface) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln("cant listet port", err)
	}

	server := grpc.NewServer()

	RegisterSubscribeOnSportsLinesServer(server, newServer(storage))

	logrus.WithFields(logrus.Fields{
		"port": port,
	}).Info("Starting grpc server")
	err = server.Serve(lis)
	if err != nil {
		logrus.WithError(err).Info("grps server error")
	}
}

type subscribeServer struct {
	storage storage.DatabaseInterface
}

func newServer(storage storage.DatabaseInterface) *subscribeServer {
	return &subscribeServer{
		storage: storage,
	}
}

func (s *subscribeServer) SubscribeOnSportsLines(stream SubscribeOnSportsLines_SubscribeOnSportsLinesServer) error {
	requestChan := make(chan *Request)
	errChan := make(chan error)

	go func() {
		for {
			request, err := stream.Recv()
			if err == io.EOF {
				errChan <- err
				return
			}
			if err != nil {
				errChan <- err
				return
			}

			requestChan <- request
			logrus.WithFields(logrus.Fields{
				"request": request,
			}).Info("Get request")
		}
	}()

	var request *Request
	var lines map[string]float32

	for {
		select {
		case err := <-errChan:
			return err
		case newRequest := <-requestChan:
			if request != nil {
				if reflect.DeepEqual(request.Sports, newRequest.Sports) {
					request.Interval = newRequest.Interval
					break
				}
			}
			request = newRequest
			lines = make(map[string]float32, len(request.Sports))
		default:
			if request != nil {
				for _, sport := range request.Sports {
					line, err := s.storage.GetLastLine(sport)
					if err != nil {
						return err
					}
					lines[sport] = line - lines[sport]
				}

				err := stream.Send(&Response{
					Lines: lines,
				})
				if err != nil {
					logrus.WithError(err).Info("can't send response")
				}

				logrus.WithFields(logrus.Fields{
					"lines": lines,
				}).Info("Send lines")
				time.Sleep(time.Duration(request.Interval) * time.Second)
			}
		}
	}
}
