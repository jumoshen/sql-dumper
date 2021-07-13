#!/bin/sh
backupdir=$(pwd)/backup
time=$7
if [[ $1&&$2 ]]; then
	$1 -u$2 -P$3 -h$4  -p$5 $6 | gzip > $backupdir/myblog$time.sql.gz
	find $backupdir -name "myblog*.sql.gz" -type f -mtime +30 -exec rm {}  \; > /dev/null 2>&1
else
	echo "${time} :mysqldump or host missed\r\n"
fi