package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/compute/metadata"
)

type Config struct {
	CredentialsFile  string
	Project          string
	Region           string
	Zone             string
	InstanceGroup    string
	InstanceTemplate string
	Instance         string
	OnShutdown       string
	Timeout          time.Duration
}

func readConfig() (*Config, error) {
	c := &Config{}

	// try to get config from metadata server
	c.Project, _ = metadata.ProjectID()
	c.InstanceGroup, _ = metadata.InstanceAttributeValue("created-by")
	c.InstanceTemplate, _ = metadata.InstanceAttributeValue("instance-template")
	c.Zone, _ = metadata.Zone()
	c.Instance, _ = metadata.InstanceName()

	// convert zone "us-central1-b" into region "us-central1"
	if len(c.Zone) > 2 {
		c.Region = c.Zone[:len(c.Zone)-2]
	}

	if c.InstanceTemplate != "" {
		// parse `projects/xxx/global/instanceTemplates/xxx`
		c.InstanceTemplate = lastPathElement(c.InstanceTemplate)
	}

	if c.InstanceGroup != "" {
		// parse `projects/xxx/regions/us-central1/instanceGroupManagers/xxx`
		c.InstanceGroup = lastPathElement(c.InstanceGroup)
	}

	// parse flags
	credentialsFileFlag := flag.String("creds", "", "")
	projectFlag := flag.String("project", "", "")
	zoneFlag := flag.String("zone", "", "")
	regionFlag := flag.String("region", "", "")
	instanceGroupFlag := flag.String("instance-group", "", "")
	instanceTemplateFlag := flag.String("instance-template", "", "")
	instanceFlag := flag.String("instance", "", "")
	onShutdownFlag := flag.String("on-shutdown", "", "")
	timeoutFlag := flag.String("timeout", "", "")
	flag.Parse()

	c.CredentialsFile = *credentialsFileFlag

	if isFlagSet("project") {
		c.Project = *projectFlag
	}

	if isFlagSet("region") {
		c.Region = *regionFlag
	}

	if isFlagSet("zone") {
		c.Zone = *zoneFlag
	}

	if isFlagSet("instance-group") {
		c.InstanceGroup = *instanceGroupFlag
	}

	if isFlagSet("instance-template") {
		c.InstanceTemplate = *instanceTemplateFlag
	}

	if isFlagSet("instance") {
		c.Instance = *instanceFlag
	}

	c.OnShutdown = *onShutdownFlag

	if isFlagSet("timeout") {
		d, err := time.ParseDuration(*timeoutFlag)
		if err != nil {
			return nil, fmt.Errorf("-timeout: %w", err)
		}
		c.Timeout = d
	} else {
		c.Timeout = 10 * time.Minute
	}

	c.Project = strings.TrimSpace(c.Project)
	c.Region = strings.TrimSpace(c.Region)
	c.InstanceGroup = strings.TrimSpace(c.InstanceGroup)
	c.InstanceTemplate = strings.TrimSpace(c.InstanceTemplate)

	if c.Project == "" {
		return nil, fmt.Errorf("missing -project")
	}

	if c.Region == "" {
		return nil, fmt.Errorf("missing -region")
	}

	if c.Zone == "" {
		return nil, fmt.Errorf("missing -zone")
	}

	if c.InstanceGroup == "" {
		return nil, fmt.Errorf("missing -instance-group")
	}

	if c.InstanceTemplate == "" {
		return nil, fmt.Errorf("missing -instance-template")
	}

	if c.Instance == "" {
		return nil, fmt.Errorf("missing -instance")
	}

	if c.OnShutdown == "" {
		return nil, fmt.Errorf("missing -on-shutdown")
	}

	return c, nil
}

func isFlagSet(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func lastPathElement(str string) string {
	x := strings.Split(str, "/")
	if len(x) == 1 {
		return str
	}

	return x[len(x)-1]
}
