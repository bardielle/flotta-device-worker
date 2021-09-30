package deregistration

import (
	"git.sr.ht/~spc/go-log"
	"github.com/jakub-dzon/k4e-device-worker/internal/configuration"
	"github.com/jakub-dzon/k4e-operator/models"
	"github.com/jakub-dzon/k4e-device-worker/internal/workload"
)

type Deregistration struct {
	workloads    *workload.WorkloadManager
	config       *configuration.Manager
}

func NewDeregistration(workloadsManager *workload.WorkloadManager, configManager *configuration.Manager) (*Deregistration, error) {
	deregstration := Deregistration{
		workloads:                   workloadsManager,
		config:                      configManager,
	}
	return &deregstration, nil
}

func (d *Deregistration) Update(configuration models.DeviceConfigurationMessage) error {
	err := d.workloads.RemoveAllWorkloads()
	if err != nil {
		log.Errorf("failed to remove workloads: %v", err)
		return err
	}

	// after removing all the workloads the manifets directory can be removed
	err = d.workloads.DeleteManifestsDir()
	if err != nil {
		log.Errorf("failed to delete manifets directory: %v", err)
		return err
	}
	return nil
}
