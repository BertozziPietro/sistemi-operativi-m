#!/bin/bash

CONFIG_FILE="test-config.txt"
set_to=false

while IFS= read -r line; do
    if [[ $line == FILE_OUTPUT=* ]]; then
        FILE_OUTPUT="${line#FILE_OUTPUT=}"
    elif [[ $line == FILE_RESULTS=* ]]; then
        FILE_RESULTS="${line#FILE_RESULTS=}"
    elif [[ $line == TERMINAL_* ]]; then
        updated_line=$(echo "$line" | sed 's/=[^=]*$/=false/')
        sed -i "s/^$line$/$updated_line/" "$CONFIG_FILE"
        set_to=true
    elif $set_to; then
        updated_line=$(echo "$line" | sed 's/=[^=]*$/=true/')
        sed -i "s/^$line$/$updated_line/" "$CONFIG_FILE"
    fi
done < "$CONFIG_FILE"

> "$FILE_OUTPUT"
> "$FILE_RESULTS"

echo "Reset completato con successo!"
