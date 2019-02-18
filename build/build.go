package build

import (
	"context"

	"github.com/QMSTR/go-qmstr/service"
)

// QmstrBuildClient is used to send build information to the QMSTR master server
type QmstrBuildClient struct {
	buildService service.BuildServiceClient
}

// NewQmstrBuildClient returns a new instance of QmstrBuildClient to communicate with the build service
func NewQmstrBuildClient(srv service.BuildServiceClient) *QmstrBuildClient {
	return &QmstrBuildClient{buildService: srv}
}

// BuildFiles sends the given FileNode instances to the QMSTR master server
func (qbc *QmstrBuildClient) BuildFiles(fileNodes ...*service.FileNode) error {
	build, err := qbc.buildService.Build(context.Background())
	if err != nil {
		return err
	}

	for _, fn := range fileNodes {
		build.Send(fn)
	}
	return nil
}

// PackageFiles adds the given FileNode instances to the current package
func (qbc *QmstrBuildClient) PackageFiles(fileNodes ...*service.FileNode) error {
	pkg, err := qbc.buildService.Package(context.Background())
	if err != nil {
		return err
	}

	for _, fn := range fileNodes {
		pkg.Send(fn)
	}
	return nil
}
