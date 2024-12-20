package main

//ben impostato ma qui c'è un grande problema ovvero le risorse condivise sono autonome e devi sempre usare indici ovunque
//esempio: come faccio a sapere se non posso entrare nel tunnel da cliente... dipende quale

import (
	"fmt"
	"math/rand"
	"time"
)

type req struct {
	id int
	ack chan int
}

type tunnel struct {
	addetto int
	cliente int
}

const (
	NUM_CLIENTI = 100
	NUM_ADDETTI = 10
	NUM_TUNNELS  = 10
	NUM_PI      = 30
	
	MAX      = 30
	MAX_BUFF = 100
	
	TEMPO_MIN     = 500
	TEMPO_INIZIO  = 500
	TEMPO_CLIENTE  = 1000
	TEMPO_ADDETTO = 3000
)

var (
	entraAddetto        = make(chan req, MAX_BUFF)
	esceAddetto         = make(chan req, MAX_BUFF)
	filtratoEsceAddetto = make(chan req, MAX_BUFF)
	terminaFiltro       = make(chan chan bool)

	entraCliente = [2]chan req{make(chan req, MAX_BUFF), make(chan req, MAX_BUFF)}
	esceCliente  = [2]chan req{make(chan req, MAX_BUFF), make(chan req, MAX_BUFF)}

	finito              = make(chan bool, MAX_BUFF)
	bloccaAutolavaggio  = make(chan bool)
	terminaAutolavaggio = make(chan bool)
)

func when(b bool, c chan req) chan req {
	if !b {
		return nil
	}
	return c
}

func autolavaggio() {
	fmt.Println("[AUTOLAVAGGIO] inizio")

	var (
		addetti = 0
		clientiTunnel = 0
		clientiPI = 0
		
		fine = false
	)
	
	var tunnels = [NUM_TUNNELS]tunnel{}
	for i := 0; i < NUM_TUNNELS; i++ {
		tunnels[i] = tunnel{addetto: -1, cliente: -1}
	}
	
	trova := func(id int) int {
		for i := 0; i < NUM_TUNNELS; i++ {
			if (tunnels[i].addetto != -1 && tunnels[i].cliente == id) {
				return i
			}
		}
		return -1
	}
	libero := func() bool {
 		return clientiTunnel + clientiPI < MAX
	}
	
	filtraEsceAddetto := func() {
		for {
		    select {
		    case r := <-esceAddetto: {
		        if tunnels[r.id].cliente == -1 {
		            filtratoEsceAddetto <- r
		        }
		    }
		    case ack:= <-terminaFiltro: {
		        ack <- true
		        return
		    }}
		}
	}
	go filtraEsceAddetto()

	for {
		fmt.Printf("[AUTOLAVAGGIO] Addetti: %1d, Clienti Tunnel: %03d, Clienti PI: %03d, EntraA: %1d, EsceA: %1d, EntraCT: %03d, EsceCT: %03d, EntraCPI: %03d, EsceCPI: %03d\n", addetti, clientiTunnel, clientiPI, len(entraAddetto), len(esceAddetto), len(entraCliente[0]), len(esceCliente[0]), len(entraCliente[1]), len(esceCliente[1]))
		select {
		case r := <-when(!fine,
		entraAddetto): {
			tunnels[r.id].addetto = r.id
			addetti++
			r.ack <- r.id
		}
		case r := <-when(fine,
		entraAddetto): {
			r.ack <- -1
		}
		case r := <-filtratoEsceAddetto: {
			tunnels[r.id].addetto = -1
			addetti--
			r.ack <- r.id
		}
		case r := <-when(len(esceAddetto) == 0 && len(entraCliente[1]) == 0 && libero(),
		entraCliente[0]): {
			var id = trova(-1)
			tunnels[id].cliente = r.id
			clientiTunnel++
			r.ack <- id
		}
		case r := <-esceCliente[0]: {
			var id = trova(r.id)
			tunnels[id].cliente = -1
			clientiTunnel--
			r.ack <- id
		}
		case r := <-when(libero(),
		entraCliente[1]): {
			clientiPI++
			r.ack <- r.id
		}
		case r := <-esceCliente[1]: {
			clientiPI--
			r.ack <- r.id
		}
		case <-bloccaAutolavaggio: {
			fine = true
		}
		case <-terminaAutolavaggio: {
			var ack = make(chan bool)
			terminaFiltro <- ack
			<-ack
			finito <- true
			fmt.Println("[AUTOLAVAGGIO] fine")
			return
		}}
	}
}

func addetto(id int) {
	fmt.Printf("[ADDETTO %03d] inizio\n", id)

	// preparazione
	r := req{id: id, ack: make(chan int)}
	
	// comportamento
	var continua int 
	for {
		time.Sleep(time.Duration(rand.Intn(TEMPO_INIZIO)+TEMPO_MIN) * time.Millisecond)
		fmt.Println("[ADDETTO %03d] mi metto in coda per usare il tunnel", id) 
		entraAddetto <- r
		continua = <-r.ack
		if continua > 0 {
			finito <- true
			fmt.Printf("[ADDETTO %03d] fine\n", id)
			return
		}
		fmt.Println("[ADDETTO %03d] è il mio turno di usare il tunnel", id)
		time.Sleep(time.Duration(rand.Intn(TEMPO_ADDETTO)+TEMPO_MIN) * time.Millisecond)
		fmt.Println("[ADDETTO %03d] mi metto in coda per uscire dal tunnel", id) 
		esceAddetto <- r
		<-r.ack
		fmt.Println("[ADDETTO %03d] è il mio turno di uscire dal tunnel", id)
	}
}

func cliente(id int, tipo int) {
	time.Sleep(time.Duration(rand.Intn(TEMPO_INIZIO)+TEMPO_MIN) * time.Millisecond)
	fmt.Printf("[CLIENTE %03d] inizio\n", id)
	
	// preparazione
	r := req{id: id, ack: make(chan int)}
	
	// comportamento
	var desideri = [2]string {"usare un tunnel", "usare un area PI"}
	var risorsa int
	for i := 0; i < 2; i++ {
		if tipo == 2 || tipo == i {
			fmt.Printf("[CLIENTE %03d] mi metto in coda per %s\n", id, desideri[i])
			entraCliente[i] <- r
			risorsa = <-r.ack
			fmt.Printf("[CLIENTE %03d] è il mio turno di %s numero %03d\n", id, desideri[i], risorsa)
			time.Sleep(time.Duration(rand.Intn(TEMPO_CLIENTE)+TEMPO_MIN) * time.Millisecond)
			fmt.Printf("[CLIENTE %03d] mi metto in coda per %s \n", id, desideri[i])
			esceCliente[i] <- r
			risorsa = <-r.ack
			fmt.Printf("[CLIENTE %03d] è il mio turno di %s numero %03d\n", id, desideri[i], risorsa)
		}
	}
	finito <- true
	fmt.Printf("[CLIENTE %03d] fine\n", id)
}

func main() {
	fmt.Println("[MAIN] inizio")
	rand.Seed(time.Now().Unix())
	
	// inizio goroutines
	go autolavaggio()
	for i := 0; i < NUM_ADDETTI; i++ {
		go addetto(i)
	}
	for i := 0; i < NUM_CLIENTI; i++ {
		go cliente(i, i % 3)
	}
	
	// fine goroutines
	for i := 0; i < NUM_CLIENTI; i++ {
		<-finito
	}
	bloccaAutolavaggio <- true
	for i := 0; i < NUM_ADDETTI; i++ {
		<-finito
	}
	terminaAutolavaggio <- true
	<-finito
	
	fmt.Println("[MAIN] fine")
}
