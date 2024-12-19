#!/bin/bash

#configurazione

CONFIG_FILE="test-config.txt"
if [[ ! -f $CONFIG_FILE ]]; then
    echo "[CONFIG_FILE] File di configurazione $CONFIG_FILE non trovato!"
    exit 1
fi
if [[ ! -r $CONFIG_FILE ]]; then
    echo "[CONFIG_FILE] Il file di configurazione $CONFIG_FILE non è leggibile!"
    exit 1
fi
while IFS= read -r line; do
    if [[ -z "$line" || "$line" =~ ^# ]]; then
        continue
    fi
    if [[ ! "$line" =~ ^[A-Z_]+=([-a-z0-9_\.]+)$ ]]; then
        echo "[CONFIG_FILE] '$line' non è nel formato corretto (PARAMETRO=valore)!"
        exit 1
    fi
done < "$CONFIG_FILE"
source "$CONFIG_FILE"

if [[ -z $FILE_GO ]]; then
    echo "[FILE_GO] La variabile FILE_GO non è definita in $CONFIG_FILE!"
    exit 1
fi
if [[ ! -f $FILE_GO ]]; then
    echo "[FILE_GO] File $FILE_GO non trovato!"
    exit 1
fi
if [[ "${FILE_GO##*.}" != "go" ]]; then
    echo "[FILE_GO] Il file $FILE_GO non ha estensione .go!"
    exit 1
fi

if [[ $FILE_OUTPUT != *.txt ]]; then
    echo "[FILE_OUTPUT] Il file $FILE_OUTPUT non ha estensione .txt!"
    exit 1
fi

if [[ $FILE_RESULTS != *.txt ]]; then
    echo "[FILE_RESULTS] Il file $FILE_RESULTS non ha estensione .txt!"
    exit 1
fi

if [[ -z $TERMINAL_OUTPUT ]]; then
    echo "[TERMINAL_OUTPUT] La variabile TERMINAL_OUTPUT non è definita in $CONFIG_FILE!"
    exit 1
fi
if [[ $TERMINAL_OUTPUT != "true" && $TERMINAL_OUTPUT != "false" ]]; then
    echo "[TERMINAL_OUTPUT] Il parametro $TERMINAL_OUTPUT non è booleano!"
    exit 1
fi

if [[ -z $TERMINAL_RESULTS ]]; then
    echo "[TERMINAL_RESULTS] La variabile TERMINAL_RESULTS non è definita in $CONFIG_FILE!"
    exit 1
fi
if [[ $TERMINAL_RESULTS != "true" && $TERMINAL_RESULTS != "false" ]]; then
    echo "[TERMINAL_RESULTS] Il parametro $TERMINAL_RESULTS non è booleano!"
    exit 1
fi

if [[ -z $DEADLOCK_ALERTER ]]; then
    echo "[DEADLOCK_ALERTER] La variabile DEADLOCK_ALERTER non è definita in $CONFIG_FILE!"
    exit 1
fi
if [[ $DEADLOCK_ALERTER != "true" && $DEADLOCK_ALERTER != "false" ]]; then
    echo "[DEADLOCK_ALERTER] Il parametro $DEADLOCK_ALERTER non è booleano!"
    exit 1
fi

if [[ -z $MESSAGE_COUNTER ]]; then
    echo "[MESSAGE_COUNTER] La variabile MESSAGE_COUNTER non è definita in $CONFIG_FILE!"
    exit 1
fi
if [[ $MESSAGE_COUNTER != "true" && $MESSAGE_COUNTER != "false" ]]; then
    echo "[MESSAGE_COUNTER] Il parametro $MESSAGE_COUNTER non è booleano!"
    exit 1
fi

if [[ -z $CHARACTER_COUNTER ]]; then
    echo "[CHARACTER_COUNTER] La variabile CHARACTER_COUNTER non è definita in $CONFIG_FILE!"
    exit 1
fi
if [[ $CHARACTER_COUNTER != "true" && $CHARACTER_COUNTER != "false" ]]; then
    echo "[CHARACTER_COUNTER] Il parametro $CHARACTER_COUNTER non è booleano!"
    exit 1
fi

if [[ -z $PROTAGONIST_IDENTIFIER ]]; then
    echo "[PROTAGONIST_IDENTIFIER] La variabile PROTAGONIST_IDENTIFIER non è definita in $CONFIG_FILE!"
    exit 1
fi
if [[ $PROTAGONIST_IDENTIFIER != "true" && $PROTAGONIST_IDENTIFIER != "false" ]]; then
    echo "[PROTAGONIST_IDENTIFIER] Il parametro $PROTAGONIST_IDENTIFIER non è booleano!"
    exit 1
fi

if [[ -z $ENDING_SEEKER ]]; then
    echo "[ENDING_SEEKER] La variabile ENDING_SEEKER non è definita in $CONFIG_FILE!"
    exit 1
fi
if [[ $ENDING_SEEKER != "true" && $ENDING_SEEKER != "false" ]]; then
    echo "[ENDING_SEEKER] Il parametro $ENDING_SEEKER non è booleano!"
    exit 1
fi

if [[ -z $CONVERSATION_SYNTHESIZER ]]; then
    echo "[CONVERSATION_SYNTHESIZER] La variabile CONVERSATION_SYNTHESIZER non è definita in $CONFIG_FILE!"
    exit 1
fi
if [[ $CONVERSATION_SYNTHESIZER != "true" && $CONVERSATION_SYNTHESIZER != "false" ]]; then
    echo "[CONVERSATION_SYNTHESIZER] Il parametro $CONVERSATION_SYNTHESIZER non è booleano!"
    exit 1
fi

#esecuzione

go run "$FILE_GO" > "$FILE_OUTPUT" 2>&1
if [[ $TERMINAL_OUTPUT == "true" ]]; then
    echo -e "\n*** RISULTATO DELL'ESECUZIONE ***\n"
    cat "$FILE_OUTPUT"
    echo -e "\n"
fi
if [[ $? -ne 0 ]]; then
    echo "[$FILE_GO] Errore durante l'esecuzione di $FILE_GO ..."
    exit 1
fi

#analisi

if [[ $DEADLOCK_ALERTER == "true" ]]; then
    DEADLOCK_ALERTER=$(grep -Fxq "fatal error: all goroutines are asleep - deadlock!" "$FILE_OUTPUT" && echo "Deadlock Fatale" || echo "Nessun Deadlock")
    if [[ $? -ne 0 ]]; then
        echo "[[DEADLOCK_ALERTER]]\n\nErrore durante l'avvertimento dell'eventuale deadlock" 
        exit 1
    fi
    echo -e "[[DEADLOCK_ALERTER]]\n\n$DEADLOCK_ALERTER\n" >> "$FILE_RESULTS"
fi

if [[ $MESSAGE_COUNTER == "true" ]]; then
    MESSAGE_COUNTER=$(grep -o '\[.*\]' "$FILE_OUTPUT" | wc -l)
    if [[ $? -ne 0 ]]; then
        echo "[[MESSAGE_COUNTER]]\n\nErrore durante il conteggio dei messaggi" 
        exit 1
    fi
    echo -e "[[MESSAGE_COUNTER]]\n\n$MESSAGE_COUNTER\n" >> "$FILE_RESULTS"
fi

if [[ $CHARACTER_COUNTER == "true" ]]; then
    CHARACTER_COUNTER=$(grep -o '\[.*\]' "$FILE_OUTPUT" | sort | uniq | wc -l)
    if [[ $? -ne 0 ]]; then
        echo "[[CHARACTER_COUNTER]]\n\nErrore durante il conteggio dei personaggi" 
        exit 1
    fi
    echo -e "[[CHARACTER_COUNTER]]\n\n$CHARACTER_COUNTER\n" >> "$FILE_RESULTS"
fi

if [[ $PROTAGONIST_IDENTIFIER == "true" ]]; then
    PROTAGONIST_IDENTIFIER=$(grep -o '\[.*\]' "$FILE_OUTPUT" | sort | uniq -c | sort -nr | awk '{count=$1; $1=""; printf "%-20s %s\n", $0, count}' | sed 's/^ //')
    if [[ $? -ne 0 ]]; then
        echo "[[PROTAGONIST_IDENTIFIER]]\n\nErrore durante l'identificazione dei protagonisti"
        exit 1
    fi
    echo -e "[[PROTAGONIST_IDENTIFIER]]\n\n$PROTAGONIST_IDENTIFIER\n" >> "$FILE_RESULTS"
fi

if [[ $ENDING_SEEKER == "true" ]]; then	
	ENDING_SEEKER=$(inizio_count=$(grep -o "inizio" "$FILE_OUTPUT" | wc -l); fine_count=$(grep -o "fine" "$FILE_OUTPUT" | wc -l); echo "Inizio: $inizio_count, Fine: $fine_count, Differenza: $((inizio_count - fine_count))"; awk '/\[.*\] inizio/ {match($0, /\[([^\]]+)\] inizio/, arr); if (arr[1] != "") {inizio[arr[1]] = 1}} /\[.*\] fine/ {match($0, /\[([^\]]+)\] fine/, arr); if (arr[1] != "") {delete inizio[arr[1]]}} END {for (name in inizio) print "[" name "]"}' "$FILE_OUTPUT" | sort)


    if [[ $? -ne 0 ]]; then
        echo "[[ENDING_SEEKER]]\n\nErrore durante la ricerca delle conclusioni" 
        exit 1
    fi
    echo -e "[[ENDING_SEEKER]]\n\n$ENDING_SEEKER\n" >> "$FILE_RESULTS"
fi

if [[ $CONVERSATION_SYNTHESIZER == "true" ]]; then
    CONVERSATION_SYNTHESIZER=$(grep -o '\[.*\]' "$FILE_OUTPUT" | uniq -c | awk '{count=$1; $1=""; printf "%-20s %s\n", $0, count}' | sed 's/^ //')
    if [[ $? -ne 0 ]]; then
        echo "[[CONVERSATION_SYNTHESIZER]]\n\nErrore durante la sintesi della conversazione"
        exit 1
    fi
    echo -e "[[CONVERSATION_SYNTHESIZER]]\n\n$CONVERSATION_SYNTHESIZER\n" >> "$FILE_RESULTS"
fi

if [[ $TERMINAL_RESULTS == "true" ]]; then
    echo -e "\n*** RISULTATO DELL'ANALISI ***\n"
    cat "$FILE_RESULTS"
fi
