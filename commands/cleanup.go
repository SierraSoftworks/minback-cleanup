package commands

import (
	"fmt"
	"net/url"

	"github.com/SierraSoftworks/minback-cleanup/tools"
	log "github.com/Sirupsen/logrus"
	minio "github.com/minio/minio-go"
	"github.com/urfave/cli"
)

var Cleanup = cli.Command{
	Name: "cleanup",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "server",
			EnvVar: "MINIO_SERVER",
		},
		cli.StringFlag{
			Name:   "access-key",
			EnvVar: "MINIO_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   "secret-key",
			EnvVar: "MINIO_SECRET_KEY",
		},
		cli.StringFlag{
			Name:   "bucket",
			Value:  "backups",
			EnvVar: "MINIO_BUCKET",
		},
		cli.StringFlag{
			Name:  "db",
			Usage: "The name of the database backup files (my-db-2017-12-19.backup would use 'my-db')",
		},
		cli.StringSliceFlag{
			Name:  "keep, k",
			Usage: "@7d/1d will keep a backup every 1d for all backups 7d old or older",
		},
	},
	Action: func(c *cli.Context) error {
		logger := log.
			WithField("server", c.String("server")).
			WithField("bucket", c.String("bucket")).
			WithField("db", c.String("db"))

		specs := []tools.SampleSpec{}
		for _, spec := range c.StringSlice("keep") {
			s, err := tools.NewSampleSpec(spec)
			if err != nil {
				log.WithField("spec", spec).WithError(err).Error("could not parse keep spec")
				return err
			}

			specs = append(specs, s)
		}

		server, err := url.Parse(c.String("server"))
		if err != nil {
			logger.WithError(err).Error("failed to parse server URL")
			return err
		}

		client, err := minio.New(server.Host, c.String("access-key"), c.String("secret-key"), server.Scheme == "http")
		if err != nil {
			logger.WithError(err).Error("failed to create client")
			return err
		}

		if exists, err := client.BucketExists(c.String("bucket")); !exists {
			if err == nil {
				err = fmt.Errorf("bucket does not exist")
				logger.WithError(err).Error("bucket does not exist")
			} else {
				logger.WithError(err).Error("failed to determine if bucket exists")
			}

			return err
		}

		doneCh := make(chan struct{})
		defer close(doneCh)

		objs := client.ListObjectsV2(
			c.String("bucket"),
			fmt.Sprintf("%s-", c.String("db")),
			false,
			doneCh,
		)

		for obj := range objs {
			if obj.Err != nil {
				log.WithError(err).Error("failed to enumerate bucket objects")
				return err
			}

			logger := logger.WithFields(log.Fields{
				"object.key":          obj.Key,
				"object.lastModified": obj.LastModified,
			})

			logger.Debug("enumerated bucket object")

			ts, err := tools.ParseFilename(obj.Key, fmt.Sprintf("%s-", c.String("db")))
			if err != nil {
				logger.WithError(err).Warn("failed to parse filename, skipping")
				continue
			}

			if tools.Matches(ts, specs) {
				logger.Info("keeping backup file")
			} else {
				logger.Warn("removing backup file")
				if err := client.RemoveObject(c.String("bucket"), obj.Key); err != nil {
					logger.WithError(err).Error("failed to remove backup file")
					return err
				}
			}
		}

		return nil
	},
}
