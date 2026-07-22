// Reinforced falsetto trainer.
//
// Plays a random pitch (E, F, F#, G in the octave above middle C)
// paired with a random vowel sound ("ee", "eh", "ah", "oh", "oo"), one
// combination at a time, without repeating a combo within a run.
//
// Press Return/Enter to mark a trial as a success. Press any other key
// to mark it as a failure. After all 30 combinations have been played,
// prints the number of successes, failures, and the success percentage.
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"

	"github.com/ebitengine/oto/v3"
	"golang.org/x/term"
)

const (
	sampleRate          = 44100
	toneDurationSeconds = 1.0
	greenCheck          = "\033[92m✓\033[0m"
	redX                = "\033[91m✗\033[0m"
)

type pitch struct {
	name string
	freq float64
}

// Octave above middle C (C5), equal temperament, A4 = 440 Hz.
var pitches = []pitch{
	{"E", 659.25},
	{"F", 698.46},
	{"F#", 739.99},
	{"G", 783.99},
}

var vowels = []string{"ee", "eh", "ah", "oh", "oo"}

type combo struct {
	pitch pitch
	vowel string
}

func makeTone(frequency float64) []byte {
	numSamples := int(sampleRate * toneDurationSeconds)
	buf := new(bytes.Buffer)
	buf.Grow(numSamples * 2)
	for i := 0; i < numSamples; i++ {
		t := float64(i) / sampleRate
		sample := math.Sin(2 * math.Pi * frequency * t)
		binary.Write(buf, binary.LittleEndian, int16(sample*32767))
	}
	return buf.Bytes()
}

func playTone(ctx *oto.Context, frequency float64) *oto.Player {
	audio := makeTone(frequency)
	player := ctx.NewPlayer(bytes.NewReader(audio))
	player.Play()
	return player
}

// getKeypress blocks until a single key is pressed and returns it
// without requiring Enter.
func getKeypress() (string, error) {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return "", err
	}
	defer term.Restore(fd, oldState)

	buf := make([]byte, 1)
	if _, err := os.Stdin.Read(buf); err != nil {
		return "", err
	}
	return string(buf), nil
}

func isEnter(key string) bool {
	return key == "\r" || key == "\n"
}

func center(s string, width int) string {
	if len(s) >= width {
		return s
	}
	total := width - len(s)
	left := total / 2
	right := total - left
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
}

func printResultsTable(results map[combo]bool) {
	colWidth := 0
	for _, p := range pitches {
		if len(p.name) > colWidth {
			colWidth = len(p.name)
		}
	}
	colWidth += 2

	vowelWidth := 0
	for _, v := range vowels {
		if len(v) > vowelWidth {
			vowelWidth = len(v)
		}
	}

	header := strings.Repeat(" ", vowelWidth)
	for _, p := range pitches {
		header += center(p.name, colWidth)
	}
	fmt.Println(header)

	for _, v := range vowels {
		row := v + strings.Repeat(" ", vowelWidth-len(v))
		for _, p := range pitches {
			mark := redX
			if results[combo{p, v}] {
				mark = greenCheck
			}
			// Marks contain ANSI escape codes, which don't take up display
			// width, so widen the field to compensate (mirrors the Python
			// version's `col_width + len(mark) - 1`).
			row += center(mark, colWidth+len(mark)-1)
		}
		fmt.Println(row)
	}
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	op := &oto.NewContextOptions{
		SampleRate:   sampleRate,
		ChannelCount: 1,
		Format:       oto.FormatSignedInt16LE,
	}
	ctx, ready, err := oto.NewContext(op)
	if err != nil {
		return fmt.Errorf("initializing audio: %w", err)
	}
	<-ready

	var combos []combo
	for _, p := range pitches {
		for _, v := range vowels {
			combos = append(combos, combo{p, v})
		}
	}
	rand.Shuffle(len(combos), func(i, j int) {
		combos[i], combos[j] = combos[j], combos[i]
	})

	successes := 0
	failures := 0
	results := make(map[combo]bool)

	fmt.Println("Press Return for a success, any other key for a failure. Ctrl+C to quit early.")
	fmt.Println()

	// Keep played sounds alive until the run ends, since oto.Player must
	// not be garbage-collected while audio is still playing.
	var players []*oto.Player

	for i, c := range combos {
		fmt.Printf("%d. %s\n", i+1, c.vowel)
		players = append(players, playTone(ctx, c.pitch.freq))
		key, err := getKeypress()
		if err != nil {
			return fmt.Errorf("reading keypress: %w", err)
		}
		success := isEnter(key)
		results[c] = success
		if success {
			successes++
		} else {
			failures++
		}
	}

	total := successes + failures
	percent := 0.0
	if total > 0 {
		percent = float64(successes) / float64(total) * 100
	}
	fmt.Println()
	fmt.Println("Done!")
	fmt.Printf("Successes: %d\n", successes)
	fmt.Printf("Failures: %d\n", failures)
	fmt.Printf("Success rate: %.1f%%\n", percent)
	fmt.Println()
	printResultsTable(results)

	return nil
}
