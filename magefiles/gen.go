//go:build mage

package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/coder/coder/coderd/database/migrations"
	"github.com/coder/coder/coderd/database/postgres"
	"github.com/magefile/mage/mg"
)

type Gen mg.Namespace

func (Gen) DumpSQL() error {
	dst := filepath.Join(
		"coderd", "database", "dump.sql",
	)

	if destNewer(dst, sourceFilter{"coderd/database", []string{
		`migrations`,
		`\.sql$`,
	}}) {
		return nil
	}

	connection, closeFn, err := postgres.Open()
	if err != nil {
		return err
	}
	defer closeFn()

	db, err := sql.Open("postgres", connection)
	if err != nil {
		return err
	}

	err = migrations.Up(db)
	if err != nil {
		return err
	}
	const minimumPostgreSQLVersion = 13
	hasPGDump := false
	if _, err = exec.LookPath("pg_dump"); err == nil {
		out, err := exec.Command("pg_dump", "--version").Output()
		if err == nil {
			// Parse output:
			// pg_dump (PostgreSQL) 14.5 (Ubuntu 14.5-0ubuntu0.22.04.1)
			parts := strings.Split(string(out), " ")
			if len(parts) > 2 {
				version, err := strconv.Atoi(strings.Split(parts[2], ".")[0])
				if err == nil && version >= minimumPostgreSQLVersion {
					hasPGDump = true
				}
			}
		}
	}

	cmdArgs := []string{
		"pg_dump",
		"--schema-only",
		connection,
		"--no-privileges",
		"--no-owner",

		// We never want to manually generate
		// queries executing against this table.
		"--exclude-table=schema_migrations",
	}

	if !hasPGDump {
		cmdArgs = append([]string{
			"docker",
			"run",
			"--rm",
			"--network=host",
			fmt.Sprintf("postgres:%d", minimumPostgreSQLVersion),
		}, cmdArgs...)
	}
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...) //#nosec
	cmd.Env = append(os.Environ(), []string{
		"PGTZ=UTC",
		"PGCLIENTENCODING=UTF8",
	}...)
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	for _, sed := range []string{
		// Remove all comments.
		"/^--/d",
		// Public is implicit in the schema.
		"s/ public\\./ /g",
		"s/::public\\./::/g",
		"s/'public\\./'/g",
		// Remove database settings.
		"s/SET .* = .*;//g",
		// Remove select statements. These aren't useful
		// to a reader of the dump.
		"s/SELECT.*;//g",
		// Removes multiple newlines.
		"/^$/N;/^\\n$/D",
	} {
		cmd := exec.Command("sed", "-e", sed)
		cmd.Stdin = bytes.NewReader(output.Bytes())
		output = bytes.Buffer{}
		cmd.Stdout = &output
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	dump := fmt.Sprintf("-- Code generated by 'mage gen:dumpSQL'. DO NOT EDIT.\n%s", output.Bytes())
	err = os.WriteFile(
		dst,
		[]byte(dump), 0o600,
	)
	if err != nil {
		return err
	}
	return nil
}

func (Gen) GoQuerier() error {
	mg.Deps((Gen).DumpSQL)

	dest := filepath.Join("coderd", "database", "querier.go")
	if destNewer(dest, sourceFilter{"coderd/database", []string{
		`dump.sql$`,
		`sqlc.yaml$`,
		`\.sql$`,
	}}) {
		return nil
	}

	return shell("./coderd/database/generate.sh").run()
}

func (Gen) ProvisionerProto() error {
	var (
		dest   = filepath.Join("provisionersdk", "proto", "provisioner.pb.go")
		source = filepath.Join("provisionersdk", "proto", "provisioner.proto")
	)
	if destNewer(
		dest, sourceFilter{source, nil},
	) {
		return nil
	}

	return shell(`
		protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-drpc_out=. \
		--go-drpc_opt=paths=source_relative \
		%s \
	`, source).run()
}

func (Gen) ProvisionerdProto() error {
	var (
		dest   = filepath.Join("provisionerd", "proto", "provisionerd.pb.go")
		source = filepath.Join("provisionerd", "proto", "provisionerd.proto")
	)
	if destNewer(
		dest, sourceFilter{source, nil},
	) {
		return nil
	}

	return shell(`
		protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-drpc_out=. \
		--go-drpc_opt=paths=source_relative \
		%s \
	`, source).run()
}

func (Gen) TypesGenerated() error {
	dest := filepath.Join("site/src/api/typesGenerated.ts")
	if destNewer(
		dest, sourceFilter{"codersdk", []string{`\.go$`}},
	) {
		return nil
	}

	destFi, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer destFi.Close()

	c := goRun("./scripts/apitypings/main.go")
	c.Stdout = destFi
	c.Stderr = os.Stderr
	err = c.run()
	if err != nil {
		return err
	}

	return shell("yarn format:types").cd("site").run()
}

func (Gen) All() {
	mg.Deps(
		(Gen).GoQuerier, (Gen).DumpSQL,
		(Gen).ProvisionerProto, (Gen).ProvisionerdProto,
		(Gen).TypesGenerated,
	)
}
