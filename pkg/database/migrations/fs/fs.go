package fs

import "embed"

// Migrations contains the files to initialize the database.
//
//go:embed *.sql
var Migrations embed.FS

var MigrationsPath = "."
