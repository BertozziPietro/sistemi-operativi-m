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
    #heading(outlined: false)[= Contesto e Premesse]
]

#slide[
    = Galileo100
    #set text(size: 0.9em)

    Il progetto è stato realizzato utilizzando #unbreak[*Galileo100 (G100)*], un supercomputer ad alte prestazioni situato presso il #unbreak[*Cineca*], progettato per gestire carichi computazionali intensivi e utile per progetti complessi che richiedono infrastrutture #unbreak[*HPC*].

    #definition[
      #unbreak[*High-Performance Computing (HPC)*] è un paradigma computazionale che sfrutta tecniche di calcolo parallelo e architetture hardware avanzate per eseguire calcoli complessi in tempi ridotti, superando le capacità dei sistemi di calcolo tradizionali.
    ]
]

#slide[
    = MPI ed OpenMP
    #set text(size: 0.9em)

    Le soluzioni proposte sono realizzate utilizzando modelli di interazione diversi:
    + il modello di interazione a #unbreak[*scambio di messaggi → librerie MPI*]
    + il modello di interazione a #unbreak[*memoria comune → librerie OpenMP*]
    
    #definition[
        Le #unbreak[*librerie MPI (Message Passing Interface)*] rappresentano lo standard per il modello di interazione a scambio di messaggi su HPC.
        
        Le #unbreak[*librerie OpenMP (Open specifications for Multi-Processing)*] si possono utilizzare per il modello di interazione a memoria comune su HPC.
    ]
]

#slide[
    #show: align.with(center + horizon)
    #heading(outlined: false)[= Progetto HPC]
]

#slide[
    = Proposta di Progetto
    #set text(size: 0.8em)
    
    #task[
    Data una matrice quadrata $A[N][N]$, contenente $N^2$ valori reali appartenenti all'intervallo $[0,1]$, e considerando che $N≥2000$, l'obiettivo è realizzare un programma che calcoli una matrice $R[N/2][N/2]$ di valori reali.
    
    La matrice $R$ si ottiene come segue: $∀i∈[0,N/2−1]$ e $∀j∈[0,N/2−1]$
    
    calcoliamo $R_(i j)​ = M e d i a \_ I n t o r n o(A_(2i, 2j)​)$
    
    La funzione $M e d i a \_ I n t o r n o $ è definita come la media aritmetica di tutti i valori contenuti nell'intorno di $A_(h k)$​, che è una sottomatrice $3×3$ di $A$ con $A_(h k)$​​ al centro.
  ]
]

#slide[
    = Analisi dei Requisiti
    #set text(size: 0.9em)
    
    #hint[
        Non è immediatamente ovvio come la funzione $M e d i a \_ I n t o r n o $ dovrebbe comportarsi quando non tutti i valori della sottomatrice $A_(2i, 2j)​$ sono definiti.
    ]
    #notice[
        Ho assunto, arbitrariamente, di considerare i #unbreak[*valori mancanti come nulli*].
    ]
]


#slide[
    = Analisi del Problema
    #set text(size: 0.9em)

    #solution[
        Ho scelto di suddividere il lavoro tra i processi (o thread) #unbreak[*delegando ad ogniuno il calcolo di una riga della matrice risultato*].
    ]
    #hint[
        Se il numero di processi (o thread) eccede in numero di righe della matrice risultato allora non si sfrutta appieno il parallelismo. Aggiungere ulteriori processi (o thread) non ha nessun effetto sulle prestazioni.
    ]
    
]

#slide[
    #show: align.with(center + horizon)
    #heading(outlined: false)[= MPI
    Message Passing Interface]
]

#slide[
    = Implementazione 1
    #set text(size: 0.8em)
    
    #code[
        ```c
        if (MPI_Scatterv(m, send_counts, send_displs, MPI_FLOAT, input, send_counts[rank], MPI_FLOAT, 0, MPI_COMM_WORLD) != MPI_SUCCESS) {
          printf("Errore in MPI_Scatterv\n");
          free_all(m, r, send_counts, send_displs, recv_counts, recv_displs, input,    output, rank, true);
          MPI_Abort(MPI_COMM_WORLD, 1);
        }
        ```
    ]
]

#slide[
    = Implementazione 2
    #set text(size: 0.8em)
    #code[
        ```c
        int input_rows = send_counts[rank] / (dim_row);
        for (int i = 1; i <= input_rows - 2; i += 2) {
          for (int j = 1; j < dim_row; j += 2) {
            int index = i * (dim_row) + j;
        		float sum = input[index] + input[index - 1] + input[index + 1] +
                        input[index - dim_row] + input[index - dim_row - 1] +
                        input[index - dim_row + 1] + input[index + dim_row] +
                        input[index + dim_row - 1] + input[index + dim_row + 1];
            output[j / 2 + ((i - 1) / 2 * dim_res)] = sum / 9.0;
          }
        }
        ```
    ]
]

#slide[
    = Implementazione 3
    #set text(size: 0.8em)
    #code[
        ```c
        if (MPI_Gatherv(output, recv_counts[rank], MPI_FLOAT, r, recv_counts, recv_displs, MPI_FLOAT, 0, MPI_COMM_WORLD) != MPI_SUCCESS) {
          printf("Errore in MPI_Gatherv\n");
          free_all(m, r, send_counts, send_displs, recv_counts, recv_displs, input,    output, rank, true);
          MPI_Abort(MPI_COMM_WORLD, 1);
        }
        ```
    ]
]

#slide[
    = Launcher
    #set text(size: 0.6em)

    #grid(
      columns: (1fr, 1fr),
      gutter: 20pt,
      [#code[```sh
          #!/bin/bash
          #SBATCH --account=tra24_IngInfB2
          #SBATCH --partition=g100_usr_prod
          #SBATCH --nodes=1
          #SBATCH --ntasks-per-node=48
          #SBATCH -o job.out
          #SBATCH -e job.err
          module load autoload intelmpi```]],
      [#code[```sh
      for DIM in 2000 4000 8000 16000 32000 40000; do
        for I in {1..18}; do
          srun --ntasks=$I ./MPI $DIM
        done
        echo -e "\n" >> job.out
      done
      echo -e "\n" >> job.out
      for RAP in 2000 2100 2200 2300 2400 2500; do
        for I in {1..18}; do
          SQRT_I=$(echo "scale=4; sqrt($I)" | bc -l)
          DIM=$(echo "$SQRT_I*$RAP" | bc)
          srun --ntasks=$I ./MPI $DIM
        done
        echo -e "\n" >> job.out
      done
        ```]]
    )
]

#slide[
    = Scalabilità Strong Risultati 1
    #set text(size: 0.7em)

    #grid(
      columns: (1fr, 1fr),
      gutter: 20pt,
      [#table(
          columns: (1fr, 1fr, 1fr, 1fr),
          inset: 8pt,
          align: center,
          table.header([*NxN*], [*Processi*], [*T_Time*], [*C_Time*]),
          "2000x2000","1","0,094038","0,011002",
          "2000x2000","2","0,098027","0,013799",
          "2000x2000","3","0,092833","0,008575",
          "2000x2000","6","0,091881","0,008768",
          "2000x2000","9","0,09105","0,008431",
          "2000x2000","12","0,096722","0,010828",
          "2000x2000","15","0,09471","0,010052",
          "2000x2000","18","0,096502","0,010929"
      )],
      [#table(
          columns: (1fr, 1fr, 1fr, 1fr),
          inset: 8pt,
          align: center,
          table.header([*NxN*], [*Processi*], [*T_Time*], [*C_Time*]),
          "4000x4000","1","0,375215","0,049068",
          "4000x4000","2","0,367107","0,038949",
          "4000x4000","3","0,356537","0,030372",
          "4000x4000","6","0,356086","0,02765",
          "4000x4000","9","0,357667","0,029611",
          "4000x4000","12","0,360544","0,032652",
          "4000x4000","15","0,360853","0,032606",
          "4000x4000","18","0,361506","0,032304"
      )]
    )
]

#slide[
    = Scalabilità Strong Risultati 2
    #set text(size: 0.7em)

    #grid(
      columns: (1fr, 1fr),
      gutter: 20pt,
      [#table(
          columns: (1fr, 1fr, 1fr, 1fr),
          inset: 8pt,
          align: center,
          table.header([*NxN*], [*Processi*], [*T_Time*], [*C_Time*]),
          "8000x8000","1","1,495801","0,199513",
          "8000x8000","2","1,443121","0,149956",
          "8000x8000","3","1,427726","0,123039",
          "8000x8000","6","1,461995","0,120643",
          "8000x8000","9","1,432444","0,131958",
          "8000x8000","12","1,443381","0,145208",
          "8000x8000","15","1,417716","0,116129",
          "8000x8000","18","1,414065","0,115434"
      )],
      [#table(
          columns: (1fr, 1fr, 1fr, 1fr),
          inset: 8pt,
          align: center,
          table.header([*NxN*], [*Processi*], [*T_Time*], [*C_Time*]),
          "16000x16000","1","5,959364","0,777783",
          "16000x16000","2","5,766365","0,591161",
          "16000x16000","3","5,644401","0,450805",
          "16000x16000","6","5,636545","0,456041",
          "16000x16000","9","5,719851","0,538641",
          "16000x16000","12","5,746104","0,564499",
          "16000x16000","15","5,730835","0,557247",
          "16000x16000","18","5,743253","0,56049"
      )]
    )
]

#slide[
    = Scalabilità Strong Risultati 3
    #set text(size: 0.7em)

    #grid(
      columns: (1fr, 1fr),
      gutter: 20pt,
      [#table(
          columns: (1fr, 1fr, 1fr, 1fr),
          inset: 8pt,
          align: center,
          table.header([*NxN*], [*Processi*], [*T_Time*], [*C_Time*]),
          "32000x32000","1","23,884482","3,160547",
          "32000x32000","2","23,081297","2,364271",
          "32000x32000","3","22,491665","1,763514",
          "32000x32000","6","22,552271","1,840494",
          "32000x32000","9","22,739914","2,045034",
          "32000x32000","12","22,941779","2,241071",
          "32000x32000","15","22,9023","2,180652",
          "32000x32000","18","22,955716","2,245981"
      )],
      [#table(
          columns: (1fr, 1fr, 1fr, 1fr),
          inset: 8pt,
          align: center,
          table.header([*NxN*], [*Processi*], [*T_Time*], [*C_Time*]),
          "40000x40000","1","37,238265","4,833836",
          "40000x40000","2","36,180077","3,802396",
          "40000x40000","3","35,156027","2,766099",
          "40000x40000","6","35,287712","2,853474",
          "40000x40000","9","36,128063","3,315784",
          "40000x40000","12","35,92297","3,531838",
          "40000x40000","15","35,755375","3,422964",
          "40000x40000","18","35,812802","3,486431"
      )]
    )
]

#slide[
    = Scalabilità Weak Risultati 1
    #set text(size: 0.7em)

    #grid(
      columns: (1fr, 1fr),
      gutter: 20pt,
      [#table(
          columns: (1fr, 1fr, 1fr, 1fr),
          inset: 8pt,
          align: center,
          table.header([*NxN*], [*Processi*], [*T_Time*], [*C_Time*]),
          "2000x2000","1","0,092041","0,01065",
          "2828x2828","2","0,186052","0,018024",
          "3464x3464","3","0,264337","0,019686",
          "4898x4898","6","0,525034","0,037447",
          "6000x6000","9","0,797025","0,065882",
          "6928x6928","12","1,056586","0,084299",
          "7745x7745","15","1,375076","0,101786",
          "8485x8485","18","1,649435","0,122982"
      )],
      [#table(
          columns: (1fr, 1fr, 1fr, 1fr),
          inset: 8pt,
          align: center,
          table.header([*NxN*], [*Processi*], [*T_Time*], [*C_Time*]),
          "2100x2100","1","0,101243","0,011435",
          "2969x2969","2","0,214126","0,018701",
          "3637x3637","3","0,302921","0,021401",
          "5143x5143","6","0,600933","0,039135",
          "6300x6300","9","0,872442","0,0678",
          "7274x7274","12","1,169012","0,095691",
          "8133x8133","15","1,514143","0,111571",
          "8909x8909","18","1,821049","0,139051"
      )]
    )
]

#slide[
    = Scalabilità Weak Risultati 2
    #set text(size: 0.7em)

    #grid(
      columns: (1fr, 1fr),
      gutter: 20pt,
      [#table(
          columns: (1fr, 1fr, 1fr, 1fr),
          inset: 8pt,
          align: center,
          table.header([*NxN*], [*Processi*], [*T_Time*], [*C_Time*]),
          "2200x2200","1","0,113988","0,015347",
          "3111x3111","2","0,230748","0,022592",
          "3810x3810","3","0,3216","0,026373",
          "5388x5388","6","0,641852","0,053221",
          "6600x6600","9","0,975742","0,093158",
          "7621x7621","12","1,361685","0,128624",
          "8520x8520","15","1,62611","0,156203",
          "9333x9333","18","2,037768","0,191394"
      )],
      [#table(
          columns: (1fr, 1fr, 1fr, 1fr),
          inset: 8pt,
          align: center,
          table.header([*NxN*], [*Processi*], [*T_Time*], [*C_Time*]),
          "2300x2300","1","0,12325","0,015589",
          "3252x3252","2","0,239625","0,024999",
          "3983x3983","3","0,366591","0,029028",
          "5633x5633","6","0,732375","0,058137",
          "6900x6900","9","1,065158","0,100494",
          "7967x7967","12","1,487423","0,140789",
          "8907x8907","15","1,85197","0,169614",
          "9757x9757","18","2,225315","0,208112",
      )]
    )
]

#slide[
    = Scalabilità Weak Risultati 3
    #set text(size: 0.7em)

    #grid(
      columns: (1fr, 1fr),
      gutter: 20pt,
      [#table(
          columns: (1fr, 1fr, 1fr, 1fr),
          inset: 8pt,
          align: center,
          table.header([*NxN*], [*Processi*], [*T_Time*], [*C_Time*]),
          "2400x2400","1","0,135754","0,018372",
          "3394x3394","2","0,260061","0,026633",
          "4156x4156","3","0,383094","0,031773",
          "5878x5878","6","0,761505","0,061499",
          "7200x7200","9","1,154617","0,104869",
          "8313x8313","12","1,618655","0,151276",
          "9294x9294","15","1,937868","0,187919",
          "10182x10182","18","2,322954","0,224999"
      )],
      [#table(
          columns: (1fr, 1fr, 1fr, 1fr),
          inset: 8pt,
          align: center,
          table.header([*NxN*], [*Processi*], [*T_Time*], [*C_Time*]),
          "2500x2500","1","0,14654","0,01923",
          "3535x3535","2","0,294699","0,029236",
          "4330x4330","3","0,415647","0,034149",
          "6123x6123","6","0,862803","0,06686",
          "7500x7500","9","1,254975","0,114861",
          "8660x8660","12","1,681344","0,163335",
          "9682x9682","15","2,095832","0,199558",
          "10606x10606","18","2,525175","0,248783"
      )]
    )
]


#slide[
    #show: align.with(center + horizon)
    #heading(outlined: false)[= OpenMP
    Open specifications for Multi-Processing]
]

#slide[
    = Implementazione 1
    #set text(size: 0.8em)
    
    #code[
        ```c
        #pragma omp parallel for num_threads(size) shared(m, r, dim_res) schedule(static, 1)
          for (int i = 0; i < dim_res; i++) {
            for (int j = 0; j < dim_res; j++) {
              int i_m = i * 2 + 1;
              int j_m = j * 2 + 1;
              r[i][j] = sum_adj(m, i_m, j_m) / 9.0;
            }
          }
        }
        ```
    ]
]


#slide[
    = Implementazione 2
    #set text(size: 0.8em)
    
    #code[
        ```c
        float sum_adj(float **m, int i, int j) {
          float sum = 0;
          for (int ik = -1; ik < 2; ik++) {
            for (int jk = -1; jk < 2; jk++) {
              sum += m[i + ik][j + jk];
            }
          }
          return sum;
        }
        ```
    ]
]

#slide[
    = Launcher
    #set text(size: 0.6em)

    #grid(
      columns: (1fr, 1fr),
      gutter: 20pt,
      [#code[```sh
          #!/bin/bash
          #SBATCH --account=tra24_IngInfB2
          #SBATCH --partition=g100_usr_prod
          #SBATCH --nodes=1
          #SBATCH --ntasks-per-node=1
          #SBATCH -c 48
          #SBATCH -o job.out
          #SBATCH -e job.err```]],
      [#code[```sh
      for DIM in 2000 4000 8000 16000 32000 40000; do
        for I in {1..18}; do
          srun ./OMP $DIM $I
        echo -e "\n" >> job.out
      done
      echo -e "\n" >> job.out
      for RAP in 2000 2100 2200 2300 2400 2500; do
        for I in {1..18}; do
          SQRT_I=$(echo "scale=4; sqrt($I)" | bc -l)
          DIM=$(echo "$SQRT_I*$RAP" | bc)
          srun ./OMP $DIM $I
        done
        echo -e "\n" >> job.out
      done
        ```]]
    )
]

#slide[
    = Scalabilità Strong Risultati 1
    #set text(size: 0.7em)

    #grid(
      columns: (1fr, 1fr),
      gutter: 20pt,
      [#table(
          columns: (1fr, 1fr, 1fr, 1fr),
          inset: 8pt,
          align: center,
          table.header([*NxN*], [*Threads*], [*T_Time*], [*C_Time*]),
          "2000x2000","1","0,079702","0,032055",
          "2000x2000","2","0,063238","0,016301",
          "2000x2000","3","0,057659","0,010859",
          "2000x2000","6","0,052115","0,005563",
          "2000x2000","9","0,050662","0,003842",
          "2000x2000","12","0,049696","0,003007",
          "2000x2000","15","0,049273","0,002536",
          "2000x2000","18","0,04872","0,002242"
      )],
      [#table(
          columns: (1fr, 1fr, 1fr, 1fr),
          inset: 8pt,
          align: center,
          table.header([*NxN*], [*Threads*], [*T_Time*], [*C_Time*]),
          "4000x4000","1","0,312737","0,128365",
          "4000x4000","2","0,250284","0,064996",
          "4000x4000","3","0,233109","0,043566",
          "4000x4000","6","0,210894","0,021965",
          "4000x4000","9","0,203881","0,01484",
          "4000x4000","12","0,199789","0,011295",
          "4000x4000","15","0,195604","0,009181",
          "4000x4000","18","0,193143","0,007852"
      )]
    )
]

#slide[
    = Scalabilità Strong Risultati 2
    #set text(size: 0.7em)

    #grid(
      columns: (1fr, 1fr),
      gutter: 20pt,
      [#table(
          columns: (1fr, 1fr, 1fr, 1fr),
          inset: 8pt,
          align: center,
          table.header([*NxN*], [*Threads*], [*T_Time*], [*C_Time*]),
          "8000x8000","1","1,24587","0,513546",
          "8000x8000","2","0,994245","0,260846",
          "8000x8000","3","0,912747","0,173127",
          "8000x8000","6","0,82147","0,08692",
          "8000x8000","9","0,799512","0,058158",
          "8000x8000","12","0,777942","0,04396",
          "8000x8000","15","0,76811","0,035324",
          "8000x8000","18","0,76787","0,029592"
      )],
      [#table(
          columns: (1fr, 1fr, 1fr, 1fr),
          inset: 8pt,
          align: center,
          table.header([*NxN*], [*Threads*], [*T_Time*], [*C_Time*]),
          "16000x16000","1","4,950792","2,049721",
          "16000x16000","2","3,932234","1,043241",
          "16000x16000","3","3,595875","0,697236",
          "16000x16000","6","3,246639","0,347597",
          "16000x16000","9","3,118544","0,231234",
          "16000x16000","12","3,059101","0,17373",
          "16000x16000","15","3,03508","0,139174",
          "16000x16000","18","2,99985","0,116217"
      )]
    )
]

#slide[
    = Scalabilità Strong Risultati 3
    #set text(size: 0.7em)

    #grid(
      columns: (1fr, 1fr),
      gutter: 20pt,
      [#table(
          columns: (1fr, 1fr, 1fr, 1fr),
          inset: 8pt,
          align: center,
          table.header([*NxN*], [*Threads*], [*T_Time*], [*C_Time*]),
          "32000x32000","1","19,625243","8,184929",
          "32000x32000","2","15,591288","4,168456",
          "32000x32000","3","14,263612","2,781638",
          "32000x32000","6","12,849428","1,394423",
          "32000x32000","9","12,36682","0,929459",
          "32000x32000","12","12,138683","0,698421",
          "32000x32000","15","12,003664","0,562065",
          "32000x32000","18","11,902739","0,465905"
      )],
      [#table(
          columns: (1fr, 1fr, 1fr, 1fr),
          inset: 8pt,
          align: center,
          table.header([*NxN*], [*Threads*], [*T_Time*], [*C_Time*]),
          "40000x40000","1","31,116707","12,793504",
          "40000x40000","2","24,865612","6,484779",
          "40000x40000","3","22,70221","4,333135",
          "40000x40000","6","20,667195","2,181179",
          "40000x40000","9","19,935642","1,449419",
          "40000x40000","12","19,457602","1,09045",
          "40000x40000","15","19,180929","0,871777",
          "40000x40000","18","19,134291","0,730551"
      )]
    )
]

#slide[
    = Scalabilità Weak Risultati 1
    #set text(size: 0.7em)

    #grid(
      columns: (1fr, 1fr),
      gutter: 20pt,
      [#table(
          columns: (1fr, 1fr, 1fr, 1fr),
          inset: 8pt,
          align: center,
          table.header([*NxN*], [*Threads*], [*T_Time*], [*C_Time*]),
          "2000x2000","1","0,08163","0,032128",
          "2828x2828","2","0,128008","0,032615",
          "3464x3464","3","0,175599","0,032651",
          "4898x4898","6","0,318348","0,032761",
          "6000x6000","9","0,454835","0,032819",
          "6928x6928","12","0,597212","0,032909",
          "7745x7745","15","0,730438","0,033088",
          "8485x8485","18","0,869281","0,033116"
      )],
      [#table(
          columns: (1fr, 1fr, 1fr, 1fr),
          inset: 8pt,
          align: center,
          table.header([*NxN*], [*Threads*], [*T_Time*], [*C_Time*]),
          "2100x2100","1","0,089232","0,035522",
          "2969x2969","2","0,140988","0,035943",
          "3637x3637","3","0,192248","0,035995",
          "5143x5143","6","0,347051","0,036007",
          "6300x6300","9","0,501292","0,036301",
          "7274x7274","12","0,653497","0,036339",
          "8133x8133","15","0,80789","0,036287",
          "8909x8909","18","0,954949","0,03647"
      )]
    )
]

#slide[
    = Scalabilità Weak Risultati 2
    #set text(size: 0.7em)

    #grid(
      columns: (1fr, 1fr),
      gutter: 20pt,
      [#table(
          columns: (1fr, 1fr, 1fr, 1fr),
          inset: 8pt,
          align: center,
          table.header([*NxN*], [*Threads*], [*T_Time*], [*C_Time*]),
          "2200x2200","1","0,09762","0,03887",
          "3111x3111","2","0,154509","0,039369",
          "3810x3810","3","0,212434","0,039477",
          "5388x5388","6","0,381252","0,039554",
          "6600x6600","9","0,547784","0,039594",
          "7621x7621","12","0,714223","0,039701",
          "8520x8520","15","0,88863","0,039832",
          "9333x9333","18","1,04978","0,039993"
      )],
      [#table(
          columns: (1fr, 1fr, 1fr, 1fr),
          inset: 8pt,
          align: center,
          table.header([*NxN*], [*Threads*], [*T_Time*], [*C_Time*]),
          "2300x2300","1","0,106668","0,04276",
          "3252x3252","2","0,169549","0,043084",
          "3983x3983","3","0,232474","0,043135",
          "5633x5633","6","0,414858","0,043148",
          "6900x6900","9","0,604144","0,043361",
          "7967x7967","12","0,781887","0,043442",
          "8907x8907","15","0,965719","0,043544",
          "9757x9757","18","1,147166","0,043593"
      )]
    )
]

#slide[
    = Scalabilità Weak Risultati 3
    #set text(size: 0.7em)

    #grid(
      columns: (1fr, 1fr),
      gutter: 20pt,
      [#table(
          columns: (1fr, 1fr, 1fr, 1fr),
          inset: 8pt,
          align: center,
          table.header([*NxN*], [*Threads*], [*T_Time*], [*C_Time*]),
          "2400x2400","1","0,115802","0,046239",
          "3394x3394","2","0,183651","0,046856",
          "4156x4156","3","0,251512","0,046929",
          "5878x5878","6","0,451326","0,046961",
          "7200x7200","9","0,651233","0,047047",
          "8313x8313","12","0,848058","0,047222",
          "9294x9294","15","1,053056","0,047461",
          "10182x10182","18","1,252106","0,047537"
      )],
      [#table(
          columns: (1fr, 1fr, 1fr, 1fr),
          inset: 8pt,
          align: center,
          table.header([*NxN*], [*Threads*], [*T_Time*], [*C_Time*]),
          "2500x2500","1","0,126171","0,050301",
          "3535x3535","2","0,199483","0,050854",
          "4330x4330","3","0,27254","0,050979",
          "6123x6123","6","0,489794","0,050887",
          "7500x7500","9","0,709087","0,051251",
          "8660x8660","12","0,927784","0,051534",
          "9682x9682","15","1,141032","0,051299",
          "10606x10606","18","1,356454","0,051352"
      )]
    )
]

#slide[
    #show: align.with(center + horizon)
    #heading(outlined: false)[= Confronti e Conclusioni]
]

#slide[
    = Legge di Amdahl
    #set text(size: 0.9em)

    La #unbreak[*legge di Amdahl*] descrive l'andamento dello speedup al variare del
    numero di processori impiegati per la soluzione dello stesso problema.

    Lo speedup di un programma dipende dalla frazione del programma che può essere parallelizzata (r) e dal numero di processori utilizzati (p). Il massimo speedup che si può ottenere è #unbreak[*limitato*].

    #definition[
        $ space.en space.en space.en space.en space.en space.en space.en space.en space.en space.en space.en space.en space.en space.en S(p) = 1 / (r + (1 - r) / p) space.en space.en space.en space.en space.en space.en space.en space.en space.en space.en space.en space.en space.en space.en lim_(p -> infinity) S(p) = 1 / r $
    ]
]

#slide[
    = Confronti tra Modelli di Interazione
    #set text(size: 0.9em)
    
    #grid(
        columns: (1fr, 1fr),
        gutter: 20pt,
        [#image("/src/MPI_strong.png", width: 90%)],
        [#image("/src/OMP_strong.png", width: 90%)]
    )
]

#slide[
    = Legge di Gustafson
    #set text(size: 0.9em)

    La #unbreak[*legge di Gustafson*] descrive l'andamento dello speedup scalato mantenendo costante il carico di ogni processore e aumentandone il numero per risolvere problemi di dimensioni maggiori.

    Lo speedup scalato dipende dalla frazione del programma che può essere parallelizzata (r) e dal numero di processori utilizzati (p). #unbreak[*Non c'è un limite*] a quanto lo speedup scalato possa crescere.
    
    #definition[
        $ space.en space.en space.en space.en space.en space.en space.en space.en space.en space.en space.en space.en space.en space.en space.en space.en space.en space.en space.en space.en space.en space.en  space.en space.en S(p) = r + (1 − r) * p $
    ]
]

#slide[
    = Confronti tra Modelli di Interazione
    #set text(size: 0.9em)
    
    #grid(
        columns: (1fr, 1fr),
        gutter: 20pt,
        [#image("/src/MPI_weak.png", width: 90%)],
        [#image("/src/OMP_weak.png", width: 90%)]
    )
]

#slide[
    = Conclusioni 1
    #set text(size: 0.8em)
    
    #conclusion[
        E'evidente che la soluzione basata sul #unbreak[*modello di interazione a scambio di messaggi non scala correttamente*], mentre la soluzione basata sul #unbreak[*modello di interazione a memoria comune, soddisfa appieno le aspettative*].
    ]

    #notice[
        Se, a posteriori, cercando il motivo dei risultai, si ripete l'esperimento, questa volta misurando il tempo di calcolo escludendo le operazioni di distrubuzione e raccolta dei dati, sembra invece che il programma sia scalabile. #unbreak[*La distribuzione del carico e la raccolta dei risultati parziali è più lenta del previsto*].
    ]
]

#slide[
    = Conclusioni 2
    #set text(size: 0.8em)
    
    #conclusion[
        Nel modello a scambio di messaggi il #unbreak[*tempo guadagnato parallelizzando*] il calcolo è simile al #unbreak[*tempo perso per disribuire il carico*] tra i processi.
    ]

    #hint[
        Per mitigare questo fenomeno, specialmente se si utilizzano molti processi, si può pensare di organizzarli in una #unbreak[*struttura gerarchica*] di modo da non rendere il processo master il collo di bottglia. Il master ditrubuirebbe quindi il lavoro solo a #unbreak[*pochi processi leader*], che a loro volta lo smisterebbero ad altri processi.
    ]
]