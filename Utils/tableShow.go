package Utils

import (
	"github.com/gookit/color"
	"github.com/olekukonko/tablewriter"
	"os"
	"sync"
)

var mu sync.Mutex

func TableShow(keys []string, values [][]string) {
	mu.Lock()
	defer mu.Unlock()
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetHeader(keys)
	table.AppendBulk(values)
	table.Render()

}
func DomainTableShow(keys []string, values [][]string, Gong string) {
	mu.Lock()
	defer mu.Unlock()
	color.RGBStyleFromString("205,155,29").Println("\n", Gong, " 查询")
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetHeader(keys)
	table.AppendBulk(values)
	table.Render()

}
func PortTableShow(keys []string, values [][]string, Gong string) {
	mu.Lock()
	defer mu.Unlock()
	color.RGBStyleFromString("205,155,29").Println("\n", Gong)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetHeader(keys)
	table.AppendBulk(values)
	table.Render()

}
