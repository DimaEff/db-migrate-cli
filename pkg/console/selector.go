package console

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

const inputSize = 3

const defaultSelectedIndexList int = 0

const (
	selectedMark     = ">"
	clearCode        = "\033[H\033[2J"
	moveCursorUpCode = "\033[%dA"
)

const (
	ctrlC                 = 0x03
	enter                 = 0x0D
	specialSymbolStart    = 0x1B
	specialSymbolContinue = 0x5B
	arrowUp               = 0x41
	arrowDown             = 0x42
)

type OptionValue int

// rawInput is fixed size array because terminal special inputs never exceed 3 bytes
type rawInput [inputSize]byte

type SelectListOption struct {
	Value OptionValue
	Text  string
}

type SelectList struct {
	options       []SelectListOption
	selectedIndex int
}

func NewSelectList(options []SelectListOption) (*SelectList, error) {
	err := validateOptions(options)
	if err != nil {
		return nil, err
	}

	return &SelectList{
		options:       options,
		selectedIndex: defaultSelectedIndexList,
	}, nil
}

func validateOptions(options []SelectListOption) error {
	if len(options) == 0 {
		return fmt.Errorf("new select list: options len must be > 0")
	}

	values := make(map[int]string)
	for _, o := range options {
		if _, exists := values[int(o.Value)]; exists {
			return fmt.Errorf("all option`s values must be unique")
		}
		values[int(o.Value)] = o.Text
	}

	return nil
}

// Run blocks until user selects an option or cancels selection with Ctrl+C
func (sl *SelectList) Run() (*SelectListOption, error) {
	prevState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return nil, fmt.Errorf("select list run, make raw: %w", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), prevState)

	sl.render()

	input := rawInput{}
	for {
		_, err := os.Stdin.Read(input[:])
		if err != nil {
			return nil, fmt.Errorf("select list run, read input: %w", err)
		}

		switch {
		case input.isCtrlC():
			return nil, fmt.Errorf("selection canceled")
		case input.isEnter():
			return &sl.options[sl.selectedIndex], nil
		case input.isArrowUp():
			if sl.isValidIndex(sl.selectedIndex - 1) {
				sl.selectedIndex--
				sl.render()
			}
		case input.isArrowDown():
			if sl.isValidIndex(sl.selectedIndex + 1) {
				sl.selectedIndex++
				sl.render()
			}
		}
	}
}

func (sl *SelectList) render() {
	sl.moveCursorToStart()
	for i, opt := range sl.options {
		prefix := " "
		if sl.selectedIndex == i {
			prefix = selectedMark
		}
		fmt.Printf("\033[2K\r%s%2d. %s\n", prefix, i+1, opt.Text)
	}
}

func clearConsole() {
	fmt.Print(clearCode)
}

func (sl *SelectList) moveCursorToStart() {
	fmt.Printf(moveCursorUpCode, len(sl.options))
}

func (sl *SelectList) isValidIndex(index int) bool {
	return index >= 0 && index < len(sl.options)
}

func (i *rawInput) isCtrlC() bool {
	return i[0] == ctrlC
}

func (i *rawInput) isEnter() bool {
	return i[0] == enter
}

func (i *rawInput) isArrowUp() bool {
	return i.isSpecialInput() && i[2] == arrowUp
}

func (i *rawInput) isArrowDown() bool {
	return i.isSpecialInput() && i[2] == arrowDown
}

func (i *rawInput) isSpecialInput() bool {
	return i[0] == specialSymbolStart && i[1] == specialSymbolContinue
}
