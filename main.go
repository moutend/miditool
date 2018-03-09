package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	midi "github.com/moutend/go-midi"
	"github.com/moutend/go-midi/constant"
	"github.com/moutend/go-midi/event"
)

const (
	SOLO = false
	MUTE = true
)

var (
	Version   = "develop"
	Revision  = "latest"
	DebugFlag bool
)

type NotesFlag struct {
	Notes []constant.Note
}

func (n *NotesFlag) String() string {
	return ""
}
func (n *NotesFlag) Set(v string) error {
	note, err := constant.ParseNote(v)
	if err != nil {
		return err
	}

	n.Notes = append(n.Notes, note)

	return nil
}

func HasNotes(notes []constant.Note, note constant.Note) bool {
	for _, n := range notes {
		if n == note {
			return true
		}
	}
	return false
}

func main() {
	if err := run(os.Args); err != nil {
		log.New(os.Stderr, "error: ", 0).Fatal(err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) < 2 {
		return HelpCommand(args)
	}
	switch args[1] {
	case "version":
		return VersionCommand(args)
	case "help":
		return HelpCommand(args)
	}

	f := flag.NewFlagSet(fmt.Sprintf("%s %s", args[0], args[1]), flag.ExitOnError)
	f.BoolVar(&DebugFlag, "debug", false, "Enable debug output.")
	f.Parse(args[1:])
	args = f.Args()
	command := args[0]
	args = args[1:]

	switch command {
	case "solo":
		return SoloCommand(args)
	case "mute":
		return MuteCommand(args)
	case "velocity":
		return PlusCommand(args)
	default:
		return fmt.Errorf("%s is not a subcommand", command)
	}
}

func SoloCommand(args []string) error {
	return SoloOrMuteCommand(SOLO, args)
}

func MuteCommand(args []string) error {
	return SoloOrMuteCommand(MUTE, args)
}

func SoloOrMuteCommand(mode bool, args []string) error {
	var notesFlag NotesFlag
	var outputFilenameFlag string

	f := flag.NewFlagSet("solo", flag.ExitOnError)
	f.Var(&notesFlag, "note", "specify notes which play")
	f.Var(&notesFlag, "n", "alias of --note")
	f.StringVar(&outputFilenameFlag, "output", "output.mid", "specify output file name")
	f.StringVar(&outputFilenameFlag, "o", "output.mid", "alias of --output")

	err := f.Parse(args)
	if err != nil {
		return err
	}
	args = f.Args()
	if len(args) < 1 {
		return nil
	}

	file, err := ioutil.ReadFile(args[0])
	if err != nil {
		return err
	}

	m, err := midi.NewParser(file).Parse()
	if err != nil {
		return err
	}

	for _, t := range m.Tracks {
		for _, e := range t.Events {
			switch e.(type) {
			case *event.NoteOnEvent:
				noteOnEvent := e.(*event.NoteOnEvent)
				if mode == HasNotes(notesFlag.Notes, noteOnEvent.Note()) {
					noteOnEvent.SetVelocity(0)
				}
			}
		}
	}

	return ioutil.WriteFile(outputFilenameFlag, m.Serialize(), 0644)
}

func PlusCommand(args []string) error {
	var notesFlag NotesFlag
	var outputFilenameFlag string
	var globalVelocityFlag int

	f := flag.NewFlagSet("solo", flag.ExitOnError)
	f.IntVar(&globalVelocityFlag, "v", 0, "specify velocity")
	f.Var(&notesFlag, "note", "specify notes which play")
	f.Var(&notesFlag, "n", "alias of --note")
	f.StringVar(&outputFilenameFlag, "output", "output.mid", "specify output file name")
	f.StringVar(&outputFilenameFlag, "o", "output.mid", "alias of --output")

	err := f.Parse(args)
	if err != nil {
		return err
	}
	args = f.Args()
	if len(args) < 1 {
		return nil
	}

	file, err := ioutil.ReadFile(args[0])
	if err != nil {
		return err
	}

	m, err := midi.NewParser(file).Parse()
	if err != nil {
		return err
	}

	for _, t := range m.Tracks {
		for _, e := range t.Events {
			switch e.(type) {
			case *event.NoteOnEvent:
				noteOnEvent := e.(*event.NoteOnEvent)
				if HasNotes(notesFlag.Notes, noteOnEvent.Note()) {
					velocity := noteOnEvent.Velocity() + uint8(globalVelocityFlag)
					if velocity > 127 {
						velocity = 127
					}
					if velocity < 0 {
						velocity = 0
					}
					noteOnEvent.SetVelocity(velocity)
				}
			}
		}
	}

	return ioutil.WriteFile(outputFilenameFlag, m.Serialize(), 0644)
}

func HelpCommand(args []string) error {
	return nil
}

func VersionCommand(args []string) error {
	return nil
}
