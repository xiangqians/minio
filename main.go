// Copyright (c) 2015-2021 MinIO, Inc.
//
// This file is part of MinIO Object Storage stack
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main // import "github.com/minio/minio"

//go:generate go install tool

import (
	"github.com/minio/minio/internal/ext/thumb"
	"log"
	"os"

	// MUST be first import.
	_ "github.com/minio/minio/internal/init"

	minio "github.com/minio/minio/cmd"
)

// Program arguments: server "D:\tmp\minio\data" --address ":9000" --console-address ":9001"

func main() {
	err := thumb.Startup()
	if err != nil {
		log.Fatalf("Error starting thumb: %v", err)
	}
	defer thumb.Shutdown()

	minio.Main(os.Args)
}
