# Argomenti Principali

## Virtualizzazione

## Protezione

1. Quali sono le definizioni di protezione e di sicurezza, nel contesto dei sistemi operativi?
   
   La protezione consiste nell'insieme di attività volte a garantire il controllo dell'accesso alle risorse logiche e fisiche da parte degli utenti, mentre la sicurezza riguarda l'insieme delle tecniche con le quali regolamentare l'accesso degli utenti al sistema di elaborazione. La sicurezza impedisce accessi non autorizzati al sistema e i conseguenti tentativi dolosi di alterazione e distruzione dei dati.

2. Quali sono i meccanismi di sicurezza?
   
   I meccanismi di sicurezza sono: identificazione (chi sei?), autenticazione (come faccio a sapere che sei chi dici di essere?), autorizzazione (cosa puoi fare?).

3. Quali concetti chiave è necessario introdurre per descrivere il controllo degli accessi ad un sistema?
   
   E' necessario introdurre i concetti di: modelli, politiche e meccanismi. Un modello di protezione definisce i soggetti (attivi), gli oggetti (passivi) e i diritti di accesso dei soggetti (su oggetti o soggetti). Le politicha possono essere Discretional Access Control (DAC), Mandatory Access Control (MAC) o Role Based Access Control (RABC). I meccanismi di protezione sono gli strumenti che permettono di imporre una determinata politica e i principi sono la flessibilità del sistema di protezione e la separazione tra i meccanismi e le politiche.

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
    
    Garantisce: mediazione completa, applicando le regole di sicurezza ad ogni accesso e non solo, isolamento proteggendo reference monitor e base di dati da modifiche non autorizzate,  e infine verificabilità delle precedenti.

15. Quali sono le principali categorie dell'Orange Book?
    
    D (Minimal Protection), C (Discretional Protection), B (Mandatory Protection), A (Verified Protection).

## Programmzione Concorrente

## Modello a Memoria Comune

## Nucleo Sistema Memoria Comune

## Modello Scambio di Messaggi

## Comunicazione Sincronizzazione Estesa

## Implementazioni Concorrenza

## Algoritm Sincronizzazione Distribuiti

## Introduzione HCP

## Programmazione HPC

## Programmazione Parallela
