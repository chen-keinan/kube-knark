package ui

import (
	"bytes"
	"strings"

	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

const (
	// EmptyCell represents an empty cell in the table.
	EmptyCell = ""
)

const (
	// number of rows to skip based on the
	// widgets.Table.Draw method.
	skipRows = 2
)

// TableRows represent the rows of the table.
type TableRows [][]string

// Table is extending the Table in the termui/v3/widgets/ package
// to support scrolling/filtering.
type Table struct {
	*widgets.Table
	// out are the rows that are going to be rendered.
	out TableRows
	//Rows are the rows of the table.
	Rows TableRows
	// offset is the offset from the first row of the table.
	offset int
	// visibleRows is the number of rows that will be displayed.
	visibleRows int
	// current height of the table.
	height int
	// curr represents the position in the table within the bounds of [offset, visibleRows).
	curr int
	// prev is the previous position curr held.
	prev int
	// filter to be applied on the table rows
	filter *bytes.Buffer
	// column on which the filter should be applied
	filterColumn int
	// number of rows per entry in the table
	rowsPerEntry int

	// colors which will be used to paint the table rows.
	Colors struct {
		// default color of each row
		Text termui.Color
		// text color of the selected row
		SelectedRowFg termui.Color
		// row color of the selected row
		SelectedRowBg termui.Color
	}
}

// NewTable returns a default instance of xtui.Table.
func NewTable(lightTheme bool) *Table {
	t := &Table{
		Table:        widgets.NewTable(),
		out:          nil,
		Rows:         nil,
		offset:       0,
		visibleRows:  0,
		height:       0,
		curr:         0,
		prev:         0,
		filter:       new(bytes.Buffer),
		filterColumn: -1,
		rowsPerEntry: 1,
	}
	// Default colors
	if lightTheme {
		t.Colors.Text = termui.ColorBlack
	} else {
		t.Colors.Text = termui.ColorWhite
	}
	t.Colors.SelectedRowFg = termui.ColorBlack
	t.Colors.SelectedRowBg = termui.ColorGreen
	return t
}

// paintActiveRows paints the active row in the
// specified table.
func (t *Table) paintActiveRow() {
	t.RowStyles[t.prev] = termui.NewStyle(t.Colors.Text)
	t.RowStyles[t.curr] = termui.NewStyle(t.Colors.SelectedRowFg, t.Colors.SelectedRowBg, termui.ModifierBold)
}

// SetRect resize the table, and correctly sets the height of the table.
// x,y : top left corner
// x1,y2 : bottom right corner
func (t *Table) SetRect(x, y, x1, y2 int) {
	t.Table.SetRect(x, y, x1, y2)

	// new number of visible rows is less than
	// the old number of visible rows.
	if y2-y-skipRows < t.visibleRows {
		t.prev = t.curr
		t.curr--
		if t.curr < 0 {
			t.prev = t.curr
			t.curr = 0
		}
	}

	t.height = y2 - y
	t.visibleRows = t.height - skipRows

	t.paintActiveRow()
}

// ScrollUp scrolls the table one row up
func (t *Table) ScrollUp() {
	if t.curr > 0 {
		t.prev = t.curr
		t.curr--
	} else {
		if t.offset > 0 {
			t.offset--
		}
	}

	t.paintActiveRow()
}

// ScrollDown scrolls the table one row down
func (t *Table) ScrollDown() {
	if t.curr < t.visibleRows-1 {
		t.prev = t.curr
		t.curr++
	} else {
		if t.offset+t.visibleRows < len(t.out) {
			t.offset++
		}
	}

	t.paintActiveRow()
}

// reCalcView recalculates the view into the table, handling any out of bounds errors.
func (t *Table) reCalcView() {
	if len(t.out) == 0 {
		return
	}
	// Adjust the visible rows based on the available
	// height.
	if t.visibleRows < t.height-skipRows {
		t.visibleRows = t.height - skipRows
	}
	// Avoid overflow if the number of displayed rows
	// is greater than the number of available rows.
	if t.offset+t.visibleRows > len(t.out) {
		t.visibleRows = len(t.out) - t.offset
	}
	// Avoid underflow if the table height is less than the top left
	// corner of the table.
	if t.visibleRows < 0 {
		t.visibleRows = 0
	}
	t.Table.Rows = t.out[t.offset : t.offset+t.visibleRows]
}

// Draw extends the method Draw from tui.Table to also include filtering.
func (t *Table) Draw(buf *termui.Buffer) {
	if t.filter.String() != "" && t.filterColumn >= 0 {
		var filteredRows [][]string
		for i := 0; i < len(t.Rows); i += t.rowsPerEntry {
			if strings.Contains(t.Rows[i][t.filterColumn], t.filter.String()) {
				for r := 0; r < t.rowsPerEntry; r++ {
					filteredRows = append(filteredRows, t.Rows[i+r])
				}
			}
		}

		// if no match against the filter
		if len(filteredRows) == 0 {
			// make an empty table based on the number of columns
			// of the last render.
			// NOTE: if the number of columns changes might produce
			// unwanted behavior
			if len(t.Table.Rows) != 0 {
				columns := len(t.Table.Rows[0])
				filteredRows = [][]string{
					make([]string, columns),
				}
			}
		}
		t.out = filteredRows
	} else {
		t.out = t.Rows
	}

	t.reCalcView()
	// Avoid panic in the termui/table draw method, if no rows are supplied by the user.
	if len(t.Table.Rows) == 0 {
		return
	}

	t.paintActiveRow()
	t.Table.Draw(buf)
}
