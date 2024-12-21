//anche qui le risorse sono identificate e singole perchè le macchine possono stare nei posti dei camion quindi devi stare li a capire chi deve stare dove
//per ora la questione è ignorata e le macchine stanno solo nel posto delle macchine e i camion solo nel posto dei camion

package main

import (
	"fmt"
	"math/rand"
	"time"
	"strings"
)

const (
	NS       = 10
	NM       = 10
	MAX_BUFF = 30

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
	canaliSpazzaneve = [AZIONI_SPAZZANEVE]chan chan bool{
					 make(chan chan bool),
					 make(chan chan bool),
					 make(chan chan bool),
					 make(chan chan bool)}
	canaliCamper     = [AZIONI_CAMPER]chan chan bool{
					 make(chan chan bool, MAX_BUFF),
					 make(chan chan bool, MAX_BUFF),
					 make(chan chan bool, MAX_BUFF),
					 make(chan chan bool, MAX_BUFF)}
	canaliAuto       = [AZIONI_AUTO]chan chan bool{
					 make(chan chan bool, MAX_BUFF),
					 make(chan chan bool, MAX_BUFF),
					 make(chan chan bool, MAX_BUFF),
					 make(chan chan bool, MAX_BUFF)}

	finito         = make(chan bool, MAX_BUFF)
	bloccaStrada   = make(chan bool)
	terminaStrada  = make(chan bool)
)

func lunghezze(canali []chan chan bool) []int {
	lunghezze := make([]int, len(canali))
	for i, c := range canali { lunghezze[i] = len(c) }
	return lunghezze
}

func when(b bool, c chan chan bool) chan chan bool {
	if !b { return nil }
	return c
}

func strada() {
	const nome = "STRADA"
	var spazi = strings.Repeat(" ", len(nome)+3)
	fmt.Printf("[%s] inizio", nome)

	var (
		liberiM = NS
		liberiS = NM
		spazzaneve = false
		camper = 0
		auto = 0
		direzione = 0
		fine = false
	)
	
	vuota := func() bool {
 		return camper + auto == 0
	}
	
	priorità := func(canali ...chan chan bool) bool {
		for _, c := range canali { if len(c) > 0 { return false } }
		return true
	}

	for {
		var (
			lunghezzeSpazaneve = lunghezze(canaliSpazzaneve[:])
			lunghezzeCamper    = lunghezze(canaliCamper[:])
			lunghezzeAuto      = lunghezze(canaliAuto[:])
		)
		fmt.Printf("[%s] LiberiM: %03d, LiberiS: %03d, Spazzaneve: %5t, Camper: %03d, Auto: %03d, Direzione: %03d, Fine: %5t\n%sCanaliSpazzaneve: %v, CanaliCamper: %v, CanaliAuto: %v\n",
		nome, liberiM, liberiS, spazzaneve, camper, auto, direzione, fine, spazi, lunghezzeSpazaneve, lunghezzeCamper, lunghezzeAuto)
		
		select {
		case ack := <-when(!fine && vuota(),
		canaliSpazzaneve[0]): {
			spazzaneve = true
			ack <- true
		}
		case ack := <-when(fine,
		canaliSpazzaneve[0]): {
			ack <- false
		}
		case ack := <-canaliSpazzaneve[1]: {
			spazzaneve = false
			ack <- true
		}
		case ack := <-when(vuota() &&
		priorità(canaliSpazzaneve[0], canaliCamper[2], canaliAuto[2], canaliCamper[0], canaliAuto[0]),
		canaliSpazzaneve[2]): {
			spazzaneve = true
			ack <- true
		}
		case ack := <-canaliSpazzaneve[3]: {
			spazzaneve = false
			ack <- true
		}
		case ack := <-when(!spazzaneve && direzione != -1 && liberiM > 0 &&
		priorità(canaliSpazzaneve[0], canaliCamper[2], canaliAuto[2]),
		canaliCamper[0]): {
			liberiM--
			camper++
			ack <- true
		}
		case ack := <-canaliCamper[1]: {
			camper--
			ack <- true
		}
		case ack := <-when(!spazzaneve && direzione != 1 &&
		priorità(canaliSpazzaneve[0]),
		canaliCamper[2]): {
			liberiM++
			camper++
			ack <- true
		}
		case ack := <-canaliCamper[3]: {
			camper--
			ack <- true
		}
		case ack := <-when(!spazzaneve && direzione != -1 && liberiS > 0 &&
		priorità(canaliSpazzaneve[0], canaliCamper[2], canaliAuto[2], canaliCamper[0]),
		canaliAuto[0]): {
			liberiS--
			auto++
			ack <- true
		}
		case ack := <-canaliAuto[1]: {
			auto--
			ack <- true
		}
		case ack := <-when(!spazzaneve && direzione != 1 &&
		priorità(canaliSpazzaneve[0], canaliCamper[2]),
		canaliAuto[2]): {
			liberiS++
			auto++
			ack <- true
		}
		case ack := <-canaliAuto[3]: {
			auto--
			ack <- true
		}
		case <-bloccaStrada: {
			fine = true
		}
		case <-terminaStrada: {
			finito <- true
			fmt.Printf("[%s] fine", nome)
			return
		}}
	}
}

func spazzaneve(id int) {
	const nome = "SPAZZANEVE"
	fmt.Printf("[%s %03d] inizio\n", nome, id)

	var (
		ack = make(chan bool)
		azioni = [AZIONI_SPAZZANEVE]string {"scendere per la strada", "sostare al bar", "salire per la strada", "sostare al castello"}
		continua bool
	)
	for {
		for i, azione := range azioni {
			time.Sleep(time.Duration(rand.Intn(TEMPO_SPAZZANEVE)+TEMPO_MINIMO) * time.Millisecond)
			fmt.Printf("[%s %03d] mi metto in coda per %s\n", nome, id, azione)
			canaliSpazzaneve[i] <- ack
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

func camper(id int) {
	const nome = "CAMPER"
	fmt.Printf("[%s %03d] inizio\n", nome, id)
	
	var (
		ack = make(chan bool)
		azioni = [AZIONI_CAMPER]string{"salire per la strada", "sostare al castello", "scendere per la strada", "tornare a casa"}
	)
	for i, azione := range azioni {
		time.Sleep(time.Duration(rand.Intn(TEMPO_CAMPER)+TEMPO_MINIMO) * time.Millisecond)
		fmt.Printf("[%s %03d] mi metto in coda per %s\n", nome, id, azione)
		canaliCamper[i] <- ack
		<-ack
		fmt.Printf("[%s %03d] è il mio turno di %s\n", nome, id, azione)
	}
	
	finito <- true
	fmt.Printf("[%s %03d] fine\n", nome, id)
}

func auto(id int) {
	const nome = "AUTO"
	fmt.Printf("[%s %03d] inizio\n", nome, id)
	
	var (
		ack = make(chan bool)
		azioni = [AZIONI_AUTO]string {"salire per la strada", "sostare al castello", "scendere per la strada", "tornare a casa"}
	)
	for i, azione := range azioni {
		time.Sleep(time.Duration(rand.Intn(TEMPO_AUTO)+TEMPO_MINIMO) * time.Millisecond)
		fmt.Printf("[%s %03d] mi metto in coda per %s\n", nome, id, azione)
		canaliAuto[i] <- ack
		<-ack
		fmt.Printf("[%s %03d] è il mio turno di %s\n", nome, id, azione)
	}
	
	finito <- true
	fmt.Printf("[%s %03d] fine\n", nome, id)
}

func main() {
	fmt.Println("[MAIN] inizio")
	rand.Seed(time.Now().Unix())
	
	go strada()
	go spazzaneve(0)
	for i := 0; i < NUM_CAMPER; i++ {
		go camper(i)
	}
	for i := 0; i < NUM_AUTO; i++ {
		go auto(i)
	}
	
	for i := 0; i < NUM_CAMPER+NUM_AUTO; i++ {
		<-finito
	}
	bloccaStrada <- true
	<-finito
	terminaStrada <- true
	<-finito
	
	fmt.Println("[MAIN] fine")
}
