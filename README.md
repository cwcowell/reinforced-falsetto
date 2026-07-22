# reinforced-falsetto

An interactive ear-training drill for reinforced falsetto (a.k.a. mixed voice).
It plays a random pitch (E, F, F#, G in the octave above middle C)
paired with a random vowel ("ee", "eh", "ah", "oh", "oo"). Press **Return**
if you hit it successfully, or any other key if you didn't. Every one of the
20 pitch/vowel combinations is used exactly once per run, then a results
table and success rate are printed.

## Prerequisites

- [Go](https://go.dev/dl/) 1.21 or later
- Speakers or headphones

## Running the trainer

```bash
git clone git@github.com:cwcowell/reinforced-falsetto.git
cd reinforced-falsetto
go run .
```

## Building a binary

```bash
go build -o reinforced-falsetto .
./reinforced-falsetto
```

On Windows this produces `reinforced-falsetto.exe`.

## How it works

- Press **Return/Enter** to mark the current pitch/vowel combo a success.
- Press **any other key** to mark it a failure.
- The drill ends automatically after all 20 combinations have played once.
- Results print as a table (pitches across the top, vowels down the side)
  with a green check for successes and a red X for failures, followed by
  your success count, failure count, and success percentage.
