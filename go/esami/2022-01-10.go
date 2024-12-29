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

	MAXS     = 30
	MAX_BUFF = 100
	
	NUM_UFFICI = 10

	NUM_CONSULENTI     = 10
	NUM_AMMINISTRATORE = 50
	NUM_PROPRIETARIO   = 100
	
	TEMPO_MINIMO         = 500
	TEMPO_CONSULENTE     = 3000
	TEMPO_AMMINISTRATORE = 1000
	TEMPO_PROPRIETARIO   = 1000
	
	AZIONI_CONSULENTE = 2
	AZIONI_AMMINISTRATORE = 3
	AZIONI_PROPRIETARIO = 3
	
	TIPI_PROPRIETARIO = 2
)

var (
	canaliConsulente = [AZIONI_CONSULENTE]chan richiesta{
					 make(chan richiesta, MAX_BUFF),
					 make(chan richiesta, MAX_BUFF)}
					 
	canaliAmministratore = [AZIONI_AMMINISTRATORE]chan richiesta{
					 make(chan richiesta, MAX_BUFF),
					 make(chan richiesta, MAX_BUFF),
					 make(chan richiesta, MAX_BUFF)}
					
	canaliProprietario = [TIPI_PROPRIETARIO][AZIONI_PROPRIETARIO]chan richiesta{
					{make(chan richiesta, MAX_BUFF), make(chan richiesta, MAX_BUFF), make(chan richiesta, MAX_BUFF)},
					{make(chan richiesta, MAX_BUFF), make(chan richiesta, MAX_BUFF), make(chan richiesta, MAX_BUFF)}}

	finito         = make(chan bool, MAX_BUFF)
	bloccaFiliale  = make(chan bool)
	terminaFiliale = make(chan bool)
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

func filiale() {
	const nome = "FILIALE"
	var spazi = strings.Repeat(" ", len(nome) + 3)
	fmt.Printf("[%s] inizio\n", nome)

	var (
		soli = 0
		accompagnati = 0
		fine = false
	)
	
	var (
		uffici [NUM_UFFICI]int
		posticipatiACK [NUM_CONSULENTI]chan int
	)
	for i := range uffici { uffici[i] = -2 }
	
	disponibile := func() int {
		for i := range uffici { if (uffici[i] == -1) { return i } }
		return -1
	}
	
	libera := func(aggiunti int) bool { return soli + 2 * accompagnati + aggiunti <= MAXS }

	for {
		var canaliProprietarioSlice = make([][]chan richiesta, len(canaliProprietario))
		for i, row := range canaliProprietario { canaliProprietarioSlice[i] = append([]chan richiesta(nil), row[:]...) }
		var (
			lunghezzeConsulente     = lunghezzeCanaliInVettore(canaliConsulente[:])
			lunghezzeAmministratore = lunghezzeCanaliInVettore(canaliAmministratore[:])
			lunghezzeProprietario   = lunghezzeCanaliInMatrice(canaliProprietarioSlice)
		)
		
		fmt.Printf("[%s] soli: %03d, accompagnati: %03d, fine: %5t\n%sVettoreUffici: %v\n%sCanaliConsulente: %v, CanaliAmministratore: %v, CanaliProprietario: %v\n",
		nome, soli, accompagnati, fine, spazi, uffici, spazi, lunghezzeConsulente, lunghezzeAmministratore, lunghezzeProprietario)

	select {
		case r := <-when(!fine,
		canaliConsulente[0]): {
			uffici[r.soggetto] = -1
			r.ack <- 1
		}
		case r := <-when(fine,
		canaliConsulente[0]): {
			r.ack <- -1
		}
		case r := <-canaliConsulente[1]: {
			uffici[r.soggetto] = -2
			if uffici[r.soggetto] >= 0 { posticipatiACK[r.soggetto] = r.ack } else { r.ack <- 1 }
		}
		case r := <-when(libera(1),
		canaliAmministratore[0]): {
			soli++
			r.ack <- 1
		}
		case r := <-when(disponibile() >= 0,
		canaliAmministratore[1]): {
			var oggetto = disponibile()
			uffici[oggetto] = r.soggetto
			soli--
			r.ack <- oggetto
		}
		case r := <-canaliAmministratore[2]: {
			if posticipatiACK[r.oggetto] != nil {
				uffici[r.oggetto] = -2
				posticipatiACK[r.oggetto] <- 1
				posticipatiACK[r.oggetto] = nil
			} else { uffici[r.oggetto] = -1 }
			r.ack <- 1
		}
		case r := <-when(libera(1) &&
		priorità(canaliAmministratore[0]),
		canaliProprietario[0][0]): {
			soli++
			r.ack <- 1
		}
		case r := <-when(disponibile() >= 0 &&
		priorità(canaliAmministratore[1]),
		canaliProprietario[0][1]): {
			var oggetto = disponibile()
			uffici[oggetto] = r.soggetto
			soli--
			r.ack <- oggetto
		}
		case r := <-canaliProprietario[0][2]: {
			if posticipatiACK[r.oggetto] != nil {
				uffici[r.oggetto] = -2
				posticipatiACK[r.oggetto] <- 1
				posticipatiACK[r.oggetto] = nil
			} else { uffici[r.oggetto] = -1 }
			r.ack <- 1
		}
		case r := <-when(libera(2) &&
		priorità(canaliAmministratore[0], canaliProprietario[0][0]),
		canaliProprietario[1][0]): {
			soli += 2
			r.ack <- 1
		}
		case r := <-when(disponibile() >= 0 &&
		priorità(canaliAmministratore[1], canaliProprietario[0][1]),
		canaliProprietario[1][1]): {
			var oggetto = disponibile()
			uffici[oggetto] = r.soggetto
			soli -= 2
			r.ack <- oggetto
		}
		case r := <-canaliProprietario[1][2]: {
			if posticipatiACK[r.oggetto] != nil {
				uffici[r.oggetto] = -2
				posticipatiACK[r.oggetto] <- 1
				posticipatiACK[r.oggetto] = nil
			} else { uffici[r.oggetto] = -1 }
			r.ack <- 1
		}
		case <-bloccaFiliale: {
			fine = true
		}
		case <-terminaFiliale: {
			finito <- true
			fmt.Printf("[%s] fine\n", nome)
			return
		}}
	}
}

func consulente(id int) {
	const nome = "CONSULENTE"
	fmt.Printf("[%s %03d] inizio\n", nome, id)

	var (
		r = richiesta { soggetto: id, oggetto: NON_SIGNIFICATIVO, ack: make(chan int) }
		azioni = [AZIONI_CONSULENTE]string {"entrare in servizio", "riposare"}
		continua int
	)
	for {
		for i, azione := range azioni {
			time.Sleep(time.Duration(rand.Intn(TEMPO_CONSULENTE) + TEMPO_MINIMO) * time.Millisecond)
			fmt.Printf("[%s %03d] mi metto in coda per %s\n", nome, id, azione)
			canaliConsulente[i] <- r
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

func amministratore(id int) {
	const nome = "AMMINISTRATORE"
	fmt.Printf("[%s %03d] inizio\n", nome, id)
	
	var (
		r = richiesta { soggetto: id, oggetto: NON_SIGNIFICATIVO, ack: make(chan int) }
		azioni = [AZIONI_AMMINISTRATORE]string {"entrare nell'ufficio", "tornare a casa"}
	)
	for i := 0; i < AZIONI_AMMINISTRATORE; i++ {
		time.Sleep(time.Duration(rand.Intn(TEMPO_AMMINISTRATORE) + TEMPO_MINIMO) * time.Millisecond)
		fmt.Printf("[%s %03d] mi metto in coda per %s\n", nome, id, azioni[i])
		canaliAmministratore[i] <- r
		r.oggetto = <-r.ack
		if r.oggetto < 0 {
			finito <- true
			fmt.Printf("[%s %03d] fine\n", nome, id)
			return
		}
		fmt.Printf("[%s %03d] è il mio turno di %s numero %d\n", nome, id, azioni[i], r.oggetto)
	}
	
	finito <- true
	fmt.Printf("[%s %03d] fine\n", nome, id)
}

func proprietario(id int, tipo int) {
	const nome = "PROPRIETARIO"
	fmt.Printf("[%s %03d] inizio\n", nome, id)
	
	var (
		r = richiesta { soggetto: id, oggetto: NON_SIGNIFICATIVO, ack: make(chan int) }
		azioni = [AZIONI_PROPRIETARIO]string {"entrare nell'ufficio", "tornare a casa"}
	)
	for i := 0; i < AZIONI_PROPRIETARIO; i++ {
		time.Sleep(time.Duration(rand.Intn(TEMPO_PROPRIETARIO) + TEMPO_MINIMO) * time.Millisecond)
		fmt.Printf("[%s %03d] mi metto in coda per %s\n", nome, id, azioni[i])
		canaliProprietario[tipo][i] <- r
		r.oggetto = <-r.ack
		if r.oggetto < 0 {
			finito <- true
			fmt.Printf("[%s %03d] fine\n", nome, id)
			return
		}
		fmt.Printf("[%s %03d] è il mio turno di %s numero %d\n", nome, id, azioni[i], r.oggetto)
	}
	
	finito <- true
	fmt.Printf("[%s %03d] fine\n", nome, id)
}

func main() {
	fmt.Println("[MAIN] inizio")
	rand.Seed(time.Now().Unix())
	
	go filiale()
	for i := 0; i < NUM_CONSULENTI; i++ { go consulente(i) }
	for i := 0; i < NUM_AMMINISTRATORE; i++ { go amministratore(i) }
	for i := 0; i < NUM_PROPRIETARIO; i++ { go proprietario(i, i % 2) }
	
	for i := 0; i < NUM_AMMINISTRATORE + NUM_PROPRIETARIO; i++ { <-finito }
	bloccaFiliale <- true
	for i := 0; i < NUM_CONSULENTI; i++ { <-finito }
	terminaFiliale <- true
	<-finito
	
	fmt.Println("[MAIN] fine")
}
