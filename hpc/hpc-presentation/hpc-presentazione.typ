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
    = Philosophische Argumente

    #definition[
        Ein Argument ist eine Ansammlung von Aussagesätzen, von denen behauptet wird, dass die einen (die #unbreak[*Annahmen*] bzw. #unbreak[*Prämissen*]) einen anderen (die #unbreak[*Konklusion*]) in der Art stützen würden, dass es rational wäre, anzunehmen, die Konklusion wäre wahr, wenn man annimmt, dass die Prämissen wahr sind.
    ]
]

#slide[
    #set text(size: 0.95em)
    #hint[
        Argumente stellen Begründungszusammenhänge dar. Die #unbreak[*Annahmen*] sollen eine gemeinsame Wissensgrundlage sein, auf die sich der Philosoph beruft, die #unbreak[*Konklusion*] das Ergebnis, was sich aus dieser Grundlage ergeben soll.
    ]

    #hint[
        Die *Annahmen* werden getrennt über den Strich geschrieben, die Konklusion darunter. Ein Beispiel folgt.
    ]
]

#slide[
    #example[
        1. Die Gesellschaft steht vor dem Untergang und muss gerettet werden.
        2. Die Gesellschaft kann nur gerettet werden, wenn die Populisten an die Macht kommen.
        #box(line(length: 100%))
        3. Die Populisten müssen an die Macht kommen.#footnote[Vgl. David Lanius: _Wie argumentieren Rechtspopulisten? Eine Argumentationsanalyse am Beispiel des AfD-Wahlprogramms._ https://davidlanius.de/de/wie-argumentieren-rechtspopulisten/ (05.03.2024, 08:30 Uhr).]
    ]
]

#slide[
    = Contesto e Premesse
    #set text(size: 0.9em)

    Il progetto è stato realizzato utilizzando #unbreak[*Galileo100 (G100)*], un supercomputer ad alte prestazioni situato presso il #unbreak[*Cineca*], progettato per gestire carichi computazionali intensivi e utile per progetti complessi che richiedono infrastrutture #unbreak[*HPC*].

    #definition[
      High-Performance Computing (HPC) è un paradigma computazionale che sfrutta tecniche di calcolo parallelo e architetture hardware avanzate per eseguire calcoli complessi in tempi ridotti, superando le capacità dei sistemi di calcolo tradizionali.
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
        Ho assunto, arbitrariamente, di considerare i valori mancanti come nulli.
    ]
]


#slide[
    = Analisi del Problema
    #set text(size: 0.8em)

    #hint[
        Si può dividere il calcolo di un qualsiasi $R_(i j)$ in due fasi:
            + la somma dei valori della relativa sottomatrice $A_(i, j)$
            + la divisione della somma per 9.0
    ]

    #hint[
        Conoscere il risultato della prima fase del calcolo di un valore della matrice risultato accelera la prima fase per i valori adiacenti, in tutte le direzioni. Questo perché 6 dei 9 valori sono condivisi tra le sottomatrici adiacenti.
    ]
]

#slide[
    = Analisi del Problema
    #set text(size: 0.8em)

    #hint[
        E' possibile utilizzare una sliding window quadrata; eseguendo per ogni prima fase di ogni nuovo valore, 3 somme e 3 sottrazioni, invece che 9 somme. Conoscendo un valore si può trovare un valore adiacente in maniera più efficiente, e il procedimento si può iterare, a patto di calcolare sempre un valore adiacente ad uno noto.
    ]

    #solution[
        Ogni soluzine ottimale del problema suddivide il carico tra i processi di modo che
            + sia quanto più bilanciato possibile
            + le celle assegnate ad ogni processo siano adiacienti tra loro
    ]
]


#slide[
    = OpenMP ed MPI
    #set text(size: 0.9em)

    Le due soluzioni proposte sono state realizzate utilizzando rispettivamente il modello di interazione OpenMP, ed il modello di interazione MPI.
    #definition[
      
    ]
    #definition[
      
    ]
]

#slide[
    = OpenMP: Progettazione
    #set text(size: 0.9em)
    
    #solution[
      
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
      [ciao/*#image("immagine.png", width: 100%)*/],
    )
]

#slide[
    = OpenMP: Scalabilità Strong
]

#slide[
    = MPI: Progettazione
    #set text(size: 0.9em)
    
    #solution[
      
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
    = Confronto e Conclusioni
    #set text(size: 0.9em)
    
    #conclusion[
      
    ]
]