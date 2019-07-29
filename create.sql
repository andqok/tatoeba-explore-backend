create table sentences (
  sentence_id serial primary key,
  sentence_number integer,
  sentence_text text,
  lang varchar(4)
);

create table links (
  link_id serial primary key,
  link_1 integer,
  link_2 integer
);

INSERT INTO sentences (lang, sentence_number, sentence_text)
    SELECT 'cmn', '1', '我們試試看'
    WHERE NOT EXISTS (
      SELECT * FROM sentences WHERE lang = 'cmn' and sentence_number = 1 and sentence_text = '我們試試看'
    ) ;

select * from sentences as s
  left join links as l  on l.link_1 = s.sentence_number
  left join sentences as s2 on s2.sentence_number = link_2
  where s.sentence_number = 110;

CREATE TABLE words (
  word text,
  sentence_number integer,
  lang varchar(4)
);

SELECT word, COUNT(*) AS frequency
FROM words
GROUP BY word
ORDER BY COUNT(*) DESC;
