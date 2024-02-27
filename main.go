package main

import (
	"bufio"
	"fmt"
	"image/color"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Theme customization
type myTheme struct{}

var _ fyne.Theme = (*myTheme)(nil)

func (m myTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == "button" {
		return color.RGBA{R: 26, G: 140, B: 255, A: 255}
	}
	return theme.DefaultTheme().Color(name, variant)
}
func (m myTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}
func (m myTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}
func (m myTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == "text" {
		return 30
	}
	return theme.DefaultTheme().Size(name)
}

const GAME_LENGTH_MINS = 3

func main() {
	words := NewTrie()
	LoadWordsListFromFile(words)

	a := app.New()
	a.Settings().SetTheme(&myTheme{})

	res, _ := fyne.LoadResourceFromPath("resources/icon.png")
	a.SetIcon(res)

	w := a.NewWindow("BoggleGo")
	wSeed := a.NewWindow("BoogleGo - Seed")

	// Seed input
	var seed int64
	lbl := widget.NewLabel("Enter game seed (blank for random)")
	entrySeed := widget.NewEntry()
	btn := widget.NewButton("Go!", func() {
		n, err := strconv.ParseInt(entrySeed.Text, 10, 64)
		if err != nil {
			n = time.Now().UnixNano()
		}
		seed = n
		println("Using seed: ", n)
		w.SetTitle(fmt.Sprintf("BoogleGo (%d)", seed))
		StartGame(seed, w, words)
		wSeed.Close()
	})

	entrySeed.Validator = validation.NewRegexp("^[0-9]*$", "Invalid input")

	lastText := entrySeed.Text
	entrySeed.OnChanged = func(s string) {
		if entrySeed.Validate() == nil {
			lastText = entrySeed.Text
			return
		}
		entrySeed.Text = lastText
	}

	entrySeed.OnSubmitted = func(s string) {
		btn.Tapped(nil)
	}

	inputSeed := container.NewVBox(lbl, entrySeed, btn)
	wSeed.SetContent(inputSeed)
	wSeed.Canvas().Focus(entrySeed)
	wSeed.SetFixedSize(true)
	wSeed.Show()

	a.Run()
}

func LoadWordsListFromFile(trie *Trie) {
	file, _ := os.Open("resources/palabras_esp.txt")
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		trie.Insert(scanner.Text())
	}
}

func end(window fyne.Window, words map[string]bool) {
	scores := map[int]int{
		4: 1,
		5: 2,
		6: 3,
		7: 4,
		8: 11,
	}

	points := 0
	for w := range words {
		points += scores[min(len(w), 8)]
	}

	d := dialog.NewInformation("Game over!", fmt.Sprintf("Puntaje: %d", points), window)
	d.SetOnClosed(func() {
		os.Exit(0)
	})
	d.Show()
}

func StartGame(seed int64, w fyne.Window, words *Trie) {
	// Buttons grid (dice)
	var buttons []fyne.CanvasObject
	diceSlice := ThrowDice(seed)
	for _, v := range diceSlice {
		if v == "Q" {
			v = "Qu"
		}
		btn := widget.NewButton(v, nil)
		buttons = append(buttons, btn)
	}
	diceGrid := container.New(layout.NewGridLayout(4), buttons...)

	// Guessed words list
	guessedMap := map[string]bool{}
	guessedWords := []string{}
	guessedList := widget.NewList(
		func() int {
			return len(guessedWords)
		},
		func() fyne.CanvasObject {
			txt := canvas.NewText("templatewithlong", color.White)
			txt.TextSize = 15
			return txt
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*canvas.Text).Text = guessedWords[i]
		})

	// Initialize board
	board := NewBoard(diceSlice)

	//Entry (input)
	entry := widget.NewEntry()
	entry.PlaceHolder = "Ingresa la palabra..."
	entry.OnSubmitted = func(s string) {
		runes := []rune(strings.ReplaceAll(strings.ToLower(s), "qu", "q"))
		found := WordExists(runes, *board) && words.Search(s)
		_, exists := guessedMap[s]
		if found && !exists {
			guessedWords = append(guessedWords, s)
			sort.Strings(guessedWords)
			guessedMap[s] = true
		}
		board.Reset()

		guessedList.Refresh()
		guessedList.ScrollToBottom()
		entry.SetText("")
	}

	// Timer
	remainingTime := time.Date(0, 0, 0, 0, GAME_LENGTH_MINS, 0, 0, time.UTC)
	remainingTimeString := remainingTime.Format("4:05")
	timer := widget.NewLabel(remainingTimeString)
	timer.Alignment = fyne.TextAlignCenter
	timer.Importance = widget.SuccessImportance
	go func() {
		for range time.Tick(time.Second) {
			remainingTime = remainingTime.Add(time.Second * -1)
			if remainingTime.Minute() == 0 && remainingTime.Second() <= 30 {
				timer.Importance = widget.DangerImportance
			}
			remainingTimeString = remainingTime.Format("4:05")
			timer.SetText(remainingTimeString)
			if remainingTime.Minute() == 0 && remainingTime.Second() == 0 {
				end(w, guessedMap)
				break
			}
		}
	}()
	listAndTimer := container.NewBorder(layout.NewSpacer(), timer, layout.NewSpacer(), layout.NewSpacer(), guessedList)

	// Found words title
	headerLabel := canvas.NewText("Encontradas:", color.RGBA{
		R: 0,
		G: 157,
		B: 255,
		A: 255,
	})
	headerLabel.TextSize = 16
	headerLabel.Alignment = fyne.TextAlignCenter
	headerLabel.TextStyle.Bold = true
	rightPanel := container.NewBorder(headerLabel, layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(), listAndTimer)

	// Build GUI and start
	content := container.NewBorder(layout.NewSpacer(), entry, layout.NewSpacer(), rightPanel, diceGrid)
	w.SetContent(content)
	w.Resize(fyne.NewSize(520, 500))
	w.SetFixedSize(true)
	w.Canvas().Focus(entry)

	w.SetMaster()
	w.Show()
}
