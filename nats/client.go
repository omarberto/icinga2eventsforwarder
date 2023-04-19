//go:generate ../../../tools/readme_config_includer/generator
package nats

import (
	//"wp/tls"
	//_ "embed"
	"fmt"
    "sync"
	//"strings"

	"github.com/nats-io/nats.go"

	"wp/icinga2eventsforwarder/config"
)

type NATS struct {

	conf *config.NATS

	conn       *nats.Conn
}

var lock = &sync.Mutex{}

var NATSclient *NATS

func GetNATSclient() *NATS {
    if NATSclient == nil {
        lock.Lock()
        defer lock.Unlock()
        if NATSclient == nil {
            fmt.Println("Creating single instance now.")
            NATSclient = &NATS{}
        } else {
            fmt.Println("Single instance already created.")
        }
    }

    return NATSclient
}

func (n *NATS) Connect(c *config.NATS) error {
	n.conf = c

	var err error

	opts := []nats.Option{
		nats.MaxReconnects(-1),
	}

	tlsConfig, err := n.conf.ClientConfig.TLSConfig()
	if err != nil {
		return err
	}

	opts = append(opts, nats.Secure(tlsConfig))

	// try and connect
	n.conn, err = nats.Connect(n.conf.Server, opts...)

	return err
}

func (n *NATS) Close() error {
	n.conn.Close()
	return nil
}

func (n *NATS) Write(message []byte) error {

	err := n.conn.Publish(n.conf.Subject, message)
	if err != nil {
		return fmt.Errorf("FAILED to send NATS message: %w", err)
	}

	return nil
}