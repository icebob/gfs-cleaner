# GFS Cleaner
It is a GFS (Grandfather-Father-Son) backup files cleaner.

## Install

**Linux x64**
```bash
sudo wget https://github.com/icebob/gfs-cleaner/releases/download/v0.1.0/gfs-cleaner_0.1.0_Linux_x86_64.tar.gz
sudo tar -xvzf gfs-cleaner_0.1.0_Linux_x86_64.tar.gz -C /usr/local/bin/ gfs-cleaner
sudo chmod +x /usr/local/bin/gfs-cleaner
```

## Usage

**Clean old backup files from `backup` folder**
```bash
$ gfs-cleaner clean backup
```
Default retentions:   
- 7 daily backups
- 4 weekly backups
- 12 monthly backups
- 10 yearly backups

**Clean old backup files with custom retention values**
```bash
$ gfs-cleaner clean backup --daily 5 --weekly 2 --monthly 10 --yearly 10
```

## License
gfs-cleaner is available under the [MIT license](https://tldrlegal.com/license/mit-license).

## Contact

Copyright (C) 2019 Icebob

[![@icebob](https://img.shields.io/badge/github-icebob-green.svg)](https://github.com/icebob) [![@icebob](https://img.shields.io/badge/twitter-Icebobcsi-blue.svg)](https://twitter.com/Icebobcsi)
