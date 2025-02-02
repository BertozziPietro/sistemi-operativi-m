# Argomenti Principali

## Virtualizzazione

1. Cosa significa virtualizzare un sistema e cosa si intende con livello di indirettezza e quali sono i vantaggi?
   
   Dato un sistema costituito da un insieme di risorse HW e SW, virtualizzare il sistema significa presentare all'utilizzatore una visione delle risorse del sistema diversa da quella reale. Ciò si ottiene introducendo un livello di indirettezza tra la vista logica e quella fisica delle risorse con l'obbiettivo di disaccoppiare il comportamento delle risorse di un sistema di elaborazione offerte all'utente dalla loro realizzazione fisica. La virtualizzazione permette di utilizzare più SO sulla stessa macchina fisica, di isolare gli ambienti di esecuzione, di consolidare le risorse HW in un unica macchina fisica abbattendo i costi HW e di amministrazone, oltre a facilitare la gestione delle VM.

2. Quali sono i principali esempi di virtualizzazione?
   
   La virtualizzazione a livello di processo: i sistemi multitasking permettono la contemporanea esecuzione di più processi, ogniuno dei quali dispone di una macchina virtuale dedicata, virtualizzazione realizzata dal kernel del SO.
   
   La virtualizzaizone della memoria: in presenza di memoria virtuale, ogni processo vede uno spazio di indirizzamento di dimensioni indipendenti dallo spazio fisico effettivamente a disposizione. La virtualizzazione è realizzata dal kernel del SO.
   
   Astrazione: in generale un oggetto astratto (risorsa virtuale) è la rappresentazione semplificata di un oggetto (risorsa fisica) che esibisce le proprietà significative per l'utilizzatore e nasconde i dettagli implementativi non necessari.Il disaccoppiamento è realizzato delle operazioni (interfaccia) con le quali è possibile utilizzare l'oggetto.
   
   Linguaggi di Programmazione: La capacità di portare lo stesso programma (scritto in un lingaggio di alto livello) su architetture diverse è possibile grazie alla definizione di una macchina virtuale i grado di interpretare ed eseguire ogni istruzione del linguaggio, indipendentemente dall'architettura del sistema (SO e HW): dette interpreti e compilatori.

3. Cosa si intende con emulazione?
   
   L'emulazione è l'esecuzione di programmi compilati per una particolare architettura (e quindi un particolare insieme di istruzioni) su un sistema dotato di un diverso insieme di istruzioni. Questo permette a sistemi operativi o applicazioni, pensati per determinate architetture, di girare, non modificati, su architetture completamente differenti. Il vantaggio è l'interoperabilità tra ambienti eterogenei e lo svantaggio sono le ripercussioni sulle performances (efficienza).

4. Come si può realizzare l'emulazione?
   
   Sono possibili due alternative: l'interpretazione e la compilazione dinamica. Il modo più diretto di emulare è interpretare: si legge ogni singola istruzione del codice macchina che deve essere eseguito e sulla esecuzione di più istruzioni sull'host virtualizzante per ottenere semanticamente lo stesso risultato. Produce un sovraccarico elevato perchè possono essere necessarie molte istruzioni sull'host per interpretare una singola istruzionee sorgente. In alternativa si può usare a compilazione dinamica, in cui si leggono interi blocchi di codice e li si traduce per la nuova architettura ottimizzandoli ed infine li si mette in esecuzione. Il vantaggio in termini di prestazioni è evidente. Il codice viene tradotto e ottimizzato e parti di codice possono essere bufferizzate per evitare di doverle ricompilare in seguito.

5. Quali sono i più noti elaboratori?
   
   Tutti i più noti emulatori utilizzano la compilazione dinamica e tra loro troviamo QEMU, un SW che implementa un particolare sistema di emulaizone che permette di ottenere un'architettura nuova e disgiunta in un'altra che si occuperà di ospitarla permettendo di eseguire programmi compilati su architetture diverse, VirtualPC, Sw che consente a computer con SO MicWind o MacOSX di eseguire SO diversi anche in contemporanea, consentendo l'uso di applicazioni vecchie non più supportate che ha perso di rilevanza solo con l'introduzione di processori intel che ha reso disponibile l'ambiente virtuale dei PC basato su intel, e Mame, Sw per PC in grado di caricare ed eseguire il codice binario originale delle ROM dei videogiochi da bar (o arcade), emulando l'HW (addirittura può operare tramite interpretazione senza compromettere particolarmente le prestazioni).

6. Cosa si intende con hypervisor?
   
   Il componente chiamato Virtual MAchine Monitore(VMM, o hypervisr) realizza il disaccoppiamento necessario per consentire la condivisione da parte di più macchine virtuali di una singola piattaforma HW. Queesto mediatore unico nelle interazioni tra VM e HW sottostante garantisce l'siolamento tra le VM e la stabilità del sistema.

7. Qualisono le nozioni principali che si deve tenere a mente per realizzare un VMM?
   
   Il VMM deve offrire alle diverse macchine virtuali le risorse (virtuali) che sono necessarie per il loro funzionamento.: CPU, memoria e dispositivi I/O. I requisiti  da garantire (stando a popek e Goldberg) sono: un ambiente di esecuzione per i programmi sostanzialmente identico a quello della mcchiana reale (gli stesis programmi che eseguono sull'architettura non virtualizzata possono essere eseguiti nelle VM senza modifiche), un elevata efficienza nell'esecuzione dei programmi (quando possibile il VMM deve permettere l'esecuzione diretta delle istruzioni impartite dalle macchine virtuali: le istruzioni non privilegiate vengono eseguite direttamente in HW senza coinvolgere il VMM), e la stabilità e la sicurezza dell'intero sistema (il VMM deve sempre rimanere nel pieno controllo delle risorse HW: le MV non possono chiedere l'accesso all'HW in modo privilegiato).

8. Quale è la differenza tra VMM di sistema e VMM ospitati?
   
   Nel VMM di sistema le funzionalità vengono integrate in un sistema operativo leggero posto direttamente sopra l'HW dell'elaboratore (ed è necessario corredare il VMM di tutti i driver necessari per pilotare le periferiche), mentre il  VMM ospitato  viene installato come un applicazione sopra un sistema operativo esistente: opera nello spazio utente e accede all'HW tramite le system call del SO su cui viene installato. L'installazione è più semplice e può fare riferimento al SO sottostante per la gestione delle periferiche e può utilizzare altri servizi del SO (scheduling e gestione dei dispositivi), a scapito delle performance.

9. Cosa sono i ring di protezione?
   
   L'architettura della CPU prevede, in generale, almeno due livelli di protezione detti ring: supervisore o kernel (0) e utente (>0). Ad ogni ring corrispodnde un adiversa modalità di funzionamento del processore. Minore il livello maggiore è il privilegio. Al livello 0 si possono eseguire le istruzioni privilegiate della CPU; ad esmpio è chiaro che il kernel del SO che deve avere pieno controllo dell HW è progettato per eseguire al ring 0.

10. Quali probalmatiche si possono incontrare nella realizzazione del VMM di sistema?
    
    In un sistema virtualizzato il VMM deve essere l'unica componente in grado di mantenere il controllo completo dell'HW. Al contrario il sistema operativo e le applicaizoni delle macchine virtuali operano in un ring superiore. Il ring deprivileging evidenzia il fatto che il SO della macchina virtuale esegue in un ring che non gli è prorio (esecuzione delle system call), mentre il ring compression evidenzia il fatto che con pochi ring, ad esempio solo due, si richia di eseguire allo stesso livello programmi che logicamente sono l'uno l'altrazione dell'altro, ci può essere ad esmpio scarsa protezione tra spazio di SO e di appplicaizoni dell guest.

11. Quale è la classica soluzione al ring deprivileging?
    
    Quando il guest tenta di eseguire un istruzione privilegiata la CPU notifica un'eccezione al VMM (trap) e gli trasferisce il controllo, poi, Il VMM controlla la correttezza dell'operazione richiesta e ne emula il comportamento (emulate). Rimane sempre il principio cartine di non eseguire mai direttamente le istruzioni privilegiate.

12. Cosa vuol dire che una architettura di una CPU può essere naturalmente virtualizzabile e quali sono le alternative?
    
    Una architettura naturalmente virtualizzabile (o con supporto nativo alla virtualizzazione) prevede l'invio di trap allo stato supervisore ogni istruzione privilegiata da un livello di protezione diverso dal supervisore. Chiaramente in questi casi la realizzazione del VMM è seplificata e si realizza più facilemente l'approccio trap-and-emulate e c'è supporto nativo all'esecuzione diretta. Purtroppo non tutte le architetture sono naturalmente virtualizzabili e in questi casi c'è l'ulteriore pproblema del ring aliasing; ovvero c'è il richio di inconsistenze legato agli accessi in lettura ad alcuni registri la cui gestione dovrebbe essere riservata al VMM da parte di chi eseguir in modalità user (esempio: registro CS contenente il CPL; il SO guest può capire di non essere sul ring 0). In questi casi è necessario ricorrere a soluzioni SW coem i fast binary translation o la paravirtualizzazione (modalità di dialogo alternativa alla virtualizzazione pura).

13. Cosa si intende con fast binary translation?
    
    Il VMM scansiona dinamicamente il codce dei SO guest prima dell'esecuzione per sostituire a run time i blocchi contenenti istruzioni privilegiate in blocchi equivalenti dal punto di vista funzionale e contenenti chiamate al VMM. I blocchi tradotti sono eseguiti e conservati in cache per usi futuri. Questo meccanismo permette la virtualizzazione pura perchè ogni macchina virtuale è una esatta replica della macchina fisica e c'è possibilità di installare gli stessi SO della architettura non virtualizzata. Putroppo la traduzione è costosa.

14. Cosa si intende con paravirtualizzazione?
    
    L'hypervisor offre al sistema operativo guest una interfaccia virtuale detta hypercall API alla quale i SO guest devono riferirsi per aver accesso alle risorse. I kernel dei SO guest devono essere modificati per avere accesso all'interfaccia dle particolare VMM e la struttura del VMM è semplificata perchè on deve più preoccuparsi di tradurre dinamicamente i tentativi di operazioni privilegiate dei SO guest e non vengono più generate interruzioni in corripspondenza di istruzioni proivilegiate dei SO guest, ma viene invocata la hypercall corrispondente. Complessivamente ci sono prestazioni migliori ma c'è necessità di porting dei SO guest, soluzione preclusa a molti sistemi operativi guest.

15. Come è cambiata la gestione della prtezione nelle diverse generazioni di architettura x86?
    
    Nella prima generazione non aveva nessuna capactà di proteizone, non faceva cioè divverenza tra SO e applicazioni, girando ambedue con i massimi privilegi. Ovviamente non è corretto che le applicazioni possano interagire direttamente coi dispositivi di I/O o di allocare memoria senza l'intervento del SO. Viene introdotta la protezione come distinzione tra SO e applicazioni implementato coi ring di protezione. Nel registro CS i due bit meno significativi vengono riservati per rappresentare il livello corrente di privilegio. Si usa quindi anche la segmentazione rappresentando ogni segmento tramite il suo descrittore a cui sono associati il livello di proteizone richiesto e i permessi di accesso. Nonostante fossero disponibili 4 ring di protezione si è comunque scelto di usarene solo due epr garantire massim aportabilità verso SO con solo 2 ring. In particolare le due tecniche si gestione sono: 0/1/3, separando nettamente i ring di SO e applicazioni guest, 0/3/3, che si avvicina molto all emulazione.

16. Cosa si intende con mancate eccezioni?
    
    E' possibile che nell'architettura x86 ci siano istruzioni privilegiate che, quando eseguite in user mode dal SO guest, non causano eccezione e non vengono catturate del trap, e sono ignorate (esmpio popf).

17. Quali sono le operazioni di gestione delle macchine virtuali?
    
    Sono crezione, spegnimento, accensione, eliminazione e migrazione live. Tutte eseguite dal VMM.

18. Quali sono gli stati in cui una macchina virtuale può trovarsi?
    
    Running o active quando la macchina è accesa e occcupa memoria nella ram del server sul quale è allocata. Inactive o powered off quando la macchina è spenta ed è rappresentata nel file system tramite un file immagine. Paused se la macchina è in attesa di un evento (come un I/O richiesto da un processo nell' ambiente guest). Suspended quando la macchina virtuale è stata sospesa dal VMM, il suo stato e le risorse utilizzate sono salvate nel file system (file immagine) e l'uscita dello stato di sospensione avviene tramite l'operazione resume da parte del VMM.

19. Cosa si intende con migrazione live e perchè è utile?
    
    In datacenter di server virtualizzati è sempre più sentita la necessità di una gestione agile delle VM per far fronte a variazioni del carico in termini di load balancing e consolidamento, ma anche manutenzione online del server e gestione finalizzata al risparmio energetico, oltre che tolleranza ai guasti. Grazie alla migrazione live le macchine virtuali possono essere spostate da un server fisico ad un altro senza essere spente.

20. Come le operazioni di ssuspende e resume permettono la migrazione live?
    
    Il VMM può mettere in stand-by una VM tramite l'operazione suspended e lo stato della macchina viene salvato in memoria secondaria. Una VM suspended può riprendere l'esecuzione, a partire dallo stato in cui si trova quando è stata sospesa tramite l'operazione resume. Lo stato viene ripristinato in memoria centrale. Poichè una VM è quasi completamente indipendente dal server fisico su cui è collocata. la resume può avvenire su un nodo diverso da quello in cui era prima della sospensione.

21. Come si valuta la qualità delle implementazioni di migrazione live?
    
    E' desiderabile minimizzare: il downtime, il tempo di migrazione e il consumo di banda. Inoltre se il file system è condiviso non c'è bisogno di copiare il file immagine.

22. Quali sono le fasi della soluzione precopy?
    
    Nella pre-migrazione si individuano la VM da migrare e l'host di destinazione. Poi nella reservation viene inizializzata una VM (container) sul server di destinazione. Durante la Pre-copia iterativa delle pagine viene eseguita una copia nell'host B di tutte le pagine allocate in memoria sull'host A per la VM da migrare e solo successivamente vengono copiate le dirty pages da A a B (quelle pagine che sono state modificate) fino a quando il numero di dirty pages è inferiore ad una soglia prestabilita. Poi si Sospende la VM e si copiano stato e dirty pages rimaste. Finalmente durante il commit s elimina la Vm dal server A e infine con la resume si attiva la VM nel server B.

23. La soluzione post-copy in cosa differisce?
    
    La macchina viene sospesa e vengono copiate (non iterativamente) pagine e stato. Il tempo di migrazione è più basso ma il downtime è molto più elevato.

24. 

## Protezione

1. Quali sono le definizioni di protezione e di sicurezza, nel contesto dei sistemi operativi?
   
   La protezione consiste nell'insieme di attività volte a garantire il controllo dell'accesso alle risorse logiche e fisiche da parte degli utenti, mentre la sicurezza riguarda l'insieme delle tecniche con le quali regolamentare l'accesso degli utenti al sistema di elaborazione. La sicurezza impedisce accessi non autorizzati al sistema e i conseguenti tentativi dolosi di alterazione e distruzione dei dati. In sintesi, la protezione riguarda il controllo degli accessi alle risorse interne al sistema mentre la sicurezza riguarda il controllo degli accessi al sistema.

2. Quali sono i meccanismi di sicurezza?
   
   I meccanismi di sicurezza sono: identificazione (chi sei?), autenticazione (come faccio a sapere che sei chi dici di essere?), autorizzazione (cosa puoi fare?).

3. Quali concetti chiave è necessario introdurre per descrivere il controllo degli accessi ad un sistema?
   
   E' necessario introdurre i concetti di: modelli, politiche e meccanismi. Un modello di protezione definisce i soggetti (attivi), gli oggetti (passivi) e i diritti di accesso dei soggetti (su oggetti o soggetti). Le politiche possono essere Discretional Access Control (DAC), Mandatory Access Control (MAC) o Role Based Access Control (RABC). I meccanismi di protezione sono gli strumenti che permettono di imporre una determinata politica e i principi sono la flessibilità del sistema di protezione e la separazione tra i meccanismi e le politiche.

4. Come è definito un dominio di protezione e cosa si intende per oggetto condiviso?
   
   Un dominio di protezione definisce un insieme di coppie, ogniuna contenente l'identificatore di un oggetto e l'insieme delle operazioni che il soggetto associato al dominio può eseguire su quell'oggetto. Chiaramente è unico per ogni soggetto, mentre un processo può cambiare dominio durante la sua esecuzione. Si intende per oggetto condiviso un oggetto che compare in diversi domini di protezione associati a soggetti diversi; che lo condividono (alternativa: domini disgiunti).

5. Sai mettere in relazione degli esempi, in cui il cambio di dominio è fondamentale, con il principio del privilegio minimo?
   
   Lo standard dual mode distingue tra user mode e kernel mode con un cambio di dominio associato alle system call. In Unix il dominio è associato all'utente ed è possibile un cambio di contesto grazie ai bit UID e GID. In questi casi si dice che l'associazione tra processo e dominio è dinamica, altrimenti è statica. Il principio del privilegio minimo, secondo cui ai soggetti sono garantiti solo gli oggetti strettamente necessari per la sua esecuzione, si sposa perfettamente col concetto di dominio di protezione dinamico.

6. Quali sono i pro e i contro dei due modi di propagazione di un diritto?
   
   Un diritto può esssere propagato per trasferimento o per copia e la differenza è che nel primo caso chi trasferisce perde il diritto. La propagazione per trasferimento permette di rispettare i vincoli legali o tecnici garantendo l'esclusività della licenza o del diritto. La possibilità di copiare un diritto è determinata dal copy flag (*).

7. Sai spiegare il significato dei diritti owner, control, switch?
   
   Chi possiede il diritto owner può assegnare o revocare qualunque diritto di accesso sull'oggetto di cui è proprioetario. Chi possiede il diritto control può revocare qualunque diritto di accesso su qualunque oggetto al soggetto che controlla. Chi possiete il diritto switch su un altro soggetto può commutare il proprio dominio nel dominio del soggetto in questione (permesso anche il ritorno).

8. Quale sarebbe il probelma se la matrice degli accessi non ci fosse?
   
   Non sarebbe possibile verificare se un accesso è consentito o meno, nè cambiare il dominio dinamicamente. Non sarebbe possibile cambiare lo stato di protezione in modo controllato e non sarebbero noti i soggetti e gli oggetti del sistema.

9. Quali sono i pregi e i difetti delle diverse implementazioni della matrice degli accessi in un sistema operativo?
   
   Se si implementa come access control list (ACL) allora è semplice aggiungere o rimuovere oggetti ma è impegnativo aggiungere o togliere soggetti. Nel caso delle Capability List (CL) è il contrario. La revoca dei diritti di accesso può essere generale o selettiva, parziale o totale, temporanea o permanente; con ACL risulta semplice, con CL risulta compessa. La soluzoine ideale è quella mista: la ACL è persistente e quando un soggetto accede ad un oggetto, la si controlla una sola volta e si aggiorna la CL del soggetto, che viene distrutta solo dopo l'ultimo accesso del soggetto.

10. Nel contesto dei sistemi multilivello, quale è il significato di clarance levels, sensitivity levels, security levels e category set?
    
    Rispettivamente sono i livelli per i soggetti, i livelli per gli oggetti, i livelli per la classificazione gerarchica della riservatezza (non classificato, confidenziale, segreto, top secret) e un insieme di categorie dipendenti dall'applicazione in cui i dati sono usati.

11. Perchè il modello Biba è debole all'attacco del cavallo di Troia mentre il modello  Bell-La Padula invece no?
    
    I due modelli sono caratterizzati dalle rispettive proprietà di sicurezza semplice e di integrità *. Nel Bell-La Padula il flusso delle informazioni è verso l'alto dal momento che si può leggere verso il basso e scrivere verso l'alto mentre nel Biba è il contrario. Nel attacco del cavallo di Troia l'attaccante tenta con l'inganno di far eseguire a chi ha più privilegi ed informazioni di lui un programma che scrive queste informazioni in un file non classificato di modo da poterle leggere, ma non funziona sempre; il flusso di informazioni verso l'alto garantisce confidenzialità, mentre il flusso verso il basso garantisce integrità.

12. Perchè tutti i sistemi ad alta sicurezza sono sistemi fidati?
    
    Per i sistemi fidati è possibile definire formalmente dei requisiti di sicurezza e la presenza dei requisiti è necessaria per poterli classificare come "elevata sicurezza".

13. Che legame c'è tra TCB e l'audit file?
    
    Il Trusted Computing Base contiene i privilegi di sicurezza di ogni soggetto e gli attributi di protezione di ogni oggetto, mentre nell'audit file vengono mantenuti gli eventi importanti come i tentativi di violazione della sicurezza o le modifiche autorizzate al TCB stesso.

14. Nei sistemi ad elevata sicurezza, cosa garantisce il reference monitor?
    
    Garantisce: mediazione completa, applicando le regole di sicurezza ad ogni accesso e non solo, isolamento proteggendo reference monitor e base di dati da modifiche non autorizzate, e infine verificabilità delle precedenti.

15. Quali sono le principali categorie dell'Orange Book?
    
    D (Minimal Protection), C (Discretional Protection), B (Mandatory Protection), A (Verified Protection).

## Programmzione Concorrente

1. Cosa si intende con programmazione concorrente?
   
   La programmazione concorrente è l'insieme di delle tecniche, metodologie e strumenti per il supporto all'esecuzione di sistemi software composti da insiemi di attività svolte simultaneamente.

2. Come funziona la classificazione di Flynn?
   
   In un sistema ci può essere parallelismo a livello di istruzioni o parallelismo a livello di dati. I sistemi vengono divisi in SISD(single instruction single data) come ad esmpio i traduzionali computer di Von Neumann, MISD(multiple instruction single data) come i computer in pipeline, SIMD(single instruction multiple data) come array sistolici o processori vettoriali e MIMD(multiple instruction multiple data) come i multicomputers (MPP, COW) e i multiprocessors (UMA, NUMA).

3. Quali sono i diversi tipi di architettura?
   
   Single processor (o monoprocessore), shered-memory multiprocessor (o multiprocessore) e distributed memory system (o sistemi a memoria distribuita). Da notare che i nodi di un sistema a memoria distribuita possono essere monoprocessori o multiprocessori.

4. Quale è la differenza tra UMA e NUMA?
   
   Nei sistemi UMA (UniformMemory Access) il tempo di accesso è uniforme da ogni processore ad ogni locazione di memoria (detti anche SMP), la rete di iinterconnessione è realizzata da un memory bus o da crossbar switch enormalmente il numeor di processori è ridotto (20-30). Nei sistemi NUMA (Non Uniform Acess Memory) il tempo di accesso dipende dalla distanza tra processore e memoria. La rete di interconnessione è un insieme di switch e memorie strutturato ad albero e ogni processore ha memorie più vicine e più lontane. La memoria è quindi organizzata gerarchicamente  per evitare la congestione del bus e normalment il numero dei processori è elevato (anche centinaia). 

5. Quali sono i principali tipi di sistemi a memoria distribuita?
   
   Sono i multicomputer e i network system. I primi sono tightly coupled in cui i processori e la rete sono fisicamente vicini e i secondi sono loosley coupled in cui i nodi sono collegati ad una rete locale o geografica (Ethernet/Internet).

6. A grandi linee come possiamo classificare le applicazioni concorrenti?
   
   In multithreaded, multitasking(sistemi distribuiti) e parallele. I primi sono il classico caso di concorrenza in cui i processi vengono schedulati e eseguiti indipendentemente e sono più dei processori. Nei sistemi distribuiti le componenti (task) vengono eseguite su nodi collegati tramite opportuni mezzi di interconnessione detti canali e possono essere organizzati in diversi modi (client-server, perr to peer, puplisher-substriber). Le applicaizoni parallele risolvono problemi complessi sfrutttando il parallelismo HW eseguendo su appositi modellli paralleli.

7. Quali sono le definizioni di algoritmo, programma, processo, elaboratore ed evento?
   
   L'algoritmo è un procedimento logico che deve essere eseguito per risolvere un determinato problema. Il programma è la descrizione dell'algoritmo mediante un opportuno formalismo (linguaggio di programmazione), che rende possibile l'esecuzione dell'algoritmo da parte di un particolare elaboratore. Il processo è l'insieme ordinato di eventi cui da luogo un elaboratore quando opera sotto il controllo di un programma. Un elaboratore è un entità astratta realizzata in HW e in SW, in grado di eseguire programmi (descritti in un dato linguaggio). L'evento è l'esecuzione di un'operazione tra quelle appartenenti all'insieme che l'elaboratore sa riconoscere ed eseguire; ogni evento determina una transizione di stato dell'elaboratore. In conclusione Più processi possono essere associati allo stesso programma e ciascuno rappresenta l'esecuzione dello stesso codice con dati in ingresso diversi.

8. Quali sono le possibili interazioni tra processi?
   
   Ammesso che i processi siano tra loro interagenti e non indipendenti; possono essere in cooperazione, interazione prevista e desiderata, competizione, interazione prevedibile ma non desiderata e non grave, o interferenza, non prevista, non desiderata e quasi certamente non inevitabile.

9. Quali soluzioni forniscono i modelli linguistici per la concorrenza?
   
   SOno il fork/join e il cobegin/coend. Con fork si può creare e attivare un processo che esegue in parallelo col chiamante, mentre join consente di determinare quando un processo creato tramite la fork ha terminato il suo compito, sincronizzandosi con tale evento: c'è quindi la possibilità di denotare in modo esplicito il processo sulla cui terminazione ci si vuole sincronizzare. In ogni blocco cobegin/coend invece ogni istruzione viene eseguita in parallelo e in ogni istruzione si può inserire un ulteriore blocco. Da mìnotare che non tutti i grafici di precedenza non possono essere espressi tramite cobegin/coend.

10. Quali sono le proprietà che caratterizzano i processi sequenziali e quali sono le proprietà che caratterizzano i processi concorrenti?
    
    I nodi del grafo di precedenza rappresentano i singoli eventi del processo mentre gli archi orientati identificano le precedenza temporali tra tali eventi. Nei processi sequenziali il grafo di processi è ad ordinamento totale mentre nei programmmi concorrenti il grafo di precedenza è ad ordinamento parziale.

11. Quali sono le proprietà che caratterizzano i programmi sequenziali e quali sono le proprietà che caratterizzano i programmi concorrenti?
    
    La traccia dell'esecuzione (o storia) è la sequenza degli stati attraversati dal sistema di elaborazione durante l'esecuzione del programma (variabili esplicite ed implicite). Nei programmi sequenziali la traccia è costante ad ogni esecuzione, mentre nei programmi concorrenti la traccia può variare ad ogni esecuzione (non determinismo).

12. Quali esempio possono aiutare a capire cosa si intende con safety e liveness?
    
    Sono due proprietà dei programmi. La safety garantisce che durante l'esecuzione non si entrerà mai in uno stato errato (stato in cui le variabili assumono valori indesiderati), mentre la liveness garantisce che durante l'esecuzione, prima o poi si entrerà in uno stato corretto (stato in cui le variabili assumono valori desiderati). Alcuni esempi di safety sono la correttezza dello stato finale, la mutua esclusione dell'accesso a risorse condivise e l'assenza di deadlock. Alcuni esempi di liveness sono la terminazione e l'assenza di starvation.

## Modello a Memoria Comune

1. Quali sono le premesse del sistema a memoria comune?
   
   L'ambiente è globale e il sistema è visto come un insieme di processi attivi e risorse o oggetti passive/i. Le interazioni tra processi possono esssere di competizione o di cooperazione. Il modello rappresenta la maturale astrazione del funzionamento di un sistema in multiprogrammaizone costruito da uno o più processori che hanno accesso ad una memoria comune. Ad ogni processore può essere associata una memoria privata ma ogni interaizone avviene tramite oggetti in memoria comune. 

2. Cosa si intende con risorsa?
   
   E' un qualunque oggetto, fisico o logico, di cui un processo necessità per portare a termine il suo compito. Le risorse sono raggruppate in classi; una classe identifica l'insieme di tutte e sole le operaizoni che un processo può eseguire per operare su risorse di quella classe. Ogni risorsa si identifica con una struttura dati allocata in memoria comune (anche risorse fisiche: descrittori del dispositivo).

3. Cosa è il gestore di una risorsa e quali sono i suoi compiti?
   
   Per ogni risprsa R, il suo gestore definisce, in ogni istante t, l'insieme SR(t) dei processi che, in tale istante, hanno il diritto di operare su R. Si possono classificare in dedicate o condivise, alllocate staticamente o dinamicamente e private o comuni. Il gestore mantiene aggornato l'insieme SR(t) e cioè lo stato di allocazione di una risorsa, fornisce i meccanismi che un processo può utilizzare per acquisire il diritto di operare sulla risorsa, entrando a far parte dell'insieme SR(t), e per rilasciare tale diritto quando non è più necessario e implementare la strategia di allocaizone della risorsa e cioè definire quando, a chi, e per quanto tempo allocare la risorsa. Il gestore di una risorsa è una risorsa condivisa nel modello a memoria comune.

4. Come si accede ad una risorsa allocata dinamicamente?
   
   E' necessario prevedere un gestore che implementa le funzioni di richiesta e rilascio della risorsa, rispettivamente prima e dopo aver eseguito le operazioni che necessitavano della risorsa.

5. Come si accede ad una risorsa condivisa?
   
   E'necessario assocurarsi che gli accessi alla risorsa avvengano in modo non divisibile: le funzioni di accesso alla risorsa devono quindi essere programmate come classe di sezioni critiche (9utilizzando i meccanismi di sincronizzazione).

6. Cosa si intende con  regione critica condizionale?
   
   E' un formalismo che permette di esprimere la specifica di un qualunque vincolo di sincronizzaizone. Il corpo della regione rappresenta un'operaizone da eseguire sulla risorsa condivisa e quindi costituise una sezione critica che deve essere eseguita in mutua esclusione con le altre operazioni definite su R.

7. Quali sono i casi particolari di regioni critiche?
   
   ```c
   //semplice mutua esclusione
   region R << S; >>
   
   //semplice vincolo di sincronizzazione
   region R << when(C) >>
   
   // specifica dello stato della risorsa per poter eseguire l'operazione
   region R << when(C) S; >>
   ```

8. Cosa si intende con mutua esclusione e cosa sono le sezioni critiche?
   
   La condizione tale per cui le operazione con le quali i processi accedono alle variabili comuni non si sovrappongono nel tempo. La sequenza di istruzioni con le quali un processo accede e modifica un insieme di variabili comuni prende il nome di sezione critica. Ad un insieme di variabili comuni ppossono essere associate una sola sezione critica (usata da tutti i processi) o più sezioni critiche (classe di sezioni critiche). La regola di mutua esclusione stabilisce che: sezioni critiche appartenenti alla stessa classe devono escludersi mutualmente nel tempo oppure ad ogni istante può essere "in esecuzione" al più una sezione critica di ogni classe. La mutua esclusione prevede un prologo e un epilogo nl codice durante i quali la risorsa viene acquisita e liberata.

9. Quali sono le possibili soluzioni al problema della mutua esclusione?
   
   Possono essere algoritmiche: come gli algoritmi di dekker e peterson o l'algoritmo del fornaio. Queste soluzioni implementano una attesa attiva. Possono essere HW based con disabilitazione delle interruzioni o con lock e unlock. Infine possono esssere di natura SW, realizzate coi semafori e coi derivati, in cui si sospendono effettivamente i processi in attesa.

## Nucleo Sistema Memoria Comune

1. Come si mettono in evidenza le proprietà logiche di comunicazone e sincronizzazione tra processi senza doversi preoccupare degli aspetti implementativi delle particolari caratteristiche del processore fisico?
   
   In un sistema multiprogrammato vengono offerte tante unità di elaborazione astratte (macchine virtuali) quanti sono  i processi e ogni macchina possiede come set di istruzioni elementari quelle corrispondenti all'unità centrale reale più le istruzioni relative alla crezione ed eliminazione dei processi, al meccanismo con i dispositivi di I/O visti come processori esterni).

2. Di cosa si occupa il kernel?
   
   Si chiama kernel il modulo (o insieme di funzioni) realizzato in SW, HW o FW che supporta il concetto di processo e realizza gli strumenti necessari per la gestione dei processi. Costituisce il livello più interno di un qualunque sistema basato su processi concorrenti. E' il livello più elementare di un sistema operativo multiprogrammato e fornisce il supporto a tempo di esecuzione per un linguaggio per la programmazione concorrente. Il nucleo è il solo modulo conscio della esistenza delle interruzioni: i processi vengono sospesi grazie a specifiche primitive nel nucleo e viene risvegliato dopo che il nucleo ha ricevuto il segnale di interruzione da dispositivo e poi CPU. La gestione delle interruzioni è quindi invisibile ai processi ed ha come unico effetto rilevabile di rallentare la loro esecuzione sulle rispettive macchine virtuali.

3. Cosa si intende con contessto,  salvataggio e ripristino del contesto?
   
   Il contesto di un processo è l'insieme delle informazioni contenute nei registri del processore, quando esso opera sotto il controllo del processo. Il salvataggio avviene quando il processo perde il controllo del processore, e il contenuto dei registri del processore (ovvero i contesto) viene salvato in una struttura dati associata al processo, chiamata descrittore. Infine il ripristino avviene quando un processo viene schedulato e i valori salvati nel suo descrittore vengono caricati nei registri del processore.

4. Quali sono le funzioni del nucleo?
   
   Il compito fndamentale del nucleo di un sistema a processi è gestire le transizioni di stato dei processi. In particolare deve gestire il salvataggio e il ripristino dei contesti dei processi, scegliere a quale tra i processi pronti assegnare l'unità di elaborazione (scheduling CPU), gestire le interruzioni dei dispositivi esterni e infine realizzare i meccanismi di sincronizzazione dei processi.

5. Quali sono le caratterictiche desiderabili del nucleo?
   
   Sono efficienza ottenuta anche con soluzioni HW e mediante microprogrammi (condiziona l'intera struttura), dimensioni ridotte anche frazie alla semplicità delle funzioni, e separazione tra politiche e meccanismi.

6. Cosa è e cosa contiene un descrittore di un processo?
   
   Contiene l'identificatore del processo, lo stato del processo e la modalità di servizio (parametri di scheduling), ma anche il contesto del processo in temrini di contatore del programma, registri e indirizzo all'area di memoria privata del processo e infine contiene il riferimento a code (puntatore ad indirizzo successivo).

7. Cosa sono le code dei processi pronti?
   
   Ci possono essere una o più code dei processi pronti( scheduling con priorità). Dato che non è detto che vi sia sempre almeno un processo pronto , la coda contiene sempre almeno un processo fittizzio (dummy process) che va in esecuzione solo quando tutte le altre code sono vuote. Ha priorità minima ed è sempre pronto. Ci sono puntatori ad entrambe le estremità del ìla coda per facilitare le due operazioni di inserimento e prelievo dei descrittori.

8. Cosa è la coda dei descrittori liberi?
   
   E' la coda nella quale sono concatenati i descrittori disponibili per la creazione di nuovi processi e nella quale sono re-inseriti i descrittori dei processi terminati.

9. Come fa il nucleo a sapere quale processo è in esecuzione?
   
   Questa informazione, rappresentata dall'indice del processo, viene contenuta in un aparticolare variabile del nucleo (spesso, un registro del processore).

10. In che senso il nucleo è strutturato in due livelli?
    
    Il livello superiore contiene tutte le informazione direttamente utilizzabili dai processi sia interni isia esterni (dispositivi I/O); come le primitive per la creazione, eliminazione e sincronizzazione dei processi e le funzioni di risposta ai segnali di interruzione. Nel livello inferiore sono relizzate le funzionalità di cambio di contesto: salvataggio dello stato, sceduling o assegnazione della cpu e ripristino dello stato.

11. Quali sono le accortezze a livello di sicurezza di nucleo?
    
    Le funzioni del nucleo sono le sole che possono operare s lle strutture dati che rappresentano lo stato del sistema e che possono utilizzare istruzioni privilegiate (comunicaizoni dispositivi e gestione interruzioni). Pertanto nucleo e processi eseguono in ambienti separati: il nucleo in kernel mode nel ring 0 e i processi in user mode nei ring positivi.

12. Come si realizza il cambio ring e di privilegio?
    
    Nel caso di chiamate da processi esterni (dispositivi) si utilizza il meccanismo di risposta al segnale di interruzione (interruzione esterne), mentre nel caso di funzioni chiamate da processi interni, il passaggio è ottenuto tramite l'esecuzione di system calls o chiamate al supervisore (interruzioni esterne).

13. Come viene gestito il temporizzatore?
    
    Per consentire la modalità di servizio a divisione di tempo è necessario che il nucleo gestisca un dispositivo temporizzatore tramite un'apposita procedura che ad intervalli di tempo fissati, provveda a sospendere il processo in esecuzione ed assegnare l'unità di elaborazione ad un altro processo (cambio di contesto).

14. Come si può realizzare un semaforo in un sistema monoprocessore?
    
    Può essere implementato tramite una variabile intera che rappresenta il suo valore non negativo e da una coda di descrittori di processi in attesa sul semaforo. La coda viene gestita con un apolitica FIFO: i processi risultano oridnati secondo il loro tempo di arrivo nella coda associata al semaforo. IL descrittore di un processo viene inserito nella coda del semaforo come conseguenza di una primitiva p non passante; e poi prelevato per effetto di una v.

15. Quali sono i due modelli di organizzaizone interna possibili nei sistemi operativi multiprocessore?
    
    Se il sistema operativo esegue su un architettura multiprocessore allora deve gestire una molteplicità di CPU, ogniuna delle quali può accedere alla stessa memoria condivisa. I due modelli sono: il modello SMP (Simmetric Multi Processing) e il modello a nuclei distinti.

16. Cosa caratterizza il modello SMP?
    
    Nel modello SMP c'è un unica copia del nucleo del sistema operativo allocata nella memoria comuneche si occupa della gestione di tutte le risorse disponibili, comprese le CPU. Ogni processo può essere allocato su qualunque CPU ed è possibile che processi che eseguono su CPU diverse richiedano contemporaneamente funzioni del nucleo (System Call), e questa competizione tra le CPU per eseguire nel nucleo necessità di sincronizzazione.

17. Come si può sincronizzare l'accesso al nucleo delle diverse CPU?
    
    Con la soluzione ad un solo lock, si garantisce la mutua esclusione in modo semplice, ma si limita il grado di parallelismo, escludento a priori ogni possibilità di esecuzione contemporanea. Con la soluzione a più lock invece si individuano all'interno del nucleo diverse classi  di sezioni critiche, ogniuna associata ad una struttura dai separata e sufficientemente indipendente dalle altre, e ad ogniuna viene associato un lock distinto.

18. E' sempre la scelta migliore, in sistemi SMP, schedulare ogni processo su uno qualunque dei processori, per massimizzare il bilanciamento del carico?
    
    Può essere conveniente assegnare un processo ad un determinato processore perchè possono accedere velocemente alla loro memoria privata e potrebbe convenire schedulare il processo sul processore la cui memoria privata già contiene il suo codice o perchè in sistemi NUMA l'accesso alla memoria più vicina è più rapido e conviene quindi schedulare il processo sul processore più vicino alla memoria dove è allocato il suo spazio di indirizzamento. I processori hanno anche memoria cache e per ridurre l' overhead conviene schedulare un processo nel processore sul quale era stato precedentemente eseguito. Inoltre la scelta sulla politica di scheduling ha un impatto sul numero di code da gestire:  da una in tutto a una per CPU (nodo in questo caso).

19. Cosa caratterizza il modello a nuclei distinti?
    
    Ogni nucleo è  dedicato alla gestione di una diversa specifica CPU. L'assunzione di base è che l'insieme dei processi che esseguiranno nel sistema sia partizionabile in sottoinsieme (nodi virtuali) lascamente connessi, cioè con un ridotto numero di interazioni reciproche. Ciascun nodo virtuale è assegnato ad un nodo fisico e tutte le strutture dati relative al nodo virtuale come i descrittori dei processi,  i semafori locali e le codee dei processi pronti, vengono allocate sulla memoria privata del nodo fisico.In questo modo tutte le interazioni solcali al nodo virtuale possono avvenire indipendentemente e concorrentemente a quelle di altri nodi virtuali, facendo riferimento al nucleo del nodo. E' fondamentale che solo le interazioni tra processi appartenenti a  nodi diversi utilizzanno la memoria comune.

20. Quali sono i vantaggi di SMP e quali quelli di nuclei distiniti?
    
    Il primo permette un gestione ottimale delle risorse (miglior bilanciamento) mentre il secondoaumenta il grado di parallelismo (disaccoppiamento più basso e maggiore scalabilità).

21. Quale è la differenza tra semafori privati e semafori condivisi?
    
    In un sistema multiprogrammato multiprocessore con modello a nuclei distinti i semafori privati riguardano un singolo nodo e vengono realizzati come nel caso monoprocessore, mentre i semafori condivisi tra i nodi sono utilizzati da processi appartenenti a nodi virtuali diversi. La memoria comune dovrà contenere tutte le informazioni relative ai semafori condivisi.

22. Cosa è il rappresentante del processo?
    
    E' l'insieme minimo di informazioni sufficienti per identificare sia il nodo fisico su cui il processo opera, sia il descrittore contenuto nella memoria provata del processore.

23. Come vengono realizzati i semafori condivisi?
    
    Ogni semaforo condiviso viene rappresentato nella memmoria comune da un interno non negativo e l'accesso al semaforo e protetto da un lock, anch'esso in memoria comune. Per ogni semaforo condiviso sono mantenute varie code: su ogni nodo è mantenuta una coda locale contenente i descrittori dei processi locali sospesi, e un ultima coda globale dei rappresentati di tutti i processi sospesi, accessibile solo dal nucleo.

24. Perchè si parla di comunicazione tempestiva tra i nuclei?
    
    La comunicazione tra i nuclei deve essere tempestiva. Il nucleo che esegue la v deve interrospere qualsiasi cosa stia facendo per mandare un segnale di interruzioen al nucleo che si era sospeso con una p. Utilizza un buffer in memoria comune dove è presente la coda dei rappresentanti di S. Estrae quindi il descrittore dalla coda locale portando il processo nello stato ready.

## Modello a Scambio di Messaggi

1. Quali sono le caratteristiche del modello a scambio di messaggi?
   
   Ogni processo può accedere esclusivamente alle risorse allocate nella propria memoria locale. Ogni risorsa del sistema è accessibile direttamente ad un solo processo (gestore della risorsa). Se una risorsa è necessaria a più processi applicativi, ciascuno di questi processi clienti dovrà delegare l'unico processo che può operare sulla risorsa, processo gestore, all'esecuzione delle operazioni richieste. il gestore della risorsa coincide quindi con un processo server.

2. Cosa è un canale di comunicazione e da quali parametri è caratterizzato?
   
   E' un collegamento logico mediante il quale due o più processi comunicano. Il nucleo della macchina concorrente realizza l'astrazione canale come meccanismo primitivo per lo scambio di informazioni. Poi il linguaggio di programmazione offre gli strumenti linguistici di alto livello per specificare i canali di comunicazione e utilizzarli per esprimere le interazioni tra i processi. I parametri che caratterizzano il concetto di canale sono: la direzione del flusso dei dati che un canale può trasferire, la designazione del canale e dei processi origine e destinatario e il tipo di sincronizzaizone.

3. Come si possono classificare i canali?
   
   Possono essere monodirezionali o bidirezionali, possono essere link, port o mailbox,  possono favorire una comunicazione asincrona, sincrona o con sincronizzazione estesa.

4. Quali sono le caratteristiche della comunicazione asincrona?
   
   Il processo mittente continua la sua esecuzione immediatamente dopo l'invio del messaggio: così facendo le informazioni ricevute non possono essere attribuite allo stato attuale del mittente. L'invio del messaggio non è un punto di sincronizzazione. C'è carenza espressiva e difficoltà nella verifica dei programmmi ma l'assenza di vincoli di sincronizzaizone favorisce il grado di concorrenza. Da un punto di vista realizzativo, sarebbe necessario un buffer di capacità illimitata, ma ovviamente un l'invio su un canale pieno è sospensivo.

5. Quali sono le caratteristiche della comunicazione sincrona semplice?
   
   Il primo processo ad eseguire invio o ricezine si sospende in attesa che l'altro sia pronto ad eseguire l'operazione corrispondente: così facendo l'inzio è un punto di sincronizzaizone e ogni messaggio ricevuto contiene informazioni attribuibili allo stato attuale del processo mittente; sempificando scrittura e verifica dei programmi. Non è necessario bufferizzare.

6. Quali sono le principali differenze tra le caratteristiche tra le comunicazioni asincrona e sincrona?
   
   La prima consente maggiore parallelismo, a scabilto della semploicità d'uso mentre le seconda è più espressiva ma le sospensioni possono peggiorare le performance.

7. Come si può realizzare una comunicazione sincrona usando primitive asincrone?
   
   Mandando un messaggio di ack dopo dal processo ricevente al mittente.

8. Quali sono le caratteristiche della comunicazione con sincronizzazione estesa?
   
   Si parte dal presupposto che ogni messaggio inviato rappresenta una richiesta al destinatario dell'esecuzione di una certa azione. Il processo mittente rimane in atetsa fino a che il ricevente non ha terminato di svolgere l'azione richiesta. Il punto di sincronizzaizone semplifica la verifica dei programmi anche se c'è riduzione del parallelismo. Il modello rimane client-servitore e c'è analogia semantica con la chiamata di procedura. Da notare  che in generale il server ha la possibilità di offrire diversi servizi e di conseguenza gestire diversi canali di ingresso. Nella sincronizzazione estesa quindi la recive non è bloccante anche se queso può introdurre attesa attiva.

9. Quale è il meccanismo di ricezione ideale nella comunicaione con sincronizzaizone estesa?
   
   Deve consentire al processo server di verificare contemporaneamente la dispobilità di messaggi su più canali, abilitare la ricezione si un messaggio da un qualunque canale contenente messaggi e quando tutti i canali sono vuoti, bloaccar il processo in attesa che arrivi il  messaggio, qualunque sia il canale su cui arriva. Questo meccanismo è relizzabile tramite comandi id guardia?

10. Cosa si intende con comando con guardia?
    
    ```c 
    <espressione_booleana>; <recive> -> <istruzione>
    ```
    
    La valutazione della guardia può fornire tre diversi valori: guardia fallita, guardia ritardata e guardia valida. Rispettivamente quando l'espressione booleana è falsa, quando l'espressione booleana ha valore true ma nel canale su cui viene eseguita non ci sono messaggi e quando l'espressione booleana ha valore true e nel canale c'è almeno un mesaggio (le receive esegue senza ritardi).

11. Cosa si intende con comando con guardia alternativo?
    
    ```c
    select {
        [] <guardia_1> -> <istruzione_1>;
        ...
        [] <guardia_n> -> <istruzione_n>;
    }
    ```
    
    Il comando con guardia alternativo (select) racchiude un numero arbitrario di comandi con guardia semplice. Se una o più guardie sono valide viene scelto in maniera on deterministica uno dei rami con la guardia valida, se tutte le guardie non fallite sono ritardate il processo in esecuzione si sospende, mentre se tutte le guardie sono fallite allora il comando termina.

## Sincronizzazione Estesa

## Implementazioni Concorrenza

1. Cosa è una goroutine?
   
   E' l'unità di esecuzione concorrente: una funziona che esegue concorrentemente ad altre nello stesso spazio di indirizzamento. Un programma go è costituito da una o più goroutine concorrenti. In generale ci sono più goroutine per thread di SO e le lingole goroutine possono essere estremamente leggere. La costante di ambiente GOMAXPROCS determina il numero di thread.
   
   $$
   GOMAXPROCS = \frac{goroutine}{thread}
   $$

2. Quyali filosofia di fondo di go determina le scelte nell'implementazione delle comunicazioni tra goroutine?
   
   "Do not communicate by sharing memory. Instead, share memory by
   communicating". I canali sono oggetti di prima classe e possono essere sia simmetrici che asimmetrici, 1-1, 1-m, m-m, la comunnicazione può essere sincrona o asincrona e sia bidirezionale che monodirezionale.

3. Quali sono le particolarità del costrutto select in go?
   
   E' una istruzione analogoa al comando con guardia alternativo e la selezione è non deterministica tra i rami di guardia valida. Le guardie sono tutte semolici e il linguaggio non prevede la guardia logica: le guardi possono essere solo valide o ritardate. A questo si può in qualche modo ovviare utilizzando when:
   
   ```go
   func when(b bool, c chan int) chan int {
       if !b return nil
       return c
   }
   ```

4. Cosa si intende per unità di concorrenza in ada?
   
   L'unità di concorrenza in ada è il task e il lingiaggio adotta come metodo di interazione tra i processi il randevous esteso con comunicazione asimmetrica. Se un processo prende il nome di task allora i programmi concorrenti prendono il nome di procedure. Il messanismo di comunicazione prevede che ogni task possa definire delle entry visibili agli altri task e che questi ultimi possono chiamare dall'esterno. Un task che chiama una entry di un altro si mette in attesa del termine. Da notare che ad una stessa entry possono esserre associate più accept e ad esse possono corrispondere azioni diverse a seconda dalla fase di esecuzione del task. Un Task (server) può esporre più operazioni e accettare le richieste attraverso il comando con guardia alternativo select.

5. Quali sono gli obbiettivi della progettazione di ada e come questi influenzano le sue caratteristiche?
   
   E' progettato per ridurre i costi di sviluppo e manutenzione, prevenire i bugs rilevandoli il prima possibile (tempo di compilazione), favorire riutilizzo e sviluppo in team e semplificare la manutenzione in termini di leggibilità e autodocumentazione. Ada è  quindi un linguaggio fortemente e statcamente tipato (due  tipi diversi non possono essere confrontati). La keyword access rende i puntatori più semplici da usare e ci sono attributi che arricchiscono i tipi in relazione al loro dominio. Complessivamente ada soddisfa i più alti standard di sicurezza del software.

## Algoritmi Sincronizzazione Distribuiti

1. Quale è il legame tra modello a scambio di messaggi e sistema distribuito?
   
   Il modello a scambio di messaggi è la naturale astrazione di un sistema
   distribuito, nel quale processi distinti eseguono su nodi fisicamente separati,
   collegati tra di loro attraverso una rete. Nel sistema distribuito non ci sono risorse condivise e non c'è un glock globale. Si passa dla concorrente al distribuito e si introduce la possibilità di malfunzionamenti indipendenti.

2. Quali sono le proprietà desiderabili nel distribuito?
   
   Salabilità e tolleranza ai guasti. La scalabilità garantisce che nell’applicazione distribuita le prestazioni aumentano al crescere del numero di nodi utilizzati. La tolleranza ai guasti garantisce che l’applicazione è in grado di funzionare anche in
   presenza di guasti (es. crash di un nodo, problemi sulla rete, ecc.).

3. Come si misurano le prestazioni di un sistema distribuito?
   
   Lo speedup e l'efficienza sono indicatori usati e idelmente hanno rispettivamente valore $n$ ed $1$.
   
   $$
   Speedup(n)=\frac{Tempo(1)}{Tempo(n)} \newline
Efficienza(n)=\frac{Speedup(n)}{n}
   $$

4. Quali spossono essere i tipi di guasto e cosa si intende come si implementa la tolleranza ai guasti?
   
   I guasti possono essere transienti, intermittenti o persistenti. Si possono implementare tecniche di ridondanza e sono necessari meccanismi di rilevazione (fault detection) e di riprestino (recovery).

5. Come è organizzata la gestione del tempo in un sistema distribuito?
   
   Ogni nodo è dotato di un suo orologio. Se gli orologi locali di due nodi non sono sincronizzati, è possibile che due eveti in due nodi diversi siano associati a due istanti temporali che fanno semprare che uno sia precedente all'altro quando in realtà è il contrario. Nel caos isa necessario un riferimento temporale unico e si può usare un orologio universale fisico (algoritmi di Berkley e Cristian) o un orologio logico; che permette di associare ad ogni eventi un istant e logico (timestamp) la cui relazione coi timestamp di altri eventi sia coerente con l'ordine in cui essi si verificano.

6. Come si possono caratterizzare i eventi concorrenti?
   
   E' possibile definire la relazione di precedenza tra eventi (Heppende Before, $->$). Data una coppia di eventi $a$ e $b$, allora se $a$ precede $b$ si indica con $a->b$, se invece $a$ precede $a$, si indica con $b->a$, e infine se nessuna delle due è valida allora sono concorrenti.

7. Come funziona l'algoritmo di Lamport?
   
   Ogni processo mantiene localmente un contatore del tempo logico e ogni nuovo evento all'interno del processo provoca un incremento del valore del contatore. Inoltre Ogni volta che il processo vuole inviare un messaggio, dopo aver incrementato il contatore, quest'ultimo viene allegato al messaggio. Quando si riceve un messaggio si assegna al proprio contatore il massimo tra il valore del contatore allegato e il valore attuale, e successivamene lo si incrementa. Usualmente implementato dal middleware che interfaccia i precessi alla rete.

8. Come possiamo classificare le soluzoini volte a garantire che due o più processi non possano eseguire contemporaneamente alcune prestabilite attività?
   
   Per risolvere la mutua esclusione distribuita si può ricorrere a soluzoni token-based o permission-based e queste ultime a loro volta possono essere centralizzate oppure decentralizate.

9. Quali sono vantaggi e svantaggi della soluzioni permission-based centralizzata?
   
   L'algoritmo è equo quindi non c'è starvation. E' anche semplice perchè prevede solo 3 messaggi: richiesta, autorizzazione e rilascio. Purtroppo non è né scalabile né tollerante ai guasti. Un procesos che non riceve autorizzazione non può sapere se non è stata concessa o se il gestore è guasto.

10. Quale dell'algoritmo permission-based è scalabile e cosa possiamo dire sulla sua tolleranza ai guasti?
    
    L'algoritmo Ricard-Agrawala prevede $2*(N-1)$ messaggi per ogni sezione critica. Inoltre la tolleranza ai guasti è pessima perchè è sufficiente che ci sia un guasto su un nodo e nessuno sarà più autorizzato a fare nulla. Si può fare una piccola modifica introducento i messaggi di accesso negato. Una volta ricevuti ci si mette di nuovo in attesa ma si può impostare un timeout per rilevare i guasti e eventualmente escluderlo dal gruppo.

11. Quaii osno gli aspetti cruciali dell'algoritmo token ring?
    
    L'intero sistema è costruito da un insieme di processi in competizione collegati tra loro in una topologia ad anello e i processi conoscono i loro vicini. Un messaggio, detto token, circola attraverso l'anello, nel verso relativo all'ordine dei processi nella topologia. Chi deve eseguire la sezione critica tiene il token fino al rilascio. E' scalabile ma ci possono essere moltissimi messaggi per ogni sezione critica. Inoltre ci sono N punti di fallimento e un crash può fare perdere il token.

12. Cosa è un algoritmo di elezione?
    
    In alcuni algoritmi è necessario che un processo svolga il ruolo di coordinatore. La disegazione può essere statica o dinamica. Nel secondo caso, per scegliere, si usa un algoritmo di elezione.

13. Quali sono le differenze tra gli algoritmi di elezione Bully ed ad Anello?
    
    Nel primo il processo che avvia l'elezione invia l'aposito messaggio a tutti i processi con l'id più alto del suo e chi non è guasto risponde positivamente. Poi se c'è stata almeno una risposta, tutti quelli che hanno risposto avviano un elezione a loro volta. Nel secondo caso quando un processo si rende conto che il coordinatore è guasto inizia un elezione mandando un messaggio col priprio id e chi lo riceve aggiunge il proprio id e lo riinvia a sua volta. Quando si riceve un messaggio con proprio id si cambia il contenuto del messaggio e si invia ora l'identità del nuovo coordinatore, ovvero l'id più alto tra tutti.

## Introduzione HCP

1. Quali sono le differenze tra il calcolo concorrente e il calcolo parallelo?
   
   In entrambi i casi si da luogo ad un insieme di attività. Sono concorrenti se sono contemporaneamente in progress, ovvero iniziate ma non temrinate, sono parallele se effettivamente le attività multiple eseguono in contemporanea.
   
   Nel primo caso il numero di processori è maggiore del numero delle CPU, nel secondo invece no. 

2. Perchè si esegue in parallelo?
   
   Per aumentare le performance in temini di complessità dei problemi che si possono risolvere e di tempo necessario.

3. Cosa ci dice la Legge di Moore sull'evoluzione dei sistemi di calcolo nel tempo?
   
   Fino ai primi anni 2000 l'evoluzione dei sistemi di calcolo seguivano un andamento preciso: in numero di transistori in ogni 18 mesi. Quando poi si sono raggiunti i limiti fisici legati all' effeto Joule e non è stato più possibile aumentare la frequenza di clock è stato necessario aumentare la capacità di calcolo a parità di frequenza. Il parallelismo in questo senso è diventato una forma di accellerazione dell'hardware.

4. Cosa si intende con Von Neumann Bottleneck?
   
   La velocità di fetching di istruzioni e dati diepden dalla velocità di trasmissione del Bus è una limitazione della velcità di esecuzione. Il modello di Von Neumann è stato quindi esteso con l'introduzione di memorie cache e di paralllelismo di basso livello, come Instruction-level parallelism (ILP) e HW multithreading (TLP). Il modello Von Neumann esteso è trasparente per lo sviluppatore.

5. La cache che tipo di memoria è?
   
   E' una memoria associativa ad accesso veloce e di capacità limitata che risiede sul chip del processore e si colooca ad un livello intermedio tra i registri e la memoria centrale. Viene gestita con criteri sul principio di località spaziale e temporale (cache hit/miss ed hit-rate).

6. Cosa si intende con parallelismo a livello di istruzione?
   
   L'esecuzione di ogni istruzione viene attuata attraverso una sequenza di fasi. Ogni fase può essere affidata ad un unità funzionale indipendente che opera in parallelo alle altre. Si possono mettere in pipelining collegando tutte le unità funzionali tra loro eseguendo fasi diverse di istruzioni diverse in parallelo. In alternativa ci possono essere più istanze di ogni unità funzionale.

7. Cosa si intende con hardware multi-threading?
   
   Permette a due thread di condividere la stessa CPU (core), utilizzando una tecnica di sovrapposizione. Ciò è reso possibile dalla duplicazione dei registri che mantengono lo stato di ogni thread (PC, IR, ecc) e da un meccanismo HW che implementa il context switch tra un thread ed un altro in modo molto efficiente. Sono possibili 2 approcci: a grana fine e a grana grossa. Nel primo caso viene eseguito context switch dopo ogni istruzione e nel secondo viene eseguito context switch quando il thread corrente è in una situaizone di attesa.

8. Come si realizza la parallelizzazione esplicita?
   
   SI possono usare 2 modelli: scambio di messaggi (MPI) o memoria condivisa (OpenMP). Il parallelismo si ottiene distribuendo task diversi a processi diversi; ogni processo è assegnato a una CPU a sua completa disposizione. Normalmente su utilizza il paradigma SPMD (single program multiple data) sfruttando il branching condizionale. Solo pochi programmi sono embarassingly parallel, nella maggior parte dei casi le iterazioni non sono indipendenti tra loro ed è necessaria che i processi siano sincronizzati. Quindi prima s i divide il lavoro e ppoi ci si occupa di sincronizzazione e comunicazione.

9. Cosa sono i petaFLOPS e gli exaFLOPS?
   
   Sono unità di misura delle prestazioni ovvero i floating point operations per second. Corrispondono rispettivamente a $10^{15}$ e a $10^{18}$ FLOPS.

10. Quale è il legame tra speedup ed efficienza?
    
    Lo speedup misura quanto è più veloce la versione parallela rispetto alla versione sequenziale, ovvero il guadagno della parallelizzazione e nel caso ideale vale 1. Nei casi non ideali c'è overhead dovuto alla creazione e allocazione dei processi, alla cominicazione e alla sincronizzazione e anche alla distribuzione non bilanciata del lavoro. L'efficienza serve a misurare la scalabilità: quanto più rimane costate tanto più un programma è scalabile.

11. Quando si misura la scalabilità quali opzioni si hanno?
    
    Si può misurare la scalabilità strong, ovvero quanto si può guadagnare nella soluzione di uno stesso problema di dimensione fissata aumentando il numero di processori, oppure la scalabilità weak, ovvero se è possibile risolvere lo stesso problema di dimensioni maggiori nello stesso tempo.

12. Cosa dice la Legge di Amdahl e cosa ci dice sulla saclabilità strong?
    
    Partendo dalla premessa che non tutto un programma è parallelizzabile, si calcola lo speedup in funzione di $r$ ovvero della frazione del tempo totale di esecuzione spesa nella parte non parallelizzabile e si ottinene $\lim_{p \to \infty} S = \frac{1}{r}$ ; ovvero se $r$ e diverso da 0, lo speedup non può crescere all'infinito. Similmente vale $\lim_{p \to \infty} E = 0$ il che ci conferma che solo nel caso ideale al crescere di p, l'efficienza si mantiene costante.

13. Come si valuta la scalabilità weak?
    
    Si usano efficienza scalata e speedup scalato:
    
    $$
    E_s(p, N) = \frac{T(N, 1)}{T(pN, p)} \newline
S_s(p, N) = E_sp
    $$

14. Cosa implica la Legge di Gustafson?
    
    La Legge di Gustafson presuppone che il probelma sia di dimensione variabile e si concentra sulla scalabilità weak. implica che assegnando ad ogni processore un workload costante (1-r), lo speedup cresce linearmente con il numero dei processori.
    
    $$
    S(p, pN) = \frac{T(1, pN)}{T(p, pN)} \newline
da \space cui \newline
S(p, pN) = r + (1 - r)p
    $$

## Programmazione Parallela con MPI

1. In quale ambito lo standard è rappreesentato dalle librerie MPI?
   
   Se i nodi dell'architettura non condividono memoria e lo sviluppo dei programmi paralleli si fonda sul modello a scambio di  messaggi (esempio Cluster HPC).

2. Quali sono le caratteristiche principali dello standard Message Passing Interface?
   
   E' basato sul paradigma SPMD con molteplici istanze dello stesso programma, ogniuna in esecuzione contemporanea su un nodo distinto. Ogni istanza rappresenta un processo MPI. Offre un ricco set di funzioni per esprimere comunicazione tra processi sia punto-punto che collective, con semantiche sia sincrone che asincrone. Offre inoltre potenti strumenti per data partitioning e data collecting e gestisce i processi in maniera statica e implicita definendo il grado di parallelismo a tempo di caricamento.

3. Cosa sono i comunicator?
   
   Sono astrazioni che definiscono un dominio di comunicazione, ovvero un insieme di processi che possono comunicare tra loro; due processi possono scambiarsi messaggi solo se appartengono allo stesso comunicator. Esiste un comunicator di default detto MPI_COMM_WORLD. A partire da questo è possibile crearne di altri.

4. Quali sono i limiti della soluzioni centralizzate e cosa si può fare per superarli?
   
   Al crescere del numero dei nodi il master potrebbe rappresentare un collo di bottiglia, dovendo ricevere molti messaggi. Per mitigari, si può distribuire il carico di comunicazione tra più nodi utilizzando degli schemi di comunicaizone gerarchici che coinvolgono tutti i nodi. Ogni nodo dell’albero riceve messaggi dai nodi figli e manda un messaggio «cumulativo» al padre. In questo modo il master viene alleggerito.

5. Quali sono le intestazioni delle funzioni principali di MPI?

```c
int MPI_Init(int* argc, char*** argv);

int MPI_Finalize(void);

int MPI_Comm_create(MPI_Comm comm, MPI_Group group, MPI_Comm* new_comm);

int MPI_Comm_spawn(const char* command, char* argv[],
int maxprocs, MPI_Info info, int root, MPI_Comm comm,
MPI_Comm* intercomm, int array_of_errcodes[]);

int MPI_Comm_size(MPI_Comm comm, int* size);

int MPI_Comm_rank(MPI_Comm comm, int* rank);

int MPI_Send(const void* buffer, int count, MPI_Datatype datatype,
int dest, int tag, MPI_Comm comm);

int MPI_Ssend(const void* buffer, int count, MPI_Datatype datatype,
int dest, int tag, MPI_Comm comm);

int MPI_Isend(const void* buffer, int count, MPI_Datatype datatype,
int dest, int tag, MPI_Comm comm, MPI_Request* request);

int MPI_Bsend(const void* buffer, int count, MPI_Datatype datatype,
int dest, int tag, MPI_Comm comm);

int MPI_Wait(MPI_Request* request, MPI_Status* status);

int MPI_Test(MPI_Request* request, int* flag, MPI_Status* status);

int MPI_Recv(void* buffer, int count, MPI_Datatype datatype, int source,
int tag, MPI_Comm comm, MPI_Status* status);

int MPI_Irecv(void* buffer, int count, MPI_Datatype datatype, int source,
int tag, MPI_Comm comm, MPI_Request* request);

int MPI_Reduce(const void* send_buffer, void* receive_buffer, int count,
MPI_Datatype datatype, MPI_Op operation, int root, MPI_Comm comm);

int MPI_Ireduce(const void* send_buffer, void* receive_buffer,
int count, MPI_Datatype datatype, MPI_Op operation, int root,
MPI_Comm comm, MPI_Request* request);

int MPI_Allreduce(const void* sendbuf, void* recvbuf, int count,
MPI_Datatype datatype, MPI_Op operation, MPI_Comm comm);

int MPI_Bcast(void* buffer, int count, MPI_Datatype datatype,
int emitter_rank, MPI_Comm communicator);

int MPI_Scatter(const void* sendbuf, int count_send,
MPI_Datatype datatype_send, void* recvbuf, int count_recv,
MPI_Datatype datatype_recv, int root, MPI_Comm comm);

int MPI_Gather(void* sendbuf, int count_send,
MPI_Datatype datatype_send, void* recvbuf, int count_recv,
MPI_Datatype datatype_recv, int root, MPI_Comm communicator);

int MPI_Barrier(MPI_Comm comm);

double MPI_Wtime(void);
```

## Programmazione Parallela con OpenMP

1. Quali modelli si possono utilizzare in caso di memoria condivisa tra i processi?
   
   Si possono usare si il modello a scambio di messaggi, ad esempio con MPI, che il modello a memoria condivisa. In questo ultimo caso si possono usare diverse tecnologie: i sistemi multicore o multiprocessor come OpenMP o pthreads oppure le GP-GPU per le quali esistono librerie specifiche che consentono lo sviluppo di programmi destinati ad eseguire su GPU come CUPA (libreria proprietaria NVidia) oppure openCL.

2. Quali sonole caratteristiche principali di OpenMP?
   
   E' una libreria che permette di parallelizzare il codice di programmi C, utilizzando un approccio dichiarativo. Offre strumenti per gestire i thread paralleli, ottenere/impostare informazioni sull'ambiente di esecuzione, definire la visibilità delle variabli rispetto ai thread paralleli e di sincronizzare i thread tra loro con sezioni critiche o bariere di sincronizzazione.

3. Quali sono le clausole principali della direttiva parallel?
   
   Sono diverse. Con num_threads(N) si imposta il numero di thread paralleli. Usando shared private e firstprivate (ogni processo utilizza una copia privata inizializzata al valore che aveva prima di pragma) si può specificare il campo di visibilità (mentre di default la variabili sono private solo se definite internamente al blocco parallel, private altrimenti). Con reduction è possibile utilizzare una variabile di appoggio permettendo di valutare l'espessione in parallelo (aggiornamento di var in mutua esclusione). Con if si rende condizionale la parallelizzazione, se false allora l'esecuzione è sequenziale. Con for è possibile parallelizzare un ciclo e usando schedule ci si può assicurare di bilanciare il carico che può essere static o dynamic.

4. Come si possono gestire le informazioni sui thread in esecuzione?
   
   ```c
   #include <omp.h>
   
   #pragma omp parallel
   {
       // Restituisce il numero di thread paralleli
       // da utilizzare all'interno di un blocco parallelo
       int num_threads = omp_get_num_threads();
   
       // Restituisce il rank del thread che lo invoca (0 è il master)
       int thread_id = omp_get_thread_num();
   
       // Imposta a n il numero di thread paralleli
       // nei successivi blocchi paralleli
       omp_set_num_threads(4);
   }
   ```

5. A cosa servono le direttive master e single?
   
   La direttiva master indica che solo il thread master (tipicamente il thread che ha un ID pari a 0)mentre la direttiva single crea una barierea di sincronizzazione implicita, garantendo che solo un thread, scelto arbitrariamente, esegua il blocco di codice.

6. Quali clausole hanno come scopo principale la sincronizzazione?
   
   La prima è la direttiva critical: il blocco di istruzioni immediatamente successivo alla direttiva viene eseguito un solo processo alla volta. Anche la direttiva barrier è utile ed implmenta la classica barriera di sincronizzazione in un team di threads. Inoltre è possibile usare i lock per reealizzare schemi di sincronizzaizone ad hoc.
   
   ```c
   //esattamente equivalente a critical
   #include <stdio.h> #include <omp.h>
   omp_lock_t my_lock;
   int main() {
       omp_init_lock(&my_lock);
       #pragma omp parallel num_threads(4)
       {
           int i, j, t = omp_get_thread_num( );
           for (i = 0; i < 5; ++i) {
               omp_set_lock(&my_lock); //prologo sezione critica
               printf("Thread %d – inizio sezione critica %d\n", t,i);
               printf("Thread %d - fine sezione critica %d\n", t, i);
               omp_unset_lock(&my_lock); //epilogo sezione critica
           }
       }
       omp_destroy_lock(&my_lock);
   }
   ```

7. Come si può misurare il tempo?
   
   ```c
   double omp_get_wtime(void);
   ```

8. Quali osservazioni possiamo fare paragonando OpenMP ad altre librerie come pthread oppure MPI?
   
   Pthread utilizza un paradigma MPMD ed un modello di creazione fork-join. Mette a disposizione un ampio set di politiche per la sinconizzazione specifiche (mutex, semafori, condition) e risulta particolarmente adatto per algoritmi task-parallel. OpenMp a confronto utilizza un approccio di più alto livello, basato su SPMD ed un modello di crreazione cobegin-coend. La sincronizzazione avviene tramite direttive e clausole che implemntano schemi predefiniti (barrier, critical, reduction), o anche ad hoc coi lock. Ideale per la modellazione di algoritmi data-parallel.
   
   Se mettiamo a confronto OpenMP con MPI è evidente che il primo è ben più semplice da utilizzare (vedi bilalnciamento del carico) e che il secondo, insieme alla complessità di utilizzo, ha tra le sue proprietà una maggiore scalabilità e portabilità. Da notare che è possibile combinare i due e beneficiare dei vantaggi di entrambi: si parla di Hybridization.
