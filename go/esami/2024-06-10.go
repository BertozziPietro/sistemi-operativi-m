package main

import (
	"fmt"
	"math/rand"
	"time"
)

const MAXPROC = 100 // massimo numero clienti
const M = 2         // numero di tunnel
const N = 5         // numero aree per la pulizia interni (PI)
const MAX = 6       // max numero di tunnel e PI contemporaneamente usati

const TUNNEL = 0
const PI = 1

//stati addetto 0 OUT, 1 IN, 2 EXITING
const OUT = 0
const IN = 1
const EXITING = 2

const DIMBUFF = 100

var acquisizione [2]chan int
var rilascio [2]chan int
var inizio_presidio = make(chan int, DIMBUFF)
var termina_presidio = make(chan int, DIMBUFF)
var ACK_CLI [MAXPROC]chan int //canali ack clienti
var ACK_ADD [M]chan int       //canali ack addetti

var done = make(chan bool)
var termina = make(chan bool)
var stop = make(chan bool)

func when(b bool, c chan int) chan int {
	if !b {
		return nil
	}
	return c
}

func cliente(id int) {
	washed := false // indica se è stato eseguito il lavaggio esterno (evito che un cliente non chieda alcun servizio)
	time.Sleep(time.Duration(rand.Intn(2)) * time.Second)

	if rand.Intn(2) == TUNNEL {
		washed = true
		// Acquisizione di un tunnel per il lavaggio carrozzeria
		fmt.Printf("[CLIENTE %d] richiede un TUNNEL per il lavaggio della carrozzeria\n", id)
		acquisizione[TUNNEL] <- id
		tun := <-ACK_CLI[id]
		fmt.Printf("[CLIENTE %d] usa il TUNNEL n. %d\n", id, tun)
		// Lavaggio
		time.Sleep(time.Duration(1) * time.Second)

		// Rilascio tunnel
		rilascio[TUNNEL] <- tun
		<-ACK_CLI[id]
		fmt.Printf("[CLIENTE %d] rilascia il TUNNEL n. %d\n", id, tun)
	}
	if rand.Intn(2) == PI || washed == false { //pulizia interni
		// Acquisizione area di pulizia
		fmt.Printf("[CLIENTE %d] richiede un'AREA PULIZIA INTERNI\n", id)
		acquisizione[PI] <- id
		area := <-ACK_CLI[id]
		fmt.Printf("[CLIENTE %d] usa l'AREA PULIZIA INTERNI n. %d\n", id, area)
		// pulizia auto..
		time.Sleep(time.Duration(2) * time.Second)

		// Rilascio area di pulizia
		rilascio[PI] <- area
		<-ACK_CLI[id]
		fmt.Printf("[CLIENTE %d] ha rilasciato l'AREA PULIZIA INTERNI n. %d\n", id, area)
	}
	fmt.Printf("[CLIENTE %d] ho finito e me ne vado \n", id)
	// Terminazione
	done <- true
}

//	ADDETTO

func addetto(id int) {
	fmt.Printf("[ADDETTO %d] creato\n", id)

	for {
		// Inizia presidio del proprio tunnel
		inizio_presidio <- id
		fine := <-ACK_ADD[id]

		if fine < 0 {
			// Terminazione
			fmt.Printf("[ADDETTO %d] termino\n", id)
			done <- true
			return
		}

		// presidio tunnel
		time.Sleep(time.Duration(8) * time.Second)

		// Termina presidio del proprio tunnel
		fmt.Printf("[ADDETTO %d] vuole riposarsi\n", id)
		termina_presidio <- id
		<-ACK_ADD[id]

		// Riposo:
		time.Sleep(time.Duration(1) * time.Second)
	}
}

// AUTOLAVAGGIO

func autolavaggio() {
	var tunnel_utilizzati = 0
	var tunnel_disponibili = 0 // numero di tunnel con addetto presente e non impegnati da auto
	var aree_pulizia_usate = 0
	var stato_add [M]int  //0 OUT, 1 IN, 2 EXITING
	var auto_in_T [M]int  //AUTO_IN_T [i]=id cliente (se tunnel i occupato), altrimenti -1
	var auto_in_PI [N]int //AUTO_IN_PI [i]=id cliente (se PI i occupata), altrimenti -1
	var fine = false

	for i := 0; i < M; i++ {
		stato_add[i] = OUT // inizialmente non c'è l'addetto
		auto_in_T[i] = -1  // tunnel inizialmente libero
	}
	for i := 0; i < N; i++ {
		auto_in_PI[i] = -1
	}

	for { // calcolo i tunnel disponibili nell'iterazione corrente:
		tunnel_disponibili = 0
		for i := 0; i < M; i++ {
			if auto_in_T[i] == -1 && stato_add[i] == IN { //tunnel vuoto e addetto presente
				tunnel_disponibili++
			}
		}
		fmt.Printf("[AUTOLAVAGGIO] sono  disponibili  %d tunnel e %d aree PI\n", tunnel_disponibili, N-aree_pulizia_usate)

		select {
		// Acquisizione di un tunnel
		case id := <-when(tunnel_disponibili > 0 && tunnel_utilizzati < M && tunnel_utilizzati+aree_pulizia_usate < MAX && len(acquisizione[PI]) == 0, acquisizione[TUNNEL]):

			found := false
			assegnato := -1
			for i := 0; i < M && !found; i++ {
				if auto_in_T[i] == -1 && stato_add[i] == IN {
					found = true
					assegnato = i
				}
			}
			auto_in_T[assegnato] = id
			tunnel_utilizzati++
			fmt.Printf("[AUTOLAVAGGIO] Cliente %d acquisisce il tunnel %d\n", id, assegnato)
			ACK_CLI[id] <- assegnato

		// Acquisizione di un'area di pulizia interni
		case id := <-when(aree_pulizia_usate < N && tunnel_utilizzati+aree_pulizia_usate < MAX, acquisizione[PI]):
			assegnato := -1
			found := false
			i := 0
			for i = 0; i < N && !found; i++ {
				if auto_in_PI[i] == -1 {
					found = true
					assegnato = i
				}
			}
			auto_in_PI[assegnato] = id
			aree_pulizia_usate++
			ACK_CLI[id] <- assegnato
			fmt.Printf("[AUTOLAVAGGIO] Cliente %d acquisisce l'area per la pulizia interni n. %d\n", id, assegnato)

		// Rilascio di un tunnel (t)
		case t := <-rilascio[TUNNEL]:

			tunnel_utilizzati--
			id := auto_in_T[t]
			auto_in_T[t] = -1
			ACK_CLI[id] <- 1
			fmt.Printf("[AUTOLAVAGGIO] Cliente %d rilascia il tunnel %d\n", id, t)

			// Eventuale risveglio dell'addetto sospeso
			if stato_add[t] == EXITING {
				stato_add[t] = OUT
				ACK_ADD[t] <- 1
				fmt.Printf("[AUTOLAVAGGIO] L'addetto %d termina il presidio del tunnel %d\n", t, t)
			}

		// Rilascio di un'area di pulizia interni
		case area := <-rilascio[PI]:
			aree_pulizia_usate--
			id := auto_in_PI[area]
			auto_in_PI[area] = -1
			ACK_CLI[id] <- 1
			fmt.Printf("[AUTOLAVAGGIO] Cliente %d rilascia l'area di pulizia interni n. %d\n", id, area)

		// Addetto inizia il presidio del proprio tunnel
		case id := <-when(!fine, inizio_presidio):
			stato_add[id] = IN
			ACK_ADD[id] <- 1
			fmt.Printf("[AUTOLAVAGGIO] L'addetto %d inizia il presidio del tunnel %d\n", id, id)

		case id := <-when(fine, inizio_presidio):
			ACK_ADD[id] <- -1

		// Addetto termina il presidio del proprio tunnel
		case id := <-termina_presidio:
			if auto_in_T[id] == -1 { //tunnel non utilizzato
				stato_add[id] = OUT
				ACK_ADD[id] <- 1
				fmt.Printf("[AUTOLAVAGGIO] L'addetto %d termina il presidio del tunnel %d\n", id, id)
			} else { // Sospensione dell'addetto se il suo tunnel è ancora utilizzato da un cliente
				stato_add[id] = EXITING
				fmt.Printf("[AUTOLAVAGGIO] L'addetto %d si sospende in attesa della terminazione dell'utilizzo del tunnel %d\n", id, id)
			}

		case <-termina:
			fine = true
			fmt.Printf("Clienti terminati, Autolavaggio in chiusura..\n")

		// Terminazione di tutte le goroutine
		case <-stop:
			done <- true
			return

		}

	}

}

func main() {

	var n_clienti int

	// Inizializzazione canali
	for i := 0; i < 2; i++ {
		acquisizione[i] = make(chan int, DIMBUFF)
		rilascio[i] = make(chan int, DIMBUFF)
	}

	for i := 0; i < M; i++ {
		ACK_ADD[i] = make(chan int, DIMBUFF)
	}

	fmt.Printf("\nQuanti clienti (al massimo %d)? ", MAXPROC)
	fmt.Scanf("%d", &n_clienti)

	for i := 0; i < n_clienti; i++ {
		ACK_CLI[i] = make(chan int, DIMBUFF)
	}

	rand.Seed(time.Now().Unix())

	go autolavaggio()
	fmt.Printf("\n\n*** AUTOLAVAGGIO APERTO ***\n [ci sono %d tunnel e %d aree PI]\n\n", M, N)

	for i := 0; i < M; i++ {
		go addetto(i)
	}

	for i := 0; i < n_clienti; i++ {
		go cliente(i)
	}

	// Attesa terminazione clienti e addetti
	for i := 0; i < n_clienti; i++ {
		<-done
	}
	termina <- true // tutti i clienti sono terminati

	for i := 0; i < M; i++ {
		<-done
	}
	// tutti gli addetti sono terminati
	stop <- true //termino il server
	<-done
	fmt.Printf("*** CHIUSURA AUTOLAVAGGIO *** \n\n")
}
