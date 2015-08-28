package client

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	Cli "github.com/docker/docker/cli"
	"io/ioutil"
	"strings"
)

// CmdMount Mount layers from a created container.
//
// Usage:
// 	docker mount [OPTIONS] CONTAINER
func (cli *DockerCli) CmdMount(args ...string) error {
	cmd := Cli.Subcmd(
		"mount",
		[]string{"[CONTAINER]"},
		strings.Join([]string{
			"Mount layers from a container.\n",
		}, ""),
		true,
	)
	var (
		empty     = false
		flMntOpts = cmd.String([]string{"o", "-mount-opts"}, "", "Mount options comma separated")
		//flMntPoint = cmd.String([]string{"d", "-mount-point"}, "", "Mount point for merged layers")
	)
	//cmd.Require(flag.Exact, 1)
	cmd.ParseFlags(args, true)

	if cmd.Arg(0) == "" {
		empty = true
	}

	if empty {
		// List mounted containers
		fmt.Printf("deadbeef\t/var/lib/docker/storagedriver/deadbeef\tro,nosuid\n")
	} else {
		container := cmd.Arg(0)
		servResp, err := cli.call("POST", "/containers/"+container+"/mount", nil, nil)
		if err != nil {
			logrus.Fatal(err.Error())
		}

		body, err := ioutil.ReadAll(servResp.body)
		if err != nil {
			logrus.Fatal(err.Error())
		}

		logrus.Infof("container=%s mount_opts=%s server_resp=%s\n", container, *flMntOpts, body)
	}
	return nil
}
