package display

import (
	"errors"
	"fmt"
	"github.com/bgartzi/uhhmm/filters"
	"github.com/bgartzi/uhhmm/host"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
	"strconv"
	"strings"
)

type HostDisplayer struct {
	filters filters.FilterChain
	columns []string
	table   table.Writer
}

// FKA ID: Hosts don't have ID written on theirselves. Need a way of adding that
// extra field name somehow so it's "not hardcoded" everywhere
const HostPositionColName string = "Index"

func DefaultHostDisplayerConfig() HostDisplayer {
	return HostDisplayer{
		columns: []string{HostPositionColName, "Address", "NickName", "Info", "Labels"},
		table:   table.NewWriter(),
	}
}

func (displayer *HostDisplayer) AddFilter(filter filters.HostFilter) {
	displayer.filters.Append(filter)
}

func (displayer *HostDisplayer) applyFilters(hosts []host.Host) []host.Host {
	return displayer.filters.Apply(hosts)
}

func (displayer *HostDisplayer) ConfigColumns(columns []string) {
	displayer.columns = columns
}

func (displayer *HostDisplayer) EmptyColumnsConfig() {
	displayer.columns = []string{}
}

func (displayer *HostDisplayer) AddColumn(column string) {
	displayer.columns = append(displayer.columns, column)
}

func (displayer *HostDisplayer) getHeader() ([]string, error) {
	if len(displayer.columns) == 0 {
		return []string{}, errors.New("Empty list of columns to display")
	}
	return displayer.columns, nil
}

func (displayer *HostDisplayer) configOutputTable() error {
	displayer.table.SetOutputMirror(os.Stdout)
	displayer.table.SetStyle(table.Style{
		Name: "empty",
		Box: table.BoxStyle{
			PaddingRight: "   ",
		},
		Options: table.Options{
			DrawBorder:      false,
			SeparateColumns: true,
			SeparateFooter:  false,
			SeparateHeader:  false,
			SeparateRows:    false,
		},
	})
	return nil
}

func rowFromSlice(fields []string) table.Row {
	ret := make([]interface{}, len(fields))
	for i, field := range fields {
		ret[i] = field
	}
	return table.Row(ret)
}

func (displayer *HostDisplayer) buildTableHeader() {
	displayer.table.AppendHeader(rowFromSlice(displayer.columns))
}

func (displayer *HostDisplayer) toFields(host host.Host, pos int) ([]interface{}, error) {
	fields := make([]interface{}, len(displayer.columns))
	for i_field, fieldName := range displayer.columns {
		switch fieldName {
		case HostPositionColName:
			fields[i_field] = strconv.Itoa(pos)
		case "Address":
			fields[i_field] = host.Address
		case "Port":
			fields[i_field] = host.Port
		case "NickName":
			fields[i_field] = host.NickName
		case "Info":
			fields[i_field] = host.Info
		case "User":
			fields[i_field] = host.User
		case "Labels":
			fields[i_field] = strings.Join(host.Labels, ",")
		default:
			return nil, fmt.Errorf("Unsupported host field %s", fieldName)
		}
	}
	return fields, nil
}

func (displayer *HostDisplayer) buildTableData(hosts []host.Host) error {
	for pos, host := range hosts {
		fields, err := displayer.toFields(host, pos)
		if err != nil {
			return err
		}
		displayer.table.AppendRow(fields)
	}
	return nil
}

func (displayer *HostDisplayer) renderTable() {
	displayer.table.Render()
}

func (displayer *HostDisplayer) Display(hosts []host.Host) error {
	err := displayer.configOutputTable()
	if err != nil {
		return err
	}
	displayer.buildTableHeader()
	hosts = displayer.applyFilters(hosts)
	err = displayer.buildTableData(hosts)
	if err != nil {
		return err
	}
	displayer.renderTable()
	return nil
}
