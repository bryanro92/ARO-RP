package main

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	utillog "github.com/Azure/ARO-RP/pkg/util/log"
	"github.com/sirupsen/logrus"
)

var (
	certDir       = "secrets"
	certExtension = ".pem"
)

func run(ctx context.Context, l *logrus.Entry) error {
	certList, err := walkCerts(ctx, l)
	if err != nil {
		l.Fatal(err)
	}
	return checkCertExpiry(ctx, l, certList)

}

func walkCerts(ctx context.Context, l *logrus.Entry) ([]string, error) {
	var f []string
	err := filepath.Walk(certDir, func(file string, i os.FileInfo, err error) error {
		if strings.EqualFold(filepath.Ext(file), certExtension) {
			l.Infof("found cert, %s", file)
			f = append(f, file)
		}
		return err
	})
	return f, err
}

func checkCertExpiry(ctx context.Context, l *logrus.Entry, list []string) error {
	var err error
	var r []byte
	for i := range list {
		r, err = ioutil.ReadFile(list[i])
		if err != nil {
			l.Fatal(err)
		}
		// loops over blocks in a pem file
		// if block is a private key, skip and get next block
		// if block != private key, its a cert we can check expiry of
		block, rest := pem.Decode(r)
		for block.Type == "PRIVATE KEY" {
			block, rest = pem.Decode(rest)
		}
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			l.Fatal(err)
		}
		daysRemaining := cert.NotAfter.Sub(time.Now()).Hours() / 24
		if daysRemaining < 30 {
			l.Warnf("%s expires in %v days", cert.Subject.CommonName, daysRemaining)
		}
	}
	return err
}

func main() {
	log := utillog.GetLogger()

	if err := run(context.Background(), log); err != nil {
		log.Fatal(err)
	}
}
