package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/crshao/grpc-graphql-gateway/nsq/handler"

	"github.com/crshao/grpc-graphql-gateway/nsq/client"
	"github.com/nsqio/go-nsq"
	"github.com/segmentio/go-queue"
	"github.com/tj/docopt"
	"github.com/tj/go-gracefully"
	"gopkg.in/yaml.v2"
)

var Version = "0.0.1"

const Usage = `
  Usage:
    nsq_to_postgres --config file
    nsq_to_postgres -h | --help
    nsq_to_postgres --version

  Options:
    -c, --config file   configuration file path
    -h, --help          output help information
    -v, --version       output version

`

type Config struct {
	Postgres *client.Config
	Nsq      map[string]interface{}
}

type messageHandler struct{}
type Message struct {
	Name      string
	Timestamp string
}

// HandleMessage implements the Handler interface.
func (h *messageHandler) HandleMessage(m *nsq.Message) error {
	//Process the Message
	var request Message

	if err := json.Unmarshal(m.Body, &request); err != nil {
		log.Println("Error when Unmarshaling the message body, Err : ", err)
		// Returning a non-nil error will automatically send a REQ command to NSQ to re-queue the message.
		return err
	}
	//Print the Message
	log.Println("Message")
	log.Println("--------------------")
	log.Println("Name : ", request.Name)
	log.Println("Timestamp : ", request.Timestamp)
	log.Println("--------------------")
	log.Println("")

	// Returning a non-nil error will automatically send a REQ command to NSQ to re-queue the message.
	return nil
}

func main() {
	args, err := docopt.Parse(Usage, nil, true, Version, false)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	log.Printf("starting nsq_to_postgres version %s", Version)

	// Read config
	file := args["--config"].(string)
	b, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("error reading config: %s", err)
	}

	// Unmarshal config
	config := new(Config)
	err = yaml.Unmarshal(b, config)
	if err != nil {
		log.Fatalf("error unmarshalling config: %s", err)
	}

	// Validate config
	err = config.Postgres.Validate()
	if err != nil {
		log.Fatalf("configuration error: %s", err)
	}

	// Apply nsq config
	c := queue.NewConsumer("", "nsq_to_postgres")

	for k, v := range config.Nsq {
		c.Set(k, v)
	}

	// Connect
	log.Printf("connecting to postgres")
	db, err := client.New(config.Postgres)
	if err != nil {
		log.Fatalf("error connecting: %s", err)
	}

	// Bootstrap with table
	err = db.Bootstrap()

	if err != nil {
		log.Printf("error bootstrapping: %s", err)
	}

	// Start consumer
	log.Printf("starting consumer")

	c.Start(handler.New(db))

	gracefully.Shutdown()
	log.Printf("stopping consumer")
	c.Stop()

	log.Printf("bye :)")
}
