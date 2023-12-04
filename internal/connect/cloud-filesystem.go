package connect

import (
	"fmt"
	
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	
	"github.com/creamsensation/cp/internal/config"
	"github.com/creamsensation/cp/internal/style"
)

func CloudFilesystem(c config.Filesystem) *minio.Client {
	client, err := minio.New(
		c.Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(c.AccessKey, c.SecretKey, ""),
			Secure: true,
		},
	)
	if err != nil {
		fmt.Printf(
			"Filesystem [cloud:%s]: %s\n",
			style.GoldColor.Render(c.StorageName),
			style.RedColor.Render("ERROR"),
		)
		fmt.Printf("%s\n", style.RedColor.Render(err.Error()))
	}
	if err == nil {
		fmt.Printf(
			"Filesystem [cloud:%s]: %s\n",
			style.GoldColor.Render(c.StorageName),
			style.EmeraldColor.Render("CONNECTED"),
		)
	}
	return client
}
