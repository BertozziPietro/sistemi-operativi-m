package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	NUM_CITTADINI = 90
	
	M1       = 10
	M2       = 20
	N 		 = 30
	MAX_BUFF = 50
	
	TEMPO_MIN       = 200
	TEMPO_INIZIO    = 200
	TEMPO_CITTADINO = 500
	TEMPO_ADDETTO   = 1000
)

var (
	entraAddetto = make(chan chan bool, MAX_BUFF)
	esceAddetto  = make(chan chan bool, MAX_BUFF)

	entraCittadino = [2]chan chan bool{make(chan chan bool, MAX_BUFF), make(chan chan bool, MAX_BUFF)}
	esceCittadino  = [2]chan chan bool{make(chan chan bool, MAX_BUFF), make(chan chan bool, MAX_BUFF)}

	finito           = make(chan bool, MAX_BUFF)
	bloccaCasaAcqua  = make(chan bool)
	terminaCasaAcqua = make(chan bool)
)

func when(b bool, c chan chan bool) chan chan bool {
	if !b {
		return nil
	}
	return c
}

func casaAcqua() {
	fmt.Println("[CASA_ACQUA] inizio")

	var (
		monetine = 0
		monetone = 0
		litri = N
		
		dentro = 0
		
		fine = false
	)
	
	piena := func() bool {
 		return monetine == M1 || monetone == M2
	}
	vuoto := func() bool {
 		return litri == 0
	}

	for {
		fmt.Printf("[CASA_ACQUA] Monetine: %03d, Monetone: %03d, Litri: %03d, Dentro: %03d, Fine: %5t, EntraA: %1d, EsceA: %1d, EntraCP: %03d, EsceCP: %03d, EntraCG: %03d, EsceCG: %03d\n", monetine, monetone, litri, dentro, fine, len(entraAddetto), len(esceAddetto), len(entraCittadino[0]), len(esceCittadino[0]), len(entraCittadino[1]), len(esceCittadino[1]))
		
		select {
		case ack := <-when(!fine && dentro == 0 && (piena() || vuoto()),
		entraAddetto): {
			monetine = 0
			monetone = 0
			litri = N
			dentro = 2
			ack <- true
		}
		case ack := <-when(fine,
		entraAddetto): {
			ack <- false
		}
		case ack := <-esceAddetto: {
			dentro = 0
			ack <- true
		}
		case ack := <-when(dentro == 0 && !(piena() || vuoto()),
		entraCittadino[0]): {
			monetine++
			litri--
			dentro = 1
			ack <- true
		}
		case ack := <-esceCittadino[0]: {
			dentro = 0
			ack <- true
		}
		case ack := <-when(dentro == 0 && !(piena() || vuoto()) && len(entraCittadino[0]) == 0,
		entraCittadino[1]): {
			monetone++
			litri -= 2
			dentro = 1
			ack <- true
		}
		case ack := <-esceCittadino[1]: {
			dentro = 0
			ack <- true
		}
		case <-bloccaCasaAcqua: {
			fine = true
		}
		case <-terminaCasaAcqua: {
			finito <- true
			fmt.Println("[CASA_ACQUA] fine")
			return
		}}
	}
}

func addetto(id int) {
	fmt.Println("[ADDETTO] inizio")

	// preparazione
	var ack = make(chan bool)
	
	// comportamento
	var continua bool 
	for {
		time.Sleep(time.Duration(rand.Intn(TEMPO_INIZIO)+TEMPO_MIN) * time.Millisecond)
		fmt.Println("[ADDETTO] mi metto in coda per entrare nella casa dell'acqua") 
		entraAddetto <- ack
		continua = <-ack
		if !continua {
			finito <- true
			fmt.Println("[ADDETTO] fine")
			return
		}
		fmt.Println("[ADDETTO] è il mio turno di entrare nella casa dell'acqua")
		time.Sleep(time.Duration(rand.Intn(TEMPO_ADDETTO)+TEMPO_MIN) * time.Millisecond)
		fmt.Println("[ADDETTO] mi metto in coda per entrare nella casa dell'acqua") 
		esceAddetto <- ack
		<-ack
		fmt.Println("[ADDETTO] è il mio turno di entrare nella casa dell'acqua")
	}
}

func cittadino(id int, tipo int) {
	time.Sleep(time.Duration(rand.Intn(TEMPO_INIZIO)+TEMPO_MIN) * time.Millisecond)
	fmt.Printf("[CITTADINO %03d] inizio\n", id)
	
	// preparazione
	var ack = make(chan bool)
	
	// comportamento
	fmt.Printf("[CITTADINO %03d] mi metto in coda per entrare nella casa dell'acqua\n", id)
	entraCittadino[tipo] <- ack
	<-ack
	fmt.Printf("[CITTADINO %03d] è il mio turno di entrare nella casa dell'acqua\n", id)
	time.Sleep(time.Duration(rand.Intn(TEMPO_CITTADINO)+TEMPO_MIN) * time.Millisecond)
	fmt.Printf("[CITTADINO %03d] mi metto in coda per entrare nella casa dell'acqua\n", id)
	esceCittadino[tipo] <- ack
	<-ack
	fmt.Printf("[CITTADINO %03d] è il mio turno di entrare nella casa dell'acqua\n", id)
	
	finito <- true
	fmt.Printf("[CITTADINO %03d] fine\n", id)
}

func main() {
	fmt.Println("[MAIN] inizio")
	rand.Seed(time.Now().Unix())
	
	ruotaTipo := func(i int) int {
		if i%3 == 0 { return 1 } else { return 0 }
	}
	
	// inizio goroutines
	go casaAcqua()
	go addetto(0)
	for i := 0; i < NUM_CITTADINI; i++ {
		go cittadino(i, ruotaTipo(i))
	}
	
	// fine goroutines
	for i := 0; i < NUM_CITTADINI; i++ {
		<-finito
	}
	bloccaCasaAcqua <- true
	<-finito
	terminaCasaAcqua <- true
	<-finito
	
	fmt.Println("[MAIN] fine")
}
