package views

import (
	"fmt"
	"log"
	"sync"

	"github.com/jroimartin/gocui"
	"github.com/zerodoctor/homeapi-front/channel"
	"github.com/zerodoctor/homeapi-front/consumer"
	"github.com/zerodoctor/homeapi-front/model"
)

var (
	treeView *gocui.View

	root = model.NewLabel(
		"root", "root", true, false, 0, 0, nil,
	)
)

// SetTreeView :
func SetTreeView(g *gocui.Gui, maxX int, maxY int) error {
	if v, err := g.SetView("tree", 0, 0, maxX/5, (maxY-(maxY/15))-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "tree"
		v.Wrap = false
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack

		treeView = v
		//InTreeChan <- NewData(0, 0, false, "", model.EmptyFileFolder())
		if _, err := SetCurrentViewOnTop(g, "tree"); err != nil {
			return err
		}

		channel.InTreeChan <- channel.NewData("init", 0, false, "", nil)
	}

	return nil
}

// currentBuffer :
//	 could be very large (memory wise) but hopefully user doesn't need to keep every folder open
// 	 They'll probably will
var currentBuffer []model.Label

// PrintTreeView :
func PrintTreeView(g *gocui.Gui, wg *sync.WaitGroup) {
	defer wg.Done()
	for data := range channel.InTreeChan {
		if treeView != nil {
			err := showTreeView(g, data)
			if err != nil {
				log.Panicln(err)
			}
			channel.InStatusChan <- Logging("", "Done!", false)
		}
	}
}

// not sure if we need to send Data struct in here
func showTreeView(g *gocui.Gui, data channel.Data) error {
	if len(currentBuffer) <= 0 {
		channel.InStatusChan <- Logging("Info", "initalizing tree buffer... ", true)
		currentBuffer = append(currentBuffer, root)
		printTreeView(g)
		return nil
	}
	if data.Integer >= len(currentBuffer) || data.Integer <= -1 {
		channel.InStatusChan <- Logging("Error", "index is outside the current tree buffer", true)
		return nil
	}

	parent := currentBuffer[data.Integer]

	switch data.Type {
	case "open":
		return openTreeView(g, parent, data)
	case "refresh":
		return refreshTreeView(g, parent, data)
	case "view":
	}

	return nil

}

func refreshTreeView(g *gocui.Gui, parent model.Label, data channel.Data) error {

	_, t := getFileData(parent, data)
	if t != 0 {
		return nil
	}
	printTreeView(g)

	if data.Type == "refresh" {
		content := consumer.GetFolderContent(parent.ID)
		channel.InScreenChan <- channel.NewData("open", data.Integer, false, parent.ID, content)
	}

	return nil
}

func getFileData(parent model.Label, data channel.Data) (model.Label, int) {
	channel.InStatusChan <- Logging("Info", "fecthing "+parent.Name+" children... ", true)
	//InStatusChan <- Logging("Sending Request: " + (time.Now()).String(), true)
	err := consumer.CheckFolderContent(parent.ID)
	//InStatusChan <- Logging("Recevice Response: " + (time.Now()).String(), true)
	if err != nil {
		channel.InStatusChan <- Logging("Warning", "couldn't fetch "+parent.ID+" from db ", true)
		return parent, 1
	}

	reponse := consumer.GetLabelContent(parent.ID)
	if len(reponse) <= 0 {
		channel.InStatusChan <- Logging("Info", parent.ID+" is empty", true)
		channel.InScreenChan <- channel.NewData("open", data.Integer, false, parent.ID, nil)
		return parent, 1
	}
	parent.Children = reponse
	currentBuffer[data.Integer] = parent

	return parent, 0
}

func openTreeView(g *gocui.Gui, parent model.Label, data channel.Data) error {
	parent.Open = !parent.Open
	currentBuffer[data.Integer] = parent

	if parent.Open {
		channel.InStatusChan <- Logging("Info", "opening "+parent.Name+"... ", true)
		if parent.Children == nil || data.Type == "refresh" {
			result, t := getFileData(parent, data)
			if t != 0 {
				return nil
			}
			parent = result
		}

		content := consumer.GetFolderContent(parent.ID)
		channel.InScreenChan <- channel.NewData("open", data.Integer, false, parent.ID, content)

		count := 0
		for _, child := range parent.Children {
			if child.Folder {
				child.Depth = parent.Depth + 1
				child.Index = data.Integer + count + 1
				currentBuffer = model.Insert(currentBuffer, child, child.Index)
				count++
			}
		}
		currentBuffer[data.Integer] = parent
		printTreeView(g)
		updateScreenView(g, parent.Name)
	} else {
		channel.InStatusChan <- Logging("Info", "closing "+parent.Name+"... ", true)

		total := 0
		parent.Index = data.Integer
		for i := parent.Index + 1; i < len(currentBuffer); i++ {
			if parent.Depth >= currentBuffer[i].Depth {
				break
			}
			total++
		}
		currentBuffer = model.RemoveFrom(currentBuffer, parent.Index+1, parent.Index+total+1)
		currentBuffer[data.Integer] = parent
		printTreeView(g)
	}
	return nil
}

// could move this to the screen.go file
func updateScreenView(g *gocui.Gui, title string) {
	g.Update(func(g *gocui.Gui) error {
		screenView, err := g.View("screen")
		if err == nil {
			screenView.Title = title
		}
		return nil
	})
}

func printTreeView(g *gocui.Gui) {
	g.Update(func(g *gocui.Gui) error {
		treeView.Clear()
		var str string
		for _, label := range currentBuffer {
			for i := 0; i < label.Depth; i++ {
				str += "| "
			}
			if label.Folder && !label.Open {
				str += "+ "
			} else if label.Folder && label.Open {
				str += "- "
			} else {
				str += "  "
			}
			str += label.Name + "\n"
		}
		fmt.Fprint(treeView, str)

		return nil
	})
}
