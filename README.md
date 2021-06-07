# genlayout

## Usage
`./genkey command argument`
### Commands
#### [load]
Reads `text.txt` in the same directory and outputs its bigram, trigram, and skipgram data into `data.json`. This data is used for every other command.
#### [r]ank
Outputs a sorted list of all the layouts with their scores.
#### [a]nalyze layout
Outputs detailed analysis of a given layout.
#### [g]enerate
Attempts to find the optimal layout according to the scoring algorithm.
#### [sfbs] layout
Lists the sfb frequency and most frequent sfbs.
#### [dsfbs] layout
Lists the dsfb frequency and most frequent dsfbs.
#### [speed] layout
Lists each finger and its unweighted speed for the layout.
