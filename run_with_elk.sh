#!/bin/bash
# Run error-simulator with Logstash forwarding (alternative to LOGSTASH_HOST env)
# When using go run main.go directly, set: LOGSTASH_HOST=localhost LOGSTASH_PORT=5001

set -e
cd "$(dirname "$0")"

export LOGSTASH_HOST="${LOGSTASH_HOST:-localhost}"
export GITHUB_REPOSITORY="${GITHUB_REPOSITORY:-dhruvkuchhal-pbhealth/error-simulator}"
export LOGSTASH_PORT="${LOGSTASH_PORT:-5001}"
export ELASTIC_APM_SERVICE_NAME="${ELASTIC_APM_SERVICE_NAME:-error-simulator}"
export ELASTIC_APM_SERVER_URL="${ELASTIC_APM_SERVER_URL:-http://localhost:8200}"

# Check if Logstash is reachable
if ! nc -z "$LOGSTASH_HOST" "$LOGSTASH_PORT" 2>/dev/null; then
  echo "⚠️  Logstash not reachable at $LOGSTASH_HOST:$LOGSTASH_PORT"
  echo "   Start ELK stack: cd elk-stack && docker compose up -d"
  echo ""
fi

echo "Starting error-simulator (logs → Logstash, APM → Kibana)..."
echo "Trigger errors: curl http://localhost:\${SERVER_PORT:-8092}/error/nil-pointer"
echo "APM/Logs: http://localhost:5601/app/apm/services/error-simulator/overview?rangeFrom=now-1h&rangeTo=now"
echo "  (Use 'Last 1 hour' time range in Kibana if you see 'No data to display')"
echo ""

exec go run main.go
