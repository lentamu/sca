#!/usr/bin/env sh

# wait-for-it.sh

set -e

TIMEOUT=30
WAIT_INTERVAL=2
HOST=$1
PORT=$2
shift 2

# Usage message
usage() {
    echo "Usage: wait-for-it.sh host:port [timeout] [command]"
    exit 1
}

# Check if the host and port are provided
if [ -z "$HOST" ] || [ -z "$PORT" ]; then
    usage
fi

# Wait for the specified host and port to be available
echo "Waiting for $HOST:$PORT..."

for i in $(seq 1 $TIMEOUT); do
    if nc -z "$HOST" "$PORT"; then
        echo "$HOST:$PORT is available."
        break
    fi
    echo "Waiting for $HOST:$PORT to be available..."
    sleep $WAIT_INTERVAL
done

if ! nc -z "$HOST" "$PORT"; then
    echo "Timeout reached while waiting for $HOST:$PORT."
    exit 1
fi

# Execute the provided command
exec "$@"
