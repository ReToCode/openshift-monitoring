package client

import (
	"net"
	"github.com/cenkalti/rpc2"
	"github.com/SchweizerischeBundesbahnen/openshift-monitoring/models"
	"os"
)

func StartDaemon(h string, dt string, ns string) *rpc2.Client {
	// Local state
	host, _ := os.Hostname()
	d := models.Daemon{Hostname: host,
		Namespace: ns,
		DaemonType: dt,
		StartedChecks: 0,
		FailedChecks:0,
		SuccessfulChecks:0}

	dc := &models.DaemonClient{Daemon: d,
		Quit: make(chan bool),
		ToHub: make(chan models.CheckResult)}

	// Register on hub
	conn, _ := net.Dial("tcp", h)
	dc.Client = rpc2.NewClient(conn)
	dc.Client.Handle("startChecks", func(client *rpc2.Client, checks *models.Checks, reply *string) error {
		startChecks(dc, checks)
		*reply = "ok"
		return nil
	})
	dc.Client.Handle("stopChecks", func(client *rpc2.Client, stop *bool, reply *string) error {
		stopChecks(dc)
		*reply = "ok"
		return nil
	})
	// Start handling from & to hub
	go dc.Client.Run()
	go handleCheckResultToHub(dc)

	registerOnHub(h, dc)

	return dc.Client
}

func StopDaemon(c *rpc2.Client) {
	unregisterOnHub(c)
}
