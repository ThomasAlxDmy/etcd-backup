
Repository moved to [FanTV repository](https://github.com/fanhattan/etcd-backup)

<br/><br/><br/>


# etcd-backup

etcd-backup is a simple, efficient and lightweight Command line utility to backup and restore [etcd](https://github.com/coreos/etcd) keys.

## Dependencies

etcd-backup has only one dependency: go-etcd [the golang offical library for ETCD](https://github.com/coreos/go-etcd)

## Installation

  Installation composed of 3 steps:

* Install go [offical documentation](http://golang.org/doc/install/source)
* Download the project `git clone git@github.com:ThomasAlxDmy/etcd-backup.git`
* Download the dependency `go get github.com/coreos/go-etcd/etcd`
* Build the binary `cd etcd-backup` and then  `go install`

## Dumping

### Usage

    $ etcd-dump dump

This is the easiest way to dump the whole `etcd` keys. Results will be stored in a json file `etcd-dump.json`
file in the directory where you executed the command.

The default Backup strategy for dumping is dump all keys and conserve the order : `keys:["/"], recursive:true, sorted:true`
The backup strategy can be overwritten in the etcd-backup configuration file. See fixtures/backup-configuration.json

### Command line options and default values

  `-config` Mandatory etcd-backup configuration file location, default value: "backup-configuration.json". See [Configuration section](#config) for more information.<br/>
  `-concurent-request` Number of concurent resquest that will be executed during the restore(restore mode only), default value is 10.<br/>
  `-retries` Number of retires that will be executed if the requests fail, default value is 5.<br/>
  `-etcd-config` Mandatory etcd configuration file location, default value: "etcd-configuration.json". See fixtures folder for an example.<br/>
  `-dump` Location of the dump file data will be store in (in case of a dump) or load from (in case of a restore), default value: "etcd-dump.json".<br/>


    $ etcd-dump -config=myBackupConfig.json -retries=2 -etcd-config=myClusterConfig.json -dump=result.json dump

### <a name="config"/>Configuration

The `dump.keys` supports differents configurations:

  {
    "key": "/",
    "recursive": true
  }

Recursively dump all the keys inside `/`.

  {
    "key": "/myKey"
  }

Dump only the key `/myKey`.


### Dump File structure

Dumped keys are stored in an array of keys, the key path is the absolute path. By design non-empty directories are not saved in the dump file, and empty directories does not contain the `value` key:

    [{ "key": "/myKey", "value": "value1" },{ "key": "/dir/myKey/object", "value": "test" }, {"key": "/dir/mydir"}]

## Restoring

### Usage

    $ etcd-dump restore

Restore the keys from the `etcd-dump.json` file.

### Command line options and default values

  `-config` Mandatory etcd-backup configuration file location, default value: "backup-configuration.json". See [Configuration section](#config) for more information.<br/>
  `-concurent-request` Number of concurent resquest that will be executed during the restore(restore mode only), default value is 10.<br/>
  `-retries` Number of retires that will be executed if the requests fail, default value is 5.<br/>
  `-etcd-config` Mandatory etcd configuration file location, default value: "etcd-configuration.json". See fixtures folder for an example.<br/>
  `-dump` Location of the dump file data will be store in (in case of a dump) or load from (in case of a restore), default value: "etcd-dump.json".<br/>

    $ etcd-dump -config=myBackupConfig.json -retries=2 -etcd-config=myClusterConfig.json -dump=dataset.json -concurent-request=100 restore

