#!/bin/bash
sudo apt-get update -y
apt install unzip -y
wget -O v2ray-linux-64.zip https://github.com/BattleRoach/v2ray-core/releases/latest/download/v2ray-linux-64.zip
sudo unzip -o v2ray-linux-64.zip -d /usr/local/bin/v2ray
sudo chmod +x /usr/local/bin/v2ray/v2ray
sudo cp /usr/local/bin/v2ray/systemd/system/v2ray.service /lib/systemd/system
echo "{\"server\":\"$1\",\"node\":$2,\"token\":\"$3\"}"
echo "{\"server\":\"$1\",\"node\":$2,\"token\":\"$3\"}" > /usr/local/bin/v2ray/v2board.json
systemctl daemon-reload
service v2ray restart
