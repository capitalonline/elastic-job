#!/bin/bash
# Author:hukey
command_linebin='/mongodb/bin/mongo'
port=27017

if [ ! -d "/mongodb/backup/mongodbOplog_bak/mongo-$port" ];then
    mkdir -p /mongodb/backup/mongodbOplog_bak/mongo-$port
fi

if [ ! -d "/mongodb/backup/mongodbOplog_bak/log-$port" ];then
    mkdir -p /mongodb/backup/mongodbOplog_bak/log-$port
fi

bkdatapath=/mongodb/backup/mongodbOplog_bak/mongo-$port
bklogpath=/mongodb/backup/mongodbOplog_bak/log-$port

logfilename=$(date +"%Y%m%d")

echo "===MongoDB 端口为" $port "的差异备份开始，开始时间为" $(date -d today +"%Y%m%d%H%M%S")

paramBakEndDate=$(date +%s)
echo "===本次备份时间参数中的结束时间为：" $paramBakEndDate

diffTime=$(expr 65 \* 60)
echo "===备份设置的间隔时间为：" $diffTime

paramBakStartDate=$(expr $paramBakEndDate - $diffTime)
echo "===本次备份时间参数中的开始时间为：" $paramBakStartDate

diffTime=$(expr 61 \* 60)
paramAfterBakRequestStartDate=$(expr $paramBakEndDate - $diffTime)
echo "===为保证备份的连续性,本次备份后,oplog中的开始时间需小于：" $paramAfterBakRequestStartDate

bkfilename=$(date -d today +"%Y%m%d%H%M%S")

command_line="${command_linebin} 192.168.118.16:27017"

opmes=$(/bin/echo "db.printReplicationInfo()" | $command_line --quiet)

echo $opmes > /tmp/opdoctime$port.tmplog
opbktmplogfile=/tmp/opdoctime$port.tmplog
opstartmes=$(grep "oplog first event time" $opbktmplogfile | awk -F 'CST' '{print $1}' | awk -F 'oplog first event time: '  '{print $2}' | awk -F ' GMT' '{print $1}'  )
oplogRecordFirst=$(date -d "$opstartmes"  +%s)
echo "===oplog集合记录的开始时间为[格式化]：" $oplogRecordFirst
if [ $oplogRecordFirst -le $paramBakStartDate ]; then
    echo "Message --检查设置备份时间合理。备份参数的开始时间在oplog记录的时间范围内。"
else
    echo "Fatal Error --检查设置的备份时间不合理合理。备份参数的开始时间不在oplog记录的时间范围内。请调整oplog size或调整备份频率。本次备份可以持续进行，但还原时数据完整性丢失。"
fi

/mongodb/bin/mongodump -h 192.168.118.16 --port $port  -d local -c oplog.rs  --query '{ts:{$gte:Timestamp('$paramBakStartDate',1),$lte:Timestamp('$paramBakEndDate',9999)}}' -o $bkdatapath/mongodboplog$bkfilename


opmes=$(/bin/echo "db.printReplicationInfo()" | $command_line --quiet)
echo $opmes > /tmp/opdoctime$port.tmplog
opbktmplogfile=/tmp/opdoctime$port.tmplog
opstartmes=$(grep "oplog first event time" $opbktmplogfile | awk -F 'CST' '{print $1}' | awk -F 'oplog first event time: '  '{print $2}' | awk -F ' GMT' '{print $1}'  )
oplogRecordFirst=$(date -d "$opstartmes"  +%s)
echo "===执行备份后,oplog集合记录的开始时间为[时间格式化]:" $oplogRecordFirst

if [ $oplogRecordFirst -le $paramAfterBakRequestStartDate ]; then
    echo "Message --备份后，检查oplog集合中数据的开始时间，即集合中最早的一笔数据，时间不小于61分钟的时间（即参数 paramAfterBakRequestStartDate）。这样可以保证每个增量备份含有最近一个小时的全部op操作，满足文件的持续完整性，逐个还原无丢失数据风险。"
else
    echo "Fatal Error --备份后，检查oplog集合的涵盖的时间范围过小（小于61min）。设置的备份时间不合理合理，备份后的文件不能完全涵盖最近60分钟的数据。请调整oplog size或调整备份频率。本次备份可以持续进行，但还原时数据完整性丢失。"
fi

if [ -d "$bkdatapath/mongodboplog$bkfilename" ]
then
    echo "Message --检查此次备份文件已经产生.文件信息为:" $bkdatapath/mongodboplog$bkfilename >> $bklogpath/$logfilename.log
else
    echo "Fatal Error --备份过程已执行，但是未检测到备份产生的文件，请检查！" >> $bklogpath/$logfilename.log
fi

keepbaktime=$(date -d '-3 days' "+%Y%m%d%H")*
if [ -d $bkdatapath/mongodboplog$keepbaktime ]; then
    rm -rf $bkdatapath/mongodboplog$keepbaktime
    echo "Message -- $bkdatapath/mongodboplog$keepbaktime 删除完毕" >> $bklogpath/$logfilename.log
fi

echo "===MongoDB 端口为" $port "的差异备份结束，结束时间为：" $(date -d today +"%Y%m%d%H%M%S")