--------------------------------------------------------------------------------
--  *  Prog name ponte.adb
--  *  Project name ponte
--  *  PONTE A SENSO UNICO ALTERNATO
--------------------------------------------------------------------------------

with Ada.Text_IO, Ada.Integer_Text_IO;
use Ada.Text_IO, Ada.Integer_Text_IO;

procedure ponte is
   type cliente_ID is range 1..10;
   type dir_ID is  (NORD, SUD);
   task type cliente (ID: cliente_ID);
   type ac is access cliente;
   task type server is
      entry entraNORD (ID: in cliente_ID );
      entry esceNORD(ID: in cliente_ID );
      entry entraSUD (ID: in cliente_ID );
      entry esceSUD(ID: in cliente_ID );
   end server;

   S: server;

   task body server is
      MAX  : constant INTEGER := 5;
      sulponte: Integer;
      utenti: array(dir_ID'Range) of Integer;
   begin
      Put_Line ("SERVER iniziato!");
      sulponte := 0;
      for i in dir_ID'Range loop
         utenti(i) := 0;
      end loop;
      loop
         select
            when sulponte < MAX and utenti(SUD) = 0 =>
               accept entraNORD (ID: in cliente_ID ) do
                  Put_Line("sta entrando da nord il cliente "& cliente_ID'Image(ID) &" !");
                  utenti(NORD) := utenti(NORD) + 1;
                  sulponte := sulponte + 1;
                  Put ("..ora ci sono ");
                  Put(utenti(NORD));
                  Put("utenti entrati da NORD e");
                  Put(sulponte,10);
                  Put(" utenti in totale!");
                  New_Line;
               end entraNORD;
            or
               when sulponte < MAX and utenti(NORD) = 0 =>
                  accept entraSUD (ID: in cliente_ID ) do
                     Put_Line("sta entrando da sud il cliente "& cliente_ID'Image(ID) &" !");
                     utenti(SUD) := utenti(SUD) + 1;
                     sulponte := sulponte + 1;
                     Put ("..ora ci sono ");
                     Put(utenti(SUD));
                     Put("utenti entrati da SUD e");
                     Put(sulponte,10);
                     Put(" utenti in totale!");
                     New_Line;
                  end entraSUD;
            or
               accept esceNORD (ID: in cliente_ID ) do
                  Put_Line("sta uscendo da sud il cliente "& cliente_ID'Image(ID) &" !");
                  utenti(NORD) := utenti(NORD) - 1;
                  sulponte := sulponte - 1;
                  Put ("..ora ci sono ");
                  Put(utenti(NORD));
                  Put("utenti entrati da NORD e");
                  Put(sulponte,10);
                  Put(" utenti in totale!");
                  New_Line;
               end esceNORD;
            or
               accept esceSUD (ID: in cliente_ID ) do
                  Put_Line("sta uscendo da nord il cliente "& cliente_ID'Image(ID) &" !");
                  utenti(SUD) := utenti(SUD) - 1;
                  sulponte := sulponte - 1;
                  Put ("..ora ci sono ");
                  Put(utenti(SUD));
                  Put("utenti entrati da SUD e");
                  Put(sulponte,10);
                  Put(" utenti in totale!");
                  New_Line;
               end esceSUD;
         end select;
      end loop;
   end;

   task body cliente is
   begin
      Put_Line ("cliente" & cliente_ID'Image(ID) & " iniziato!");
      S.entraSUD(ID);
      delay 1.0;
      S.esceSUD(ID);
      delay 1.0;
      S.entraNORD(ID);
      delay 1.0;
      S.esceNORD(ID);
   end;

   New_client: ac;

begin
   for I in cliente_ID'Range loop
      New_client := new cliente (I);
   end loop;
end ponte;
