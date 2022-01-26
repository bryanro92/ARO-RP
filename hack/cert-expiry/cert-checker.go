package main

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	utillog "github.com/Azure/ARO-RP/pkg/util/log"
	"github.com/sirupsen/logrus"
)

var (
	certList      []string
	expiryList    []string
	certDir       = "secrets"
	certExtension = ".crt"
)

func run(ctx context.Context, l *logrus.Entry) error {
	certList = walkCerts(ctx, l)
	return checkCertExpiry(certList)

}

func walkCerts(ctx context.Context, l *logrus.Entry) []string {
	var f []string
	err := filepath.Walk(certDir, func(file string, i os.FileInfo, err error) error {
		if strings.EqualFold(filepath.Ext(file), certExtension) {
			l.Infof("found cert, %s", file)
			f = append(f, file)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return f
}

func checkCertExpiry(list []string) error {
	fmt.Println(certList)
}

func main() {
	log := utillog.GetLogger()

	if err := run(context.Background(), log); err != nil {
		log.Fatal(err)
	}

}
