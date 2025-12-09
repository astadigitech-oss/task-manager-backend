#!/bin/bash
pkill -f "./app"
nohup ./app > app.log 2>&1 &
echo "App restarted!"