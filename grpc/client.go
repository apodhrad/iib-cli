package grpc

import (
	"context"
	"io"

	"github.com/operator-framework/operator-registry/pkg/api"
	"google.golang.org/grpc"
)

type Service struct {
	Name    string
	Methods []string
}

type Channel struct {
	Name    string `json:"name"`
	CsvName string `json:"csvName"`
}

type Package struct {
	Name               string    `json:"name"`
	Channels           []Channel `json:"channels"`
	DefaultChannelName string    `json:"defaultChannelName"`
}

type Bundle struct {
	CsvName     string `json:"csvName"`
	PackageName string `json:"packageName"`
	ChannelName string `json:"channelname"`
	BundlePath  string `json:"bundlePath"`
	Version     string `json:"version"`
	Replaces    string `json:"replaces,omitempty"`
}

type Client struct {
	Registry api.RegistryClient
	Conn     *grpc.ClientConn
}

func NewClient(address string) (*Client, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := &Client{
		Registry: api.NewRegistryClient(conn),
		Conn:     conn,
	}
	return c, err
}

func (c *Client) Close() error {
	if c.Conn == nil {
		return nil
	}
	return c.Conn.Close()
}

func (c *Client) GetPackageNames() ([]*api.PackageName, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := c.Registry.ListPackages(ctx, &api.ListPackageRequest{})
	if err != nil {
		// return []Package{}, err
		return []*api.PackageName{}, err
	}

	// packages := []Package{}
	packages := []*api.PackageName{}
	for next, err := stream.Recv(); err != io.EOF; next, err = stream.Recv() {
		if err != nil && err != io.EOF {
			// return []Package{}, err
			return []*api.PackageName{}, err
		}
		// nextPackage := Package{Name: next.Name}
		packages = append(packages, next)
	}

	return packages, nil
}

func (c *Client) GetPackage(name string) (*api.Package, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return c.Registry.GetPackage(ctx, &api.GetPackageRequest{Name: name})
}
