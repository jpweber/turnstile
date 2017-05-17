/*
* @Author: jpweber
* @Date:   2017-04-27 10:47:35
* @Last Modified by:   jpweber
* @Last Modified time: 2017-04-27 10:59:48
 */

package main

import (
	"fmt"
	"os"
	"time"

	cli "gopkg.in/urfave/cli.v1"
)

func main() {

	app := cli.NewApp()
	app.Name = "turnstile"
	app.Version = "0.1.0"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "James Weber",
			Email: "jpweber@gmail.com",
		},
	}
	app.HelpName = "turnstile"
	app.Usage = "Manages vm orchestration for running short lived vm jobs"
	// app.UsageText = "turnstile - cloud vm orchestration"
	// app.ArgsUsage = "[args and such]"
	app.Commands = []cli.Command{
		cli.Command{
			Name: "test",
			Action: func(c *cli.Context) error {
				jobid := pseudoUUID()

				instanceInfo := Instance{}
				ec2Svc := newSession()
				instanceInfo.launchInstance(ec2Svc)
				db := openDB()
				fmt.Println(instanceInfo)

				//write to data store
				writeToDB(db, "instanceID", instanceInfo.InstanceID, jobid)

				// read from data store
				result := readFromDB(db, "instanceID", jobid)
				fmt.Printf("JobID: %s\nInstance ID: %s\n", jobid, result)

				// descirbe(ec2Svc, instanceInfo.InstanceID)
				instanceInfo.state(ec2Svc)

				return nil
			},
		},
		cli.Command{
			Name:    "serve",
			Aliases: []string{""},
			// Category:    "server",
			Usage: "start up in server mode",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		cli.Command{
			Name:    "state",
			Aliases: []string{""},
			// Category:    "insance",
			Usage:       "get state of an instance",
			UsageText:   "doo - does the dooing",
			Description: "Queries based on Job ID or Instance ID for the state of an instance",
			// ArgsUsage:   "[arrgh]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "instance, i",
					Usage: "Instance to get state for",
				},
				cli.StringFlag{
					Name:  "job, j",
					Usage: "Job ID to get state for",
				},
			},
			Action: func(c *cli.Context) error {
				id := c.String("instance")
				i := Instance{InstanceID: id}
				ec2Svc := newSession()
				i.state(ec2Svc)
				return nil
			},
		},
	}

	app.Run(os.Args)

}
