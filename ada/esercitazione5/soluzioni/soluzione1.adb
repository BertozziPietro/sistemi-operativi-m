with Ada.Text_IO, Ada.Integer_Text_IO;
use Ada.Text_IO, Ada.Integer_Text_IO;
with Ada.Numerics.Discrete_Random;

procedure ufficio is

   N : constant Integer := 5;
   K : constant Integer := 10;

   type tema is (tur, eve);
   type utente_ID is range 1 .. 50;
   type sportello_ID is range 1 .. N;

   task type URP is
      entry acquisizione(tema)(ID : in utente_ID; SP : out sportello_ID);
      entry rilascio(ID : in utente_ID; SP : in sportello_ID);
   end URP;

   U : URP;

   task type utente (ID : utente_ID; T : tema);

   task body utente is
      sportello : sportello_ID;
   begin
      U.acquisizione(T)(ID, sportello);
      delay 2.0;
      U.rilascio(ID, sportello);
   end utente;

   task body URP is
      sportelli : Integer;
      prio      : tema;
      contatore : Integer;
      ris       : sportello_ID;
      libero    : array (sportello_ID) of Boolean;
   begin
      sportelli := N;
      for I in sportello_ID'Range loop
         libero(I) := True;
      end loop;
      contatore := 0;
      prio      := tur;

      loop
         select
            when (sportelli > 0) and (prio = tur or (prio = eve and acquisizione(eve)'Count = 0)) =>
               accept acquisizione(tur)(ID : in utente_ID; SP : out sportello_ID) do
                  for I in sportello_ID'Range loop
                     if libero(I) = True then
                        libero(I) := False;
                        ris       := I;
                        exit;
                     end if;
                  end loop;
                  Put_Line("SERVER: ho assegnato all'Utente TUR " & utente_ID'Image(ID) & " lo sportello n. " & sportello_ID'Image(ris));
                  sportelli := sportelli - 1;
                  if prio = tur then
                     contatore := contatore + 1;
                     if contatore = K then
                        Put_Line("cambio prio TUR->EVE!");
                        prio      := eve;
                        contatore := 0;
                     end if;
                  end if;
                  SP := ris;
               end acquisizione;

            or when (sportelli > 0) and (prio = eve or (prio = tur and acquisizione(tur)'Count = 0)) =>
               accept acquisizione(eve)(ID : in utente_ID; SP : out sportello_ID) do
                  for I in sportello_ID'Range loop
                     if libero(I) = True then
                        libero(I) := False;
                        ris       := I;
                        exit;
                     end if;
                  end loop;
                  Put_Line("SERVER: ho assegnato all'Utente EVE " & utente_ID'Image(ID) & " lo sportello n. " & sportello_ID'Image(ris));
                  sportelli := sportelli - 1;
                  if prio = eve then
                     contatore := contatore + 1;
                     if contatore = K then
                        prio := tur;
                        Put_Line("cambio prio EVE->TUR!");
                        contatore := 0;
                     end if;
                  end if;
                  SP := ris;
               end acquisizione;

            or
               accept rilascio(ID : in utente_ID; SP : in sportello_ID) do
                  Put_Line("Utente " & utente_ID'Image(ID) & " rilascia sportello n." & sportello_ID'Image(SP));
                  sportelli   := sportelli + 1;
                  libero(SP) := True;
               end rilascio;
         end select;
      end loop;
   end URP;

   type acU is access utente;
   new_utente : acU;

begin
   for I in utente_ID'Range loop
      new_utente := new utente(I, tur);
      new_utente := new utente(I, eve);
   end loop;

end ufficio;
