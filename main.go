package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/spf13/viper"
)

func readConfig() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	viper.AddConfigPath(home)
	viper.SetConfigName("./bookholder")
	viper.SetConfigType("yaml")

	viper.SetEnvPrefix("BOOKHOLDER")
	viper.AutomaticEnv()
}

func setDefaults() {
	viper.SetDefault("Server", "localhost")
	viper.SetDefault("Port", "8080")
	viper.SetDefault("User", "admin")
	viper.SetDefault("Password", "admin")
}

func getData(route string) (string, error) {
	res, err := http.Get("http://" + viper.GetString("Server") + ":" + viper.GetString("Port") + route)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	bodyB, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	body := string(bodyB)

	return body, nil
}

func main() {
	readConfig()
	setDefaults()

	functionList := []string{"Overview", "Capture", "Analysis", "Report", "Settings", "Help"}

	a := app.New()
	w := a.NewWindow("Bookholder")

	w.Resize(fyne.NewSize(1200, 800))

	listView := widget.NewList(
		func() int {
			return len(functionList)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(functionList[i])
		},
	)

	var content = container.NewVBox(widget.NewLabel("Select a function from the list on the left"))
	listView.OnSelected = func(id widget.ListItemID) {
		content.Objects = nil

		switch functionList[id] {
		case "Overview":
			content.Add(widget.NewLabel("Overview"))
		case "Capture":
			content.Add(widget.NewLabel("Capture"))
		case "Analysis":
			headings := []string{"Amount", "Debit", "Offset Account", "Date", "Description"}
			var tableContent [][]string

			tableContent = append(tableContent, headings)

			table := widget.NewTable(
				func() (int, int) {
					return len(tableContent), len(tableContent[0])
				},
				func() fyne.CanvasObject {
					return widget.NewLabel("Cell")
				},
				func(i widget.TableCellID, o fyne.CanvasObject) {
					o.(*widget.Label).SetText(tableContent[i.Row][i.Col])
				},
			)

			// Set column widths (adjust values as needed)
			table.SetColumnWidth(0, 100) // Amount
			table.SetColumnWidth(1, 100) // Debit
			table.SetColumnWidth(2, 150) // Offset Account
			table.SetColumnWidth(3, 120) // Date
			table.SetColumnWidth(4, 200) // Description

			// Wrap in scroll container to avoid truncation
			tableContainer := container.NewScroll(container.NewMax(table))
			tableContainer.SetMinSize(fyne.NewSize(900, 600))

			accountNumber := widget.NewEntry()
			accountNumber.PlaceHolder = "Account number"

			accountNumber.OnSubmitted = func(s string) {
				print("Submitted: " + s)
			}

			canvas := container.NewBorder(
				accountNumber,
				tableContainer,
				nil,
				nil,
				nil,
			)
			content.Add(canvas)
		case "Report":
			content.Add(widget.NewLabel("Report"))
		case "Settings":
			value1 := widget.NewEntry()
			value1.PlaceHolder = viper.GetString("Server")
			value2 := widget.NewEntry()
			value2.PlaceHolder = viper.GetString("Port")
			value3 := widget.NewEntry()
			value3.PlaceHolder = viper.GetString("User")
			value4 := widget.NewPasswordEntry()
			value4.PlaceHolder = viper.GetString("Password")
			form := &widget.Form{
				Items: []*widget.FormItem{
					{Text: "Server", Widget: value1},
					{Text: "Port", Widget: value2},
					{Text: "User", Widget: value3},
					{Text: "Password", Widget: value4},
				},
				OnSubmit: func() {
					viper.Set("Server", value1.Text)
					viper.WriteConfig()
				},
			}
			content.Add(form)
		case "Help":
			content.Add(widget.NewLabel("Help"))
		default:
			content.Add(widget.NewLabel("Unknown function"))
		}

		content.Refresh()
	}

	listView.Select(0)

	split := container.NewHSplit(
		listView,
		content,
	)

	split.Offset = 0.2

	w.SetContent(split)

	w.ShowAndRun()
}
