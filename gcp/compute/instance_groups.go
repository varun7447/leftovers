package compute

import (
	"fmt"

	gcpcompute "google.golang.org/api/compute/v1"
)

type instanceGroupsClient interface {
	ListInstanceGroups(zone string) (*gcpcompute.InstanceGroupList, error)
	DeleteInstanceGroup(zone, instanceGroup string) error
}

type InstanceGroups struct {
	client instanceGroupsClient
	logger logger
	zones  map[string]string
}

func NewInstanceGroups(client instanceGroupsClient, logger logger, zones map[string]string) InstanceGroups {
	return InstanceGroups{
		client: client,
		logger: logger,
		zones:  zones,
	}
}

func (s InstanceGroups) Delete() error {
	var groups []*gcpcompute.InstanceGroup
	for _, zone := range s.zones {
		l, err := s.client.ListInstanceGroups(zone)
		if err != nil {
			return fmt.Errorf("Listing instance groups for zone %s: %s", zone, err)
		}
		groups = append(groups, l.Items...)
	}

	for _, i := range groups {
		n := i.Name
		// n := s.clearerName(i)

		proceed := s.logger.Prompt(fmt.Sprintf("Are you sure you want to delete instance group %s?", n))
		if !proceed {
			continue
		}

		zoneName := s.zones[i.Zone]
		if err := s.client.DeleteInstanceGroup(zoneName, n); err != nil {
			s.logger.Printf("ERROR deleting instance group %s: %s\n", n, err)
		} else {
			s.logger.Printf("SUCCESS deleting instance group %s\n", n)
		}
	}

	return nil
}

// func (s Instances) clearerName(i *gcpcompute.Instance) string {
// 	extra := []string{}
// 	if i.Tags != nil && len(i.Tags.Items) > 0 {
// 		for _, tag := range i.Tags.Items {
// 			extra = append(extra, tag)
// 		}
// 	}

// 	if len(extra) > 0 {
// 		return fmt.Sprintf("%s (%s)", i.Name, strings.Join(extra, ", "))
// 	}

// 	return i.Name
// }
