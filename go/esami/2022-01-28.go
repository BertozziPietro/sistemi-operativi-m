package main

import (
	"fmt"
	"math/rand"
	"time"
	"strings"
)

const (
	NF       = 20
	NC       = 20
	LF		 = 4
	LC		 = 4
	LM 		 = 2
	MAX_BUFF = 50

	NUM_FORNITORI = 2
	NUM_ADDETTI   = 90
	
	TEMPO_MINIMO    = 500
	TEMPO_FORNITORE = 1000
	TEMPO_ADDETTO   = 1000
	
	AZIONI_FORNITORE = 2
	AZIONI_ADDETTO   = 2
	
	TIPI_FORNITORE = 2
	TIPI_ADDETTO   = 3
)

var (
	canaliFornitore = [TIPI_FORNITORE][AZIONI_FORNITORE]chan chan bool{
					{make(chan chan bool, MAX_BUFF), make(chan chan bool, MAX_BUFF)},
					{make(chan chan bool, MAX_BUFF), make(chan chan bool, MAX_BUFF)}}
	
	canaliAddetto = [TIPI_ADDETTO][AZIONI_ADDETTO]chan chan bool{
					{make(chan chan bool, MAX_BUFF), make(chan chan bool, MAX_BUFF)},
					{make(chan chan bool, MAX_BUFF), make(chan chan bool, MAX_BUFF)},
					{make(chan chan bool, MAX_BUFF), make(chan chan bool, MAX_BUFF)}}
					
	finito           = make(chan bool, MAX_BUFF)
	bloccaMagazzino  = make(chan bool)
	terminaMagazzino = make(chan bool)
)

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

func magazzino() {
	const nome = "MAGAZZINO"
	var spazi = strings.Repeat(" ", len(nome) + 3)
	fmt.Printf("[%s] inizio\n", nome)

	var (
		mascherineFFP2 = 0
		mascherineChirurgiche = 0
		addettiFFP2 = 0
		addettiChirurgiche = 0
		scaffaleFFP2 = 0
		scaffaleChirurgiche = 0
		fine = false
	)
	
	piùFFP2 := func() bool { return mascherineFFP2 > mascherineChirurgiche }
	
	rifornimentoFFP2 := func() bool { return mascherineFFP2 < NF / 2 }
	rifornimentoChirurgiche := func() bool { return mascherineChirurgiche < NC / 2 }
	
	prelievoFFP2 := func() bool { return mascherineFFP2 >= LF }
	prelievoChirurgiche := func() bool { return mascherineChirurgiche >= LC }
	prelievoMisto := func() bool { return mascherineFFP2 >= LM && mascherineChirurgiche >= LM }

	for {
		var canaliFornitoreSlice = make([][]chan chan bool, len(canaliFornitore))
		for i, row := range canaliFornitore { canaliFornitoreSlice[i] = append([]chan chan bool(nil), row[:]...) }
		var canaliAddettoSlice = make([][]chan chan bool, len(canaliAddetto))
		for i, row := range canaliAddetto { canaliAddettoSlice[i] = append([]chan chan bool(nil), row[:]...) }
		var (
			lunghezzeFornitore = lunghezzeCanaliInMatrice(canaliFornitoreSlice)
			lunghezzeAddetto   = lunghezzeCanaliInMatrice(canaliAddettoSlice)
		)
		fmt.Printf("[%s] mascherineFFP2: %03d, mascherineChirurgiche: %03d, addettiFFP2: %03d, addettiChirurgiche: %03d, scaffaleFFP2: %03d, scaffaleChirurgiche: %03d, fine: %5t\n%scanaliFornitore: %v, canaliAddetto: %v\n",
		nome, mascherineFFP2, mascherineChirurgiche, addettiFFP2, addettiChirurgiche, scaffaleFFP2, scaffaleChirurgiche, fine, spazi, lunghezzeFornitore, lunghezzeAddetto)
		
		select {
		case ack := <-when(!fine && !piùFFP2() && rifornimentoFFP2() && scaffaleFFP2 == 0,
		canaliFornitore[0][0]): {
			mascherineFFP2 = NF
			scaffaleFFP2 = 2
			ack <- true	 
		}
		case ack := <-when(fine,
		canaliFornitore[0][0]): {
			ack <- false
		}
		case ack := <-canaliFornitore[0][1]: {
			scaffaleFFP2 = 0
			ack <- true	 
		}
		case ack := <-when(!fine && piùFFP2() && rifornimentoChirurgiche() && scaffaleChirurgiche == 0,
		canaliFornitore[1][0]): {
			mascherineChirurgiche = NC
			scaffaleChirurgiche = 2
			ack <- true	 
		}
		case ack := <-when(fine,
		canaliFornitore[1][0]): {
			ack <- false
		}
		case ack := <-canaliFornitore[1][1]: {
			scaffaleChirurgiche = 0
			ack <- true	 
		}
		case ack := <-when(prelievoFFP2() && scaffaleFFP2 != 2 &&
		priorità(canaliAddetto[2][0]),
		canaliAddetto[0][0]): {
			mascherineFFP2 -= LF
			addettiFFP2++
			scaffaleFFP2 = 1
			ack <- true	 
		}
		case ack := <-canaliAddetto[0][1]: {
			addettiFFP2--
			if addettiFFP2 == 0 { scaffaleFFP2 = 0}
			ack <- true	 
		}
		case ack := <-when(prelievoChirurgiche() && scaffaleChirurgiche != 2 &&
		priorità(canaliAddetto[2][0], canaliAddetto[0][0]),
		canaliAddetto[1][0]): {
			mascherineChirurgiche -= LC
			addettiChirurgiche++
			scaffaleChirurgiche = 1
			ack <- true	 
		}
		case ack := <-canaliAddetto[1][1]: {
			addettiChirurgiche--
			if addettiChirurgiche == 0 { scaffaleChirurgiche = 0}
			ack <- true	 
		}
		case ack := <-when(prelievoMisto() && scaffaleFFP2 != 2 && scaffaleChirurgiche != 2,
		canaliAddetto[2][0]): {
			mascherineFFP2 -= LM
			mascherineChirurgiche -= LM
			addettiFFP2++
			addettiChirurgiche++
			scaffaleFFP2 = 1
			scaffaleChirurgiche = 1
			ack <- true	 
		}
		case ack := <-canaliAddetto[2][1]: {
			addettiFFP2--
			addettiChirurgiche--
			if addettiFFP2 == 0 { scaffaleFFP2 = 0}
			if addettiChirurgiche == 0 { scaffaleChirurgiche = 0}
			ack <- true	 
		}
		case <-bloccaMagazzino: {
			fine = true
		}
		case <-terminaMagazzino: {
			finito <- true
			fmt.Printf("[%s] fine\n", nome)
			return
		}}
	}
}

func fornitore(id int, tipo int) {
	const nome = "FORNITORE"
	fmt.Printf("[%s %03d] inizio\n", nome, id)
	
	var (
		ack = make(chan bool)
		azioni = [AZIONI_FORNITORE]string {"entrare nel magazzino", "uscire dal magazzino"}
		continua bool
	)
	for {
		for i, azione := range azioni {
			time.Sleep(time.Duration(rand.Intn(TEMPO_FORNITORE) + TEMPO_MINIMO) * time.Millisecond)
			fmt.Printf("[%s %03d] mi metto in coda per %s\n", nome, id, azione)
			canaliFornitore[tipo][i] <- ack
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

func addetto(id int, tipo int) {
	const nome = "ADDETTO"
	fmt.Printf("[%s %03d] inizio\n", nome, id)
	
	var (
		ack = make(chan bool)
		azioni = [AZIONI_ADDETTO]string {"entrare nel magazzino", "uscire dal magazzino"}
	)
	for i, azione := range azioni {
		time.Sleep(time.Duration(rand.Intn(TEMPO_ADDETTO) + TEMPO_MINIMO) * time.Millisecond)
		fmt.Printf("[%s %03d] mi metto in coda per %s\n", nome, id, azione)
		canaliAddetto[tipo][i] <- ack
		<-ack
		fmt.Printf("[%s %03d] è il mio turno di %s\n", nome, id, azione)
	}
	
	finito <- true
	fmt.Printf("[%s %03d] fine\n", nome, id)
}

func main() {
	fmt.Println("[MAIN] inizio")
	rand.Seed(time.Now().Unix())
	
	go magazzino()
	for i := 0; i < NUM_FORNITORI; i++ { go fornitore(i, i % 2) }
	for i := 0; i < NUM_ADDETTI; i++ { go addetto(i, i % 3) }
	
	for i := 0; i < NUM_ADDETTI; i++ { <-finito }
	bloccaMagazzino <- true
	for i := 0; i < NUM_FORNITORI; i++ { <-finito }
	terminaMagazzino <- true
	<-finito
	
	fmt.Println("[MAIN] fine")
}
