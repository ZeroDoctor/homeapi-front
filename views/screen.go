package views

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"strconv"
	"sync"

	"github.com/jroimartin/gocui"
	"github.com/zerodoctor/homeapi-front/channel"
	"github.com/zerodoctor/homeapi-front/consumer"
	"github.com/zerodoctor/homeapi-front/model"
	"github.com/zerodoctor/homeapi-front/views/ui"
)

var (
	screenView, emptyView, dialogView *gocui.View
	x0, y0, x1, y1                    int
	conX, conY                        int
	currentParent                     int
	currentParentID, uploadPath       string
	currentFolder                     []model.FileFolder
)

// SetScreenView :
func SetScreenView(g *gocui.Gui, maxX int, maxY int) error {

	x0 = maxX/5 + 1
	y0 = 0
	x1 = maxX - 1
	y1 = (maxY - (maxY / 15)) - 2
	conX = maxX
	conY = maxY

	if v, err := g.SetView("screen", x0, y0, x1, y1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		g.Cursor = false
		v.Title = ""
		v.Wrap = false
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		currentParent = -1

		screenView = v
	}

	return nil
}

// PrintScreenView :
func PrintScreenView(g *gocui.Gui, wg *sync.WaitGroup) {
	defer wg.Done()

	setInputKeybind(g)

	for data := range channel.InScreenChan {
		if emptyView != nil {
			g.DeleteView("empty")
			emptyView = nil
		}
		if screenView != nil {
			err := showScreenView(g, data)
			if err != nil {
				log.Panicln(err)
			}
			channel.InStatusChan <- Logging("", "Done!", false)
		}
	}
}

func showScreenView(g *gocui.Gui, data channel.Data) error {
	if data.Type == "open" {
		if data.File == nil || len(data.File) <= 0 {
			err := ui.SetInfoView(g, screenView, emptyView, "this folder is empty", conX, x0, y1)
			currentParent = -1
			currentFolder = nil
			if err != nil {
				return err
			}

			channel.InStatusChan <- Logging("Info", "this folder is empty... ", true)

			return nil
		}
		currentParent = data.Integer
		return printScreenView(g, data.File, data.String.(string))
	}

	if data.Integer < 0 || data.Integer >= (len(currentFolder)+1) {
		channel.InStatusChan <- Logging("Error", "index out of bounds "+strconv.Itoa(data.Integer)+"... ", true)
		return nil
	}

	switch data.Type {
	case "view":
		return openScreenView(g, data.Integer)
	case "download":
		return downloadScreenView(g, data.Integer)
	case "question":
		return questionScreenView(g, data)
	case "delete":
		return deleteScreenView(g, data.Integer)
	case "refresh":
		return refreshScreenView(g, data)
	case "upload":
		if dialogView == nil {
			err := ui.SetInputBox(g, dialogView, "upload to: ", conX, x0, y1)
			if err != nil {
				return err
			}
			SetCurrentViewOnTop(g, "dialog")
		}
	case "cancel":
		return cancelScreenView(g)
	}

	return nil
}

func setInputKeybind(g *gocui.Gui) {
	err := g.SetKeybinding("dialog", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		SetCurrentViewOnTop(g, "screen")
		if dialogView != nil {
			g.DeleteView("dialog")
			dialogView = nil
		}
		uploadScreenView(g, uploadPath)
		return nil
	})

	if err != nil {
		log.Panicln(err)
	}
}

// the reason why this is here and not in status.go because of the currentFolder var
// status.go can access the var but I like private vars to be accessed in one file instead of the whole folder
func questionScreenView(g *gocui.Gui, data channel.Data) error {
	file := currentFolder[data.Integer]

	SetCurrentViewOnTop(g, "status")
	channel.InStatusChan <- channel.NewData(
		data.String.(string), data.Integer, true,
		"are you sure you want to "+data.String.(string)+" "+file.FullName+" (y/n) ",
		nil,
	)

	return nil
}

func cancelScreenView(g *gocui.Gui) error {
	channel.InStatusChan <- Logging("Info", "cancelling operation... ", true)
	if dialogView != nil {
		g.DeleteView("dialog")
		dialogView = nil
	}
	return nil
}

func uploadScreenView(g *gocui.Gui, path string) error {

	return nil
}

func deleteScreenView(g *gocui.Gui, index int) error {
	file := currentFolder[index]

	channel.InStatusChan <- Logging("Info", "deleting "+file.FullName+"... ", true)
	result, err := consumer.DeleteFolder(file)
	if err != nil {
		return err
	}
	channel.InStatusChan <- Logging("Info", result+" ", true)
	channel.InTreeChan <- channel.NewData("refresh", currentParent, false, "", nil)

	return nil
}

func refreshScreenView(g *gocui.Gui, data channel.Data) error {
	return nil
}

func downloadScreenView(g *gocui.Gui, index int) error {
	file := currentFolder[index]

	channel.InStatusChan <- Logging("Info", "processing download "+file.ID+"...", true)
	path, err := consumer.DownloadFolder(file)
	if err != "" {
		channel.InStatusChan <- Logging("Warning", "Something went wrong my guy "+err, true)
		return nil
	}
	channel.InStatusChan <- Logging("Info", "the file "+file.FullName+" has been downloaded to "+path+" ", true)
	if file.Dir == 1 {
		channel.InTreeChan <- channel.NewData("refresh", currentParent, false, "", nil) // a zip file has been created so lets refresh parent
	}

	return nil
}

func openScreenView(g *gocui.Gui, index int) error {
	if currentFolder[index].Dir == 1 {
		// InStatusChan <- Logging("Info", "processing " + strconv.Itoa(currentParent + index) + "... " , true)
		channel.InTreeChan <- channel.NewData("open", currentParent+index, false, "", nil)
	}

	channel.InStatusChan <- Logging("Info", "can't open a file yet ", true)
	return nil
}

func repeatString(buffer *bytes.Buffer, str string, count int) {
	if count < 0 {
		channel.InStatusChan <- Logging("Error", "count is too long ", true)
		return
	}

	for i := 0; i < count; i++ {
		buffer.WriteString(str)
	}
}

func countDigits(num int) int {
	if num <= 0 { // don't need this anymore at least for this application
		return 1
	}
	return int(math.Floor(math.Log10(float64(num)) + 1))
}

func printScreenView(g *gocui.Gui, file []model.FileFolder, parentID string) error {
	if err := screenView.SetCursor(0, 2); err != nil {
		return err
	}
	var buffer bytes.Buffer
	var headBuffer bytes.Buffer
	var topBarBuffer bytes.Buffer

	sizeName := 0
	longestName := 12   // default to 12 spaces
	biggestFile := 9999 //  the biggest size of a file also defaults to 4 spaces

	channel.InStatusChan <- Logging("Info", "processing "+parentID+"... ", true)
	currentParentID = parentID
	for _, f := range file {
		sizeName = len(f.FullName)
		if sizeName > longestName {
			longestName = sizeName
		}
		if f.Size > biggestFile {
			biggestFile = f.Size
		}
	}

	biggestFile = countDigits(biggestFile)
	end := x1 - x0 - longestName - biggestFile - 19

	for _, f := range file {
		lengthSize := countDigits(f.Size)
		buffer.WriteString(" | " + f.FullName)
		repeatString(&buffer, " ", longestName-len(f.FullName))
		buffer.WriteString(" |  " + strconv.Itoa(int(f.Dir)))
		buffer.WriteString("  | " + strconv.Itoa(f.Size))
		repeatString(&buffer, " ", biggestFile-lengthSize)
		buffer.WriteString(" |")
		repeatString(&buffer, " ", end) // prev: biggestFile - lengthSize
		buffer.WriteString(" | \n")
	}

	headBuffer.WriteString(" | Name")
	repeatString(&headBuffer, " ", longestName-4)
	headBuffer.WriteString(" | Dir | Size")
	repeatString(&headBuffer, " ", biggestFile-4)
	headBuffer.WriteString(" |")
	repeatString(&headBuffer, " ", end) // prev: biggestFile - 4
	headBuffer.WriteString(" | \n")

	topBarBuffer.WriteString(" + ")
	repeatString(&topBarBuffer, "-", x1-x0-8) // prev: longestName + biggestFile + 9
	topBarBuffer.WriteString(" + \n")

	g.Update(func(g *gocui.Gui) error {
		screenView.Clear()
		fmt.Fprintln(screenView, topBarBuffer.String()+headBuffer.String()+buffer.String()+topBarBuffer.String())
		return nil
	})

	currentFolder = file
	return nil
}
