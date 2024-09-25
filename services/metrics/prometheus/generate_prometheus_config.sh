#!/bin/sh

# Defining the port depending on the USE_HTTPS flag
if [ "$USE_HTTPS" = "true" ]; then
  APP_PORT=443
else
  APP_PORT=${APP_PORT}
fi

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