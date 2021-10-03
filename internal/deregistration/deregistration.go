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

	err = d.workloads.DeleteManifestsDir()
	if err != nil {
		log.Errorf("failed to delete manifests directory: %v", err)
		return err
	}

	err = d.workloads.DeleteTable()
	if err != nil {
		log.Errorf("failed to delete table: %v", err)
		return err
	}

	err = d.workloads.DeleteVolumeDir()
	log.Infof("Deleting volumes directory")
	if err != nil {
		log.Errorf("failed to delete volumes directory: %v", err)
		return err
	}

	err = d.config.DeleteDeviceConfig()
	if err != nil {
		log.Errorf("failed to delete device-config file: %v", err)
		return err
	}
	return nil
}
