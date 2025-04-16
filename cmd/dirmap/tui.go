package main

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Generate list of directories in tview. There is no function to act on these things.
// The primary name of each element is the name of the directory
// Needs the root directory
func generateDirList(list *tview.List, root_dir string) error {

	if stat, err := os.Stat(root_dir); err != nil {
		return err
	} else {

		// Check if we are actually looking at a real directory
		if !stat.IsDir() {
			return errors.New(root_dir + " is not a real directory!")
		}

		// Otherwise, we have a real directory and now we can generate our list of elements in the directory

		files, err := os.ReadDir(root_dir)
		if err != nil {
			return errors.New("Cannot read directory " + root_dir)
		}

		// Add the .. Directory too
		list.AddItem("..", "", 0, nil)

		// Generate list of files and directories
		for _, node := range files {

			// Only add to the list if the element is a directory
			if node.IsDir() {
				list.AddItem(node.Name(), "", 0, nil)
			}

		}

		return nil

	}

}

func createSourceSelectionBox(root_dir string) (*tview.List, error) {

	// Create our list
	// This is the absolute path to the source directory
	// This will keep track of where we are
	// This will never be ended with a /
	abs_source_directory, err := filepath.Abs(root_dir)
	if err != nil {
		return nil, err
	}

	// Clean the path
	abs_source_directory = filepath.Clean(abs_source_directory)

	src_list := tview.NewList()
	err = generateDirList(src_list, root_dir)
	if err != nil {
		return nil, err
	}

	// Add the update function
	// This will make sure we get to select the directory
	src_list.SetSelectedFunc(func(index int, primary_name string, secondary_name string, shortcut rune) {

		// Accept the selection and move around the tree
		abs_source_directory += "/" + primary_name
		abs_source_directory = filepath.Clean(abs_source_directory)

		// Generate the new information
		src_list.Clear()
		err := generateDirList(src_list, abs_source_directory)
		if err != nil {
			panic(err)
		}

		// Set the title of the window
		src_list.SetTitle(filepath.Base(abs_source_directory))

	})

	// Set the starting title
	src_list.SetTitle(filepath.Base(abs_source_directory))
	src_list.Box.SetBorder(true)

	return src_list, nil

}

func createDestinationSelectionBox(root_dir string) (*tview.List, error) {

	// Deal with the destination directory
	abs_destination_directory, err := filepath.Abs(root_dir)
	if err != nil {
		return nil, err
	}

	// Set up the destination directory
	abs_destination_directory = filepath.Clean(abs_destination_directory)
	dst_list := tview.NewList()
	err = generateDirList(dst_list, abs_destination_directory)
	dst_list.SetSelectedFunc(func(index int, primary_name string, secondary_name string, shortcut rune) {

		// Accept the selection and move around the tree
		abs_destination_directory += "/" + primary_name
		abs_destination_directory = filepath.Clean(abs_destination_directory)

		// Update the list
		dst_list.Clear()
		err := generateDirList(dst_list, abs_destination_directory)
		if err != nil {
			panic(err.Error())
		}

		// Set the title of the window
		dst_list.SetTitle(filepath.Base(abs_destination_directory))

	})

	// Set the starting title
	dst_list.SetTitle(filepath.Base(abs_destination_directory))
	dst_list.Box.SetBorder(true)

	return dst_list, nil

}

// This will create our default form box for adding filters
func createDefaultFilerFunctionBox() *tview.Form {

	form := tview.NewForm().AddTextArea("Function", "", 0, 0, 0, nil)
	return form
}

// This will be the box which holds the filters applied to each file in the map
func createFilterBox(app *tview.Application) *tview.Flex {


	// Default box attributes
	list := tview.NewFlex()
	list.SetDirection(tview.FlexRow)
	list.Box.SetBorder(true)
	list.AddItem(createDefaultFilerFunctionBox(), 0 ,1, true)

	// Deal with input
	filter_focus_index := 0
	list.SetInputCapture(func(capture *tcell.EventKey) *tcell.EventKey {


		// If we have the add in key
		if capture.Key() == tcell.KeyF1 {
			list.AddItem(createDefaultFilerFunctionBox(), 0 ,1, true)

			// Change our focus
			ele_count := list.GetItemCount()
			app.SetFocus(list.GetItem(ele_count - 1))

		// F2 will be our delete function key
		} else if capture.Key() == tcell.KeyF2 {
			ele_count := list.GetItemCount()
			if ele_count > 0 {
				prim := list.GetItem(ele_count - 1)

				// If this is the item that has the focus
				if prim.HasFocus() {

					// If it was the last element
					if ele_count == 1 {
						app.SetFocus(list)
					} else { // This is not the last element
						app.SetFocus(list.GetItem(ele_count - 2))
					}

				}

				// Remove the element from the list
				list.RemoveItem(prim)

				// If it was the last item
				// Return the focus
				if ele_count == 1 {
					app.SetFocus(list)
				}

			}

		// This will move down the function filter list
		} else if capture.Key() == tcell.KeyDown {

			// Move down one filter
			filter_focus_index += 1
			ele_count := list.GetItemCount()
			if filter_focus_index >= ele_count {
				filter_focus_index = 0
			}

			// If we have an element to select
			if ele_count > 0 {
				app.SetFocus(list.GetItem(filter_focus_index))
			}

		}

		return capture

	})

	return list

}

// This handles the general TUI stuff and the interactions
// src is the starting source directory
// dst is the destination directory
// An error is returned when something happens
func StartTui(src, dst string) error {

	app := tview.NewApplication()

	src_list, err := createSourceSelectionBox(src)
	if err != nil {
		return err
	}

	dst_list, err := createDestinationSelectionBox(dst)
	if err != nil {
		return err
	}

	filter_box := createFilterBox(app)

	script_output := tview.NewTextArea().SetBorder(true).SetTitle("Script Output")

	// Design and run app
	root := tview.NewFlex().
			AddItem(tview.NewFlex().
				AddItem(tview.NewFlex().
					AddItem(src_list, 0, 1, true).
					AddItem(dst_list, 0, 1, true), 0, 1, true).
				AddItem(script_output, 0, 1, true).SetDirection(tview.FlexRow), 0, 1, true).
			AddItem(filter_box, 0, 1, true)



	app.SetRoot(root, true)

	focus_index := 0 // This will tell us what element of the root has focus
	focus_array := [...]tview.Primitive{src_list, dst_list, filter_box}
	app.SetInputCapture(func(capture *tcell.EventKey) *tcell.EventKey {

		// Close on F10
		if capture.Key() == tcell.KeyF10 {
			app.Stop()
		} else if capture.Key() == tcell.KeyTab { // This will change the focus on the app

			// Change to the next element
			focus_index += 1
			element_count := len(focus_array)

			// Reset the focus_index
			if focus_index >= element_count {
				focus_index = 0
			}

			// If there is an element to focus on
			if element_count > 0 {
				app.SetFocus(focus_array[focus_index])
			}

		// This will be our run the filters key
		} else if capture.Key() == tcell.KeyF9 {

			// TODO: Something here

		}

		return capture

	})


	err = app.Run()
	return err

}
