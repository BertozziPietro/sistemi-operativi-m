--------------------------------------------------------------------------------
--  *  Prog name torte.adb
--  *  Project name torte
--  *  LA FABBRIAC DI TORTE
--------------------------------------------------------------------------------

with Ada.Text_IO, Ada.Integer_Text_IO;
use Ada.Text_IO, Ada.Integer_Text_IO;

procedure torte is

   type clienteOP_ID is range 1..10;
   type clienteOC_ID is range 1..4;

   type torta is (cioccolato, marmellata);
   type confezione is (cioc, marm, family);

   task type clienteOP (ID: clienteOP_ID; T: torta);
   task type clienteOC (ID: clienteOC_ID; C: confezione);

   type acOP is access clienteOP;
   type acOC is access clienteOC;

   task type server is
      entry deposito(torta) (ID: clienteOP_ID);
      entry prelievo(confezione) (ID: clienteOC_ID);
   end server;

   S: server;

   task body server is
      MAX: constant INTEGER := 18;
      sultavolo: array(torta'Range) of Integer;
   begin
      Put_Line("SERVER iniziato!");
      for i in torta'Range loop
         sultavolo(i) := 0;
      end loop;
      delay 2.0;
      loop
         select
            when sultavolo(marmellata) + sultavolo(cioccolato) < MAX and sultavolo(marmellata) < MAX - 1 =>
               accept deposito(marmellata) (ID: in clienteOP_ID) do
                  Put_Line("deposito 1 crostata " & clienteOP_ID'Image(ID) & "!");
                  sultavolo(marmellata) := sultavolo(marmellata) + 1;
               end;
            or
            when sultavolo(marmellata) + sultavolo(cioccolato) < MAX and sultavolo(cioccolato) < MAX - 1 and deposito(marmellata)'COUNT = 0 =>
               accept deposito(cioccolato) (ID: in clienteOP_ID) do
                  Put_Line("deposito 1 torta al cioccolato " & clienteOP_ID'Image(ID) & "!");
                  sultavolo(cioccolato) := sultavolo(cioccolato) + 1;
               end;
            or
            when sultavolo(marmellata) >= 1 and sultavolo(cioccolato) >= 1 =>
               accept prelievo(family) (ID: in clienteOC_ID) do
                  Put_Line("prelievo confezione family " & clienteOC_ID'Image(ID) & "!");
                  sultavolo(marmellata) := sultavolo(marmellata) - 1;
                  sultavolo(cioccolato) := sultavolo(cioccolato) - 1;
               end;
            or
            when sultavolo(marmellata) >= 1 and prelievo(family)'COUNT = 0 =>
               accept prelievo(marm) (ID: in clienteOC_ID) do
                  Put_Line("prelievo marmellata " & clienteOC_ID'Image(ID) & "!");
                  sultavolo(marmellata) := sultavolo(marmellata) - 1;
               end;
            or
            when sultavolo(cioccolato) >= 1 and prelievo(family)'COUNT = 0 and prelievo(marm)'COUNT = 0 =>
               accept prelievo(cioc) (ID: in clienteOC_ID) do
                  Put_Line("prelievo cioccolato " & clienteOC_ID'Image(ID) & "!");
                  sultavolo(cioccolato) := sultavolo(cioccolato) - 1;
               end;
         end select;
      end loop;
   end;

   task body clienteOP is
   begin
      Put_Line("clienteOP" & clienteOP_ID'Image(ID) & " di tipo" & torta'Image(T) & " iniziato!");
      S.deposito(T)(ID);
   end;

   task body clienteOC is
   begin
      Put_Line("clienteOC" & clienteOC_ID'Image(ID) & " di tipo" & confezione'Image(C) & " iniziato!");
      S.prelievo(C)(ID);
   end;

   NewOP: acOP;
   NewOC: acOC;

begin
   for I in clienteOP_ID'Range loop
      NewOP := new clienteOP (I, cioccolato);
      NewOP := new clienteOP (I, marmellata);
   end loop;

   for I in clienteOC_ID'Range loop
      NewOC := new clienteOC (I, cioc);
      NewOC := new clienteOC (I, marm);
      NewOC := new clienteOC (I, family);
   end loop;

end torte;
