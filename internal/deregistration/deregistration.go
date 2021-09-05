package deregistration

import (
	"fmt"
	"git.sr.ht/~spc/go-log"
	pb "github.com/redhatinsights/yggdrasil/protocol"
	"github.com/jakub-dzon/k4e-operator/models"
	"github.com/jakub-dzon/k4e-device-worker/internal/workload"
	"io/ioutil"
	"os"
	"path"
	"time"
)

type Deregistration struct {
	workerClient pb.WorkerClient
	manifestsDir string
}

func NewDeregistration(workerClient pb.WorkerClient, configDir string) (*Deregistration, error) {
	manifestsDir := path.Join(configDir, "manifests")
	if err := os.MkdirAll(manifestsDir, 0755); err != nil {
		return nil, fmt.Errorf("cannot create directory: %w", err)
	}
	manager := Deregistration{
		manifestsDir: manifestsDir,
		workerClient:    workerClient,
	}

	return &manager, nil
}

func (d *Deregistration) Update(configuration models.DeviceConfigurationMessage) error {
	workloads := configuration.Workloads
	if len(workloads) == 0 {
		log.Trace("No workloads")

		// Purge all the workloads
		err := purgeWorkloads(workloads)
		if err != nil {
			return err
		}
		// Remove manifests
		err = removeManifests(d.manifestsDir)
		if err != nil {
			return err
		}
		return nil
	}
	// TODO: remove configuration files ?

	return nil
}

// TODO: refactoring
func purgeWorkloads(workloads *wrapper.workloadWrapper) error{
	podList, err := workloads.List()
	if err != nil {
		log.Errorf("Cannot list workloads: %v", err)
		return err
	}
	for _, podReport := range podList {
		err := Remove(podReport.Name)
		if err != nil {
			log.Errorf("Error removing workload: %v", err)
			return err
		}
	}
	return nil
}

// TODO: refactoring
func removeManifests(manifestsDir string) error {
	manifestInfo, err := ioutil.ReadDir(manifestsDir)
	if err != nil {
		return err
	}
	for _, fi := range manifestInfo {
		filePath := path.Join(manifestsDir, fi.Name())
		err := os.Remove(filePath)
		if err != nil {
			return err
		}
	}
	return nil
}

// TODO: refactoring
func Remove(workloadName string, ww *wrapper.workloadWrapper) error {
	id := ww.mappingRepository.GetId(workloadName)
	if id == "" {
		id = workloadName
	}
	if err := ww.workloads.Remove(id); err != nil {
		return err
	}
	if err := ww.netfilter.DeleteTable(nfTableName, workloadName); err != nil {
		log.Errorf("failed to delete chain %[1]s from %s table for workload %[1]s: %v", workloadName, nfTableName, err)
	}
	if err := ww.mappingRepository.Remove(workloadName); err != nil {
		return err
	}
	return nil
}

