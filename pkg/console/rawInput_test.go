package console

import "testing"

type rawInputTestingData struct {
	input         rawInput
	expectedValue bool
	name          string
}

func Test_isSpecialInput(t *testing.T) {
	data := []rawInputTestingData{
		{rawInput{}, false, "empty input"},
		{rawInput{0x1, specialSymbolContinue}, false, "wrong first input byte"},
		{rawInput{specialSymbolStart, 0x2}, false, "wrong second input byte"},
		{rawInput{specialSymbolStart, specialSymbolContinue}, true, "valid input"},
	}

	for _, tc := range data {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.input.isSpecialInput()
			if res != tc.expectedValue {
				logRawInputErrorMessage(t, tc.expectedValue, res)
			}
		})
	}
}

func Test_isArrowDown(t *testing.T) {
	data := []rawInputTestingData{
		{rawInput{}, false, "empty input"},
		{rawInput{0x1, 0x2, arrowDown}, false, "non special input"},
		{getSpecialRawInput(0x1), false, "last input is not the 'arrowDown'"},
		{getSpecialRawInput(arrowDown), true, "valid last input"},
	}
	for _, tc := range data {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.input.isArrowDown()
			if res != tc.expectedValue {
				logRawInputErrorMessage(t, tc.expectedValue, res)
			}
		})
	}
}

func Test_isArrowUp(t *testing.T) {
	data := []rawInputTestingData{
		{rawInput{}, false, "empty input"},
		{rawInput{0x1, 0x2, arrowDown}, false, "non special input"},
		{getSpecialRawInput(0x1), false, "last input is not the 'arrowUp'"},
		{getSpecialRawInput(arrowUp), true, "valid last input"},
	}
	for _, tc := range data {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.input.isArrowUp()
			if res != tc.expectedValue {
				logRawInputErrorMessage(t, tc.expectedValue, res)
			}
		})
	}
}

func Test_isEnter(t *testing.T) {
	data := []rawInputTestingData{
		{rawInput{}, false, "empty input"},
		{rawInput{0x1}, false, "first input is not the 'enter'"},
		{rawInput{enter}, true, "valid input"},
	}
	for _, tc := range data {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.input.isEnter()
			if res != tc.expectedValue {
				logRawInputErrorMessage(t, tc.expectedValue, res)
			}
		})
	}
}

func Test_isCtrlC(t *testing.T) {
	data := []rawInputTestingData{
		{rawInput{}, false, "empty input"},
		{rawInput{0x1}, false, "first input is not the 'ctrlC'"},
		{rawInput{ctrlC}, true, "valid input"},
	}
	for _, tc := range data {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.input.isCtrlC()
			if res != tc.expectedValue {
				logRawInputErrorMessage(t, tc.expectedValue, res)
			}
		})
	}
}

func getSpecialRawInput(lastInputElement byte) rawInput {
	return rawInput{specialSymbolStart, specialSymbolContinue, lastInputElement}
}

func logRawInputErrorMessage(t *testing.T, expectedValue bool, gotRes bool) {
	t.Errorf("expected %t, got %t", expectedValue, gotRes)
}
