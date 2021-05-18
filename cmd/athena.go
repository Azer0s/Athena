package main

import "github.com/alecthomas/kong"

var cli struct {
	Connect struct {
		ConnectionString string `arg name:"url" help:"URL of a cluster node of the targeted cluster." type:"url"`
	} `cmd help:"Connect to athena cluster."`

	Configure struct {
		ConnectionString string `arg name:"url" help:"URL of a cluster node of the targeted cluster." type:"url"`
	} `cmd help:"Configure an athena cluster."`

	ClusterInit struct {
		Port     int    `optional name:"port" help:"Port to start local instance at." type:"port"`
		Hostname string `optional name:"hostname" help:"Hostname to start local instance at." type:"hostname"`
	} `cmd help:"Init a new athena cluster."`
}

func main() {
	ctx := kong.Parse(&cli)
	switch ctx.Command() {
	case "cluster-init":
	default:
		panic(ctx.Command())
	}
}
