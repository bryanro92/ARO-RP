package main

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	utillog "github.com/Azure/ARO-RP/pkg/util/log"
	"github.com/sirupsen/logrus"
)

var (
	certList      []string
	expiryList    []string
	certDir       = "secrets"
	certExtension = ".pem"
)

func run(ctx context.Context, l *logrus.Entry) error {
	certList = walkCerts(ctx, l)
	expiryList = checkCertExpiry(ctx, l, certList)
	return nil

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

func checkCertExpiry(ctx context.Context, l *logrus.Entry, list []string) []string {
	for i := range list {
		r, err := ioutil.ReadFile(list[i])
		if err != nil {
			l.Fatal(err)
		}
		block, rest := pem.Decode(r)
		for block.Type == "PRIVATE KEY" {
			block, rest = pem.Decode(rest)
			break
		}
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			l.Fatal(err)
		}
		daysRemaining := cert.NotAfter.Sub(time.Now()).Hours() / 24
		if daysRemaining < 30 {
			expiryList = append(expiryList, cert.Subject.CommonName)
			l.Warnf("%s expires in %v days", cert.Subject.CommonName, daysRemaining)
		}
	}
	return expiryList
}

func main() {
	log := utillog.GetLogger()

	if err := run(context.Background(), log); err != nil {
		log.Fatal(err)
	}

}
