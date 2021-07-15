# genlayout
This is the layout analyzer and generator that was used to create [Semimak](https://semilin.github.io/semimak). It is centered around the novel finger-speed metric. This program was made for myself, and because of that **it's not very easy to use,** especially if you don't have any programming experience. I am currently working on a powerful, configurable, and user-friendly web-based and desktop analyzer+generator. 

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

## License
Copyright © 2021 Colin Hughes

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.
