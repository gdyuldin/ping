#!/usr/bin/env bash

set -e

MOUNTPOINT=$(pwd)/$RANDOM
sudo mkdir $MOUNTPOINT
sudo guestmount -a cirros-0.3.4-x86_64-disk.img -m /dev/sda1 --rw $MOUNTPOINT
sudo cp ping $MOUNTPOINT/usr/bin/fix_id_ping
sudo chmod a+x $MOUNTPOINT/usr/bin/fix_id_ping
sudo chmod u+s $MOUNTPOINT/usr/bin/fix_id_ping
sudo umount $MOUNTPOINT
sudo rm -rf $MOUNTPOINT
