package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-go/pkg/typast"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-go/pkg/typmock"
	"github.com/typical-go/typical-go/pkg/typrls"
	"github.com/typical-go/typical-rest-server/internal/generated/typical"
	"github.com/typical-go/typical-rest-server/pkg/mysqltool"
	"github.com/typical-go/typical-rest-server/pkg/pgtool"
	"github.com/typical-go/typical-rest-server/pkg/typcfg"
	"github.com/typical-go/typical-rest-server/pkg/typdocker"
)

var descriptor = typgo.Descriptor{
	ProjectName:    "typical-rest-server",
	ProjectVersion: "0.9.4",
	ProjectLayouts: []string{"internal", "pkg"},

	Cmds: []typgo.Cmd{
		// test
		&typgo.TestProject{},
		// compile
		&typgo.CompileProject{},
		// annotate
		&typast.AnnotateProject{
			Destination: "internal/generated/typical",
			Annotators: []typast.Annotator{
				&typapp.CtorAnnotation{},
				&typapp.DtorAnnotation{},
				&typcfg.EnvconfigAnnotation{
					DotEnv:   ".env",     // generate .env file
					UsageDoc: "USAGE.md", // generate USAGE.md
				},
			},
		},
		// run
		&typgo.RunProject{
			Before: typgo.BuildCmdRuns{"annotate", "compile"},
		},
		// mock
		&typmock.MockCmd{},
		// docker
		&typdocker.DockerCmd{},
		// pg
		&pgtool.PgTool{
			Name: "pg",
			ConfigFn: func() pgtool.Configurer {
				cfg, err := typical.LoadPostgresCfg()
				if err != nil {
					log.Fatal(err)
				}
				return cfg
			},
			DockerName:   "typical-rest-server_pg01_1",
			MigrationSrc: "file://databases/librarydb/migration",
			SeedSrc:      "databases/librarydb/seed",
		},
		// mysql
		&mysqltool.MySQLTool{
			Name: "mysql",
			ConfigFn: func() mysqltool.Configurer {
				cfg, err := typical.LoadMySQLCfg()
				if err != nil {
					log.Fatal(err)
				}
				return cfg
			},
			DockerName:   "typical-rest-server_mysql01_1",
			MigrationSrc: "file://databases/albumdb/migration",
			SeedSrc:      "databases/albumdb/seed",
		},
		// reset
		&typgo.Command{
			Name:  "reset",
			Usage: "reset the project locally (postgres/etc)",
			Action: typgo.BuildCmdRuns{
				"pg.drop", "pg.create", "pg.migrate", "pg.seed",
				"mysql.drop", "mysql.create", "mysql.migrate", "mysql.seed",
			},
		},
		// release
		&typrls.ReleaseProject{
			Before: typgo.BuildCmdRuns{"test", "compile"},
			// Releaser:  &typrls.CrossCompiler{Targets: []typrls.Target{"darwin/amd64", "linux/amd64"}},
			Publisher: &typrls.Github{Owner: "typical-go", Repo: "typical-rest-server"},
		},
	},
}

func main() {
	typgo.Start(&descriptor)
}
