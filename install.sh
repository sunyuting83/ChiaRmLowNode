#!/bin/sh

# install Chia harvester
basepath=$(cd `dirname $0`; pwd)
user=${USER}
myhome=${HOME}
userbin=$myhome/.local/bin
echo "Start..."

sudo -S install -Dm777 $basepath/rmlownode $userbin/ &&
sed -i 's/testuseruzuzuz/'$user'/g' $basepath/rmlownode.service &&
sed -i 's#testpath#'$basepath'#g' $basepath/rmlownode.service &&
sudo -S install -Dm644 $basepath/rmlownode.service /usr/lib/systemd/system/rmlownode.service &&
sudo -S systemctl enable rmlownode &&
sudo -S systemctl daemon-reload &&
sudo -S systemctl start rmlownode &&
echo "install complete"
