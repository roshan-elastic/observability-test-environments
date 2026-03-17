#!/usr/bin/env bash

echo "*/2 * * * * admin /tmp/rta-random-attack.sh > /tmp/rta-random.log 2>&1" | sudo tee /etc/cron.d/rta-random

sudo chmod ugo+rx /tmp/*.sh
