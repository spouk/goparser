package lib

import (
	"log"
	"io"
)

const (
	LOG_WORKERREQUEST_PREFIX = "[workerrequest-log]"
	LOG_WORKERREQUEST_FLAGS  = log.Ldate | log.Ltime | log.Lshortfile
)

type WorkerRequest struct {
	Name string
	p    *Parser
	//ChanStockRequest chan []SearchRequest
	//ChanStockResult  chan []SearchResult
	//ChanCommand      chan string
	log *log.Logger
}

func NewWorkerRequest(name string, logout io.Writer) *WorkerRequest {
	return &WorkerRequest{
		Name: name,
		log:  log.New(logout, LOG_WORKERREQUEST_PREFIX, LOG_WORKERREQUEST_FLAGS),
	}
}
func (w *WorkerRequest) runwork() {
	for {
		select {
		case command := <-w.p.ChanCommand:
			if command == "exit" {
				return
			}
		default:
			if len(w.p.ChanStockRequest) > 0 {
				w.p.Lock()

			} else {

			}

		}
	}
}
func (w *WorkerRequest) get(u string) {

}
