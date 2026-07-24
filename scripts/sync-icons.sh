#!/usr/bin/env bash
set -euo pipefail

ICON_SOURCE_DIR="node_modules/lucide-static/icons"
ICON_TARGET_DIR="internal/web/views/icons"

icons=(
  search
  log-out
  external-link
  chevron-down
  home
  users
  package
  clipboard-list
  shield
  plus
  receipt-text
  filter
  calendar
  sun
  moon
  monitor
  user
  key-round
)

mkdir -p "$ICON_TARGET_DIR"

for icon in "${icons[@]}"; do
  cp "$ICON_SOURCE_DIR/$icon.svg" "$ICON_TARGET_DIR/$icon.svg"
done

echo "Synced ${#icons[@]} Lucide icons."
