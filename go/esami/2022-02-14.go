package main

import (
	"fmt"
	"math/rand"
	"time"
	"strings"
)

const (
	N        = 60
	MAX_BUFF = 30

	NUM_UTENTI = 100
	
	TEMPO_MINIMO = 500
	TEMPO_UTENTE = 1000
	
	AZIONI_UTENTE = 4
	
	TIPI_UTENTE = 3
)

var (
	canaliUtente = [TIPI_UTENTE][AZIONI_UTENTE]chan chan bool{
					{make(chan chan bool, MAX_BUFF), make(chan chan bool, MAX_BUFF), make(chan chan bool, MAX_BUFF), make(chan chan bool, MAX_BUFF)},
					{make(chan chan bool, MAX_BUFF), make(chan chan bool, MAX_BUFF), make(chan chan bool, MAX_BUFF), make(chan chan bool, MAX_BUFF)},
					{make(chan chan bool, MAX_BUFF), make(chan chan bool, MAX_BUFF), make(chan chan bool, MAX_BUFF), make(chan chan bool, MAX_BUFF)}}

	finito        = make(chan bool, MAX_BUFF)
	terminaMostra = make(chan bool)
)

func lunghezzeCanaliInMatrice(canali [][]chan chan bool) [][]int {
	lunghezze := make([][]int, len(canali))
	for i, riga := range canali {
		lunghezze[i] = make([]int, len(riga))
		for j, c := range riga { lunghezze[i][j] = len(c) }
	}
	return lunghezze
}

func when(b bool, c chan chan bool) chan chan bool {
	if !b { return nil }
	return c
}

func mostra() {
	const nome = "MOSTRA"
	var spazi = strings.Repeat(" ", len(nome)+3)
	fmt.Printf("[%s] inizio\n", nome)

	var (
		mostraUtenti = 0
	
		salaSorveglianti = 0
		
		corridoioScolaresche = 0
		corridoioDirezione = 0
	)
	
	liberi := func(aggiunti int) bool {
 		return mostraUtenti + aggiunti <= N
	}
	
	priorità := func(canali ...chan chan bool) bool {
		for _, c := range canali { if len(c) > 0 { return false } }
		return true
	}

	for {
		var canaliUtenteSlice = make([][]chan chan bool, len(canaliUtente))
		for i, row := range canaliUtente { canaliUtenteSlice[i] = append([]chan chan bool(nil), row[:]...) }
		var lunghezzeUtenti = lunghezzeCanaliInMatrice(canaliUtenteSlice)
		fmt.Printf("[%s] MostraUtenti: %03d, SalaSorveglianti: %03d, CorridoioScolaresche: %03d, CorridoioDirezione: %03d\n%sCanaliUtenti: %v\n",
		nome, mostraUtenti, salaSorveglianti, corridoioScolaresche, corridoioDirezione, spazi, lunghezzeUtenti)
		
		select {
		case ack := <-when(corridoioDirezione != -1 && liberi(1) &&
		priorità(canaliUtente[1][2], canaliUtente[2][2], canaliUtente[0][2]),
		canaliUtente[0][0]): {
			mostraUtenti++
			ack <- true
		}
		case ack := <-canaliUtente[0][1]: {
			salaSorveglianti++
			ack <- true
		}
		case ack := <-when(corridoioDirezione != 1 &&
		priorità(canaliUtente[1][2], canaliUtente[2][2]),
		canaliUtente[0][2]): {
			mostraUtenti--
			salaSorveglianti--
			ack <- true
		}
		case ack := <-canaliUtente[0][3]: {
			ack <- true
		}
		case ack := <-when(corridoioDirezione != -1 && liberi(25) &&
		priorità(canaliUtente[1][2], canaliUtente[2][2], canaliUtente[0][2], canaliUtente[0][0], canaliUtente[2][0]),
		canaliUtente[1][0]): {
			mostraUtenti += 25
			corridoioScolaresche++
			corridoioDirezione = 1
			ack <- true
		}
		case ack := <-canaliUtente[1][1]: {
			corridoioScolaresche--
			if corridoioScolaresche == 0 { corridoioDirezione = 0 }
			ack <- true
		}
		case ack := <-when(corridoioDirezione != 1,
		canaliUtente[1][2]): {
			mostraUtenti -= 25
			corridoioScolaresche++
			corridoioDirezione = -1
			ack <- true
		}
		case ack := <-canaliUtente[1][3]: {
			corridoioScolaresche--
			if corridoioScolaresche == 0 { corridoioDirezione = 0 }
			ack <- true
		}		
		case ack := <-when(corridoioDirezione != -1 && liberi(1) &&
		priorità(canaliUtente[1][2], canaliUtente[2][2], canaliUtente[0][2], canaliUtente[0][0]),
		canaliUtente[2][0]): {
			mostraUtenti++
			ack <- true
		}
		case ack := <-canaliUtente[2][1]: {
			ack <- true
		}
		case ack := <-when(corridoioDirezione != 1 &&
		priorità(canaliUtente[1][2]),
		canaliUtente[2][2]): {
			mostraUtenti--
			ack <- true
		}
		case ack := <-canaliUtente[2][3]: {
			ack <- true
		}
		case <-terminaMostra: {
			finito <- true
			fmt.Printf("[%s] fine\n", nome)
			return
		}}
	}
}

func utente(id int, tipo int) {
	const nome = "UTENTE"
	fmt.Printf("[%s %03d] inizio\n", nome, id)

	var (
		ack = make(chan bool)
		azioni = [AZIONI_UTENTE]string {"percorrere il corridoio verso la mostra", "entrare nella mostra", "percorrere il corridoio verso casa", "tornare a casa"}
	)
	for i, azione := range azioni {
		time.Sleep(time.Duration(rand.Intn(TEMPO_UTENTE)+TEMPO_MINIMO) * time.Millisecond)
		fmt.Printf("[%s %03d] mi metto in coda per %s\n", nome, id, azione)
		canaliUtente[tipo][i] <- ack
		<-ack
		fmt.Printf("[%s %03d] è il mio turno di %s\n", nome, id, azione)
	}
	
	finito <- true
	fmt.Printf("[%s %03d] fine\n", nome, id)
}

func main() {
	fmt.Println("[MAIN] inizio")
	rand.Seed(time.Now().Unix())
	
	go mostra()
	for i := 0; i < NUM_UTENTI; i++ {
		go utente(i, i%3)
	}
	
	for i := 0; i < NUM_UTENTI; i++ {
		<-finito
	}
	terminaMostra <- true
	<-finito
	
	fmt.Println("[MAIN] fine")
}
