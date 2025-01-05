package main

import (
	"fmt"
	"math/rand"
	"time"
	"strings"
)

type richiesta struct {
	oggetto int
	ack chan int
}

const (
	NON_SIGNIFICATIVO = -1

	NS       = 10
	NM       = 10
	MAX_BUFF = 100

	NUM_CAMPER = 50
	NUM_AUTO   = 50
	
	TEMPO_MINIMO     = 500
	TEMPO_SPAZZANEVE = 1000
	TEMPO_CAMPER     = 1000
	TEMPO_AUTO       = 1000
	
	AZIONI_SPAZZANEVE = 4
	AZIONI_CAMPER     = 4
	AZIONI_AUTO       = 4
)

var (
	canaliSpazzaneve = [AZIONI_SPAZZANEVE]chan richiesta{
					 make(chan richiesta, MAX_BUFF),
					 make(chan richiesta, MAX_BUFF),
					 make(chan richiesta, MAX_BUFF),
					 make(chan richiesta, MAX_BUFF)}
	canaliCamper     = [AZIONI_CAMPER]chan richiesta{
					 make(chan richiesta, MAX_BUFF),
					 make(chan richiesta, MAX_BUFF),
					 make(chan richiesta, MAX_BUFF),
					 make(chan richiesta, MAX_BUFF)}
	canaliAuto       = [AZIONI_AUTO]chan richiesta{
					 make(chan richiesta, MAX_BUFF),
					 make(chan richiesta, MAX_BUFF),
					 make(chan richiesta, MAX_BUFF),
					 make(chan richiesta, MAX_BUFF)}

	finito          = make(chan bool, MAX_BUFF)
	bloccaCastello  = make(chan bool)
	terminaCastello = make(chan bool)
)

func lunghezze(canali []chan richiesta) []int {
	lunghezze := make([]int, len(canali))
	for i, c := range canali { lunghezze[i] = len(c) }
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

func castello() {
	const nome = "CASTELLO"
	var spazi = strings.Repeat(" ", len(nome) + 3)
	fmt.Printf("[%s] inizio\n", nome)

	var (
		spazzaneve = false
		camper = 0
		auto = 0
		direzione = 0
		fine = false
	)
	var parcheggiLiberi[2]int
	parcheggiLiberi[0] = NS
	parcheggiLiberi[1] = NM
	
	parcheggioAuto := func() bool { return parcheggiLiberi[0] + parcheggiLiberi[1] > 0 }
	parcheggioCamper := func() bool { return parcheggiLiberi[1] > 0 }
	
	vuota := func() bool { return camper + auto == 0 }

	for {
		var (
			lunghezzeSpazaneve = lunghezze(canaliSpazzaneve[:])
			lunghezzeCamper    = lunghezze(canaliCamper[:])
			lunghezzeAuto      = lunghezze(canaliAuto[:])
		)
		fmt.Printf("[%s] Spazzaneve: %5t, Camper: %03d, Auto: %03d, Direzione: %03d, Fine: %5t\n%sVettoreParcheggiLiberi: %v\n%sCanaliSpazzaneve: %v, CanaliCamper: %v, CanaliAuto: %v\n",
		nome, spazzaneve, camper, auto, direzione, fine, spazi, parcheggiLiberi, spazi, lunghezzeSpazaneve, lunghezzeCamper, lunghezzeAuto)
		
		select {
		case r := <-when(!fine && vuota(),
		canaliSpazzaneve[0]): {
			spazzaneve = true
			r.ack <- 1
		}
		case r := <-when(fine,
		canaliSpazzaneve[0]): {
			r.ack <- -1
		}
		case r := <-canaliSpazzaneve[1]: {
			spazzaneve = false
			r.ack <- 1
		}
		case r := <-when(vuota() &&
		priorità(canaliSpazzaneve[0], canaliCamper[2], canaliAuto[2], canaliCamper[0], canaliAuto[0]),
		canaliSpazzaneve[2]): {
			spazzaneve = true
			r.ack <- 1
		}
		case r := <-canaliSpazzaneve[3]: {
			spazzaneve = false
			r.ack <- 1
		}
		case r := <-when(!spazzaneve && direzione != -1 && parcheggioCamper() &&
		priorità(canaliSpazzaneve[0], canaliCamper[2], canaliAuto[2]),
		canaliCamper[0]): {
			parcheggiLiberi[1]--
			camper++
			r.ack <- 1
		}
		case r := <-canaliCamper[1]: {
			camper--
			r.ack <- 1
		}
		case r := <-when(!spazzaneve && direzione != 1 &&
		priorità(canaliSpazzaneve[0]),
		canaliCamper[2]): {
			parcheggiLiberi[1]++
			camper++
			r.ack <- 1
		}
		case r := <-canaliCamper[3]: {
			camper--
			r.ack <- 1
		}
		case r := <-when(!spazzaneve && direzione != -1 && parcheggioAuto() &&
		priorità(canaliSpazzaneve[0], canaliCamper[2], canaliAuto[2], canaliCamper[0]),
		canaliAuto[0]): {
			auto++
			if parcheggiLiberi[0] > 0 {
				parcheggiLiberi[0]--
				r.ack <- 0
			} else {
				parcheggiLiberi[1]--
				r.ack <- 1
			}
		}
		case r := <-canaliAuto[1]: {
			auto--
			r.ack <- 1
		}
		case r := <-when(!spazzaneve && direzione != 1 &&
		priorità(canaliSpazzaneve[0], canaliCamper[2]),
		canaliAuto[2]): {
			auto++
			parcheggiLiberi[r.oggetto]++
			r.ack <- 1
		}
		case r := <-canaliAuto[3]: {
			auto--
			r.ack <- 1
		}
		case <-bloccaCastello: {
			fine = true
		}
		case <-terminaCastello: {
			finito <- true
			fmt.Printf("[%s] fine\n", nome)
			return
		}}
	}
}

func spazzaneve(id int) {
	const nome = "SPAZZANEVE"
	fmt.Printf("[%s %03d] inizio\n", nome, id)

	var (
		r = richiesta { oggetto: NON_SIGNIFICATIVO, ack: make(chan int) }
		azioni = [AZIONI_SPAZZANEVE]string {"scendere per la strada", "sostare al bar", "salire per la strada", "sostare al castello"}
		continua int
	)
	for {
		for i, azione := range azioni {
			time.Sleep(time.Duration(rand.Intn(TEMPO_SPAZZANEVE) + TEMPO_MINIMO) * time.Millisecond)
			fmt.Printf("[%s %03d] mi metto in coda per %s\n", nome, id, azione)
			canaliSpazzaneve[i] <- r
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

func camper(id int) {
	const nome = "CAMPER"
	fmt.Printf("[%s %03d] inizio\n", nome, id)
	
	var (
		r = richiesta { oggetto: NON_SIGNIFICATIVO, ack: make(chan int) }
		azioni = [AZIONI_CAMPER]string{"salire per la strada", "sostare al castello", "scendere per la strada", "tornare a casa"}
	)
	for i, azione := range azioni {
		time.Sleep(time.Duration(rand.Intn(TEMPO_CAMPER) + TEMPO_MINIMO) * time.Millisecond)
		fmt.Printf("[%s %03d] mi metto in coda per %s\n", nome, id, azione)
		canaliCamper[i] <- r
		<-r.ack
		fmt.Printf("[%s %03d] è il mio turno di %s\n", nome, id, azione)
	}
	
	finito <- true
	fmt.Printf("[%s %03d] fine\n", nome, id)
}

func auto(id int) {
	const nome = "AUTO"
	fmt.Printf("[%s %03d] inizio\n", nome, id)
	
	var (
		r = richiesta { oggetto: NON_SIGNIFICATIVO, ack: make(chan int) }
		azioni = [AZIONI_AUTO]string {"salire per la strada", "sostare al castello", "scendere per la strada", "tornare a casa"}
	)
	for i, azione := range azioni {
		time.Sleep(time.Duration(rand.Intn(TEMPO_AUTO) + TEMPO_MINIMO) * time.Millisecond)
		fmt.Printf("[%s %03d] mi metto in coda per %s\n", nome, id, azione)
		canaliAuto[i] <- r
		if i == 0 { r.oggetto = <-r.ack } else { <-r.ack }
		fmt.Printf("[%s %03d] è il mio turno di %s\n", nome, id, azione)
	}
	
	fmt.Printf("[%s %03d] fine\n", nome, id)
	finito <- true
}

func main() {
	fmt.Println("[MAIN] inizio")
	rand.Seed(time.Now().Unix())
	
	go castello()
	go spazzaneve(0)
	for i := 0; i < NUM_CAMPER; i++ { go camper(i) }
	for i := 0; i < NUM_AUTO; i++ { go auto(i) }
	
	for i := 0; i < NUM_CAMPER + NUM_AUTO; i++ { <-finito }
	bloccaCastello <- true
	<-finito
	terminaCastello <- true
	<-finito
	
	fmt.Println("[MAIN] fine")
}
