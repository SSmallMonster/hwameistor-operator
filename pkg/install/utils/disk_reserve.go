package utils

import (
	"context"

	hwameistoroperatorv1alpha1 "github.com/hwameistor/hwameistor-operator/api/v1alpha1"
	hwameistorv1alpha1 "github.com/hwameistor/hwameistor/pkg/apis/hwameistor/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func CheckComponentsInstalledSuccessfully(clusterInstance *hwameistoroperatorv1alpha1.Cluster) bool {
	if ldmStatus := clusterInstance.Status.LocalDiskManager; ldmStatus == nil {
		return false
	} else if ldmStatus.Health != "Normal" {
		return false
	}
	if !clusterInstance.Spec.LocalStorage.Disable {
		if lsStatus := clusterInstance.Status.LocalStorage; lsStatus == nil {
			return false
		} else if lsStatus.Health != "Normal" {
			return false
		}
	}
	if !clusterInstance.Spec.AdmissionController.Disable {
		if clusterInstance.Status.AdmissionController.Health != "Normal" {
			return false
		}
	}
	if !clusterInstance.Spec.Scheduler.Disable {
		if clusterInstance.Status.Scheduler.Health != "Normal" {
			return false
		}
	}
	if !clusterInstance.Spec.Evictor.Disable {
		if clusterInstance.Status.Evictor.Health != "Normal" {
			return false
		}
	}
	if !clusterInstance.Spec.ApiServer.Disable {
		if clusterInstance.Status.ApiServer.Health != "Normal" {
			return false
		}
	}
	if !clusterInstance.Spec.Metrics.Disable {
		if clusterInstance.Status.Metrics.Health != "Normal" {
			return false
		}
	}

	return true
}

func ListLocalDisks (cli client.Client) ([]hwameistorv1alpha1.LocalDisk, error) {
	ldList := hwameistorv1alpha1.LocalDiskList{}
	if err := cli.List(context.TODO(), &ldList); err != nil {
		return nil, err
	}

	return ldList.Items, nil
}

func checkLocalDiskReserveCondition (ld *hwameistorv1alpha1.LocalDisk, diskReserveConfigurations []hwameistoroperatorv1alpha1.DiskReserveConfiguration) bool {
	for _, diskReserveConfiguration := range diskReserveConfigurations {
		if (ld.Spec.NodeName == diskReserveConfiguration.NodeName) && (ld.Spec.DiskAttributes.Type == diskReserveConfiguration.DiskType) {
			return true
		}
	}
	return false
}

func ReserveDisk (clusterInstance *hwameistoroperatorv1alpha1.Cluster, cli client.Client) error {
	localDisks, err := ListLocalDisks(cli)
	if err != nil {
		return err
	}
	
	for _, localDisk := range localDisks {
		if checkLocalDiskReserveCondition(&localDisk, clusterInstance.Spec.DiskReserveConfigurations) {
			localDisk.Spec.Reserved = true
			if err := cli.Update(context.TODO(), &localDisk); err != nil {
				return err
			}
		}
	}

	return nil
}

