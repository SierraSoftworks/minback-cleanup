# minback-cleanup
**Minio Rolling Backup Management Container**

This container provides a clean, job based, approach to rotating backups
on a schedule according to given rules. It is intended to be used in
conjunction with the various `minback` containers which place backup
files in a Minio bucket for you.

## Features
* Allows multiple backup retention periods to be selected
* Lightweight and short lived
* Simple and well tested Go implementation

## Example
In this example, let's assume we have a very aggressive backup schedule that
creates a new backup every minute. We would like to keep:

 - Every backup for the last `1h`
 - A backup every `15m` after `1h`
 - A backup every `1h` after `6h`
 - A backup every week (`7d`) after 1 week
 - A backup every 4 weeks (`4w`) after a month (`4w`)
 - A backup every year after 3 years (`156w`)

```sh
docker run --rm --env-file backup.env minback/cleanup cleanup --db my_db --keep "~1h/15m" --keep "~6h/1h" --keep "~7d/1w" --keep "~4w/4w" --keep "~52w/156w"
```

#### `backup.env`
```
MINIO_SERVER=https://play.minio.io
MINIO_ACCESS_KEY=minio
MINIO_SECRET_KEY=miniosecret
MINIO_BUCKET=backups
```

## Usage
```
NAME:
   main.exe - A new cli application

USAGE:
   main.exe [global options] command [command options] [arguments...]

VERSION:
   v1.0.0

AUTHOR:
   Benjamin Pannell <admin@sierrasoftworks.com>

COMMANDS:
     cleanup
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --log-level value  DEBUG|INFO|WARN|ERROR (default: "INFO")
   --help, -h         show help
   --version, -v      print the version

COPYRIGHT:
   Sierra Softworks © 2017
```

### `cleanup`
```
NAME:
   main.exe cleanup -

USAGE:
   main.exe cleanup [command options] [arguments...]

OPTIONS:
   --server value           [$MINIO_SERVER]
   --access-key value       [$MINIO_ACCESS_KEY]
   --secret-key value       [$MINIO_SECRET_KEY]
   --bucket value          (default: "backups") [$MINIO_BUCKET]
   --db value              The name of the database backup files (my-db-2017-12-19.backup would use 'my-db')
   --keep value, -k value  ~7d/1d will keep a backup every 1d for all backups 7d old or older
```

## Configuration
You can configure command line options using environment variables if you wish.

#### `MINIO_SERVER=https://play.minio.io`
The Minio server you wish to send backups to.

#### `MINIO_ACCESS_KEY=minio`
The Access Key used to connect to your Minio server.

#### `MINIO_SECRET_KEY=miniosecret`
The Secret Key used to connect to your Minio server.

#### `MINIO_BUCKET=backups`
The Minio bucket you wish to store your backup in.

[Kubernetes CronJob]: https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/
[Minio]: https://minio.io/
