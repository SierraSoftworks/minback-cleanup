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
docker run --rm --env-file backup.env minback/cleanup cleanup --db my_db --keep "@1h/15m~5m" --keep "@6h/1h~15m" --keep "@7d/1w~1d" --keep "@4w/4w~1w" --keep "@52w/156w~1w"
```

#### `backup.env`
```
MINIO_SERVER=https://play.minio.io
MINIO_ACCESS_KEY=minio
MINIO_SECRET_KEY=miniosecret
MINIO_BUCKET=backups
```

## Keep Specifications
A keep specification is composed of three different properties:
1. **Window** (`@`) tells us when this window starts to apply from. It defaults to "right now" and can be
   extended out with a time duration, for example `52w` tells us that the window starts at 1 year and
   continues into the past from there.

   Windows "stack", so `@52w` and `@4w` will combine to have one specification applying from `4w` to `52w`
   and the next from `52w` onwards.

2. **Spacing** (`/`) tells us which backups within a window to keep. A spacing of `15m` will keep any
   backups which fall on the `15m` boundary (according to their timestamp) - so `:00`, `:15`, `:30` and `:45`.

3. **Smudge Factor** (`~`) tells us that we should round the backup's time to the nearest provided duration.
   This is great for situations where your backup's timestamp is a lot more accurate than your window and you
   wish to smudge things to ensure it falls on a boundary.

   For example, a smudge of `~15m` will convert `11:13:18` to `11:15:00` and have it match the `/15m` spacing
   we selected previously. You can also use larger smudges like `~1d` or even `~52w` (for a year), but keep in
   mind that we don't de-duplicate things if multiple backups match because of the same smudge factor.

You define a keep specification by concatenating each property's control character with its duration and listing
them in any order you wish.

#### Example: `@2w/1d~6h`
In this example, let's assume we create a new backup every 6 hours but only want to store one backup per day
for everything older than 2 weeks (i.e. your 6 hour frequency is for short-term recovery).

1. **Window** is 2 weeks (`2w`), meaning this specification will only apply to backups older than 2 weeks,
   and until an older window is found.
2. **Spacing** is 1 day (`1d`), meaning we will keep one backup per day within this window. Backups will be selected
   based on their timestamp and only those which match (after smudging) the day boundary will be kept.
3. **Smudge Factor** is 6 hours (`6h`) which ensures that backups which occur up to 3 hours either side of the
   day boundary (midnight) will be kept.

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
   --keep value, -k value  @7d/1d will keep a backup every 1d for all backups 7d old or older
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
