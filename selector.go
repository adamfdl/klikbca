package klikbca

/*
   For table:last-child
      The span only has 2 childs which are both tables. The first table is for the heading that
      renders "INFORMASI REKENING - MUTASI REKENING". The second table is the account statement that we
      need to scrape.

   For tr:nth-child(2)
      Need to skip the first row, because the first row is a header
      Need to skip the last row, because the last row is the "Saldo awal" dan "Saldo akhir", which we dont need

   For td:nth-child(2)
       Inside the table row there are 3 table datas. The first and the third table data does not contain anything.
       What we are interested in is only inside table data number 2.

   For table
       Inside td:nth-child(2), there is another table. This table is what we are going to scrape.

   For tbody
       Inside table there is tbody. Should be self explanatory.

   For tr[bgcolor]
       Inisde <tbody> there will be multiple <tr>. The first <tr> does not have bgcolor class styling, which should be the heading
       that stores "TGL" and "KETERANGAN". The rest of the <tr> has the class bgcolor which provides
       the settlement data that we are going to scrape.

   For td:not([valign])
       For each <tr[bgcolor] there will be 3 table datas. The first one is to store the date, we are not interested in this data.
       The second one is the actual settlement information along with the amount. The third one I don't know what.
       tl;dr we are only interested in the second one, and the second one does not have [valign] class.

*/
const dailySettlementSelector = `span[class="blue"] > table:last-child > tbody > tr:nth-child(2) > td:nth-child(2) > table > tbody > tr[bgcolor] > td:not([valign])`
