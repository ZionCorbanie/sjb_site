#!/usr/bin/env sh
set -e  # Exit on error

cd /app

if [ "$ENVIRONMENT" = "development" ]; then
    echo "Starting in development mode..."

    echo "Starting dev server..."
    air &

    echo "Starting tailwind watch..."
	#pnpm build &

    echo "Starting templ watch..."
    make templ-watch &

    npx @tailwindcss/cli -i ./static/css/input.css -o ./static/css/style.css --watch

    wait
else
    echo "Starting in production mode..."
    make build
    ./bin/app
fi
