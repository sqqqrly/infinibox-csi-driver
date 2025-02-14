package storage

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"infinibox-csi-driver/api"
	"infinibox-csi-driver/common"
	"infinibox-csi-driver/helper"
	tests "infinibox-csi-driver/test_helper"
	"testing"

	"github.com/stretchr/testify/suite"
)

func (suite *StorageHelperSuite) SetupTest() {
	suite.api = new(api.MockApiService)
	suite.accessMock = new(helper.MockAccessModesHelper)
	suite.cs = &commonservice{api: suite.api, accessModesHelper: suite.accessMock}

	tests.ConfigureKlog()
}

type StorageHelperSuite struct {
	suite.Suite
	api        *api.MockApiService
	accessMock *helper.MockAccessModesHelper
	cs         *commonservice
}

func TestStorageHelperSuite(t *testing.T) {
	suite.Run(t, new(StorageHelperSuite))
}

////////////////////////// validateProtocolToNetworkSpace tests ///////////////////

// test NFS and TREEQ protocal validate with network space.
func (suite *StorageHelperSuite) Test_Network_Protocol_Match_NFS_TREEQ_Success() {
	//service := StorageHelperSuite{cs: *suite.cs}

	networkSpace := api.NetworkSpace{Service: common.NS_NFS_SVC}
	var scProtocol = common.PROTOCOL_NFS
	scNetSpace := []string{"someSpace", "someOtherSpace"}

	//networkSpaceErr := errors.New("Some error")
	suite.api.On("GetNetworkSpaceByName", mock.Anything).Return(networkSpace, nil)

	// validate NFS
	err := validateProtocolToNetworkSpace(scProtocol, scNetSpace, suite.cs.api)
	assert.Nil(suite.T(), err, "Expected Nil returned on success ")

	// validate TREEQ
	scProtocol = common.PROTOCOL_TREEQ
	err = validateProtocolToNetworkSpace(scProtocol, scNetSpace, suite.cs.api)
	assert.Nil(suite.T(), err, "Expected Nil returned on success ")

}

// validate iscsi protocol network space match
func (suite *StorageHelperSuite) Test_Network_Protocol_Match_ISCSI_Success() {
	//service := StorageHelperSuite{cs: *suite.cs}

	scNetSpace := []string{"someSpace", "someOtherSpace"}

	// validate ISCSI
	scProtocol := common.PROTOCOL_ISCSI
	iNetworkSpace := api.NetworkSpace{Service: common.NS_ISCSI_SVC}
	suite.api.On("GetNetworkSpaceByName", mock.Anything).Return(iNetworkSpace, nil)
	err := validateProtocolToNetworkSpace(scProtocol, scNetSpace, suite.cs.api)
	assert.Nil(suite.T(), err, "Expected Nil returned on success ")
}

// validate iscsi protocol with NFS service fails
func (suite *StorageHelperSuite) Test_Network_Protocol_MisMatch_ISCSI_Failure() {
	//service := StorageHelperSuite{cs: *suite.cs}

	scNetSpace := []string{"someSpace", "someOtherSpace"}

	// validate ISCSI
	scProtocol := common.PROTOCOL_ISCSI
	iNetworkSpace := api.NetworkSpace{Service: common.NS_NFS_SVC}
	suite.api.On("GetNetworkSpaceByName", mock.Anything).Return(iNetworkSpace, nil)
	err := validateProtocolToNetworkSpace(scProtocol, scNetSpace, suite.cs.api)
	assert.NotNil(suite.T(), err, "Expected iscsi to not match with NFS service ")
}

// validate NFS protocol with iSCSI service fails
func (suite *StorageHelperSuite) Test_Network_Protocol_MisMatch_NFS_Failure() {
	//service := StorageHelperSuite{cs: *suite.cs}

	scNetSpace := []string{"someSpace", "someOtherSpace"}

	// validate ISCSI
	scProtocol := common.PROTOCOL_NFS
	iNetworkSpace := api.NetworkSpace{Service: common.NS_ISCSI_SVC}
	suite.api.On("GetNetworkSpaceByName", mock.Anything).Return(iNetworkSpace, nil)
	err := validateProtocolToNetworkSpace(scProtocol, scNetSpace, suite.cs.api)
	assert.NotNil(suite.T(), err, "Expected iscsi to not match with NFS service ")
}

// validate FC protocol with iSCSI service fails
func (suite *StorageHelperSuite) Test_Network_Protocol_MisMatch_FC_ISCSI_Failure() {
	//service := StorageHelperSuite{cs: *suite.cs}

	scNetSpace := []string{"someSpace", "someOtherSpace"}

	scProtocol := common.PROTOCOL_FC
	iNetworkSpace := api.NetworkSpace{Service: common.NS_ISCSI_SVC}
	suite.api.On("GetNetworkSpaceByName", mock.Anything).Return(iNetworkSpace, nil)
	err := validateProtocolToNetworkSpace(scProtocol, scNetSpace, suite.cs.api)
	assert.NotNil(suite.T(), err, "Expected iscsi to not match with NFS service ")
}

// validate FC protocol with NFS service fails
func (suite *StorageHelperSuite) Test_Network_Protocol_MisMatch_FC_NFS_Failure() {
	//service := StorageHelperSuite{cs: *suite.cs}

	scNetSpace := []string{"someSpace", "someOtherSpace"}

	scProtocol := common.PROTOCOL_FC
	networkSpace := api.NetworkSpace{Service: common.NS_ISCSI_SVC}
	suite.api.On("GetNetworkSpaceByName", mock.Anything).Return(networkSpace, nil)
	err := validateProtocolToNetworkSpace(scProtocol, scNetSpace, suite.cs.api)
	assert.NotNil(suite.T(), err, "Expected iscsi to not match with NFS service ")
}

// validate iscsi protocol with iSCSI and NFS service fails
func (suite *StorageHelperSuite) Test_Network_Protocol_MisMatch_NAMESPACES_Failure() {
	//service := StorageHelperSuite{cs: *suite.cs}

	scNetSpace := []string{"someiscsiSpace", "someNfsSpace"}

	scProtocol := common.PROTOCOL_ISCSI
	iscsiNetworkSpace := api.NetworkSpace{Service: common.NS_ISCSI_SVC}
	nfsNetworkSpace := api.NetworkSpace{Service: common.NS_NFS_SVC}

	suite.api.On("GetNetworkSpaceByName", mock.Anything).Return(iscsiNetworkSpace, nil).Once()
	suite.api.On("GetNetworkSpaceByName", mock.Anything).Return(nfsNetworkSpace, nil).Once()
	err := validateProtocolToNetworkSpace(scProtocol, scNetSpace, suite.cs.api)
	assert.NotNil(suite.T(), err, "Expected iscsi to not match with NFS service ")
}

// validate NFS protocol with no network spaces fails
func (suite *StorageHelperSuite) Test_Network_Protocol_NFS_NO_NETWORKSPACES_Failure() {

	scNetSpace := []string{} // no network spaces
	scProtocol := common.PROTOCOL_FC
	networkSpace := api.NetworkSpace{Service: common.NS_ISCSI_SVC}

	suite.api.On("GetNetworkSpaceByName", mock.Anything).Return(networkSpace, nil)
	err := validateProtocolToNetworkSpace(scProtocol, scNetSpace, suite.cs.api)
	assert.NotNil(suite.T(), err, "Expected non-nil for empty network space list")
}
