package main

import (
	"fmt"
	"math/rand"
	"time"
	"strings"
)

const (
	MAXP     = 20
	MAXC     = 20
	TOT      = 20
	MAX_BUFF = 6

	NUM_NASTRI = 4
	NUM_ROBOT  = 2
	
	TEMPO_MINIMO = 500
	TEMPO_NASTRO = 1000
	TEMPO_ROBOT  = 1000
	TEMPO_FINE   = 2000
	
	AZIONI_NASTRO = 1
	AZIONI_ROBOT  = 2
	
	TIPI_NASTRO = 4
	TIPI_ROBOT  = 2
)

var (
	canaliNastro = [TIPI_NASTRO]chan chan bool{
				   make(chan chan bool, MAX_BUFF),
				   make(chan chan bool, MAX_BUFF),
				   make(chan chan bool, MAX_BUFF),
				   make(chan chan bool, MAX_BUFF)}
	
	canaliRobot = [TIPI_ROBOT][AZIONI_ROBOT]chan chan bool{
				  {make(chan chan bool, MAX_BUFF), make(chan chan bool, MAX_BUFF)},
				  {make(chan chan bool, MAX_BUFF), make(chan chan bool, MAX_BUFF)}}
					
	finito              = make(chan bool, MAX_BUFF)
	bloccaStabilimento  = make(chan bool)
	terminaStabilimento = make(chan bool)
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

func stabilimento() {
	const nome = "STABILIMENTO"
	var spazi = strings.Repeat(" ", len(nome) + 3)
	fmt.Printf("[%s] inizio\n", nome)

	var (
		CA = 0
		CB = 0
		PA = 0
		PB = 0
		Ap = 0
		Bp = 0
		A = 0
		B = 0
	)
	
	pienoC := func() bool { return CA + CB == MAXC }
	pienoP := func() bool { return PA + PB == MAXP }
	
	piùA := func() bool { return A > B }

	for {
		var canaliRobotSlice = make([][]chan chan bool, len(canaliRobot))
		for i, row := range canaliRobot { canaliRobotSlice[i] = append([]chan chan bool(nil), row[:]...) }
		var (
			lunghezzeNastro = lunghezzeCanaliInVettore(canaliNastro[:])
			lunghezzeRobot  = lunghezzeCanaliInMatrice(canaliRobotSlice)
		)
		fmt.Printf("[%s] CA: %03d, CB: %03d, PA: %03d, PB: %03d, Ap: %03d, Bp: %03d, A: %03d, B: %03d\n%sCanaliNastro: %v, CanaliRobot: %v,\n",
		nome, CA, CB, PA, PB, Ap, Bp, A, B, spazi, lunghezzeNastro, lunghezzeRobot)
		
		//inesatto sia in questo esercizio che nell altro con la fine strana si deve scegliere una canale per ogni personaggio e bloccare quello e non provarli tutti altrimenti si rischia di alllinearsi male
		//ma forse posi se è allineato su un altro canale si ferma sull altro canale perchè la select non è in esecuzione...
        if A + B >= TOT {
        	time.Sleep(time.Duration(TEMPO_FINE) * time.Millisecond)
            for _, c := range canaliNastro {  <-c <- false }
            for _, row := range canaliRobot { for _, c := range row { if len(c) > 0 { <-c <- false  } } }
            finito <- true
			fmt.Printf("[%s] fine\n", nome)
			return
        }

		select {
		case ack := <-when(!piùA() && !pienoC(),
		canaliNastro[0]): {
			CA++
			ack <- true
		}
		case ack := <-when(piùA() && !pienoC(),
		canaliNastro[1]): {
			CB++
			ack <- true
		}
		case ack := <-when(!piùA() && !pienoP(),
		canaliNastro[2]): {
			PA++
			ack <- true	 
		}
		case ack := <-when(piùA() && !pienoP(),
		canaliNastro[3]): {
			PB++
			ack <- true
		}
		case ack := <-when(!piùA() && CA > 0,
		canaliRobot[0][0]): {
			CA--
			ack <- true
		}
		case ack := <-when(!piùA() && PA > 0,
		canaliRobot[0][1]): {
			PA--
			Ap++
			if Ap % 4 == 0 { A++ }
			ack <- true
		}
		case ack := <-when(piùA() && CB > 0,
		canaliRobot[1][0]): {
			CB--
			ack <- true
		}
		case ack := <-when(piùA() && PB > 0,
		canaliRobot[1][1]): {
			PB--
			Bp++
			if Bp % 4 == 0 { B++ }
			ack <- true
		}}
	}
}

func nastro(id int, tipo int) {
	const nome = "NASTRO"
	fmt.Printf("[%s %03d] inizio\n", nome, id)
	
	var (
		ack = make(chan bool)
		azione = "depositare nello stabilimento"
		continua bool
	)
	for {
		time.Sleep(time.Duration(rand.Intn(TEMPO_NASTRO) + TEMPO_MINIMO) * time.Millisecond)
		fmt.Printf("[%s %03d] mi metto in coda per %s\n", nome, id, azione)
		canaliNastro[tipo] <- ack
		continua = <-ack
		if !continua {
			fmt.Printf("[%s %03d] fine\n", nome, id)
			return
		}
		fmt.Printf("[%s %03d] è il mio turno di %s\n", nome, id, azione)
	}
}

func robot(id int, tipo int) {
	const nome = "ROBOT"
	fmt.Printf("[%s %03d] inizio\n", nome, id)
	
	var (
		ack = make(chan bool)
		azioni = [AZIONI_ROBOT]string {"prelevare il cerchio", "montare il cerchio"}
		ruote = 4
		continua bool
	)
	for {
		for r := 0; r < ruote; r++ {
			for i, azione := range azioni {
				time.Sleep(time.Duration(rand.Intn(TEMPO_ROBOT) + TEMPO_MINIMO) * time.Millisecond)
				fmt.Printf("[%s %03d] mi metto in coda per %s %d\n", nome, id, azione, r)
				canaliRobot[tipo][i] <- ack
				continua = <-ack
				if !continua {
					fmt.Printf("[%s %03d] fine\n", nome, id)
					return
				}
				fmt.Printf("[%s %03d] è il mio turno di %s %d\n", nome, id, azione, r)
			}
		}
	}
}

func main() {
	fmt.Println("[MAIN] inizio")
	rand.Seed(time.Now().Unix())
	
	go stabilimento()
	for i := 0; i < NUM_NASTRI; i++ { go nastro(i, i) }
	for i := 0; i < NUM_ROBOT; i++ { go robot(i, i) }
	
	<-finito
	
	fmt.Println("[MAIN] fine")
}
