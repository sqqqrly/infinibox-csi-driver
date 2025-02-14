/*
Copyright 2022 Infinidat
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package storage

import (
	"context"
	"errors"
	"fmt"
	"infinibox-csi-driver/common"
	"strconv"
	"strings"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog"
)

func (treeq *treeqstorage) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (csiResp *csi.CreateVolumeResponse, err error) {
	klog.V(4).Infof("CreateVolume called pvName %s parameters %v", req.GetName(), req.GetParameters())

	capacity, err := nfsSanityCheck(req, map[string]string{
		common.SC_POOL_NAME:     `\A.*\z`, // TODO: could make this enforce IBOX pool_name requirements, but probably not necessary
		common.SC_NETWORK_SPACE: `\A.*\z`, // TODO: could make this enforce IBOX network_space requirements, but probably not necessary
	}, map[string]string{
		common.SC_MAX_FILESYSTEMS:           `\A\d+\z`,
		common.SC_MAX_TREEQS_PER_FILESYSTEM: `\A\d+\z`,
		common.SC_MAX_FILESYSTEM_SIZE:       `\A.*\z`, // TODO: add more specific pattern
	}, treeq.nfsstorage.cs.api)
	if err != nil {
		return nil, err
	}

	storageClassParameters := req.GetParameters()
	treeq.nfsstorage.storageClassParameters = storageClassParameters
	treeqVolumeContext, err := treeq.treeqService.IsTreeqAlreadyExist(storageClassParameters[common.SC_POOL_NAME], strings.Trim(storageClassParameters[common.SC_NETWORK_SPACE], ""), req.GetName())
	if err != nil {
		klog.Errorf("error locating existing treeq %s", err.Error())
		return nil, err
	}
	if len(treeqVolumeContext) == 0 {
		treeqVolumeContext, err = treeq.treeqService.CreateTreeqVolume(storageClassParameters, capacity, req.GetName())
		if err != nil {
			klog.Errorf("error creating treeq volume %s", err.Error())
			return nil, err
		}
	}

	treeqVolumeContext[common.SC_NFS_EXPORT_PERMISSIONS] = storageClassParameters[common.SC_NFS_EXPORT_PERMISSIONS]
	treeqVolumeContext[common.SC_STORAGE_PROTOCOL] = storageClassParameters[common.SC_STORAGE_PROTOCOL]

	volumeID := treeqVolumeContext["ID"] + "#" + treeqVolumeContext["TREEQID"]
	klog.V(4).Infof("CreateVolume final treeqVolumeMap %v volumeID %s", treeqVolumeContext, volumeID)
	return &csi.CreateVolumeResponse{
		Volume: &csi.Volume{
			VolumeId:      volumeID,
			CapacityBytes: capacity,
			VolumeContext: treeqVolumeContext,
			ContentSource: req.GetVolumeContentSource(),
		},
	}, nil
}

func getVolumeIDs(volumeID string) (filesystemID, treeqID int64, err error) {
	volproto := strings.Split(volumeID, "#")
	if len(volproto) != 2 {
		err = fmt.Errorf("volume Id %s and other details not found", volumeID)
		return 0, 0, err
	}
	if filesystemID, err = strconv.ParseInt(volproto[0], 10, 64); err != nil {
		return 0, 0, err
	}

	// volumeID example := "94148131#20000$$nfs_treeq"
	treeqdetails := strings.Split(volproto[1], "$")

	if treeqID, err = strconv.ParseInt(treeqdetails[0], 10, 64); err != nil {
		return 0, 0, err
	}

	return filesystemID, treeqID, nil
}

func (treeq *treeqstorage) DeleteVolume(ctx context.Context, req *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	klog.V(2).Infof("DeleteVolume called on volume ID %s", req.GetVolumeId())
	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume ID missing in request")
	}

	filesystemID, treeqID, err := getVolumeIDs(req.GetVolumeId())
	if err != nil {
		klog.Errorf("Invalid Volume ID %v", err)
		return nil, status.Error(codes.InvalidArgument, "Invalid volume ID")
	}
	nfsDeleteErr := treeq.treeqService.DeleteTreeqVolume(filesystemID, treeqID)
	if nfsDeleteErr != nil {
		if strings.Contains(nfsDeleteErr.Error(), "FILESYSTEM_NOT_FOUND") {
			klog.Error("treeq already delete from infinibox")
			return &csi.DeleteVolumeResponse{}, nil
		}
		return nil, nfsDeleteErr
	}
	klog.V(2).Infof("treeq ID %s successfully deleted", req.GetVolumeId())
	return &csi.DeleteVolumeResponse{}, nil
}

func (treeq *treeqstorage) ControllerPublishVolume(ctx context.Context, req *csi.ControllerPublishVolumeRequest) (*csi.ControllerPublishVolumeResponse, error) {
	return &csi.ControllerPublishVolumeResponse{}, nil
}

func (treeq *treeqstorage) ControllerUnpublishVolume(ctx context.Context, req *csi.ControllerUnpublishVolumeRequest) (*csi.ControllerUnpublishVolumeResponse, error) {
	return &csi.ControllerUnpublishVolumeResponse{}, nil
}

func (treeq *treeqstorage) CreateSnapshot(ctx context.Context, req *csi.CreateSnapshotRequest) (*csi.CreateSnapshotResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Unsupported operation for treeq")
}

func (treeq *treeqstorage) DeleteSnapshot(ctx context.Context, req *csi.DeleteSnapshotRequest) (*csi.DeleteSnapshotResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Unsupported operation for treeq")
}

func (treeq *treeqstorage) ControllerExpandVolume(ctx context.Context, req *csi.ControllerExpandVolumeRequest) (expandVolume *csi.ControllerExpandVolumeResponse, err error) {
	defer func() {
		if res := recover(); res != nil && err == nil {
			err = errors.New("Recovered from CSI ControllerExpandVolume " + fmt.Sprint(res))
		}
	}()

	maxFileSystemSize := treeq.nfsstorage.storageClassParameters[common.SC_MAX_FILESYSTEM_SIZE]
	filesystemID, treeqID, err := getVolumeIDs(req.GetVolumeId())
	if err != nil {
		klog.Errorf("Invalid Volume ID %v", err)
		return nil, status.Error(codes.InvalidArgument, "Invalid volume ID")
	}

	capacity := int64(req.GetCapacityRange().GetRequiredBytes())
	if capacity < gib {
		capacity = gib
		klog.Warning("Volume Minimum capacity should be greater 1 GB")
	}

	klog.V(4).Infof("filesystemID %d treeqID %d capacity %d maxSize %s\n", filesystemID, treeqID, capacity, maxFileSystemSize)
	err = treeq.treeqService.UpdateTreeqVolume(filesystemID, treeqID, capacity, maxFileSystemSize)
	if err != nil {
		return
	}
	return &csi.ControllerExpandVolumeResponse{
		CapacityBytes:         capacity,
		NodeExpansionRequired: false,
	}, nil
}
