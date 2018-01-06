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
		return VelocityCommand(args)
	default:
		return fmt.Errorf("%s is not a subcommand", command)
	}
}

func SoloCommand(args []string) error {
	var notesFlag NotesFlag

	f := flag.NewFlagSet("solo", flag.ExitOnError)
	f.Var(&notesFlag, "n", "specify notes which play")

	err := f.Parse(args)
	if err != nil {
		return err
	}
	args = f.Args()

	file, err := ioutil.ReadFile(args[0])
	if err != nil {
		return err
	}

	m, err := midi.NewParser(file).Parse()
	if err != nil {
		return err
	}

	for _, track := range m.Tracks {
		for _, e := range track.Events {
			switch e.(type) {
			case *event.NoteOnEvent:
				noteOnEvent := e.(*event.NoteOnEvent)
				if !HasNotes(notesFlag.Notes, noteOnEvent.Note()) {
					noteOnEvent.SetVelocity(0)
				}
			}
		}
	}
	return ioutil.WriteFile("output.mid", m.Serialize(), 0644)
}

func MuteCommand(args []string) error {
	return nil
}

func VelocityCommand(args []string) error {
	return nil
}

func HelpCommand(args []string) error {
	return nil
}

func VersionCommand(args []string) error {
	return nil
}
