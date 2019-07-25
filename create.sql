create table sentences ( sentence_id serial primary key, sentence_number integer, sentence_text text, lang varchar(3) );

INSERT INTO sentences (lang, sentence_number, sentence_text)
    SELECT 'cmn', '1', '我們試試看'
    WHERE NOT EXISTS (
      SELECT * FROM sentences WHERE lang = 'cmn' and sentence_number = 1 and sentence_text = '我們試試看'
    ) ;
