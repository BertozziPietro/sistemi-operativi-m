package main

import (
	"fmt"
	"math/rand"
	"time"
	"strings"
)

const (
	M1       = 10
	M2       = 20
	N 		 = 30
	MAX_BUFF = 50

	NUM_CITTADINI = 90
	
	TEMPO_MINIMO    = 200
	TEMPO_ADDETTO   = 1000
	TEMPO_CITTADINO = 500
	
	AZIONI_ADDETTO   = 2
	AZIONI_CITTADINO = 2
	
	TIPI_CITTADINO = 2
)

var (
	canaliAddetto   = [AZIONI_ADDETTO]chan chan bool{
					make(chan chan bool),
					make(chan chan bool)}		
	canaliCittadino = [TIPI_CITTADINO][AZIONI_CITTADINO]chan chan bool{
					{make(chan chan bool, MAX_BUFF), make(chan chan bool, MAX_BUFF)},
					{make(chan chan bool, MAX_BUFF), make(chan chan bool, MAX_BUFF)}}
					
	finito           = make(chan bool, MAX_BUFF)
	bloccaCasaAcqua  = make(chan bool)
	terminaCasaAcqua = make(chan bool)
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

func casaAcqua() {
	const nome = "CASA_ACQUA"
	var spazi = strings.Repeat(" ", len(nome)+3)
	fmt.Printf("[%s] inizio\n", nome)

	var (
		monetine = 0
		monetone = 0
		litri = N
		dentro = 0
		fine = false
	)
	
	intervento := func() bool {
 		return monetine == M1 || monetone == M2 || litri == 0
	}
	
	priorità := func(canali ...chan chan bool) bool {
		for _, c := range canali { if len(c) > 0 { return false } }
		return true
	}

	for {
		var canaliCittadinoSlice = make([][]chan chan bool, len(canaliCittadino))
		for i, row := range canaliCittadino { canaliCittadinoSlice[i] = append([]chan chan bool(nil), row[:]...) }
		var (
			lunghezzeAddetto   = lunghezzeCanaliInVettore(canaliAddetto[:])
			lunghezzeCittadino = lunghezzeCanaliInMatrice(canaliCittadinoSlice)
		)
		fmt.Printf("[%s] Monetine: %03d, Monetone: %03d, Litri: %03d, Dentro: %03d, Fine: %5t\n%sCanaliAddetto: %v, CanaliCittadino: %v\n",
		nome, monetine, monetone, litri, dentro, fine, spazi, lunghezzeAddetto, lunghezzeCittadino)
		
		select {
		case ack := <-when(!fine && dentro == 0 && intervento(),
		canaliAddetto[0]): {
			monetine = 0
			monetone = 0
			litri = N
			dentro = 2
			ack <- true
		}
		case ack := <-when(fine,
		canaliAddetto[0]): {
			ack <- false
		}
		case ack := <-canaliAddetto[1]: {
			dentro = 0
			ack <- true
		}
		case ack := <-when(dentro == 0 && !intervento(),
		canaliCittadino[0][0]): {
			monetine++
			litri--
			dentro = 1
			ack <- true
		}
		case ack := <-canaliCittadino[0][1]: {
			dentro = 0
			ack <- true
		}
		case ack := <-when(dentro == 0 && !intervento() &&
		priorità(canaliCittadino[0][0]),
		canaliCittadino[1][0]): {
			monetone++
			litri -= 2
			dentro = 1
			ack <- true
		}
		case ack := <-canaliCittadino[1][1]: {
			dentro = 0
			ack <- true
		}
		case <-bloccaCasaAcqua: {
			fine = true
		}
		case <-terminaCasaAcqua: {
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
		ack = make(chan bool)
		azioni = [AZIONI_ADDETTO]string {"entrare nella casa dell'acqua", "uscire dalla casa dell'acqua"}
		continua bool
	)
	for {
		for i, azione := range azioni {
			time.Sleep(time.Duration(rand.Intn(TEMPO_ADDETTO)+TEMPO_MINIMO) * time.Millisecond)
			fmt.Printf("[%s %03d] mi metto in coda per %s\n", nome, id, azione)
			canaliAddetto[i] <- ack
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

func cittadino(id int, tipo int) {
	const nome = "CITTADINO"
	fmt.Printf("[%s %03d] inizio\n", nome, id)
	
	var (
		ack = make(chan bool)
		azioni = [AZIONI_ADDETTO]string {"entrare nella casa dell'acqua", "uscire dalla casa dell'acqua"}
	)
	for i, azione := range azioni {
		time.Sleep(time.Duration(rand.Intn(TEMPO_CITTADINO)+TEMPO_MINIMO) * time.Millisecond)
		fmt.Printf("[%s %03d] mi metto in coda per %s\n", nome, id, azione)
		canaliCittadino[tipo][i] <- ack
		<-ack
		fmt.Printf("[%s %03d] è il mio turno di %s\n", nome, id, azione)
	}
	
	finito <- true
	fmt.Printf("[%s %03d] fine\n", nome, id)
}

func main() {
	fmt.Println("[MAIN] inizio")
	rand.Seed(time.Now().Unix())
	
	ruotaTipo := func(i int) int {
		if i%3 == 0 { return 1 } else { return 0 }
	}
	
	go casaAcqua()
	go addetto(0)
	for i := 0; i < NUM_CITTADINI; i++ {
		go cittadino(i, ruotaTipo(i))
	}
	
	for i := 0; i < NUM_CITTADINI; i++ {
		<-finito
	}
	bloccaCasaAcqua <- true
	<-finito
	terminaCasaAcqua <- true
	<-finito
	
	fmt.Println("[MAIN] fine")
}
