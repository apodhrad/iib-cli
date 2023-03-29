package cmd

import (
	"os"
	"testing"

	"github.com/apodhrad/iib-cli/utils"
	"github.com/stretchr/testify/assert"
)

const EXPECTED_API_LIST_TEXT string = `SERVICE                                   METHOD                                                         
api.Registry                              api.Registry.GetBundle                                         
api.Registry                              api.Registry.GetBundleForChannel                               
api.Registry                              api.Registry.GetBundleThatReplaces                             
api.Registry                              api.Registry.GetChannelEntriesThatProvide                      
api.Registry                              api.Registry.GetChannelEntriesThatReplace                      
api.Registry                              api.Registry.GetDefaultBundleThatProvides                      
api.Registry                              api.Registry.GetLatestChannelEntriesThatProvide                
api.Registry                              api.Registry.GetPackage                                        
api.Registry                              api.Registry.ListBundles                                       
api.Registry                              api.Registry.ListPackages                                      
grpc.health.v1.Health                     grpc.health.v1.Health.Check                                    
grpc.reflection.v1alpha.ServerReflection  grpc.reflection.v1alpha.ServerReflection.ServerReflectionInfo  
`

func TestListApi(t *testing.T) {
	utils.GrpcStartSafely()

	os.Setenv("IIB", "quay.io/apodhrad/iib-test:v0.0.1")
	utils.GrpcStartSafely()

	table, _, err := listApi()
	assert.Nil(t, err)
	assert.Equal(t, EXPECTED_API_LIST_TEXT, table)
}
