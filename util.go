package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

func runCommand(command string, timeout time.Duration) error {
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	fmt.Println(command)
	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func isActive(s *compute.Service, project, region, instanceGroup, instanceTemplate string) (bool, error) {
	igms := compute.NewRegionInstanceGroupManagersService(s)

	ig, err := igms.Get(project, region, instanceGroup).Do()
	if err != nil {
		return false, err
	}

	found := false
	for _, v := range ig.Versions {
		if v.TargetSize.Calculated > 0 {
			if strings.Contains(strings.TrimSpace(v.Name), strings.TrimSpace(instanceTemplate)) {
				found = true
				break
			}
		}
	}

	return found, nil
}

func recreateInstance(s *compute.Service, project, region, instanceGroup, zone, instance string) (err error) {
	igms := compute.NewRegionInstanceGroupManagersService(s)

	r := &compute.RegionInstanceGroupManagersRecreateRequest{
		Instances: []string{fmt.Sprintf("zones/%v/instances/%v", zone, instance)},
	}

	retry := 0
	for {
		_, err := igms.RecreateInstances(project, region, instanceGroup, r).Do()
		if err != nil && isNotReadyErr(err) {
			time.Sleep(6 * time.Second)
			retry++
			if retry > 10 {
				return fmt.Errorf("too many retries")
			}
			continue

		} else if err != nil {
			return err
		} else if err == nil {
			return nil
		}
	}
}

func isReasonErr(err error, reason string) bool {
	if e, ok := err.(*googleapi.Error); ok {
		for _, x := range e.Errors {
			if x.Reason == reason {
				return true
			}
		}
	}
	return false
}

func isNotReadyErr(err error) bool {
	return isReasonErr(err, "resourceNotReady")
}
