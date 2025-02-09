#!/bin/bash

input_file1="input1.txt"
output_file1="filtered_output1.txt"
input_file2="input2.txt"
output_file2="filtered_output2.txt"

awk -F";" '{
    # Sostituisce la virgola decimale con il carattere speciale §
    gsub(/,/, "§", $0);

    # Controlla se il secondo numero è tra quelli validi
    if ($2 == 1 || $2 == 2 || $2 == 3 || $2 == 6 || $2 == 9 || $2 == 12 || $2 == 15 || $2 == 18) {
        # Modifica il primo numero aggiungendo "x" e sé stesso
        new_first = $1 "x" $1;

        # Stampa il risultato con i numeri tra "" e aggiunge una virgola alla fine
        printf "\"%s\",\"%s\",\"%s\",\"%s\",\n", new_first, $2, $3, $4;
    }
}' "$input_file1" | sed 's/§/,/g' > "$output_file1"

awk -F";" '{
    # Sostituisce la virgola decimale con il carattere speciale §
    gsub(/,/, "§", $0);

    # Controlla se il secondo numero è tra quelli validi
    if ($2 == 1 || $2 == 2 || $2 == 3 || $2 == 6 || $2 == 9 || $2 == 12 || $2 == 15 || $2 == 18) {
        # Modifica il primo numero aggiungendo "x" e sé stesso
        new_first = $1 "x" $1;

        # Stampa il risultato con i numeri tra "" e aggiunge una virgola alla fine
        printf "\"%s\",\"%s\",\"%s\",\"%s\",\n", new_first, $2, $3, $4;
    }
}' "$input_file2" | sed 's/§/,/g' > "$output_file2"

echo "Filtro applicato"
