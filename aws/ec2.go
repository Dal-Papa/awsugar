package aws

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/tj/go-progress"
)

// Deletable provides an interface for any EC2 resource that can be deleted
type Deletable interface {
	Type() string
	Name() string
	Delete(*session.Session) error
}

// Sweetener provides an interface to do preventive cleaning before deleting
type Sweetener interface {
	Sweeten(*session.Session) error
}

// EC2Instance is a proxy for the AWS framework struct
type EC2Instance struct {
	*ec2.Instance
}

var _ = Deletable(&EC2Instance{})

// ListInstances returns the list of EC2Instance for the specific ids provided
func ListInstances(s *session.Session, ids []*string) ([]EC2Instance, error) {
	ec2C := ec2.New(s)
	res, err := ec2C.DescribeInstances(&ec2.DescribeInstancesInput{
		InstanceIds: ids,
	})
	if err != nil {
		return nil, fmt.Errorf("Couldn't list instances: %s", err)
	}
	var list []EC2Instance
	for i := range res.Reservations {
		for _, is := range res.Reservations[i].Instances {
			list = append(list, EC2Instance{is})
		}
	}
	return list, nil
}

// Type returns the EC2 type
func (e EC2Instance) Type() string { return "EC2" }

// Name returns the EC2 Instance name
func (e EC2Instance) Name() string {
	for i := range e.Tags {
		if e.Tags[i].Key == aws.String("Name") {
			return *e.Tags[i].Value
		}
	}
	return *e.InstanceId
}

// Delete the LoadBalancer
func (e EC2Instance) Delete(s *session.Session) error {
	ec2C := ec2.New(s)
	if _, err := ec2C.TerminateInstances(&ec2.TerminateInstancesInput{
		InstanceIds: []*string{e.InstanceId},
	}); err != nil {
		return fmt.Errorf("Couldn't delete EC2 instance [%s]: %s", e.Name(), err)
	}
	return nil
}

// LoadBalancer is a proxy for the AWS framework struct
type LoadBalancer struct {
	*elb.LoadBalancerDescription
}

var _ = Deletable(&LoadBalancer{})

// ListInactiveLoadBalancers returns a list of LoadBalancer that have no
// EC2Instance attached to it.
func ListInactiveLoadBalancers(s *session.Session) ([]LoadBalancer, error) {
	elbC := elb.New(s)
	res, err := elbC.DescribeLoadBalancers(&elb.DescribeLoadBalancersInput{})
	if err != nil {
		return nil, fmt.Errorf("Couldn't list load balancers: %s", err)
	}
	list := make([]LoadBalancer, 0, len(res.LoadBalancerDescriptions))
	for _, lb := range res.LoadBalancerDescriptions {
		if len(lb.Instances) == 0 {
			list = append(list, LoadBalancer{lb})
		}
	}
	return list, nil
}

// Type returns the ELB type
func (lb LoadBalancer) Type() string { return "ELB" }

// Name returns the LoadBalancer name
func (lb LoadBalancer) Name() string { return *lb.LoadBalancerName }

// Delete the LoadBalancer
func (lb LoadBalancer) Delete(s *session.Session) error {
	elbC := elb.New(s)
	if _, err := elbC.DeleteLoadBalancer(&elb.DeleteLoadBalancerInput{
		LoadBalancerName: lb.LoadBalancerName,
	}); err != nil {
		return fmt.Errorf("Couldn't delete load balancer [%s]: %s", *lb.LoadBalancerName, err)
	}
	return nil
}

// NetworkInterface is a proxy for the AWS framework struct
type NetworkInterface struct {
	*ec2.NetworkInterface
}

var _ = Deletable(&NetworkInterface{})

// ListUnattachedNetworkInterfaces returns a list of NetworkInterface
// that are currently not attached to an EC2Instance
func ListUnattachedNetworkInterfaces(s *session.Session) ([]NetworkInterface, error) {
	ec2C := ec2.New(s)
	res, err := ec2C.DescribeNetworkInterfaces(&ec2.DescribeNetworkInterfacesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("status"),
				Values: []*string{aws.String("available")},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("Couldn't list network interfaces: %s", err)
	}
	list := make([]NetworkInterface, 0, len(res.NetworkInterfaces))
	for _, ni := range res.NetworkInterfaces {
		list = append(list, NetworkInterface{ni})
	}
	return list, nil
}

// Type returns the Network Interface type
func (ni NetworkInterface) Type() string { return "Network Interface" }

// Name returns the NetworkInterface ID
func (ni NetworkInterface) Name() string { return *ni.NetworkInterfaceId }

// Delete the NetworkInterface
func (ni NetworkInterface) Delete(s *session.Session) error {
	ec2C := ec2.New(s)
	if _, err := ec2C.DeleteNetworkInterface(&ec2.DeleteNetworkInterfaceInput{
		NetworkInterfaceId: ni.NetworkInterfaceId,
	}); err != nil {
		return fmt.Errorf("Couldn't delete network interface [%s]: %s", *ni.NetworkInterfaceId, err)
	}
	return nil
}

// EBSVolume is a proxy for the AWS framework struct
type EBSVolume struct {
	*ec2.Volume
}

var _ = Deletable(&EBSVolume{})
var _ = Sweetener(&EBSVolume{})

// ListAvailableEBS returns a list of Available EBSVolume
func ListAvailableEBS(s *session.Session) ([]EBSVolume, error) {
	ec2C := ec2.New(s)
	res, err := ec2C.DescribeVolumes(&ec2.DescribeVolumesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("status"),
				Values: []*string{aws.String("available")},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("Couldn't list EBS volumes: %s", err)
	}
	list := make([]EBSVolume, 0, len(res.Volumes))
	for _, ni := range res.Volumes {
		list = append(list, EBSVolume{ni})
	}
	return list, nil
}

// Type returns the EBS type
func (v EBSVolume) Type() string { return "EBS" }

// Name returns the Volume ID
func (v EBSVolume) Name() string { return *v.VolumeId }

// Delete the EBSVolume
func (v EBSVolume) Delete(s *session.Session) error {
	ec2C := ec2.New(s)
	if _, err := ec2C.DeleteVolume(&ec2.DeleteVolumeInput{
		VolumeId: v.VolumeId,
	}); err != nil {
		return fmt.Errorf("Couldn't delete EBS volume [%s]: %s", *v.VolumeId, err)
	}
	return nil
}

// Sweeten creates a snapshot for the volume and waits for it to finish
// before the deletion of the EBSVolume
func (v EBSVolume) Sweeten(s *session.Session) error {
	name := v.VolumeId
	tags := make([]*ec2.Tag, 0, len(v.Tags))
	for i := range v.Tags {
		if v.Tags[i].Key == aws.String("Name") {
			name = aws.String(*name + "_" + *v.Tags[i].Value)
		}
		tags = append(tags, v.Tags[i])
	}
	ec2C := ec2.New(s)
	res, err := ec2C.CreateSnapshot(&ec2.CreateSnapshotInput{
		Description: name,
		VolumeId:    v.VolumeId,
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String(ec2.ResourceTypeSnapshot),
				Tags:         tags,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("Couldn't snapshot EBS volume [%s]: %s", *v.VolumeId, err)
	}
	snap := &Snapshot{res}
	return snap.Wait(s)
}

// Snapshot is a proxy for the AWS framework struct
type Snapshot struct {
	*ec2.Snapshot
}

// Wait for the Snapshot to finish before doing anything else
func (snap *Snapshot) Wait(s *session.Session) error {
	ec2C := ec2.New(s)
	fmt.Printf("Starting to monitor snapshot [%s]. This can take a few minutes...\n",
		*snap.SnapshotId)
	bar := progress.NewInt(100)
	bar.Text(*snap.SnapshotId)
	ticker := time.NewTicker(1 * time.Minute)
	var lastPercent string
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			res, err := ec2C.DescribeSnapshots(&ec2.DescribeSnapshotsInput{
				SnapshotIds: []*string{snap.SnapshotId},
			})
			if err != nil {
				return fmt.Errorf("Couldn't wait for snapshot [%s]: %s",
					*snap.SnapshotId, err)
			}
			for _, ws := range res.Snapshots {
				if *ws.State == ec2.SnapshotStateCompleted {
					fmt.Printf("\nSnapshot completed\n")
					return nil
				}
				if lastPercent != *ws.Progress {
					lastPercent = *ws.Progress
					percentInt, _ := strconv.Atoi(lastPercent[:len(lastPercent)-1])
					bar.ValueInt(percentInt)
					bar.WriteTo(os.Stdout)
				}
			}
		}
	}
}
