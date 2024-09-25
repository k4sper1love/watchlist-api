#!/bin/sh

cat <<EOF > /etc/prometheus/prometheus.yml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'api'
    static_configs:
      - targets: ['app:$APP_PORT']
EOF

exec /bin/prometheus --config.file=/etc/prometheus/prometheus.yml