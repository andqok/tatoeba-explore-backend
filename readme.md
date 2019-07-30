## Importing Tatoeba database

Download sentences and links: https://tatoeba.org/eng/downloads

Unpack.

`sudo -i -u postgres`

`cd <my/working/dir>` where sentences are located.

`psql`

`\c tatoeba-explore`

`\copy links(link_1, link_2) from 'links.csv' delimiter E'\t' csv;`

Fix issues with double quotes, which causes wrong parsing (without explicit error):

`sed "s/\"/\"\"\"\"/g" sentences.csv > sentences-repaired.csv`

`\copy sentences(sentence_number, lang, sentence_text) from 'sentences-repaired.csv' delimiter E'\t' csv`

