package connect

import (
	"fmt"
	
	"quirk"
	
	"github.com/creamsensation/cp/env"
	"github.com/creamsensation/cp/internal/config"
	"github.com/creamsensation/cp/internal/style"
)

func Database(item config.Database) *quirk.DB {
	conn, err := quirk.Connect(
		quirk.WithLog(env.Development()),
		quirk.WithDriver(item.Driver),
		quirk.WithHost(item.Host),
		quirk.WithPort(item.Port),
		quirk.WithUser(item.User),
		quirk.WithPassword(item.Password),
		quirk.WithDbname(item.Dbname),
		quirk.WithSsl(item.Ssl),
		quirk.WithCertPath(item.CertPath),
	)
	if err != nil {
		fmt.Printf(
			"Database [%s:%s]: %s\n",
			style.BlueColor.Render(item.Driver),
			style.GoldColor.Render(item.Dbname),
			style.RedColor.Render("ERROR"),
		)
		fmt.Printf("%s\n", style.RedColor.Render(err.Error()))
	}
	if err := conn.DB.Ping(); err == nil {
		fmt.Printf(
			"Database [%s:%s]: %s\n",
			style.BlueColor.Render(item.Driver),
			style.GoldColor.Render(item.Dbname),
			style.EmeraldColor.Render("CONNECTED"),
		)
	}
	return conn
}
