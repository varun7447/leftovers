package compute

import (
	"fmt"
	"strings"

	gcpcompute "google.golang.org/api/compute/v1"
)

type instancesClient interface {
	ListInstances(zone string) (*gcpcompute.InstanceList, error)
	DeleteInstance(zone, instance string) error
}

type Instances struct {
	client instancesClient
	logger logger
	zones  map[string]string
}

func NewInstances(client instancesClient, logger logger, zones map[string]string) Instances {
	return Instances{
		client: client,
		logger: logger,
		zones:  zones,
	}
}

func (s Instances) Delete() error {
	var instances []*gcpcompute.Instance
	for _, zone := range s.zones {
		l, err := s.client.ListInstances(zone)
		if err != nil {
			return fmt.Errorf("Listing instances for zone %s: %s", zone, err)
		}
		instances = append(instances, l.Items...)
	}

	for _, i := range instances {
		n := s.clearerName(i)

		proceed := s.logger.Prompt(fmt.Sprintf("Are you sure you want to delete instance %s?", n))
		if !proceed {
			continue
		}

		zoneName := s.zones[i.Zone]
		if err := s.client.DeleteInstance(zoneName, i.Name); err != nil {
			s.logger.Printf("ERROR deleting instance %s: %s\n", i.Name, err)
		} else {
			s.logger.Printf("SUCCESS deleting instance %s\n", i.Name)
		}
	}

	return nil
}

func (s Instances) clearerName(i *gcpcompute.Instance) string {
	extra := []string{}
	if i.Tags != nil && len(i.Tags.Items) > 0 {
		for _, tag := range i.Tags.Items {
			extra = append(extra, tag)
		}
	}

	if len(extra) > 0 {
		return fmt.Sprintf("%s (%s)", i.Name, strings.Join(extra, ", "))
	}

	return i.Name
}
