package main

import (
	"flag"
	"fmt"
	"github.com/nsqio/go-nsq"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

var start = time.Now()
var ops uint64 = 0
var numbPtr = flag.Int("msg", 10000, "number of messages (default: 10000)")
var lookup = flag.String("lookup", "127.0.0.1", "IP address of nsqlookupd (default: 127.0.0.1")

func main() {

	/*
	   Below are the default port settings
	   nsqd listens at port 4150 (for TCP clients), 4151 (for HTTP clients)

	   nsqlookupd listens at port 4160 (for TCP clients), 4161 (for HTTP clients)

	   nsqadmin listens at port 4171 (for HTTP clients) or
	     to be specified (for go-nsq clients) q.ConnectToNSQLookupd("127.0.0.1:4161")

	   http://tleyden.github.io/blog/2014/11/12/an-example-of-using-nsq-from-go/
	   $ nsqlookupd &
	   $ nsqd --lookupd-tcp-address=127.0.0.1:4160 &
	   $ nsqadmin --lookupd-http-address=127.0.0.1:4161 &
	*/

	flag.Parse()

	wg := &sync.WaitGroup{}

	config := nsq.NewConfig()
	q, _ := nsq.NewConsumer("test", "ch", config)

	q.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		wg.Add(1)

		//log.Printf("Got a message: %v", string(message.Body))

		atomic.AddUint64(&ops, 1)
		if ops == uint64(*numbPtr) {
			elapsed := time.Since(start)
			log.Printf("Time took %s", elapsed)
		}

		wg.Done()
		return nil
	}))

	//err := q.ConnectToNSQD("127.0.0.1:4150") - hardcoded

	//err := q.ConnectToNSQLookupd("127.0.0.1:4161") //better

	err := q.ConnectToNSQLookupd(string(*lookup) + ":4161") //much better

	if err != nil {
		log.Panic("Could not connect")
	}
	wg.Wait()

	fmt.Scanln()
}
