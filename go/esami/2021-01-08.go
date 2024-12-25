package main

import (
	"fmt"
	"math/rand"
	"time"
	"strings"
)

const (
	MAXP     = 20
	MAXM     = 20
	TOTP     = 100
	TOTM     = 100
	NL       = 6
	Q        = 3
	MAX_BUFF = 50

	NUM_APPROVVIGIONATORI = 2
	NUM_CONSUMATORI       = 90
	
	TEMPO_MINIMO            = 500
	TEMPO_APPROVVIGIONATORE = 1000
	TEMPO_CONSUMATORE       = 1000
	TEMPO_FINE              = 2000
	
	AZIONI_APPROVVIGIONATORE = 2
	AZIONI_CONSUMATORE       = 2
	
	TIPI_APPROVVIGIONATORE = 2
	TIPI_CONSUMATORE       = 3
)

var (
	canaliApprovvigionatore = [TIPI_APPROVVIGIONATORE][AZIONI_APPROVVIGIONATORE]chan chan bool{
					         {make(chan chan bool, MAX_BUFF), make(chan chan bool, MAX_BUFF)},
					         {make(chan chan bool, MAX_BUFF), make(chan chan bool, MAX_BUFF)}}
	
	canaliConsumatore = [TIPI_CONSUMATORE][AZIONI_CONSUMATORE]chan chan bool{
					    {make(chan chan bool, MAX_BUFF), make(chan chan bool, MAX_BUFF)},
					    {make(chan chan bool, MAX_BUFF), make(chan chan bool, MAX_BUFF)},
					    {make(chan chan bool, MAX_BUFF), make(chan chan bool, MAX_BUFF)}}
					
	finito           = make(chan bool)
	bloccaDeposito   = make(chan bool)
	terminaDeposito  = make(chan bool)
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

func deposito() {
	const nome = "DEPOSITO"
	var spazi = strings.Repeat(" ", len(nome) + 3)
	fmt.Printf("[%s] inizio\n", nome)

	var (
		VP          = MAXP
		VM          = MAXM
		DP          = 0
		DM          = 0
		consumatori = 0
		dentro      = 0
		turno       = true
	)
	
	approvigionamentoPfizer := func() bool { return VP + NL == MAXP}
	approvigionamentoModerna := func() bool { return VM + NL == MAXM}
	consumo := func() bool { return VP + VM >= Q && !approvigionamentoPfizer() && !approvigionamentoModerna()}

	for {
		var canaliApprovvigionatoreSlice = make([][]chan chan bool, len(canaliApprovvigionatore))
		for i, row := range canaliApprovvigionatore { canaliApprovvigionatoreSlice[i] = append([]chan chan bool(nil), row[:]...) }
		var canaliConsumatoreSlice = make([][]chan chan bool, len(canaliConsumatore))
		for i, row := range canaliConsumatore { canaliConsumatoreSlice[i] = append([]chan chan bool(nil), row[:]...) }
		var (
			lunghezzeApprovvigionatore = lunghezzeCanaliInMatrice(canaliApprovvigionatoreSlice)
			lunghezzeConsumatore      = lunghezzeCanaliInMatrice(canaliConsumatoreSlice)
		)
		fmt.Printf("[%s] VP: %03d, VM: %03d, DP: %03d, DM: %03d, consumatori: %03d, dentro: %03d, turno: %5t\n%scanaliApprovvigionatore: %v, canaliConsumatore: %v\n",
		nome, VP, VM, DP, DM, consumatori, dentro, turno, spazi, lunghezzeApprovvigionatore, lunghezzeConsumatore)
			
		 if DP >= TOTP && DM >= TOTM {
        	time.Sleep(time.Duration(TEMPO_FINE) * time.Millisecond)
            for _, row := range canaliApprovvigionatore { for _, c := range row { if len(c) > 0 { <-c <- false  } } }
            for _, row := range canaliConsumatore { for _, c := range row { if len(c) > 0 { <-c <- false  } } }
            finito <- true
			fmt.Printf("[%s] fine\n", nome)
			return
        }

		select {
		case ack := <-when(approvigionamentoPfizer() && dentro == 0,
		canaliApprovvigionatore[0][0]): {
			VP = MAXP
			dentro = 2
			ack <- true	 
		}
		case ack := <-canaliApprovvigionatore[0][1]: {
			dentro = 0
			ack <- true	 
		}
		case ack := <-when(approvigionamentoModerna() && dentro == 0,
		canaliApprovvigionatore[1][0]): {
			VM = MAXM
			dentro = 2
			ack <- true	 
		}
		case ack := <-canaliApprovvigionatore[1][1]: {
			dentro = 0
			ack <- true	 
		}
		case ack := <-when(consumo() && dentro != 2,
		canaliConsumatore[0][0]): {
			if turno { 
				DP += Q
				VP -= Q
			} else {
				DM += Q
				VM -= Q
			}
			turno = !turno
			consumatori++
			dentro = 1
			ack <- true	 
		}
		case ack := <-canaliConsumatore[0][1]: {
			consumatori--
			if consumatori == 0 { dentro = 0}
			ack <- true	 
		}
		case ack := <-when(consumo() && dentro != 2 &&
		priorità(canaliConsumatore[0][0]),
		canaliConsumatore[1][0]): {
			if turno { 
				DP += Q
				VP -= Q
			} else {
				DM += Q
				VM -= Q
			}
			turno = !turno
			consumatori++
			dentro = 1
			ack <- true	 
		}
		case ack := <-canaliConsumatore[1][1]: {
			consumatori--
			if consumatori == 0 { dentro = 0}
			ack <- true	 
		}
		case ack := <-when(consumo() && dentro != 2 &&
		priorità(canaliConsumatore[0][0], canaliConsumatore[1][0]),
		canaliConsumatore[2][0]): {
			if turno { 
				DP += Q
				VP -= Q
			} else {
				DM += Q
				VM -= Q
			}
			turno = !turno
			consumatori++
			dentro = 1
			ack <- true	 
		}
		case ack := <-canaliConsumatore[2][1]: {
			consumatori--
			if consumatori == 0 { dentro = 0}
			ack <- true	 
		}}
	}
}

func approvvigionatore(id int, tipo int) {
	const nome = "APPROVVIGIONATORE"
	fmt.Printf("[%s %03d] inizio\n", nome, id)
	
	var (
		ack = make(chan bool)
		azioni = [AZIONI_APPROVVIGIONATORE]string {"entrare nel deposito", "uscire dal deposito"}
		continua bool
	)
	for {
		for i, azione := range azioni {
			time.Sleep(time.Duration(rand.Intn(TEMPO_APPROVVIGIONATORE) + TEMPO_MINIMO) * time.Millisecond)
			fmt.Printf("[%s %03d] mi metto in coda per %s\n", nome, id, azione)
			canaliApprovvigionatore[tipo][i] <- ack
			continua = <-ack
			if !continua {
				fmt.Printf("[%s %03d] fine\n", nome, id)
				return
			}
			fmt.Printf("[%s %03d] è il mio turno di %s\n", nome, id, azione)
		}
	}
}

func consumatore(id int, tipo int) {
	const nome = "CONSUMATORE"
	fmt.Printf("[%s %03d] inizio\n", nome, id)
	
	var (
		ack = make(chan bool)
		azioni = [AZIONI_CONSUMATORE]string {"entrare nel deposito", "uscire dal deposito"}
		continua bool
	)
	for i, azione := range azioni {
		time.Sleep(time.Duration(rand.Intn(TEMPO_CONSUMATORE) + TEMPO_MINIMO) * time.Millisecond)
		fmt.Printf("[%s %03d] mi metto in coda per %s\n", nome, id, azione)
		canaliConsumatore[tipo][i] <- ack
		continua = <-ack
			if !continua {
				fmt.Printf("[%s %03d] fine\n", nome, id)
				return
			}
		fmt.Printf("[%s %03d] è il mio turno di %s\n", nome, id, azione)
	}
	
	fmt.Printf("[%s %03d] fine\n", nome, id)
}

func main() {
	fmt.Println("[MAIN] inizio")
	rand.Seed(time.Now().Unix())
	
	go deposito()
	for i := 0; i < NUM_APPROVVIGIONATORI; i++ { go approvvigionatore(i, i % 2) }
	for i := 0; i < NUM_CONSUMATORI; i++ { go consumatore(i, i % 3) }
	
	<-finito
	
	fmt.Println("[MAIN] fine")
}
