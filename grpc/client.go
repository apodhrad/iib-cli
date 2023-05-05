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
		return []*api.PackageName{}, err
	}

	packages := []*api.PackageName{}
	for next, err := stream.Recv(); err != io.EOF; next, err = stream.Recv() {
		if err != nil && err != io.EOF {
			return []*api.PackageName{}, err
		}
		packages = append(packages, next)
	}

	return packages, nil
}

func (c *Client) GetPackage(name string) (*api.Package, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return c.Registry.GetPackage(ctx, &api.GetPackageRequest{Name: name})
}

func (c *Client) GetBundles() ([]*api.Bundle, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := c.Registry.ListBundles(ctx, &api.ListBundlesRequest{})
	if err != nil {
		return []*api.Bundle{}, err
	}

	bundles := []*api.Bundle{}
	for next, err := stream.Recv(); err != io.EOF; next, err = stream.Recv() {
		if err != nil && err != io.EOF {
			return []*api.Bundle{}, err
		}
		bundles = append(bundles, next)
	}

	return bundles, nil
}

func (c *Client) GetBundle(pkg string, channel string, csv string) (*api.Bundle, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return c.Registry.GetBundle(ctx, &api.GetBundleRequest{
		PkgName:     pkg,
		ChannelName: channel,
		CsvName:     csv,
	})
}
