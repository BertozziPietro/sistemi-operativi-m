#import "/src/library.typ": slides
#import slides: *

#set text(lang: "it")
#show: slides.with(
    series: [HPC],
    title: [Progetto Esame Sistemi Operativi M],

    box-task-title: [Obbiettivo],
    box-hint-title: [Osservazione],
    box-solution-title: [Soluzione],
    box-definition-title: [Definizione],
    box-notice-title: [Attenzione],
    box-example-title: [Esempio],

    //outline-title-text: "Ablauf",
    // show-todolist: false,

    author: "Autore Pietro Bertozzi",
    lecturer: "Docente Anna Ciampolini"
)


// #hide-todos()

#slide[
    = Sie/Ihr und Ich

    + Siezen vs. Duzen?

    + Namenskärtchen

    + Stellen Sie sich bitte kurz vor, beantworten Sie Folgendes:
        - Wie heißen Sie?
        - Welche Fächerkombination studieren Sie?
        - Wenn Sie in einem fiktiven Universum leben müssten, in welchem?
]

#slide[
    == Ihre Erwartungen an mich

    #task[
        Notieren Sie als Wort oder kurzen Stichpunkt vorne am Whiteboard, was Sie von einem erfolgreichen Logik-Tutorium und einem guten Logik-Tutor erwarten!
    ]
]

#slide[
    == Meine Erwartungen an Sie
    - aktive Mitarbeit und Ergreifen von Initiative
    - Nachbereiten der Sitzungen mit Hinsicht auf Ihre Probleme mit dem Stoff
    - Offenheit gegenüber Problemen mit ...
        - ... der Organisation des Tutoriums
        - ... dem Umfang des Tutoriums
        - ... den Themen des Tutoriums
        - ... der Lernatmosphäre im Tutorium
        - ... mir als Tutor
        - ... Ihren fachlichen Schwächen und Lücken
]

#slide[
    == Hinweise für gutes Gelingen
    #set text(size: 0.95em)
    + Benutzen Sie mind. 1-2h pro Woche um modulspezifische Aufgaben zu bearbeiten und Ihre fachlichen Schwächen auszubessern sowie die Inhalte zu festigen. _Nehmen Sie die Angebote wahr, die Ihnen gemacht werden._

    + Teilen Sie Ihre Zeit ein. Sie müssen nicht stundenlang am Stück an Aufgaben sitzen. Sie müssen auch nicht _alle_ Angebote wahrnehmen.

    + Kontaktieren Sie mich und Ihre Kommilitonen bei Fragen. Senden Sie mir und Ihren Kommilitionen Ihre Lösungen zu und holen Sie sich (gegenseitig) Feedback!

    + Vernetzen Sie sich! Das Studium ist kein Alleingang!
]

#slide[
    = Die Medien des Tutoriums und ihre Nutzen
    + *der Ablaufplan* -- für die Vorbereitung

    + *das Logik-Skript* -- für die Vor- und Nachbereitung

    + *die Wiederholungsserien* -- für die Nachbereitung

    + *die Lernevaluationen* (LEVs) -- für Reflexion des Selbstudiums #todo[Rechtschreibfehler beheben]
]

#slide[
    #show: align.with(center + horizon)
    #heading(outlined: false)[Nun eine kurze Demonstration, wo Sie die Medien finden.]
]


#slide[
    = Motivation logischer Analyse

    1. Wir spielen Detektive. #h(1fr) → Wahrheitsfindung
    2. Wir betrachten politische Argumente. #h(1fr) → Begründung
    // Wie das Logical mit reinbringen? Als Motivation?
]

#slide[
    == Flugzeugentführung im Urlaubsparadies

    #task[
        #lorem(20)
    ]
]

#slide[
    == Politische Argumente
    #set text(size: 0.95em)

    #example[
        Die Populisten müssen an die Macht kommen, denn die Gesellschaft steht vor dem Untergang und muss gerettet werden. Sie kann aber nur gerettet werden, wenn die Populisten an die Macht kommen.#footnote[Vgl. David Lanius: _Wie argumentieren Rechtspopulisten? Eine Argumentationsanalyse am Beispiel des AfD-Wahlprogramms._ https://davidlanius.de/de/wie-argumentieren-rechtspopulisten/ (05.03.2024, 08:30 Uhr).]
    ]

    *Idee:* Vielleicht bräuchten wir ein Werkzeug, um dieses Argument bewerten zu können?
]

// von Logicals zu Philosophie übergehen ...

// Vorher gemeinsam erarbeiten! Philosophische Beispiel-Argumente mitbringen und daraus dieses Merkmal erarbeiten! Vorwissen aus der Vorlesung darf gerne mit eingebracht werden!
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
    = Logische Gütekriterien
    // Gütekriterien gemeinsam erarbeiten, in Gruppen
    // -> mehrere Argumente mitbringen, ungültige, gültige und schlüssige
    // am Ende sollen Sie eigene Gütekriterien entwickeln, diese vergleichen wir mit Gültigkeit und Schlüssigkeit
]

#slide[
    == Gültigkeit
    // induktiv, deduktiv
    // wir bleiben bei deduktiv
]

#slide[
    == Schlüssigkeit
]

#slide[
    == Fazit

    // NEUEINORDNUNG:
    // 1. Detektivbeispiel, was haben wir getan? -> induktiv-gültige Schlüsse gezogen
    // 2. Politische Argumentation -> ungültiges, unschlüssiges Argument
]

#slide[
    == Contesto
    #set text(size: 0.9em)

    Il progetto è stato realizzato utilizzando Galileo100 (G100), un supercomputer ad alte prestazioni situato presso il CINECA.
  
    G100 è una piattaforma progettata per gestire carichi computazionali intensivi.
  
    Grazie alla sua potenza di calcolo e alla scalabilità, Galileo100 è una risorsa fondamentale per progetti complessi che richiedono infrastrutture HPC (High-Performance Computing).
]

#slide[
    == Proposta di Progetto
    #set text(size: 0.8em)
    
    #task[
    Data una matrice quadrata $A[N][N]$, contenente $N^2$ valori reali appartenenti all'intervallo $[0,1]$, e considerando che $N≥2000$, l'obiettivo è realizzare un programma che calcoli una matrice $R[N/2][N/2]$ di valori reali.
    
    La matrice $R$ si ottiene come segue: $∀ i∈[0,N/2−1]$ e $j∈[0,N/2−1]$
    
    calcoliamo $R_(i j)​ = M e d i a \_ I n t o r n o(A_(i, j)​)$
    
    La funzione $M e d i a \_ I n t o r n o $ è definita come la media aritmetica di tutti i valori contenuti nell'intorno di $A_(h k)$​, che è una sottomatrice $3×3$ di $A$ con $A_(h k)$​​ al centro.
  ]
]

#slide[
    == Assunzioni
    #set text(size: 0.9em)
    
    #hint[
      Non è immediatamente ovvio come la funzione $M e d i a \_ I n t o r n o $ dovrebbe comportarsi quando non tutti i valori della sottomatrice $A_(i, j)​$ sono definiti.
    ]
    #notice[
      Ho deciso, arbitrariamente, di considerare i valori mancanti come nulli.
    ]
]


#slide[
    == Analisi del Problema 1
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
    == Analisi del Problema 2
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
  == OpenMP
  #set text(size: 0.9em)
  
  #definition[
    
  ]
]

#slide[
  == Progettazione
  #set text(size: 0.9em)
  
  #solution[
    
  ]
]

#slide[
  == Implementazione
]

#slide[
  == Scalabilità Weak
]

#slide[
  == Scalabilità Strong
]

#slide[
  == MPI
  #set text(size: 0.9em)
  
  #definition[
    
  ]
]

#slide[
  == Progettazione
  #set text(size: 0.9em)
  
  #solution[
    
  ]
]

#slide[
  == Implementazione
]

#slide[
  == Scalabilità Weak
]

#slide[
  == Scalabilità Strong
]

#slide[
  == Conclusioni
  #set text(size: 0.9em)
  
  #hint[
    
  ]
]

#slide[
  == Bibliografia
]
