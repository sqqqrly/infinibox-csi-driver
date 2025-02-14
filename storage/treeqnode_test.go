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
	"infinibox-csi-driver/api"
	"infinibox-csi-driver/helper"
	tests "infinibox-csi-driver/test_helper"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func (suite *TreeqNodeSuite) SetupTest() {
	rand.Seed(time.Now().UnixNano())
	suite.nfsMountMock = new(MockNfsMounter)
	suite.osHelperMock = new(helper.MockOsHelper)
	suite.storageHelperMock = new(MockStorageHelper)
	suite.api = new(api.MockApiService)
	suite.cs = &commonservice{api: suite.api, accessModesHelper: suite.accessMock}
	tests.ConfigureKlog()
}

type TreeqNodeSuite struct {
	suite.Suite
	nfsMountMock      *MockNfsMounter
	osHelperMock      *helper.MockOsHelper
	accessMock        *helper.MockAccessModesHelper
	api               *api.MockApiService
	cs                *commonservice
	storageHelperMock *MockStorageHelper
}

func TestTreeqNodeSuite(t *testing.T) {
	suite.Run(t, new(TreeqNodeSuite))
}

func (suite *TreeqNodeSuite) Test_TreeqNodePublishVolume_IsNotExist_false() {
	nfs := nfsstorage{storageHelper: suite.storageHelperMock, cs: *suite.cs, mounter: suite.nfsMountMock, osHelper: suite.osHelperMock}
	service := treeqstorage{nfsstorage: nfs}
	randomDir := RandomString(10)
	targetPath := randomDir
	fmt.Printf("creating %s\n", "/tmp/"+targetPath)
	err := os.Mkdir("/tmp/"+targetPath, os.ModePerm)
	assert.Nil(suite.T(), err)
	defer func() {
		fmt.Printf("removing %s\n", "/tmp/"+targetPath)
		err := os.RemoveAll("/tmp/" + targetPath)
		assert.Nil(suite.T(), err)
	}()

	suite.api.On("GetFileSystemByID", mock.Anything).Return(nil, nil)
	suite.storageHelperMock.On("SetVolumePermissions", mock.Anything).Return(nil)
	suite.storageHelperMock.On("GetNFSMountOptions", mock.Anything).Return([]string{}, nil)
	suite.api.On("ExportFileSystem", mock.Anything).Return(getExportResponseValue(), nil)
	exportResp := getExportResponse()
	suite.api.On("GetExportByFileSystem", mock.Anything).Return(exportResp, nil)
	suite.api.On("DeleteExportPath", mock.Anything).Return(exportResp, nil)
	suite.nfsMountMock.On("Mount", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	contex := getPublishContexMap()
	contex["csiContainerHostMountPoint"] = "/tmp/"

	req := getNodePublishVolumeRequest(targetPath, contex)
	req.VolumeId = "94148131#20000$$nfs_treeq"
	req.Secrets = make(map[string]string)
	req.Secrets["one"] = "one"
	req.Secrets["two"] = "two"
	req.Secrets["three"] = "three"
	responce, err := service.NodePublishVolume(context.Background(), req)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), responce, "empty object")
}

func (suite *TreeqNodeSuite) Test_TreeqNodePublishVolume_mount_sucess() {
	contex := getPublishContexMap()
	contex["csiContainerHostMountPoint"] = "/tmp/"
	randomDir := RandomString(10)
	targetPath := randomDir
	err := os.Mkdir("/tmp/"+targetPath, os.ModePerm)
	assert.Nil(suite.T(), err)
	defer func() {
		err := os.RemoveAll("/tmp/" + targetPath)
		assert.Nil(suite.T(), err)
	}()
	nfs := nfsstorage{storageHelper: suite.storageHelperMock, cs: *suite.cs, mounter: suite.nfsMountMock, osHelper: suite.osHelperMock}
	service := treeqstorage{nfsstorage: nfs}
	suite.storageHelperMock.On("SetVolumePermissions", mock.Anything).Return(nil)
	suite.storageHelperMock.On("GetNFSMountOptions", mock.Anything).Return([]string{}, nil)
	suite.nfsMountMock.On("Mount", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	suite.api.On("GetFileSystemByID", mock.Anything).Return(nil, nil)
	suite.api.On("ExportFileSystem", mock.Anything).Return(getExportResponseValue(), nil)
	exportResp := getExportResponse()
	suite.api.On("GetExportByFileSystem", mock.Anything).Return(exportResp, nil)
	suite.api.On("DeleteExportPath", mock.Anything).Return(exportResp, nil)
	suite.api.On("GetFileSystemByID", mock.Anything).Return(nil, nil)

	req := getNodePublishVolumeRequest(targetPath, contex)
	req.VolumeId = "94148131#20000$$nfs_treeq"
	req.Secrets = make(map[string]string)
	req.Secrets["one"] = "one"
	req.Secrets["two"] = "two"
	req.Secrets["three"] = "three"
	_, err = service.NodePublishVolume(context.Background(), req)
	assert.Nil(suite.T(), err, "empty error")
}

func (suite *TreeqNodeSuite) Test_TreeqNodePublishVolume_mount_Error() {
	contex := getPublishContexMap()
	contex["csiContainerHostMountPoint"] = "/tmp/"
	randomDir := RandomString(10)
	targetPath := randomDir
	mountErr := errors.New("mount error")
	nfs := nfsstorage{mounter: suite.nfsMountMock, storageHelper: suite.storageHelperMock, osHelper: suite.osHelperMock}
	service := treeqstorage{nfsstorage: nfs}
	suite.storageHelperMock.On("SetVolumePermissions", mock.Anything).Return(nil)
	suite.storageHelperMock.On("GetNFSMountOptions", mock.Anything).Return([]string{}, nil)
	suite.nfsMountMock.On("Mount", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mountErr)
	_, err := service.NodePublishVolume(context.Background(), getNodePublishVolumeRequest(targetPath, contex))
	assert.NotNil(suite.T(), err, "not nil error")
}

func (suite *TreeqNodeSuite) Test_TreeqNodeUnpublishVolume_NotMountPoint_error() {
	mountErr := errors.New("mount error")
	nfs := nfsstorage{mounter: suite.nfsMountMock, osHelper: suite.osHelperMock}
	service := treeqstorage{nfsstorage: nfs}
	suite.nfsMountMock.On("IsLikelyNotMountPoint", mock.Anything).Return(true, mountErr)
	suite.osHelperMock.On("IsNotExist", mountErr).Return(false)
	suite.nfsMountMock.On("Mount", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	suite.nfsMountMock.On("Unmount", mock.Anything).Return(nil)
	targetPath := "/var/lib/kublet/"
	suite.osHelperMock.On("Remove", targetPath).Return(nil)
	volumeID := "1234"
	_, err := service.NodeUnpublishVolume(context.Background(), getNodeUnPublishVolumeRequest(targetPath, volumeID))
	assert.NotNil(suite.T(), err, "not nil error")
}

func (suite *TreeqNodeSuite) Test_TreeqNodeUnpublishVolume_NotMountPoint_IsNotExist_true() {
	mountErr := errors.New("mount error")
	nfs := nfsstorage{mounter: suite.nfsMountMock, osHelper: suite.osHelperMock}
	service := treeqstorage{nfsstorage: nfs}
	suite.nfsMountMock.On("IsLikelyNotMountPoint", mock.Anything).Return(true, nil)
	suite.osHelperMock.On("IsNotExist", mountErr).Return(true)
	suite.nfsMountMock.On("Mount", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	suite.nfsMountMock.On("Unmount", mock.Anything).Return(nil)
	targetPath := "/var/lib/kublet/"
	suite.osHelperMock.On("Remove", targetPath).Return(nil)
	volumeID := "1234"
	_, err := service.NodeUnpublishVolume(context.Background(), getNodeUnPublishVolumeRequest(targetPath, volumeID))
	assert.Nil(suite.T(), err, "empty error")
}

func (suite *TreeqNodeSuite) Test_TreeqNodeUnpublishVolume_NotMountPoint_IsNotExist_false() {
	mountErr := errors.New("not exists")
	nfs := nfsstorage{mounter: suite.nfsMountMock, osHelper: suite.osHelperMock}
	service := treeqstorage{nfsstorage: nfs}
	suite.nfsMountMock.On("IsLikelyNotMountPoint", mock.Anything).Return(true, mountErr)
	suite.osHelperMock.On("IsNotExist", mountErr).Return(false)
	suite.nfsMountMock.On("Mount", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	suite.nfsMountMock.On("Unmount", mock.Anything).Return(nil)
	targetPath := "/var/lib/kublet/"
	suite.osHelperMock.On("Remove", targetPath).Return(nil)
	volumeID := "1234"
	_, err := service.NodeUnpublishVolume(context.Background(), getNodeUnPublishVolumeRequest(targetPath, volumeID))
	assert.NotNil(suite.T(), err, "not nil error")
}

func (suite *TreeqNodeSuite) Test_TreeqNodeUnpublishVolume_notMnt_true() {
	targetPath := "/var/lib/kublet/"
	volumeID := "1234"
	nfs := nfsstorage{mounter: suite.nfsMountMock, osHelper: suite.osHelperMock}
	service := treeqstorage{nfsstorage: nfs}
	suite.nfsMountMock.On("IsLikelyNotMountPoint", mock.Anything).Return(true, nil)
	suite.nfsMountMock.On("IsNotMountPoint", mock.Anything).Return(true, nil)
	suite.osHelperMock.On("Remove", targetPath).Return(nil)
	suite.nfsMountMock.On("Unmount", mock.Anything).Return(nil)

	_, err := service.NodeUnpublishVolume(context.Background(), getNodeUnPublishVolumeRequest(targetPath, volumeID))
	assert.Nil(suite.T(), err, "empty err")
}

func (suite *TreeqNodeSuite) Test_TreeqNodeUnpublishVolume_unmount_fail() {
	mountErr := errors.New("mount error")
	targetPath := "/var/lib/kublet/"
	volumeID := "1234"
	nfs := nfsstorage{mounter: suite.nfsMountMock, osHelper: suite.osHelperMock}
	service := treeqstorage{nfsstorage: nfs}
	suite.nfsMountMock.On("IsLikelyNotMountPoint", mock.Anything).Return(true, nil)
	suite.nfsMountMock.On("IsNotMountPoint", mock.Anything).Return(true, nil)
	suite.nfsMountMock.On("Mount", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	suite.nfsMountMock.On("Unmount", targetPath).Return(mountErr)
	suite.nfsMountMock.On("Mount", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	_, err := service.NodeUnpublishVolume(context.Background(), getNodeUnPublishVolumeRequest(targetPath, volumeID))

	assert.NotNil(suite.T(), err, "not nil error")
}

func (suite *TreeqNodeSuite) Test_TreeqNodeUnpublishVolume_unmount_sucess() {
	targetPath := "/var/lib/kublet/"
	volumeID := "1234"
	nfs := nfsstorage{mounter: suite.nfsMountMock, osHelper: suite.osHelperMock}
	service := treeqstorage{nfsstorage: nfs}
	suite.nfsMountMock.On("IsLikelyNotMountPoint", mock.Anything).Return(true, nil)
	suite.nfsMountMock.On("IsNotMountPoint", mock.Anything).Return(true, nil)
	suite.nfsMountMock.On("Unmount", targetPath).Return(nil)
	suite.osHelperMock.On("Remove", targetPath).Return(nil)
	_, err := service.NodeUnpublishVolume(context.Background(), getNodeUnPublishVolumeRequest(targetPath, volumeID))
	assert.Nil(suite.T(), err, "empty err")
}

func (suite *TreeqNodeSuite) Test_NodeStageVolume() {
	nfs := nfsstorage{mounter: suite.nfsMountMock, osHelper: suite.osHelperMock}
	service := treeqstorage{nfsstorage: nfs}

	_, err := service.NodeStageVolume(context.Background(), &csi.NodeStageVolumeRequest{})
	assert.Nil(suite.T(), err, "empty err")
}

func (suite *TreeqNodeSuite) Test_NodeUnstageVolume() {
	nfs := nfsstorage{mounter: suite.nfsMountMock, osHelper: suite.osHelperMock}
	service := treeqstorage{nfsstorage: nfs}

	_, err := service.NodeUnstageVolume(context.Background(), &csi.NodeUnstageVolumeRequest{})
	assert.Nil(suite.T(), err, "empty err")
}

func RandomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
