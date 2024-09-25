#!/bin/sh

# Generating the Prometheus configuration
cat <<EOF > /etc/prometheus/prometheus.yml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'api'
    static_configs:
      - targets: ['app:$APP_PORT']

  - job_name: 'node'
    static_configs:
      - targets: ['node_exporter:9100']
EOF

# Launching Prometheus with the specified configuration
exec /bin/prometheus --config.file=/etc/prometheus/prometheus.yml