#!/bin/sh

set -o pipefail

oldpath=`pwd`

backup_dir="/var/backups/viittoen"
db_name="viittoen.db"
backup_db_name="viittoen-$(date +"%F").db"

mkdir -p $backup_dir

sqlite3 $db_name "VACUUM;"
sqlite3 $db_name ".backup '$backup_dir/$backup_db_name'"

cd $backup_dir
gzip $backup_db_name 
cd $oldpath