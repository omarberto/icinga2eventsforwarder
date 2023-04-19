package httpclient

import (
	"time"
	"fmt"
    "bufio"
	"net/http"
    "crypto/tls"

	"wp/icinga2eventsforwarder/config"
	"wp/icinga2eventsforwarder/tornado"
	"wp/icinga2eventsforwarder/nats"
)


type HTTPstream struct {

	conf *config.EVENTSSTREAM

	//conn       *nats.Conn
}

func (s *HTTPstream) Init(c *config.EVENTSSTREAM) {
	s.conf = c
}

func (s *HTTPstream) Run(hostname string) {
	
    fmt.Println("UrlRequest: ", s.conf.UrlRequest)

	for {

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: s.conf.InsecureSkipVerify},
		}
		client := &http.Client{Transport: tr}

		req, err := http.NewRequest("POST", s.conf.UrlRequest, nil)
		if err != nil {
			// handle err
		}
		req.SetBasicAuth(s.conf.Username, s.conf.Password)
		req.Header.Set("Accept", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			// handle err
			fmt.Println("Error at HTTP request: ", err.Error())
			time.Sleep(1 * time.Second)
			continue
		}


		// Read the response header
		fmt.Println("Response: Server:", resp.Header.Get("Server"))
		fmt.Println("Response: Status:", resp.Status)
		fmt.Println("Response: Proto:", resp.Proto)
		fmt.Println("Response: Content-Type:", resp.Header.Get("Content-Type"))


		// Read the response body

		scanner := bufio.NewScanner(resp.Body)

		for scanner.Scan() {
			message := tornado.NewMessage()
			//tornadoMessage := message.Append(hostname, scanner.Text())
			tornadoMessage, err := message.Parse(hostname, scanner.Text())

			if err == nil {
				//fmt.Println("Writing tornado message: ", string(tornadoMessage))
				var natsClient = nats.GetNATSclient()
				natsClient.Write(tornadoMessage)
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading HTTP response: ", err.Error())
			time.Sleep(1 * time.Second)
			continue
		}   

	}
}