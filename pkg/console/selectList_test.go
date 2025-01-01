package console

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
)

var selectListOptionsData = []SelectListOption{
	{Value: OptionValue(1), Text: "First option text"},
	{Value: OptionValue(2), Text: "Second option text"},
	{Value: OptionValue(3), Text: "Third option text"},
}

type selectListTestingData struct {
	input         SelectList
	expectedValue bool
	name          string
}

func Test_isValidIndex(t *testing.T) {
	selectList := SelectList{
		selectListOptionsData,
		defaultSelectedIndexList,
	}
	data := []struct {
		index         int
		expectedValue bool
		name          string
	}{
		{-1, false, "out of range index(upper limit)"},
		{len(selectList.options) + 1, false, "out of range index(lower limit)"},
		{0, true, "boundary value(upper limit)"},
		{len(selectList.options) - 1, true, "boundary value(lower limit)"},
		{len(selectList.options) - 2, true, "index between boundaries"},
	}

	for _, tc := range data {
		t.Run(tc.name, func(t *testing.T) {
			res := selectList.isValidIndex(tc.index)
			if res != tc.expectedValue {
				t.Errorf("expected %t, got %t", tc.expectedValue, res)
			}
		})
	}
}

func Test_moveCursorToStart(t *testing.T) {
	sl := SelectList{
		selectListOptionsData,
		defaultSelectedIndexList,
	}

	t.Run("moving course to the top of list", func(t *testing.T) {
		cursorBuf := writePrintedContent(sl.moveCursorToStart)

		expected := fmt.Sprintf(moveCursorUpCode, len(sl.options))
		if res := cursorBuf.String(); res != expected {
			t.Errorf("expected %q, got %q", expected, res)
		}
	})
}

func Test_clearConsole(t *testing.T) {
	t.Run("success clear console", func(t *testing.T) {
		clearedConsoleContent := writePrintedContent(clearConsole)
		expected := fmt.Sprintf(clearCode)
		if res := clearedConsoleContent.String(); res != expected {
			t.Errorf("expected %q, got %q", expected, res)
		}
	})
}

func Test_render(t *testing.T) {
	sl := SelectList{
		[]SelectListOption{
			{OptionValue(1), "test1"},
			{OptionValue(2), "test2"},
		},
		defaultSelectedIndexList,
	}

	t.Run("render snapshot", func(t *testing.T) {
		expected := "\x1b[2A\x1b[2K\r> 1. test1\n\x1b[2K\r  2. test2\n"
		renderedListBuf := writePrintedContent(sl.render)

		if res := renderedListBuf.String(); res != expected {
			t.Errorf("expected %q, got %q", expected, res)
		}
	})
}

func Test_validateOptions(t *testing.T) {
	data := []struct {
		options           []SelectListOption
		isMustReturnError bool
		name              string
	}{
		{[]SelectListOption{}, true, "empty options list"},
		{[]SelectListOption{
			{OptionValue(1), "test1"},
			{OptionValue(1), "test1"},
		}, true, "options list with duplicates values"},
		{selectListOptionsData, false, "valid options list"},
	}

	for _, tc := range data {
		t.Run(tc.name, func(t *testing.T) {
			err := validateOptions(tc.options)
			if (tc.isMustReturnError && err == nil) || (!tc.isMustReturnError && err != nil) {
				t.Errorf("is expect error '%t', got %s", tc.isMustReturnError, err)
			}
		})
	}
}

func TestRun(t *testing.T) {
	// TODO: add DI to Run func, implement test with mock deps
}

func TestNewSelectList(t *testing.T) {
	validSL := SelectList{
		options:       selectListOptionsData,
		selectedIndex: defaultSelectedIndexList,
	}
	data := []struct {
		options       []SelectListOption
		expectedRes   *SelectList
		isExpectedErr bool
		name          string
	}{
		{[]SelectListOption{}, nil, true, "empty options list"},
		{[]SelectListOption{
			{OptionValue(1), "test1"},
			{OptionValue(1), "test1"},
		}, nil, true, "options list with no unique values"},
		{selectListOptionsData, &validSL, false, "valid options list"},
	}

	for _, tc := range data {
		t.Run(tc.name, func(t *testing.T) {
			res, err := NewSelectList(tc.options)
			if !isSelectListsEqual(res, tc.expectedRes) {
				t.Errorf("expected %q, got %q", tc.expectedRes, res)
			}

			if tc.isExpectedErr == true && err == nil {
				t.Error("expected an error, got nil")
			}

			if tc.isExpectedErr == false && err != nil {
				t.Error("dont expect error, but is not nil")
			}
		})
	}

	t.Run("default selected index is 0", func(t *testing.T) {
		res, _ := NewSelectList(selectListOptionsData)
		if res.selectedIndex != 0 {
			t.Errorf("default selected index must be 0, but got %d", res.selectedIndex)
		}
	})
}

func writePrintedContent(printFunc func()) bytes.Buffer {
	oldStdout := os.Stdout

	r, w, _ := os.Pipe()
	os.Stdout = w

	printFunc()

	w.Close()

	var buf bytes.Buffer
	io.Copy(&buf, r)

	os.Stdout = oldStdout
	return buf
}

func isSelectListsEqual(sl1, sl2 *SelectList) bool {
	if sl1 == nil && sl2 == nil {
		return true
	}

	if (sl1 == nil && sl2 != nil) || (sl1 != nil && sl2 == nil) {
		return false
	}

	if len(sl1.options) != len(sl2.options) {
		return false
	}

	if sl1.selectedIndex != sl2.selectedIndex {
		return false
	}

	res := true
	for i, o := range sl1.options {
		opt := sl2.options[i]
		if o.Value != opt.Value || o.Text != opt.Text {
			res = false
			break
		}
	}

	return res
}
