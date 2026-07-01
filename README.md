# reinforced-falsetto

An interactive ear-training drill for reinforced falsetto (a.k.a. mixed voice).
It plays a random pitch (D#, E, F, F#, G, G# in the octave above middle C)
paired with a random vowel ("ee", "eh", "ah", "oh", "oo"). Press **Return**
if you hit it successfully, or any other key if you didn't. Every one of the
30 pitch/vowel combinations is used exactly once per run, then a results
table and success rate are printed.

## Prerequisites

- Python 3.9 or later
- [uv](https://docs.astral.sh/uv/) (manages the virtual environment and dependencies automatically)
- Speakers or headphones

### Installing uv

**macOS / Linux:**

```bash
curl -LsSf https://astral.sh/uv/install.sh | sh
```

**Windows (PowerShell):**

```powershell
powershell -ExecutionPolicy ByPass -c "irm https://astral.sh/uv/install.ps1 | iex"
```

Restart your terminal afterward so `uv` is on your `PATH`.

## Running the trainer

Clone the repo and run the script with `uv run` — it will automatically
create a virtual environment and install the required dependencies
(`numpy`, `simpleaudio`) the first time you run it:

```bash
git clone git@github.com:cwcowell/reinforced-falsetto.git
cd reinforced-falsetto
uv run reinforced-falsetto.py
```

### Linux audio note

`simpleaudio` uses ALSA on Linux. If installing dependencies fails while
building `simpleaudio`, install ALSA's development headers first, then
re-run `uv run reinforced-falsetto.py`:

```bash
sudo apt install libasound2-dev   # Debian/Ubuntu
sudo dnf install alsa-lib-devel   # Fedora
```

## How it works

- Press **Return/Enter** to mark the current pitch/vowel combo a success.
- Press **any other key** to mark it a failure.
- The drill ends automatically after all 30 combinations have played once.
- Results print as a table (pitches across the top, vowels down the side)
  with a green check for successes and a red X for failures, followed by
  your success count, failure count, and success percentage.
