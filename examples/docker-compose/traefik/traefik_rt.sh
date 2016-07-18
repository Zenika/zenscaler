#!/bin/sh
curl -s http://localhost:8080/health | perl -n -e'/average_response_time_sec\":([\d|\.]+)/ && print $1'
