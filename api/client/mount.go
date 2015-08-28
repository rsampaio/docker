package client

import (
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	Cli "github.com/docker/docker/cli"
	"io/ioutil"
	"net/url"
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
		params := url.Values{}
		params.Add("filter", "{'status':'mounted'}")
		servResp, err := cli.call("GET", "/containers/json?"+params.Encode(), nil, nil)
		if err != nil {
			return err
		}

		defer servResp.body.Close()
		containers := []types.Container{}

		if err := json.NewDecoder(servResp.body).Decode(&containers); err != nil {
			return err
		}

		logrus.Infof("%+v\n", containers)
	} else {
		container := cmd.Arg(0)
		servResp, err := cli.call("POST", "/containers/"+container+"/mount", nil, nil)
		if err != nil {
			logrus.Fatal(err.Error())
		}

		body, err := ioutil.ReadAll(servResp.body)
		defer servResp.body.Close()
		if err != nil {
			logrus.Fatal(err.Error())
		}

		logrus.Infof("container=%s mount_opts=%s server_resp=%s\n", container, *flMntOpts, body)
	}
	return nil
}
