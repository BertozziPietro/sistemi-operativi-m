#import "/src/library.typ": slides
#import slides: *

#set text(lang: "it")
#show: slides.with(
    series: [HPC],
    title: [Progetto Esame Sistemi Operativi M],
    author: "Autore Pietro Bertozzi",
    lecturer: "Docente Anna Ciampolini"
)

#slide[
    #show: align.with(center + horizon)
    #heading(outlined: false)[Nun eine kurze Demonstration, wo Sie die Medien finden.]
]


#slide[
    == Politische Argumente
    #set text(size: 0.95em)

    #example[
        Die Populisten müssen an die Macht kommen, denn die Gesellschaft steht vor dem Untergang und muss gerettet werden. Sie kann aber nur gerettet werden, wenn die Populisten an die Macht kommen.#footnote[Vgl. David Lanius: _Wie argumentieren Rechtspopulisten? Eine Argumentationsanalyse am Beispiel des AfD-Wahlprogramms._ https://davidlanius.de/de/wie-argumentieren-rechtspopulisten/ (05.03.2024, 08:30 Uhr).]
    ]

    *Idee:* Vielleicht bräuchten wir ein Werkzeug, um dieses Argument bewerten zu können?
]

#slide[
    = Contesto e Premesse
    #set text(size: 0.9em)

    Il progetto è stato realizzato utilizzando #unbreak[*Galileo100 (G100)*], un supercomputer ad alte prestazioni situato presso il #unbreak[*Cineca*], progettato per gestire carichi computazionali intensivi e utile per progetti complessi che richiedono infrastrutture #unbreak[*HPC*].

    #definition[
      #unbreak[*High-Performance Computing (HPC)*] è un paradigma computazionale che sfrutta tecniche di calcolo parallelo e architetture hardware avanzate per eseguire calcoli complessi in tempi ridotti, superando le capacità dei sistemi di calcolo tradizionali.
    ]
]

#slide[
    = Proposta di Progetto
    #set text(size: 0.8em)
    
    #task[
    Data una matrice quadrata $A[N][N]$, contenente $N^2$ valori reali appartenenti all'intervallo $[0,1]$, e considerando che $N≥2000$, l'obiettivo è realizzare un programma che calcoli una matrice $R[N/2][N/2]$ di valori reali.
    
    La matrice $R$ si ottiene come segue: $∀i∈[0,N/2−1]$ e $∀j∈[0,N/2−1]$
    
    calcoliamo $R_(i j)​ = M e d i a \_ I n t o r n o(A_(i, j)​)$
    
    La funzione $M e d i a \_ I n t o r n o $ è definita come la media aritmetica di tutti i valori contenuti nell'intorno di $A_(h k)$​, che è una sottomatrice $3×3$ di $A$ con $A_(h k)$​​ al centro.
  ]
]

#slide[
    = Analisi dei Requisiti
    #set text(size: 0.9em)
    
    #hint[
        Non è immediatamente ovvio come la funzione $M e d i a \_ I n t o r n o $ dovrebbe comportarsi quando non tutti i valori della sottomatrice $A_(i, j)​$ sono definiti.
    ]
    #notice[
        Ho assunto, arbitrariamente, di considerare i #unbreak[*valori mancanti come nulli*].
    ]
]


#slide[
    = Analisi del Problema
    #set text(size: 0.8em)

    #hint[
        Si può dividere il calcolo di un qualsiasi $R_(i j)$ in due fasi:
            + la #unbreak[*somma*] dei valori della relativa sottomatrice $A_(i, j)$
            + la #unbreak[*divisione*] della somma per 9.0.
    ]

    #hint[
        Conoscere il risultato della prima fase del calcolo di un valore della matrice risultato accelera la prima fase per i valori adiacenti, in tutte le direzioni. Questo perché 6 dei 9 #unbreak[*valori*] sono #unbreak[*condivisi*] tra le sottomatrici adiacenti.
    ]
]

#slide[
    = Analisi del Problema
    #set text(size: 0.8em)

    #hint[
        E' possibile utilizzare una #unbreak[*sliding window quadrata*]; eseguendo per ogni prima fase di ogni nuovo valore, 3 somme e 3 sottrazioni, invece che 9 somme. Conoscendo un valore si può trovare un valore adiacente in maniera più efficiente, e il procedimento si può iterare, a patto di calcolare sempre un valore adiacente ad uno noto.
    ]

    #solution[
        Ogni soluzine ottimale del problema suddivide il carico tra i processi di modo che
            + sia quanto più #unbreak[*bilanciato possibile*]
            + le celle assegnate ad ogni processo siano #unbreak[*adiacienti*] tra loro.
    ]
]


#slide[
    = MPI ed OpenMP
    #set text(size: 0.8em)

    Le due soluzioni proposte sono state realizzate utilizzando rispettivamente il modello di interazione MPI, ed il modello di interazione OpenMP.
    
    #definition[
        Le #unbreak[*librerie MPI*] rappresentano lo standard per il modello di interazione a scambio di messaggi su HPC, mentre le #unbreak[*librerie OpenMP*] rappresentano lo standard per il modello di interazione a memoria comune su HPC.
    ]
    
    #notice[
        Dalle prossime slide, "matrice" fa rifermiento alla matrice risultato $R[L = N/2][L = N/2]$.
    ]
]

#slide[
    = MPI: Progettazione 1
    #set text(size: 0.8em)
    
    #solution[
        La soluzione proposta è una naturale conseguenza delle osservazioni discusse durante l'analisi. I processi si divdono le #unbreak[*righe*] della matrice in parti, per quanto possibile, uguali (come mostrato in figura nella slide succesisva).
    ]
    #hint[
        Notiamo che questa soluzione è pensata per un #unbreak[*numero di processi minore del numero delle righe della matrice*]. Aumentare ulteriormente il numero di processi non avrebbe alcun effetto sulle prestazioni. 
    ]
]

#slide[
    = MPI: Progettazione 2
    #set text(size: 0.8em)

    #grid(
        columns: (64fr, 36fr),
        gutter: 20pt,
        [
            #definition[
                Definiamo lo #unbreak[*sbilanciamento (del carico)*] come il rapporto tra:
                - la massima differenza fra il numero di elementi delegati ad un processo (o thread)
                - il numero di elementi totali. 
            ]
        ],
        [
            #image("/src/MPI.png", width: 80%)
        ],
    )

      #notice[
          Lo sbilanciamento così definito dipende dal #unbreak[*numero di processi (o thread)*] utilizzati.
      ]
    
]

#slide[
    = MPI: Porgettazione 3
    #set text(size: 0.8em)
    
    #example[
        Se per esempio, il numero di processi è uno in meno rispetto al numero di righe della matrice, lo sbilanciamento è pari al numero di elementi un una riga (come mostrato in figura nella slide precedente).
    ]
    #conclusion[
        In termini generali, lo sbilanciamento minimo di questa soluzione è nullo se il numero di processi è divisore del numero delle righe, mentre vale $L / L^2 = 1 / L$ (dove L è la lunghezza di una riga della matrice) in caso contrario. #unbreak[*Inoltre lo sbilanciamento non può assumere nessun altro valore*].
    ]
]

#slide[
    = MPI: Implementazione
    #set text(size: 0.8em)
    
    #code[
      ```c int a = 0```
    ]
]

#slide[
    = MPI: Scalabilità Weak
]

#slide[
    = MPI: Scalabilità Strong
]


#slide[
    = OpenMP: Progettazione 1
    #set text(size: 0.8em)
    
    #solution[
        La soluzione proposta è meno intuitiva della precedente. Ad ogni thread è assegnata una #unbreak[*sottomatrice*] invece che una riga (come mostrato in figura nella slide successiva).
    ]
    #hint[
        Notiamo che questa soluzione è pensata per un #unbreak[*umero di thread minore del numero dei valori della matrice*], e non delle sue righe. E' quindi possibile sfruttare un numero elevatissimo di thread in parallelo.
    ]
]

#slide[
    = OpenMP: Progettazione 2
    #set text(size: 0.8em)

    #grid(
      columns: (1fr, 1fr, 1fr),
      gutter: 20pt,
      [#image("/src/OpenMPpochi.png", width: 100%)],
      [#image("/src/OpenMP.png", width: 100%)],
      [#image("/src/OpenMPtanti.png", width: 100%)]
    )  
]

#slide[
    = OpenMP: Progettazione 3
    #set text(size: 0.8em)
    
    #example[
        Se per esempi i thread sono 5 e la dimensione della matrice è 5x5 (come mostrato il figura nella slide precedente) allora lo sbilanciamento vale solo 3.
        Se invece, per esempio, i thread sono 3 allora lo sbilanciamento varrà sempre cirac $N/8$.
    ]
    #conclusion[
        Lo sbilanciamento minimo della soluzione è zero ma lo sbilanciamento massimo può essere significativamente più alto. #unbreak[*L'andamento dello sbilanciamento in funzione del numero dei thread è altalenante*].
    ]
]

#slide[
    = OpenMP: Implementazione
    #set text(size: 0.9em)
    
    #code[
      
    ]
]

#slide[
    = OpenMP: Scalabilità Weak
    #set text(size: 0.7em)

    #grid(
      columns: (1fr, 1fr),
      gutter: 20pt,
      [#table(
          columns: (1fr, 1fr, 1fr),
          inset: 8pt,
          align: center,
          table.header([*N*], [*Threads*], [*Time*]),
          "a", "a", "a",
          "b", "b", "b",
      )],
      [ciao],//#image("/src/graph.png", width: 100%)],
    )
]

#slide[
    = OpenMP: Scalabilità Strong
]

#slide[
    = Confronto e Conclusioni
    #set text(size: 0.9em)
    
    #conclusion[
      
    ]
]