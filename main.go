package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

func main() {
	config := &Config{}
	config, err := readConfig()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Project: %v", config.Project)
	log.Printf("Region: %v", config.Region)
	log.Printf("Zone: %v", config.Zone)
	log.Printf("Instance Group: %v", config.InstanceGroup)
	log.Printf("Instance Template: %v", config.InstanceTemplate)
	log.Printf("Instance: %v", config.Instance)
	log.Printf("Timeout: %v", config.Timeout)

	options := make([]option.ClientOption, 0)
	if config.CredentialsFile != "" {
		options = append(options, option.WithCredentialsFile(config.CredentialsFile))
	}

	c, err := compute.NewService(context.Background(), options...)
	if err != nil {
		log.Fatal(err)
	}

	// wait until instance template is no longer active
	log.Printf("Waiting until instance template '%v' is no longer part of instance group '%v' ...", config.InstanceTemplate, config.InstanceGroup)
	count := 0
	for {
		active, err := isActive(c, config.Project, config.Region, config.InstanceGroup, config.InstanceTemplate)
		if err != nil {
			log.Printf("ERROR: Polling for updates failed: %v (will try again)", err)
		}

		if !active {
			count += 1
		}

		if count >= 3 {
			break
		}

		time.Sleep(5 * time.Second)
	}

	// instance template is no longer active ...
	log.Println("Instance Template has been removed from Instance Group")
	log.Println("Waiting for -on-shutdown command to finish ...")
	if err := runCommand(config.OnShutdown, config.Timeout); err != nil {
		log.Printf("WARNING: -on-shutdown command failed: %v", err)
	}

	// delete instance ...
	log.Printf("Deleting instance '%v'", config.Instance)
	if err := recreateInstance(c, config.Project, config.Region, config.InstanceGroup, config.Zone, config.Instance); err != nil {
		log.Fatalf("FATAL: Deleting instance failed: %v", err)
	}

	log.Println("Waiting for shutdown ...")
}
