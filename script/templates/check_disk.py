#!/usr/bin/env python

from subprocess import Popen, PIPE
import string


def partitions():

    command = ['lsblk -o name /dev/sd* -n -s -l ']
    output = Popen(command, shell=True, stdout=PIPE)
    output_string = output.stdout.read()
    output_string = output_string.strip().split('\n')
    myset = list(set(output_string))
    myset_list = [i.rstrip(string.digits) for i in myset]
    devices_list = []
    for item in myset_list:
        if myset_list.count(item) == 1:
            devices_list.append(item)
    return devices_list


if __name__ == '__main__':
    print(partitions())

