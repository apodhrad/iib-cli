package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var EXPECTED_SERVICES []Service = []Service{
	{Name: "api.Registry", Methods: []string{"api.Registry.GetBundle", "api.Registry.GetBundleForChannel", "api.Registry.GetBundleThatReplaces", "api.Registry.GetChannelEntriesThatProvide", "api.Registry.GetChannelEntriesThatReplace", "api.Registry.GetDefaultBundleThatProvides", "api.Registry.GetLatestChannelEntriesThatProvide", "api.Registry.GetPackage", "api.Registry.ListBundles", "api.Registry.ListPackages"}},
	{Name: "grpc.health.v1.Health", Methods: []string{"grpc.health.v1.Health.Check"}},
	{Name: "grpc.reflection.v1alpha.ServerReflection", Methods: []string{"grpc.reflection.v1alpha.ServerReflection.ServerReflectionInfo"}},
}

const EXPECTED_SERVICES_TEXT_OUTPUT string = `SERVICE                                   METHOD                                                         
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

func TestApiCmdGetServices(t *testing.T) {
	setTestIIB(t)
	defer stopTestGrpc(t)

	services, err := apiCmdGetServices()
	assert.Nil(t, err)
	assert.Equal(t, EXPECTED_SERVICES, services)
}

func TestApiCmdToText(t *testing.T) {
	out, err := apiCmdToText(EXPECTED_SERVICES)
	assert.Nil(t, err)
	assert.Equal(t, EXPECTED_SERVICES_TEXT_OUTPUT, out)
}
