package main

import (
	"fmt"
	"math/rand"
	"time"
	"strings"
)

const (
	MAX      = 50
	MAX_BUFF = 100

	NUM_DIPENDENTI = 100
	NUM_VISITATORI = 100
	
	TEMPO_MINIMO     = 500
	TEMPO_USCIERE    = 3000
	TEMPO_DIPENDENTE = 1000
	TEMPO_VISITATORE = 1000
	
	AZIONI_USCIERE    = 2
	AZIONI_DIPENDENTE = 2
	AZIONI_VISITATORE = 2
	
	TIPI_VISITATORI = 2
)

var (
	canaliUsciere    = [AZIONI_USCIERE]chan chan bool{
					 make(chan chan bool, MAX_BUFF),
					 make(chan chan bool, MAX_BUFF)}	

	canaliDipendente = [AZIONI_DIPENDENTE]chan chan bool{
					 make(chan chan bool, MAX_BUFF),
					 make(chan chan bool, MAX_BUFF)}
					 
	canaliVisitatore = [TIPI_VISITATORI][AZIONI_VISITATORE]chan chan bool{
					{make(chan chan bool, MAX_BUFF), make(chan chan bool, MAX_BUFF)},
					{make(chan chan bool, MAX_BUFF), make(chan chan bool, MAX_BUFF)}}

	finito      = make(chan bool, MAX_BUFF)
	bloccaSala  = make(chan bool)
	terminaSala = make(chan bool)
)

func lunghezzeCanaliInVettore(canali []chan chan bool) []int {
	lunghezze := make([]int, len(canali))
	for i, c := range canali { lunghezze[i] = len(c) }
	return lunghezze
}

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

func priorità(canali ...chan chan bool) bool {
	for _, c := range canali { if len(c) > 0 { return false } }
	return true
}

func sala() {
	const nome = "SALA"
	var spazi = strings.Repeat(" ", len(nome) + 3)
	fmt.Printf("[%s] inizio\n", nome)

	var (
		usciere = 0
		dipendenti = 0
		visitatori = 0
		fine = false
	)
	
	liberi := func(aggiunti int) bool { return usciere + dipendenti + visitatori + aggiunti <= MAX }

	for {
		var canaliVisitatoriSlice = make([][]chan chan bool, len(canaliVisitatore))
		for i, row := range canaliVisitatore { canaliVisitatoriSlice[i] = append([]chan chan bool(nil), row[:]...) }
		var (
			lunghezzeUsciere    = lunghezzeCanaliInVettore(canaliUsciere[:])
			lunghezzeDipendente = lunghezzeCanaliInVettore(canaliDipendente[:])
			lunghezzeVisitatore = lunghezzeCanaliInMatrice(canaliVisitatoriSlice)
		)
		fmt.Printf("[%s] Usciere: %03d, Dipendenti: %03d, Visitatori: %03d, Fine: %5t\n%sCanaliUsciere: %v, CanaliDipendenti: %v, CanaliVisitatori: %v\n",
		nome, usciere, dipendenti, visitatori, fine, spazi, lunghezzeUsciere, lunghezzeDipendente, lunghezzeVisitatore)
		
		select {
		case ack := <-when(!fine && visitatori == 0,
		canaliUsciere[0]): {
			usciere = 1
			ack <- true
		}
		case ack := <-when(fine,
		canaliUsciere[0]): {
			ack <- false
		}
		case ack := <-when(dipendenti == 0,
		canaliUsciere[1]): {
			usciere = 0
			ack <- true
		}
		case ack := <-when(usciere == 1 && liberi(1) &&
		priorità(canaliUsciere[1]),
		canaliDipendente[0]): {
			dipendenti++
			ack <- true
		}
		case ack := <-canaliDipendente[1]: {
			dipendenti--
			ack <- true
		}
		case ack := <-when(usciere == 0 && liberi(2) &&
		priorità(canaliUsciere[1], canaliDipendente[1]),
		canaliVisitatore[0][0]): {
			visitatori += 2
			ack <- true
		}
		case ack := <-canaliVisitatore[0][1]: {
			visitatori -= 2
			ack <- true
		}
		case ack := <-when(usciere == 0 && liberi(1) &&
		priorità(canaliUsciere[0], canaliDipendente[0], canaliVisitatore[0][0]),
		canaliVisitatore[1][0]): {
			visitatori++
			ack <- true
		}
		case ack := <-canaliVisitatore[1][1]: {
			visitatori--
			ack <- true
		}
		case <-bloccaSala: {
			fine = true
		}
		case <-terminaSala: {
			finito <- true
			fmt.Printf("[%s] fine\n", nome)
			return
		}}
	}
}

func usciere(id int) {
	const nome = "USCIERE"
	fmt.Printf("[%s %03d] inizio\n", nome, id)

	var (
		ack = make(chan bool)
		azioni = [AZIONI_USCIERE]string {"per entrare nella sala", "uscire dalla sala"}
		continua bool
	)
	for {
		for i, azione := range azioni {
			time.Sleep(time.Duration(rand.Intn(TEMPO_USCIERE) + TEMPO_MINIMO) * time.Millisecond)
			fmt.Printf("[%s %03d] mi metto in coda per %s\n", nome, id, azione)
			canaliUsciere[i] <- ack
			continua = <-ack
			if !continua {
				finito <- true
				fmt.Printf("[%s %03d] fine\n", nome, id)
				return
			}
			fmt.Printf("[%s %03d] è il mio turno di %s\n", nome, id, azione)
		}
	}
}

func dipendente(id int) {
	const nome = "DIPENDENTE"
	fmt.Printf("[%s %03d] inizio\n", nome, id)
	
	var (
		ack = make(chan bool)
		azioni = [AZIONI_DIPENDENTE]string {"per entrare nella sala", "uscire dalla sala"}
	)
	for i, azione := range azioni {
		time.Sleep(time.Duration(rand.Intn(TEMPO_DIPENDENTE) + TEMPO_MINIMO) * time.Millisecond)
		fmt.Printf("[%s %03d] mi metto in coda per %s\n", nome, id, azione)
		canaliDipendente[i] <- ack
		<-ack
		fmt.Printf("[%s %03d] è il mio turno di %s\n", nome, id, azione)
	}
	
	fmt.Printf("[%s %03d] fine\n", nome, id)
	finito <- true
}

func visitatore(id int, tipo int) {
	const nome = "VISITATORE"
	fmt.Printf("[%s %03d] inizio\n", nome, id)
	
	var (
		ack = make(chan bool)
		azioni = [AZIONI_VISITATORE]string {"per entrare nella sala", "uscire dalla sala"}
	)
	for i, azione := range azioni {
		time.Sleep(time.Duration(rand.Intn(TEMPO_VISITATORE) + TEMPO_MINIMO) * time.Millisecond)
		fmt.Printf("[%s %03d] mi metto in coda per %s\n", nome, id, azione)
		canaliVisitatore[tipo][i] <- ack
		<-ack
		fmt.Printf("[%s %03d] è il mio turno di %s\n", nome, id, azione)
	}
	
	fmt.Printf("[%s %03d] fine\n", nome, id)
	finito <- true
}

func main() {
	fmt.Println("[MAIN] inizio")
	rand.Seed(time.Now().Unix())
	
	go sala()
	go usciere(0)
	for i := 0; i < NUM_DIPENDENTI; i++ { go dipendente(i) }
	for i := 0; i < NUM_VISITATORI; i++ { go visitatore(i, i % 2) }
	
	for i := 0; i < NUM_DIPENDENTI + NUM_VISITATORI; i++ { <-finito }
	bloccaSala <- true
	<-finito
	terminaSala <- true
	<-finito
	
	fmt.Println("[MAIN] fine")
}
