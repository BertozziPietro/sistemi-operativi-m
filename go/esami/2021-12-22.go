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
	NM       = 10
	MAX_BUFF = 100
	
	NUM_SUPERVISIONAMENTI = 30

	NUM_COMMESSI = 10
	NUM_CLIENTI  = 100
	
	TEMPO_MINIMO    = 500
	TEMPO_COMMESSO  = 3000
	TEMPO_CLIENTE   = 1000
	TEMPO_FORNITORE = 1000
	
	AZIONI_COMMESSO = 2
	AZIONI_CLIENTE = 2
	AZIONI_FORNITORE = 1
	
	TIPI_CLIENTE = 2
)

var (
	canaliCommesso = [AZIONI_COMMESSO]chan richiesta{
					 make(chan richiesta, MAX_BUFF),
					 make(chan richiesta, MAX_BUFF)}
					 
	canaliCliente = [TIPI_CLIENTE][AZIONI_CLIENTE]chan richiesta{
					{make(chan richiesta, MAX_BUFF), make(chan richiesta, MAX_BUFF)},
					{make(chan richiesta, MAX_BUFF), make(chan richiesta, MAX_BUFF)}}
					
	canaleFornitore = make(chan chan bool, MAX_BUFF)

	finito               = make(chan bool, MAX_BUFF)
	bloccaAbbigliamento  = make(chan bool)
	terminaAbbigliamento = make(chan bool)
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

func whenBool(b bool, c chan chan bool) chan chan bool {
	if !b { return nil }
	return c
}

func priorità(canali ...chan richiesta) bool {
	for _, c := range canali { if len(c) > 0 { return false } }
	return true
}

func abbigliamento() {
	const nome = "ABBIGLIAMENTO"
	var spazi = strings.Repeat(" ", len(nome) + 3)
	fmt.Printf("[%s] inizio\n", nome)

	var (
		totali = 0
		mascherine = 0
		fine = false
	)
	
	var (
		supervisionamenti [NUM_SUPERVISIONAMENTI]int
		posticipatiACK [NUM_COMMESSI]chan int
	)
    for i := range supervisionamenti { supervisionamenti[i] = -3 }
    
    disponibile := func() int {
		for i := range supervisionamenti { if supervisionamenti[i] == -1 { return i } }
		return -1
	}
    
    libero := func() bool { return totali < MAX }
    
    mascherina := func() bool { return mascherine > 0 }
    
    for {
		var canaliClienteSlice = make([][]chan richiesta, len(canaliCliente))
		for i, row := range canaliCliente { canaliClienteSlice[i] = append([]chan richiesta(nil), row[:]...) }
		var (
			lunghezzeCommesso  = lunghezzeCanaliInVettore(canaliCommesso[:])
			lunghezzeCliente   = lunghezzeCanaliInMatrice(canaliClienteSlice)
			lunghezzaFornitore = len(canaleFornitore)
		)
		
		//print del vettore
		fmt.Printf("[%s] totali: %03d, mascherine: %03d, fine: %5t\n%sVettoreSupervisionamenti: %v, VettorePosticipatiACK: %v\n%sCanaliCommesso: %v, CanaliCliente: %v, CanaliFornitore: %1d\n",
		nome, totali, mascherine, fine, spazi, supervisionamenti, posticipatiACK, spazi, lunghezzeCommesso, lunghezzeCliente, lunghezzaFornitore)
		
		select {
		case r := <-when(!fine && libero(),
		canaliCommesso[0]): {
			totali++
			var first = r.soggetto * 3
			for i:= first; i < first + 3; i++ { supervisionamenti[i] = -1 }
			r.ack <- 1
		}
		case r := <-when(fine,
		canaliCommesso[0]): {
			r.ack <- -1
		}
		case r := <-canaliCommesso[1]: {
			var first = r.soggetto * 3
			var aspetta = false
			for i:= first; i < first + 3; i++ {
				if supervisionamenti[i] >= 0 {
					aspetta = true
					supervisionamenti[i] = -2
				} else { supervisionamenti[i] = -3 }
			} 
			if !aspetta {
				totali--
				r.ack <- 1
			} else { posticipatiACK[r.soggetto] = r.ack }
		}
		case r := <-when(libero() && mascherina() && disponibile() >= 0 &&
		priorità(canaliCommesso[0]),
		canaliCliente[0][0]): {
			mascherine--
			totali++
			var oggetto = disponibile()
			supervisionamenti[oggetto] = r.soggetto
			r.ack <- oggetto
		}
		case r := <-canaliCliente[0][1]: {
			var i = r.oggetto / 3
			if posticipatiACK[i] == nil {
				supervisionamenti[r.oggetto] = -1
				totali--
			} else {
				var aspetta = false
				supervisionamenti[r.oggetto] = -3
				var first = r.oggetto - (r.oggetto % 3)
				for i:= first; i < first + 3; i++ { if supervisionamenti[i] == -2 { aspetta = true } }
				if !aspetta {
					totali -= 2
					posticipatiACK[i] <- 1
					posticipatiACK[i] = nil
				} else { totali-- }
			} 
			r.ack <- 1
		}
		case r := <-when(libero() && mascherina() && disponibile() >= 0 &&
		priorità(canaliCommesso[0], canaliCliente[0][0]),
		canaliCliente[1][0]): {
			mascherine--
			totali++
			var oggetto = disponibile()
			supervisionamenti[oggetto] = r.soggetto
			r.ack <- oggetto
		}
		case r := <-canaliCliente[1][1]: {
			fmt.Println("ogg %d", r.oggetto)
			var i = r.oggetto / 3
			if posticipatiACK[i] == nil {
				supervisionamenti[r.oggetto] = -1
				totali--
			} else {
				var aspetta = false
				supervisionamenti[r.oggetto] = -3
				var first = r.oggetto - (r.oggetto % 3)
				for i:= first; i < first + 3; i++ { if supervisionamenti[i] == -2 { aspetta = true } }
				if !aspetta {
					totali -= 2
					posticipatiACK[i] <- 1
					posticipatiACK[i] = nil
				} else { totali-- }
			} 
			r.ack <- 1
		}
		case ack := <-whenBool(!fine,
		canaleFornitore): {
			mascherine += NM
			ack <- true
		}
		case ack := <-whenBool(fine,
		canaleFornitore): {
			ack <- false
		}
		case <-bloccaAbbigliamento: {
			fine = true
		}
		case <-terminaAbbigliamento: {
			finito <- true
			fmt.Printf("[%s] fine\n", nome)
			return
		}}
	}
}

func commesso(id int) {
	const nome = "COMMESSO"
	fmt.Printf("[%s %03d] inizio\n", nome, id)

	var (
		r = richiesta { soggetto: id, oggetto: NON_SIGNIFICATIVO, ack: make(chan int) }
		azioni = [AZIONI_COMMESSO]string {"entrare in servizio", "riposare"}
		continua int
	)
	for {
		for i, azione := range azioni {
			time.Sleep(time.Duration(rand.Intn(TEMPO_COMMESSO) + TEMPO_MINIMO) * time.Millisecond)
			fmt.Printf("[%s %03d] mi metto in coda per %s\n", nome, id, azione)
			canaliCommesso[i] <- r
			continua = <-r.ack
			if continua < 0 {
				finito <- true
				fmt.Printf("[%s %03d] fine\n", nome, id)
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
		azioni = [AZIONI_CLIENTE]string {"chiedere al commesso", "lasciare in pace il commesso"}
	)
	
	for i, azione := range azioni {
		time.Sleep(time.Duration(rand.Intn(TEMPO_CLIENTE) + TEMPO_MINIMO) * time.Millisecond)
		fmt.Printf("[%s %03d] mi metto in coda per %s\n", nome, id, azione)
		canaliCliente[tipo][i] <- r
		r.oggetto = <-r.ack
		if r.oggetto < 0 {
			finito <- true
			fmt.Printf("[%s %03d] fine\n", nome, id)
			return
		}
		fmt.Printf("[%s %03d] è il mio turno di %s\n", nome, id, azione)
	}
	
	finito <- true
	fmt.Printf("[%s %03d] fine\n", nome, id)
}

func fornitore(id int) {
	const nome = "FORNITORE"
	fmt.Printf("[%s %03d] inizio\n", nome, id)
	
	var (
		ack = make(chan bool)
		continua bool
	)
	for  {
		time.Sleep(time.Duration(rand.Intn(TEMPO_FORNITORE) + TEMPO_MINIMO) * time.Millisecond)
		fmt.Printf("[%s %03d] mi metto in coda per depositare le mascherine\n", nome, id)
		canaleFornitore <- ack
		continua = <-ack
		if !continua {
			finito <- true
			fmt.Printf("[%s %03d] fine\n", nome, id)
			return
		}
		fmt.Printf("[%s %03d] è il mio turno di depositare le mascherine\n", nome, id)
	}
	
	finito <- true
	fmt.Printf("[%s %03d] fine\n", nome, id)
}

func main() {
	fmt.Println("[MAIN] inizio")
	rand.Seed(time.Now().Unix())
	
	go abbigliamento()
	for i := 0; i < NUM_COMMESSI; i++ { go commesso(i) }
	for i := 0; i < NUM_CLIENTI; i++ { go cliente(i, i % 2) }
	go fornitore(0)
	
	for i := 0; i < NUM_CLIENTI; i++ { <-finito }
	bloccaAbbigliamento <- true
	for i := 0; i < NUM_COMMESSI + 1; i++ { <-finito }
	terminaAbbigliamento <- true
	<-finito
	
	fmt.Println("[MAIN] fine")
}
