#!/usr/bin/env python3
"""
Reinforced falsetto trainer.

Plays a random pitch (E, F, F#, G in the octave above middle C)
paired with a random vowel sound ("ee", "eh", "ah", "oh", "oo"), one
combination at a time, without repeating a combo within a run.

Press Return/Enter to mark a trial as a success. Press any other key
to mark it as a failure. After all 30 combinations have been played,
prints the number of successes, failures, and the success percentage.

Requires: numpy, simpleaudio
    pip install numpy simpleaudio
"""

import random
import sys

import numpy as np
import simpleaudio as sa

SAMPLE_RATE = 44100
TONE_DURATION_SECONDS = 1.0

# Octave above middle C (C5), equal temperament, A4 = 440 Hz.
PITCH_FREQUENCIES = {
    "E": 659.25,
    "F": 698.46,
    "F#": 739.99,
    "G": 783.99,
}

VOWELS = ["ee", "eh", "ah", "oh", "oo"]

GREEN_CHECK = "\033[92m✓\033[0m"
RED_X = "\033[91m✗\033[0m"


def make_tone(frequency):
    t = np.linspace(0, TONE_DURATION_SECONDS, int(SAMPLE_RATE * TONE_DURATION_SECONDS), endpoint=False)
    waveform = np.sin(2 * np.pi * frequency * t)
    audio = (waveform * 32767).astype(np.int16)
    return audio


def play_tone(frequency):
    audio = make_tone(frequency)
    return sa.play_buffer(audio, 1, 2, SAMPLE_RATE)


def get_keypress():
    """Block until a single key is pressed; return it without requiring Enter."""
    if sys.platform == "win32":
        import msvcrt

        ch = msvcrt.getch()
        return ch.decode("utf-8", errors="ignore")
    else:
        import termios
        import tty

        fd = sys.stdin.fileno()
        old_settings = termios.tcgetattr(fd)
        try:
            tty.setraw(fd)
            ch = sys.stdin.read(1)
        finally:
            termios.tcsetattr(fd, termios.TCSADRAIN, old_settings)
        return ch


def is_enter(key):
    return key in ("\r", "\n")


def print_results_table(results):
    pitches = list(PITCH_FREQUENCIES)
    col_width = max(len(pitch) for pitch in pitches) + 2
    vowel_width = max(len(vowel) for vowel in VOWELS)

    header = " " * vowel_width + "".join(pitch.center(col_width) for pitch in pitches)
    print(header)
    for vowel in VOWELS:
        row = vowel.ljust(vowel_width)
        for pitch in pitches:
            mark = GREEN_CHECK if results[(pitch, vowel)] else RED_X
            row += mark.center(col_width + len(mark) - 1)
        print(row)


def main():
    combos = [(pitch, vowel) for pitch in PITCH_FREQUENCIES for vowel in VOWELS]
    random.shuffle(combos)

    successes = 0
    failures = 0
    results = {}

    print("Press Return for a success, any other key for a failure. Ctrl+C to quit early.\n")

    for i, (pitch, vowel) in enumerate(combos, start=1):
        print(f"{i}. {vowel}")
        play_tone(PITCH_FREQUENCIES[pitch])
        key = get_keypress()
        success = is_enter(key)
        results[(pitch, vowel)] = success
        if success:
            successes += 1
        else:
            failures += 1

    total = successes + failures
    percent = (successes / total * 100) if total else 0.0
    print("\nDone!")
    print(f"Successes: {successes}")
    print(f"Failures: {failures}")
    print(f"Success rate: {percent:.1f}%\n")
    print_results_table(results)


if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        print("\nInterrupted.")
