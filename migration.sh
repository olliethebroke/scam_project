#!/bin/bash
source .env

sleep 2 && goose -dir "${MIGRATION_DIR}" postgres "${PG_DSN}" up -v