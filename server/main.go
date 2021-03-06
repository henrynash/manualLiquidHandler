// server/main.go: Part of the Antha language
// Copyright (C) 2015 The Antha authors. All rights reserved.
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
//
// For more information relating to the software or licensing issues please
// contact license@antha-lang.org or write to the Antha team c/o
// Synthace Ltd. The London Bioscience Innovation Centre
// 2 Royal College St, London NW1 0NH UK

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/antha-lang/antha/driver/pb/lh"
	"github.com/antha-lang/manualLiquidHandler"
	"github.com/antha-lang/manualLiquidHandler/cli"
	"google.golang.org/grpc"
)

var (
	port int
	view string
)

func main() {
	flag.IntVar(&port, "port", 50051, "port to listen on")
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var x manualLiquidHandler.ManualExecuter
	x = cli.NewReadWriterExecutor(
		os.Stdin,
		os.Stdout,
	)

	x.Init()

	manual := manualLiquidHandler.NewManual(x)

	s := grpc.NewServer()
	lh.RegisterExtendedLiquidhandlingDriverServer(s, manual)
	go func() {
		fmt.Println("Listening at", lis.Addr().String())
		s.Serve(lis)
	}()
	x.Close()
}
