package main

import (
	"fmt"
	"math/rand"
	"time"
	"strings"
)

type richiesta struct {
	soggetto int
	oggetto int
	ack chan int
}

const (
	NON_SIGNIFICATIVO = -1

	MAX      = 30
	MAX_BUFF = 100
	
	NUM_PI      = 30
	NUM_TUNNELS = 10

	NUM_ADDETTI = 10
	NUM_CLIENTI = 100
	
	TEMPO_MINIMO  = 500
	TEMPO_ADDETTO = 3000
	TEMPO_CLIENTE = 1000
	
	AZIONI_ADDETTO = 2
	AZIONI_CLIENTE = 2
	
	TIPI_CLIENTE = 2
)

var (
	canaliAddetto = [AZIONI_ADDETTO]chan richiesta{
					 make(chan richiesta, MAX_BUFF),
					 make(chan richiesta, MAX_BUFF)}
					 
	canaliCliente = [TIPI_CLIENTE][AZIONI_CLIENTE]chan richiesta{
					{make(chan richiesta, MAX_BUFF), make(chan richiesta, MAX_BUFF)},
					{make(chan richiesta, MAX_BUFF), make(chan richiesta, MAX_BUFF)}}

	finito              = make(chan bool, MAX_BUFF)
	bloccaAutolavaggio  = make(chan bool)
	terminaAutolavaggio = make(chan bool)
)

func lunghezzeCanaliInVettore(canali []chan richiesta) []int {
	lunghezze := make([]int, len(canali))
	for i, c := range canali { lunghezze[i] = len(c) }
	return lunghezze
}

func lunghezzeCanaliInMatrice(canali [][]chan richiesta) [][]int {
	lunghezze := make([][]int, len(canali))
	for i, riga := range canali {
		lunghezze[i] = make([]int, len(riga))
		for j, c := range riga { lunghezze[i][j] = len(c) }
	}
	return lunghezze
}

func when(b bool, c chan richiesta) chan richiesta {
	if !b { return nil }
	return c
}

func priorità(canali ...chan richiesta) bool {
	for _, c := range canali { if len(c) > 0 { return false } }
	return true
}

func autolavaggio() {
	const nome = "AUTOLAVAGGIO"
	var spazi = strings.Repeat(" ", len(nome) + 3)
	fmt.Printf("[%s] inizio\n", nome)

	var (
		clientiTunnel = 0
		clientiPI = 0
		fine = false
	)
	
	var (
		tunnels [NUM_TUNNELS]int
		posticipatiACK [NUM_ADDETTI]chan int
	)
    for i := range tunnels { tunnels[i] = -2 }
    
	disponibile := func() int {
		for i := range tunnels { if (tunnels[i] == -1) { return i } }
		return -1
	}
	
	libero := func() bool { return clientiTunnel + clientiPI < MAX }
	
	liberaPI := func() bool { return clientiPI < NUM_PI }

	for {
		var canaliClienteSlice = make([][]chan richiesta, len(canaliCliente))
		for i, row := range canaliCliente { canaliClienteSlice[i] = append([]chan richiesta(nil), row[:]...) }
		var (
			lunghezzeAddetto = lunghezzeCanaliInVettore(canaliAddetto[:])
			lunghezzeCliente = lunghezzeCanaliInMatrice(canaliClienteSlice)
		)
		
		//print del vettore
		fmt.Printf("[%s] clientiTunnel: %03d, clientiPI: %03d, fine: %5t\n%sCanaliAddetto: %v, CanaliCliente: %v\n",
		nome, clientiTunnel, clientiPI, fine, spazi, lunghezzeAddetto, lunghezzeCliente)
		
		select {
		case r := <-when(!fine,
		canaliAddetto[0]): {
			tunnels[r.soggetto] = -1
			r.ack <- 1
		}
		case r := <-when(fine,
		canaliAddetto[0]): {
			r.ack <- -1
		}
		case r := <-canaliAddetto[1]: {
			//attenzione a modificare una variabile prima di leggerla
			// o a mandare gli ack troppo presto
			if tunnels[r.soggetto] >= 0 { posticipatiACK[r.soggetto] = r.ack } else { r.ack <- 1 }
			tunnels[r.soggetto] = -2
		}
		case r := <-when(libero() && disponibile() >= 0 &&
		priorità(canaliAddetto[1], canaliCliente[1][0]),
		canaliCliente[0][0]): {
			var oggetto = disponibile()
			tunnels[oggetto] = r.soggetto
			clientiTunnel++
			r.ack <- oggetto
		}
		case r := <-canaliCliente[0][1]: {
			if posticipatiACK[r.oggetto] != nil {
				tunnels[r.oggetto] = -2
				posticipatiACK[r.oggetto] <- 1
				posticipatiACK[r.oggetto] = nil
			} else { tunnels[r.oggetto] = -1 }
			clientiTunnel--
			r.ack <- 1
		}
		case r := <-when(libero() && liberaPI(),
		canaliCliente[1][0]): {
			clientiPI++
			r.ack <- 1
		}
		case r := <-canaliCliente[1][1]: {
			clientiPI--
			r.ack <- 1
		}
		case <-bloccaAutolavaggio: {
			fine = true
		}
		case <-terminaAutolavaggio: {
			finito <- true
			fmt.Printf("[%s] fine\n", nome)
			return
		}}
	}
}

func addetto(id int) {
	const nome = "ADDETTO"
	fmt.Printf("[%s %03d] inizio\n", nome, id)

	var (
		r = richiesta { soggetto: id, oggetto: NON_SIGNIFICATIVO, ack: make(chan int) }
		azioni = [AZIONI_ADDETTO]string {"entrare in servizio", "riposare"}
		continua int
	)
	for {
		for i, azione := range azioni {
			time.Sleep(time.Duration(rand.Intn(TEMPO_ADDETTO) + TEMPO_MINIMO) * time.Millisecond)
			fmt.Printf("[%s %03d] mi metto in coda per %s\n", nome, id, azione)
			canaliAddetto[i] <- r
			continua = <-r.ack
			if continua < 0 {
				fmt.Printf("[%s %03d] fine\n", nome, id)
				finito <- true
				return
			}
			fmt.Printf("[%s %03d] è il mio turno di %s\n", nome, id, azione)
		}
	}
}

func cliente(id int, tipo int) {
	const nome = "CLIENTE"
	fmt.Printf("[%s %03d] inizio\n", nome, id)
	
	var (
		r = richiesta { soggetto: id, oggetto: NON_SIGNIFICATIVO, ack: make(chan int) }
		azioni = [TIPI_CLIENTE][AZIONI_CLIENTE]string {
				{"usare un tunnel", "liberare il tunnel"},
				{"usare un area PI", "liberare l'area PI"}}
		risorsa int
	)
	for t := 0; t < TIPI_CLIENTE; t++ {
		risorsa = NON_SIGNIFICATIVO
		if tipo == 2 || tipo == t {
			for i := 0; i < AZIONI_CLIENTE; i++ {
				time.Sleep(time.Duration(rand.Intn(TEMPO_CLIENTE) + TEMPO_MINIMO) * time.Millisecond)
				r.oggetto = risorsa
				fmt.Printf("[%s %03d] mi metto in coda per %s\n", nome, id, azioni[t][i])
				canaliCliente[t][i] <- r
				risorsa = <-r.ack
				if risorsa < 0 {
					fmt.Printf("[%s %03d] fine\n", nome, id)
					finito <- true
					return
				}
				fmt.Printf("[%s %03d] è il mio turno di %s numero %d\n", nome, id, azioni[t][i], risorsa)
			}
		}
	}
	
	fmt.Printf("[%s %03d] fine\n", nome, id)
	finito <- true
}

func main() {
	fmt.Println("[MAIN] inizio")
	rand.Seed(time.Now().Unix())
	
	go autolavaggio()
	for i := 0; i < NUM_ADDETTI; i++ { go addetto(i) }
	for i := 0; i < NUM_CLIENTI; i++ { go cliente(i, i % 3) }
	
	for i := 0; i < NUM_CLIENTI; i++ { <-finito }
	bloccaAutolavaggio <- true
	for i := 0; i < NUM_ADDETTI; i++ { <-finito }
	terminaAutolavaggio <- true
	<-finito
	
	fmt.Println("[MAIN] fine")
}
