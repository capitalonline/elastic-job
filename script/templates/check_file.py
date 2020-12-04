#!/usr/bin/env python
#_*_coding:utf-8_*_

import hashlib
import os
import datetime
import sys

def check_file_path(path):
    path_list = os.listdir(path)
    path_list.sort(key=lambda fn: os.path.getmtime(path + '/' + fn) if not os.path.isdir(path + '/' + fn) else 0)
    file_path = os.path.join(path, path_list[-1])
    date_time = datetime.datetime.fromtimestamp(os.path.getmtime(file_path))
    fsize = os.path.getsize(file_path)
    return file_path, date_time.strftime("%Y-%m-%d %H:%M:%S"), fsize

if __name__ == '__main__':
    # dir path
    backup_backup = sys.argv[1]
    path,time, fsize = check_file_path(backup_backup)
    print path,"\n",fsize