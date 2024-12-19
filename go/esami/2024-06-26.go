package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	NUM_VISITATORI = 100
	NUM_DIPENDENTI = 100
	
	MAX      = 50
	MAX_BUFF = 100
	
	TEMPO_MIN        = 500
	TEMPO_INIZIO     = 500
	TEMPO_VISITATORE = 1000
	TEMPO_DIPENDENTE = 1000
	TEMPO_USCIERE    = 3000
)

var (
	entraUsciere = make(chan chan bool, MAX_BUFF)
	esceUsciere  = make(chan chan bool, MAX_BUFF)

	entraDipendente = make(chan chan bool, MAX_BUFF)
	esceDipendente  = make(chan chan bool, MAX_BUFF)

	entraVisitatore = [2]chan chan bool{make(chan chan bool, MAX_BUFF), make(chan chan bool, MAX_BUFF)}
	esceVisitatore  = [2]chan chan bool{make(chan chan bool, MAX_BUFF), make(chan chan bool, MAX_BUFF)}

	finito      = make(chan int, MAX_BUFF)
	bloccaSala  = make(chan int)
	terminaSala = make(chan int)
)

func when(b bool, c chan chan bool) chan chan bool {
	if !b {
		return nil
	}
	return c
}

func sala() {
	fmt.Println("[SALA] inizio")

	// preparazione
	var (
		usciere = 0
		dipendenti = 0
		visitatori = 0
		
		fine = false
	)
	liberi := func(aggiunti int) bool {
 	   return usciere + dipendenti + visitatori + aggiunti <= MAX
	}

	// comportamento
	for {
		fmt.Printf("[SALA] Usciere: %1d, Dipendenti: %03d, Visitatori: %03d, EntraU: %1d, EsceU: %1d, EntraD: %03d, EsceD: %03d, EntraVS: %03d, EsceVS: %03d, EntraVO: %03d, EsceVO: %03d\n", usciere, dipendenti, visitatori, len(entraUsciere), len(esceUsciere), len(entraDipendente), len(esceDipendente), len(entraVisitatore[0]), len(esceVisitatore[0]), len(entraVisitatore[1]), len(esceVisitatore[1]))
		select {
		case ack := <-when(!fine && visitatori == 0,
		entraUsciere): {
			usciere = 1
			ack <- true
		}
		case ack := <-when(fine,
		entraUsciere): {
			ack <- false
		}
		case ack := <-when(dipendenti == 0,
		esceUsciere): {
			usciere = 0
			ack <- true
		}
		case ack := <-when(usciere == 1 && len(esceUsciere) == 0 && liberi(1),
		entraDipendente): {
			dipendenti++
			ack <- true
		}
		case ack := <-esceDipendente: {
			dipendenti--
			ack <- true
		}
		case ack := <-when(usciere == 0 && len(entraUsciere) == 0 && len(entraDipendente) == 0 && liberi(2),
		entraVisitatore[0]): {
			visitatori += 2
			ack <- true
		}
		case ack := <-esceVisitatore[0]: {
			visitatori -= 2
			ack <- true
		}
		case ack := <-when(usciere == 0 && len(entraUsciere) == 0 && len(entraDipendente) == 0 && len(entraVisitatore[0]) == 0 && liberi(1),
		entraVisitatore[1]): {
			visitatori++
			ack <- true
		}
		case ack := <-esceVisitatore[1]: {
			visitatori--
			ack <- true
		}
		case <-bloccaSala: {
			fine = true
		}
		case <-terminaSala: {
			finito <- 1
			fmt.Println("[SALA] fine")
			return
		}}
	}
}

func usciere(id int) {
	fmt.Println("[USCIERE] inizio")

	// preparazione
	var ack = make(chan bool)
	
	// comportamento
	var continua bool 
	for {
		time.Sleep(time.Duration(rand.Intn(TEMPO_INIZIO)+TEMPO_MIN) * time.Millisecond)
		fmt.Println("[USCIERE] mi metto in coda per entrare nella sala") 
		entraUsciere <- ack
		continua = <-ack
		if !continua {
			finito <- 1
			fmt.Println("[USCIERE] fine")
			return
		}
		fmt.Println("[USCIERE] è il mio turno di entrare nella sala")
		time.Sleep(time.Duration(rand.Intn(TEMPO_USCIERE)+TEMPO_MIN) * time.Millisecond)
		fmt.Println("[USCIERE] mi metto in coda per uscire dalla sala") 
		esceUsciere <- ack
		<-ack
		fmt.Println("[USCIERE] è il mio turno di uscire dalla sala")
	}
}

func dipendente(id int) {
	time.Sleep(time.Duration(rand.Intn(TEMPO_INIZIO)+TEMPO_MIN) * time.Millisecond)
	fmt.Printf("[DIPENDENTE %03d] inizio\n", id)
	
	// preparazione
	var ack = make(chan bool)
	
	// comportamento
	fmt.Printf("[DIPENDENTE %03d] mi metto in coda per entrare nella sala\n", id)
	entraDipendente <- ack
	<-ack
	fmt.Printf("[DIPENDENTE %03d] è il mio turno di entrare nella sala\n", id)
	time.Sleep(time.Duration(rand.Intn(TEMPO_DIPENDENTE)+TEMPO_MIN) * time.Millisecond)
	fmt.Printf("[DIPENDENTE %03d] mi metto in coda per uscire dalla sala\n", id)
	esceDipendente <- ack
	<-ack
	fmt.Printf("[DIPENDENTE %03d] è il mio turno di entrare nella sala\n", id)
	
	finito <- 1
	fmt.Printf("[DIPENDENTE %03d] fine\n", id)
}

func visitatore(id int, tipo int) {
	time.Sleep(time.Duration(rand.Intn(TEMPO_INIZIO)+TEMPO_MIN) * time.Millisecond)
	fmt.Printf("[VISITATORE %03d] inizio\n", id)
	
	// preparazione
	var ack = make(chan bool)
	
	// comportamento
	fmt.Printf("[VISITATORE %03d] mi metto in coda per entrare nella sala\n", id)
	entraVisitatore[tipo] <- ack
	<-ack
	fmt.Printf("[VISITATORE %03d] è il mio turno di entrare nella sala\n", id)
	time.Sleep(time.Duration(rand.Intn(TEMPO_VISITATORE)+TEMPO_MIN) * time.Millisecond)
	fmt.Printf("[VISITATORE %03d] mi metto in coda per uscire dalla sala\n", id)
	esceVisitatore[tipo] <- ack
	<-ack
	fmt.Printf("[VISITATORE %03d] è il mio turno di entrare nella sala\n", id)
	
	finito <- 1
	fmt.Printf("[VISITATORE %03d] fine\n", id)
}

func main() {
	fmt.Println("[MAIN] inizio")
	rand.Seed(time.Now().Unix())
	
	// inizio goroutines
	go sala()
	go usciere(0)
	for i := 0; i < NUM_VISITATORI; i++ {
		go visitatore(i, i%2)
	}
	for i := 0; i < NUM_DIPENDENTI; i++ {
		go dipendente(i)
	}
	
	// fine goroutines
	for i := 0; i < NUM_VISITATORI+NUM_DIPENDENTI; i++ {
		<-finito
	}
	bloccaSala <- 1
	<-finito
	terminaSala <- 1
	<-finito
	
	fmt.Println("[MAIN] fine")
}
