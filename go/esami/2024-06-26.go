/*
package main

import (
	"fmt"
	"math/rand"
	"time"
)

type req struct {
	id  int
	ack chan bool
	i   int
}

//costanti
const MAX = 50
const numVisitatori = 100
const numDipendenti = 100
const ORDIINARIO = 0
const SPECIALE = 1
const MAXBUF = 100

// canali:
var entraUsciere = make(chan req, MAXBUF)
var esceUsciere = make(chan bool, MAXBUF)
var entraVisitatoreOrdinario = make(chan req, MAXBUF)
var esceVisitatoreOrdinario = make(chan req, MAXBUF)
var entraVisitatoreSpeciale = make(chan req, MAXBUF)
var esceVisitatoreSpeciale = make(chan req, MAXBUF)
var entraDipendente = make(chan req, MAXBUF)
var esceDipendente = make(chan req, MAXBUF)
var ackUsciere = make(chan bool, MAXBUF) //visibile ai dipendenti
var done = make(chan int, 100)
var termina = make(chan int)
var terminaSala = make(chan int)

func when(b bool, c chan req) chan req {
	if !b {
		return nil
	}
	return c
}

func server() { //server: gestisce accessi e uscite dalla sala
	var usciere = 0
	var dipendenti = 0
	var visitatori = 0
	var chiusura = false //settata dall'usciere quando vuole chiudere la riunione (per impedire l'ingresso di nuovi dipendenti)
	var fine = false

	for {
		fmt.Println("[Stato Sala] usciere = ", usciere, "dipendenti = ", dipendenti, ", visitatori = ", visitatori, " , fine = ", fine, ", entraUsciere = ", len(entraUsciere), ", esceUsciere = ", len(esceUsciere), ", entraDipendente = ", len(entraDipendente), ", esceDipendente = ", len(esceDipendente), ", entraVisitatoreOrdinario = ", len(entraVisitatoreOrdinario), "entraVisitatoreSpeciale = ", len(entraVisitatoreSpeciale))
		select {
		case r := <-when(!fine && dipendenti+visitatori+usciere == 0, entraUsciere):
			{
				fmt.Println("[Server] l'usciere è entrato ")
				usciere = 1
				r.ack <- true
			}
		case r := <-when(usciere == 1 && dipendenti+usciere < MAX && !chiusura, entraDipendente): //il controllo della coda dell'usciere non è necessario
			{
				fmt.Println("[Server] è entrato il Dipendente ", r.id)
				dipendenti++
				r.ack <- true
			}
		case r := <-when(usciere == 0 && visitatori+2 <= MAX && len(entraUsciere) == 0 && len(entraDipendente) == 0, entraVisitatoreSpeciale):
			{
				fmt.Println("[Server] è entrato il Visitatore Speciale ", r.id)
				visitatori += 2
				r.ack <- true
			}
		case r := <-when(usciere == 0 && visitatori < MAX && len(entraUsciere) == 0 && len(entraDipendente) == 0 && len(entraVisitatoreSpeciale) == 0, entraVisitatoreOrdinario):
			{
				fmt.Println("[Server] è entrato il Visitatore Ordinario ", r.id)
				visitatori++
				r.ack <- true
			}
		case <-esceUsciere:
			{
				if dipendenti == 0 { //l'usciere può chiudere
					fmt.Println("[Server] l'Usciere è uscito")
					usciere = 0
					ackUsciere <- true
				} else { // l'usciere attende:
					chiusura = true //l'usciere chiude la sala per impedire a nuovi dipendenti di entrare
					//nessuna risposta: l'usciere aspetta che l'ultimo dipendente se ne vada
					fmt.Println("[Server] l'Usciere attende che la sala si vuoti, prima di uscire ..")
				}
			}
		case r := <-esceDipendente:
			{
				fmt.Println("[Server] esce il Dipendente ", r.id)
				dipendenti--
				if dipendenti == 0 { //se è l'ultimo ad uscire, sblocco l'usciere
					usciere = 0
					chiusura = false
					fmt.Println("[Server] Usciere sbloccato!")
					ackUsciere <- true
				}
				r.ack <- true
			}
		case r := <-esceVisitatoreSpeciale:
			{
				fmt.Println("[Server] esce il Visitatore Speciale ", r.id)
				visitatori = visitatori - 2
				r.ack <- true
			}
		case r := <-esceVisitatoreOrdinario:
			{
				fmt.Println("[Server] esce il Visitatore Ordinario ", r.id)
				visitatori--
				r.ack <- true
			}
		case r := <-when(fine, entraUsciere):
			{
				r.ack <- false
			}
		case <-terminaSala: //quando tutti i processi visitatori e dipendenti terminano, il server termina
			{
				done <- 1 //comunicazione al main: il server sta terminando
				return
			}
		case <-termina: //quando tutti i processi visitatori e dipendenti terminano, il server si prepara a terminare
			{
				fine = true
			}
		}
	}
}

func usciere(id int) {
	var r req
	r.id = id
	r.ack = make(chan bool)
	var continua bool
	for {
		time.Sleep(time.Duration(rand.Intn(1000)+500) * time.Millisecond)
		entraUsciere <- r
		continua = <-r.ack
		if !continua {
			done <- 1
			return
		}
		time.Sleep(time.Duration(rand.Intn(3000)+500) * time.Millisecond)
		esceUsciere <- true
		<-ackUsciere
	}
}

func dipendente(id int) {
	var r req
	time.Sleep(time.Duration(rand.Intn(1000)+500) * time.Millisecond)
	r.id = id
	r.ack = make(chan bool)
	entraDipendente <- r
	<-r.ack
	time.Sleep(time.Duration(rand.Intn(1000)+500) * time.Millisecond)
	esceDipendente <- r
	<-r.ack
	done <- 1
}

func visitatore(id int, tipo int) {
	var r req
	time.Sleep(time.Duration(rand.Intn(1000)+500) * time.Millisecond)
	r.id = id
	r.ack = make(chan bool)
	if tipo == ORDIINARIO {
		entraVisitatoreOrdinario <- r
		<-r.ack
		time.Sleep(time.Duration(rand.Intn(1000)+500) * time.Millisecond)
		esceVisitatoreOrdinario <- r
		<-r.ack
	} else {
		entraVisitatoreSpeciale <- r
		<-r.ack
		time.Sleep(time.Duration(rand.Intn(1000)+500) * time.Millisecond)
		esceVisitatoreSpeciale <- r
		<-r.ack
	}
	done <- 1
}

func main() {
	rand.Seed(time.Now().Unix())

	go server()

	go usciere(0)

	for i := 0; i < numVisitatori; i++ {
		go visitatore(i, i%2)
	}

	for i := 0; i < numDipendenti; i++ {
		go dipendente(i)
	}
	time.Sleep(time.Duration(rand.Intn(8000)+500) * time.Millisecond)

	for i := 0; i < numVisitatori+numDipendenti; i++ {
		<-done
		//fmt.Println("[MAIN] fine processo", i)
	}

	termina <- 1

	<-done
	fmt.Println("[MAIN] fine usciere")

	terminaSala <- 1
	<-done // attesa terminazione server
}
*/

package main

import (
	"fmt"
	"math/rand"
	"time"
)

type req struct {
	id  int
	ack chan bool
}

const (
	ORDINARIO = 0
	SPECIALE = 1
	
	NUM_VISITATORI = 100
	NUM_DIPENDENTI = 100
	
	MAX = 50
	MAX_BUFF = 100
)

var (
	entraUsciere = make(chan req, MAX_BUFF)
	esceUsciere = make(chan req, MAX_BUFF)
	
	entraVisitatoreOrdinario = make(chan req, MAX_BUFF)
	esceVisitatoreOrdinario = make(chan req, MAX_BUFF)
	
	entraVisitatoreSpeciale = make(chan req, MAX_BUFF)
	esceVisitatoreSpeciale = make(chan req, MAX_BUFF)
	
	entraDipendente = make(chan req, MAX_BUFF)
	esceDipendente = make(chan req, MAX_BUFF)
	
	finito = make(chan int, MAX_BUFF)
	bloccaSala = make(chan int)
	terminaSala = make(chan int)
)

func when(b bool, c chan req) chan req {
	if !b {
		return nil
	}
	return c
}

func main() {
	fmt.Println("[MAIN] inizio")
	
	//inizio goroutines
	rand.Seed(time.Now().Unix())
	go sala()
	go usciere(0)
	for i := 0; i < NUM_VISITATORI; i++ {
		go visitatore(i, i%2)
	}
	for i := 0; i < NUM_DIPENDENTI; i++ {
		go dipendente(i)
	}
	
	//fine goroutines
	for i := 0; i < NUM_VISITATORI+NUM_DIPENDENTI; i++ {
		<-finito
	}
	bloccaSala <- 1
	<-finito
	terminaSala <- 1
	<-finito
	
	fmt.Println("[MAIN] fine")
}
