package main

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/alexanderthegreat96/epocell-go/core"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {

	config, _ := core.LoadJsonConfig()

	epocellApp := app.NewWithID("epocell-app")

	epocellWindow := epocellApp.NewWindow("Epocell v1.4")
	icon, _ := fyne.LoadResourceFromPath("icons/epocell-icon.png")
	epocellWindow.SetIcon(icon)

	label := widget.NewLabelWithStyle("Epocell Configuration", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	/**
	Begin implementing the form data
	*/

	// Create a widget for displaying the selected file path
	filePathEntry := widget.NewEntry()
	filePathEntry.SetText(config.EpocellFile)
	filePathEntry.TextStyle = fyne.TextStyle{Bold: true}

	filePathRawDataEntry := widget.NewEntry()
	filePathRawDataEntry.SetText(config.RawDataFile)
	filePathRawDataEntry.TextStyle = fyne.TextStyle{Bold: true}

	// Create a button to open the file picker dialog
	pickerButton := widget.NewButton("Pick File", func() {
		core.PickFile(epocellWindow, filePathEntry)
	})

	// Input with label for epocell starts at row

	epocellStartsAtRowLabel := widget.NewLabel("Epocell Data starts at row:")
	epocellStartsAtRowEntry := widget.NewEntry()

	if config.EpocellStartsAtRow == 0 {
		epocellStartsAtRowEntry.SetText("20")
	} else {
		epocellStartsAtRowEntry.SetText(strconv.Itoa(config.EpocellStartsAtRow))
	}

	epocellStartsAtRowEntry.OnChanged = func(text string) {
		if _, err := strconv.Atoi(text); err != nil && text != "" {
			prev_data := epocellStartsAtRowEntry.Text
			prev_data_int, _ := strconv.Atoi(prev_data)
			// If the input is not a valid integer, revert to the previous value
			epocellStartsAtRowEntry.SetText(strconv.Itoa(prev_data_int))
		}
	}

	epocellStartsAtRowLabelContainer := container.NewHBox(epocellStartsAtRowLabel)
	epocellStartsAtRowLabelContainer.Layout = layout.NewMaxLayout()

	epocellStartsAtRowEntryContainer := container.NewVBox(epocellStartsAtRowEntry)
	epocellStartsAtRowContainer := container.New(layout.NewGridLayoutWithColumns(2), epocellStartsAtRowLabelContainer, epocellStartsAtRowEntryContainer)

	filePathEntryLabelContainer := container.NewHBox(widget.NewLabel("Epocell Excel File:"))
	filePathEntryContainer := container.NewHBox(filePathEntry)
	filePathEntryContainer.Layout = layout.NewMaxLayout()
	pickFileContainer := container.NewHBox(pickerButton)
	pickFileContainer.Layout = layout.NewMaxLayout()
	filePickerContainer := container.New(layout.NewGridLayoutWithColumns(3), filePathEntryLabelContainer, filePathEntryContainer, pickFileContainer)

	storeNameCellLetterLabel := widget.NewLabel("Store Name CELL LETTER:")
	storeNameCellLetterEntry := widget.NewEntry()

	if config.StoreNameCellLetter == "" {
		storeNameCellLetterEntry.SetText("B")
	} else {
		storeNameCellLetterEntry.SetText(config.StoreNameCellLetter)
	}

	storeNameCellLetterEntry.TextStyle = fyne.TextStyle{Bold: true}

	storeNameCellLetterLabelContainer := container.NewHBox(storeNameCellLetterLabel)
	storeNameCellLetterLabelContainer.Layout = layout.NewMaxLayout()

	storeNameCellLetterEntryContainer := container.NewHBox(storeNameCellLetterEntry)
	storeNameCellLetterEntryContainer.Layout = layout.NewMaxLayout()

	storeNameContainer := container.New(layout.NewGridLayoutWithColumns(2), storeNameCellLetterLabelContainer, storeNameCellLetterEntryContainer)

	storeSalesUnitsCellLetterLabel := widget.NewLabel("Store Sales Units CELL LETTER:")
	storeSalesUnitsCellLetterEntry := widget.NewEntry()

	if config.StoreSalesUnitsCellLetter == "" {
		storeSalesUnitsCellLetterEntry.SetText("J")
	} else {
		storeSalesUnitsCellLetterEntry.SetText(config.StoreSalesUnitsCellLetter)
	}

	storeSalesUnitsCellLetterEntry.TextStyle = fyne.TextStyle{Bold: true}

	storeSalesUnitsCellLetterLabelContainer := container.NewHBox(storeSalesUnitsCellLetterLabel)
	storeSalesUnitsCellLetterLabelContainer.Layout = layout.NewMaxLayout()

	storeSalesUnitsCellLetterEntryContainer := container.NewHBox(storeSalesUnitsCellLetterEntry)
	storeSalesUnitsCellLetterEntryContainer.Layout = layout.NewMaxLayout()

	storeSalesUnitsContainer := container.New(layout.NewGridLayoutWithColumns(2), storeSalesUnitsCellLetterLabelContainer, storeSalesUnitsCellLetterEntryContainer)

	storeSalesUnitsDeviationCellLetterLabel := widget.NewLabel("Store Sales Units Deviation CELL LETTER:")
	storeSalesUnitsDeviationCellLetterEntry := widget.NewEntry()
	storeSalesUnitsDeviationCellLetterEntry.SetText(config.StoreSalesUnitesDeviationCellLetter)

	if config.StoreSalesUnitesDeviationCellLetter == "" {
		storeSalesUnitsDeviationCellLetterEntry.SetText("N")
	} else {
		storeSalesUnitsDeviationCellLetterEntry.SetText(config.StoreSalesUnitesDeviationCellLetter)
	}

	storeSalesUnitsDeviationCellLetterEntry.TextStyle = fyne.TextStyle{Bold: true}

	storeSalesUnitsDeviationCellLetterLabelContainer := container.NewHBox(storeSalesUnitsDeviationCellLetterLabel)
	storeSalesUnitsDeviationCellLetterLabelContainer.Layout = layout.NewMaxLayout()

	storeSalesUnitsDeviationCellLetterEntryContainer := container.NewHBox(storeSalesUnitsDeviationCellLetterEntry)
	storeSalesUnitsDeviationCellLetterEntryContainer.Layout = layout.NewMaxLayout()

	storeSalesUnitsDeviationContainer := container.New(layout.NewGridLayoutWithColumns(2), storeSalesUnitsDeviationCellLetterLabelContainer, storeSalesUnitsDeviationCellLetterEntryContainer)

	checkFormStatus := widget.NewLabel("")

	// Begin Running the Program
	saveConfig := widget.NewButtonWithIcon("Save / Refresh Configuration", theme.DocumentSaveIcon(), func() {
		checkFormStatus.SetText("")

		var epocell_file string = filePathEntry.Text
		var epocell_data_starts_at string = epocellStartsAtRowEntry.Text
		var store_name_cell_letter string = storeNameCellLetterEntry.Text
		var store_sales_units_cell_letter string = storeSalesUnitsCellLetterEntry.Text
		var store_sales_units_deviation_cell_letter string = storeSalesUnitsDeviationCellLetterEntry.Text

		if epocell_file != "" &&
			epocell_data_starts_at != "" &&
			store_name_cell_letter != "" &&
			store_sales_units_cell_letter != "" &&
			store_sales_units_deviation_cell_letter != "" {

			checkFormStatus.SetText("")

			jsonStructure := make(map[string]interface{})

			jsonStructure["epocell_file"] = epocell_file
			jsonStructure["epocell_starts_at_row"], _ = strconv.Atoi(epocell_data_starts_at)
			jsonStructure["store_name_cell_letter"] = store_name_cell_letter
			jsonStructure["store_sales_units_cell_letter"] = store_sales_units_cell_letter
			jsonStructure["store_sales_units_deviation_cell_letter"] = store_sales_units_deviation_cell_letter

			// Convert the map to JSON with indentation.
			jsonData, err := json.MarshalIndent(jsonStructure, "", "\t")
			if err != nil {
				checkFormStatus.SetText("Error encoding data to JSON: " + err.Error())
				return
			}

			// Save the JSON data to a file.
			err = os.WriteFile("config/config.json", jsonData, 0644)
			if err != nil {
				checkFormStatus.SetText("Error writting config data: " + err.Error())
				return
			}

			checkFormStatus.SetText("Configuration saved in config/config.json")

			duration := 1 * time.Second
			time.Sleep(duration)
			checkFormStatus.SetText("")

			/**
			Load data from the given configuration
			*/

			epocellFilePath := filepath.Base(epocell_file)

			loadEpocellWindow := epocellApp.NewWindow(epocellFilePath)
			loadEpocellWindowIcon, _ := fyne.LoadResourceFromPath("icons/epocell-icon.png")
			loadEpocellWindow.SetIcon(loadEpocellWindowIcon)

			sheetIndex := 0                                     // Index of the sheet you want to read (0 for the first sheet).
			numChunks := 4                                      // Number of goroutines to use for reading data.
			startRow, _ := strconv.Atoi(epocell_data_starts_at) // Starting row to read from.
			columnLetters := []string{store_name_cell_letter, store_sales_units_cell_letter, store_sales_units_deviation_cell_letter}

			// Convert column letters to their corresponding indices (0-based).
			columnIndices := make([]int, len(columnLetters))
			for i, letter := range columnLetters {
				columnIndices[i] = core.ColLetterToIndex(letter)
			}

			contents, _ := core.ParseEpocell(epocell_file, sheetIndex, numChunks, startRow, columnIndices, columnLetters)

			store_cache_data := []map[string]interface{}{}
			if contents != nil {

				for chunkData := range contents {
					for _, columnData := range chunkData {
						/**
						Implement checking if theere any data
						based on the columns provided
						*/
						storeCache := map[string]interface{}{
							"store_name":            columnData[store_name_cell_letter],
							"sales_units":           columnData[store_sales_units_cell_letter],
							"sales_units_deviation": columnData[store_sales_units_deviation_cell_letter],
						}

						store_cache_data = append(store_cache_data, storeCache)
					}
				}
			}

			if store_cache_data != nil {
				// Provide your custom filename (without extension) here
				filename := "store_cache_" + core.RemoveExtension(filepath.Base(epocell_file))
				var storeCaches []core.StoreCache

				for i, columnData := range store_cache_data {

					storeName, storeNameExists := columnData["store_name"].(string)
					salesUnits, salesUnitsExists := columnData["sales_units"].(string)
					salesUnitsDeviation, salesUnitsDeviationExists := columnData["sales_units_deviation"].(string)

					if !storeNameExists || !salesUnitsExists || !salesUnitsDeviationExists {
						fmt.Println("Skipping incomplete data for store:", columnData)
						continue
					}

					if salesUnits == "" {
						salesUnits = "0"
					}

					if salesUnitsDeviation == "" {
						salesUnitsDeviation = "0"
					}

					if storeName != "" && salesUnits != "0" && salesUnitsDeviation != "0" {
						storeCache := core.StoreCache{
							Index:                i,
							Status:               false,
							StoreName:            storeName,
							SalesUnits:           salesUnits,
							SalesUnitsDeviation:  salesUnitsDeviation,
							CsvFile:              columnData["store_name"].(string) + ".csv", // Set Csv file
							Keywords:             []string{},                                 // Set keywords
							SkipKeywords:         []string{},
							SheetId:              1,   // Set the Sheet ID
							DataStartsAtRow:      1,   // Set data starts at row
							CategoryCellLetter:   "A", // Set category cell letter
							ProductCellLetter:    "B",
							SumOfSalesCellLetter: "C", // Set sum of sales cell letter
						}
						storeCaches = append(storeCaches, storeCache)
					}
				}

				//// Remove spaces and special characters from the filename
				//formattedFilename := formatStoreName(filename)

				err := core.SaveStoreCache(storeCaches, filename)
				if err != nil {
					fmt.Println("Error saving store cache:", err)
					checkFormStatus.SetText("Error saving store cache: " + err.Error())

					duration := 1 * time.Second
					time.Sleep(duration)
					checkFormStatus.SetText("")

					epocellApp.Quit()
					core.RestartApp()

					return
				}

				checkFormStatus.SetText("Store Cache Data saved in cache folder.")
				fmt.Println("Store cache data saved to JSON files.")

				duration := 1 * time.Second
				time.Sleep(duration)
				checkFormStatus.SetText("")

				epocellApp.Quit()
				core.RestartApp()
			}

		} else {
			checkFormStatus.SetText("All the fields must be filled.")
		}

	})

	saveConfigContainer := container.NewHBox(saveConfig)
	saveConfigContainer.Layout = layout.NewMaxLayout()

	/**
	Handle multi-store configuration
	*/
	storeConfigurationButton := widget.NewButtonWithIcon("Store Configuration", theme.InfoIcon(), func() {

		storeConfigWindow := epocellApp.NewWindow("Store Configuration")
		storeConfigWindowIcon, _ := fyne.LoadResourceFromPath("icons/epocell-icon.png")
		storeConfigWindow.SetIcon(storeConfigWindowIcon)

		statusLabel := widget.NewLabel("")

		/**
		Read config data from config/config.json
		*/

		config, _ := core.LoadJsonConfig()

		if config.EpocellFile != "" {

			epocell_file_path := filepath.Base(config.EpocellFile)
			cache_file_name := "store_cache_" + core.RemoveExtension(epocell_file_path) + ".json"

			storeCacheArray, err := core.ReadStoreCacheFromFile(cache_file_name)
			if err != nil {
				checkFormStatus.SetText("Error reading store cache: " + err.Error())
				duration := 1 * time.Second
				time.Sleep(duration)
				checkFormStatus.SetText("")
				return
			}

			searchEntry := widget.NewEntry()
			searchEntry.PlaceHolder = "Search store..."

			// Initialize filteredStores with all stores
			filteredStores := storeCacheArray

			storeList := widget.NewList(
				func() int {
					return len(filteredStores)
				},
				func() fyne.CanvasObject {
					buttonLabel := widget.NewLabel("")
					buttonConfigure := widget.NewButtonWithIcon("Configure", theme.FolderOpenIcon(), func() {

					})

					buttonMark := widget.NewButtonWithIcon("Mark", theme.InfoIcon(), func() {

					})
					contentContainer := container.New(layout.NewHBoxLayout(), buttonLabel, buttonConfigure, buttonMark)
					return contentContainer
				},
				func(index int, item fyne.CanvasObject) {
					existingData := storeCacheArray
					hBox := item.(*fyne.Container)
					storeNameLabel := hBox.Objects[0].(*widget.Label)
					buttonConfigure := hBox.Objects[1].(*widget.Button)
					buttonCheck := hBox.Objects[2].(*widget.Button)

					fullStoreLabel := filteredStores[index].StoreName + " - Sales Units: " + filteredStores[index].SalesUnits + " - Sales Units Deviation: " + filteredStores[index].SalesUnitsDeviation
					storeNameLabel.SetText(fullStoreLabel)
					storeNameLabel.Refresh() // Invalidate the label to ensure it updates visually

					/**
					Set the checked button state
					*/
					var isChecked bool
					originalIndex := filteredStores[index].Index
					selectedStoreIndex := originalIndex

					for idx, store := range existingData {
						if store.Index == selectedStoreIndex {
							// Modify the selected store's data
							if existingData[idx].Status {
								isChecked = true
							} else {
								isChecked = false
							}
							break
						}
					}

					if isChecked {

						buttonCheck.SetIcon(theme.CheckButtonCheckedIcon())
						buttonCheck.SetText("Done")
					} else {
						buttonCheck.SetIcon(theme.CheckButtonIcon())
						buttonCheck.SetText("Undone")
					}

					buttonConfigure.OnTapped = func(index int) func() {
						return func() {
							originalIndex := filteredStores[index].Index

							configOneStoreWindow := epocellApp.NewWindow("Configure: " + fullStoreLabel)
							configOneStoreWindowIcon, _ := fyne.LoadResourceFromPath("icons/epocell-icon.png")
							configOneStoreWindow.SetIcon(configOneStoreWindowIcon)

							mainText := widget.NewLabel("Store Information Configuration")

							outputText := widget.NewLabel("")

							csvFilePathEntry := widget.NewEntry()
							csvFilePathEntry.SetText(filteredStores[index].CsvFile)
							csvFilePathEntry.TextStyle = fyne.TextStyle{Bold: true}

							// Create a button to open the file picker dialog
							csvPickerButton := widget.NewButton("Pick File", func() {
								core.PickFileCsv(configOneStoreWindow, csvFilePathEntry)
							})

							csvFileLabelContainer := container.NewHBox(widget.NewLabel("CSV File:"))
							csvFileEntryContainer := container.NewHBox(csvFilePathEntry)
							csvFileEntryContainer.Layout = layout.NewMaxLayout()
							pickcsvFileContainer := container.NewHBox(csvPickerButton)
							pickcsvFileContainer.Layout = layout.NewMaxLayout()
							csvfilePickerContainer := container.New(layout.NewGridLayoutWithColumns(3), csvFileLabelContainer, csvFileEntryContainer, pickcsvFileContainer)

							sheetIdLabelContainer := container.NewHBox(widget.NewLabel("Sheet ID:"))
							sheetIdLabelContainer.Layout = layout.NewMaxLayout()

							sheetIdEntry := widget.NewEntry()
							sheetIdEntry.SetText(strconv.Itoa(filteredStores[index].SheetId))
							sheetIdEntry.TextStyle = fyne.TextStyle{Bold: true}

							sheetIdEntry.OnChanged = func(text string) {
								if _, err := strconv.Atoi(text); err != nil && text != "" {
									// If the input is not a valid integer, revert to the previous value
									sheetIdEntry.SetText(strconv.Itoa(filteredStores[index].SheetId))
								}
							}

							sheetIdEntryContainer := container.NewHBox(sheetIdEntry)
							sheetIdEntryContainer.Layout = layout.NewMaxLayout()

							sheetIdContainer := container.New(layout.NewGridLayoutWithColumns(2), sheetIdLabelContainer, sheetIdEntryContainer)

							keywordsLabelContainer := container.NewHBox(widget.NewLabel("Keywords:"))
							keywordsLabelContainer.Layout = layout.NewMaxLayout()

							keywordStr := strings.Join(filteredStores[index].Keywords, ",")

							keywordsEntry := widget.NewEntry()
							keywordsEntry.SetText(keywordStr)
							keywordsEntry.TextStyle = fyne.TextStyle{Bold: true}

							keywordsEntryContainer := container.NewHBox(keywordsEntry)
							keywordsEntryContainer.Layout = layout.NewMaxLayout()

							keywordsContainer := container.New(layout.NewGridLayoutWithColumns(2), keywordsLabelContainer, keywordsEntryContainer)

							skipKeywordsLabelContainer := container.NewHBox(widget.NewLabel("Skip if contains keywords:"))
							skipKeywordsLabelContainer.Layout = layout.NewMaxLayout()

							skipKeywordstr := strings.Join(filteredStores[index].SkipKeywords, ",")

							skipKeywordsEntry := widget.NewEntry()
							skipKeywordsEntry.SetText(skipKeywordstr)
							skipKeywordsEntry.TextStyle = fyne.TextStyle{Bold: true}

							skipKeywordsEntryContainer := container.NewHBox(skipKeywordsEntry)
							skipKeywordsEntryContainer.Layout = layout.NewMaxLayout()

							skipKeywordsContainer := container.New(layout.NewGridLayoutWithColumns(2), skipKeywordsLabelContainer, skipKeywordsEntryContainer)

							dataStartsAtRowLabelContainer := container.NewHBox(widget.NewLabel("Data begins at row:"))
							dataStartsAtRowLabelContainer.Layout = layout.NewMaxLayout()

							dataStartsAtRowEntry := widget.NewEntry()
							dataStartsAtRowEntry.SetText(strconv.Itoa(filteredStores[index].DataStartsAtRow))
							dataStartsAtRowEntry.TextStyle = fyne.TextStyle{Bold: true}

							/**
							Force integers up there
							*/

							dataStartsAtRowEntry.OnChanged = func(text string) {
								if _, err := strconv.Atoi(text); err != nil && text != "" {
									// If the input is not a valid integer, revert to the previous value
									dataStartsAtRowEntry.SetText(strconv.Itoa(filteredStores[index].DataStartsAtRow))
								}
							}

							dataStartsAtRowEntryContainer := container.NewHBox(dataStartsAtRowEntry)
							dataStartsAtRowEntryContainer.Layout = layout.NewMaxLayout()

							dataStartsAtRowContainer := container.New(layout.NewGridLayoutWithColumns(2), dataStartsAtRowLabelContainer, dataStartsAtRowEntryContainer)

							categoryCellLetterLabelContainer := container.NewHBox(widget.NewLabel("Category / Product Cell Letter:"))
							categoryCellLetterLabelContainer.Layout = layout.NewMaxLayout()

							categoryCellLetterEntry := widget.NewEntry()
							categoryCellLetterEntry.SetText(filteredStores[index].CategoryCellLetter)
							categoryCellLetterEntry.TextStyle = fyne.TextStyle{Bold: true}

							categoryCellLetterEntryContainer := container.NewHBox(categoryCellLetterEntry)
							categoryCellLetterEntryContainer.Layout = layout.NewMaxLayout()

							categoryCellLetterContainer := container.New(layout.NewGridLayoutWithColumns(2), categoryCellLetterLabelContainer, categoryCellLetterEntryContainer)

							productCellLetterLabelContainer := container.NewHBox(widget.NewLabel("Product Name Cell Letter:"))
							productCellLetterLabelContainer.Layout = layout.NewMaxLayout()

							productCellLetterEntry := widget.NewEntry()
							productCellLetterEntry.SetText(filteredStores[index].ProductCellLetter)
							productCellLetterEntry.TextStyle = fyne.TextStyle{Bold: true}

							productCellLetterEntryContainer := container.NewHBox(productCellLetterEntry)
							productCellLetterEntryContainer.Layout = layout.NewMaxLayout()

							productCellLetterContainer := container.New(layout.NewGridLayoutWithColumns(2), productCellLetterLabelContainer, productCellLetterEntryContainer)

							sumOfSalesCellLetterLabelContainer := container.NewHBox(widget.NewLabel("Sum of Sales Cell Letter:"))
							sumOfSalesCellLetterLabelContainer.Layout = layout.NewMaxLayout()

							sumOfSalesCellLetterEntry := widget.NewEntry()
							sumOfSalesCellLetterEntry.SetText(filteredStores[index].SumOfSalesCellLetter)
							sumOfSalesCellLetterEntry.TextStyle = fyne.TextStyle{Bold: true}

							sumOfSalesCellLetterEntryContainer := container.NewHBox(sumOfSalesCellLetterEntry)
							sumOfSalesCellLetterEntryContainer.Layout = layout.NewMaxLayout()

							sumOfSalesCellLetterContainer := container.New(layout.NewGridLayoutWithColumns(2), sumOfSalesCellLetterLabelContainer, sumOfSalesCellLetterEntryContainer)

							keywordsContent := container.NewVBox()

							keyordCacheFileName := "store_cache_keywords_" + filteredStores[index].StoreName + "_" + core.RemoveExtension(epocell_file_path)

							fetchKeywords, _ := core.ReadStoreKeywordsCache(keyordCacheFileName)
							if fetchKeywords != nil {
								searchEntryKeywords := widget.NewEntry()
								searchEntryKeywords.PlaceHolder = "Search keywords..."

								filteredKeywords := fetchKeywords

								keywordList := widget.NewList(
									func() int {
										return len(filteredKeywords)
									},
									func() fyne.CanvasObject {
										buttonLabel := widget.NewLabel("")
										buttonCopy := widget.NewButtonWithIcon("Copy", theme.ContentCopyIcon(), func() {

										})
										contentContainer := container.New(layout.NewHBoxLayout(), buttonLabel, buttonCopy)
										return contentContainer
									},
									func(index int, item fyne.CanvasObject) {
										hBox := item.(*fyne.Container)
										keywordLabel := hBox.Objects[0].(*widget.Label)
										buttonCopy := hBox.Objects[1].(*widget.Button)

										keyword := filteredKeywords[index].Keyword
										productCount := filteredKeywords[index].Count

										keywordLabel.SetText(keyword + " : " + strconv.Itoa(productCount) + "")

										buttonCopy.OnTapped = func() {
											// Append keyword to input and separate with commas
											input := keywordsEntry // Replace with your input field reference
											currentText := keywordsEntry.Text
											if currentText != "" {
												currentText += ", "
											}
											currentText += keyword
											input.SetText(currentText)

											// Provide user feedback if needed
											buttonCopy.SetText("Copied")
											buttonCopy.Disable()
										}
									},
								)

								searchEntryKeywords.OnChanged = func(query string) {
									if query == "" {
										// If the query is empty, show all stores
										filteredKeywords = fetchKeywords
									} else {
										// Otherwise, filter the stores based on the query
										filteredKeywords = core.FilterKeywords(query, fetchKeywords)
									}

									// Define a function that returns the length of the filtered stores
									getFilteredKeywordsCount := func() int {
										return len(filteredKeywords)
									}

									// Assign the function to storeList.Length
									keywordList.Length = getFilteredKeywordsCount

									// Refresh the list
									keywordList.Refresh()
								}

								scrollContainerKeywords := container.NewVScroll(keywordList)
								scrollContainerKeywords.SetMinSize(fyne.NewSize(1000, 250)) // Adjust the height as needed

								searchKeywordsContainer := container.NewHBox(searchEntryKeywords)
								searchKeywordsContainer.Layout = layout.NewMaxLayout()

								layoutContainer := container.NewVBox(widget.NewLabelWithStyle("Suggested Keywords:", fyne.TextAlign(0), fyne.TextStyle{Bold: true}), searchKeywordsContainer, scrollContainerKeywords)

								//keywordScrollListContainer := container.NewVBox(searchKeywordsContainer, scrollContainer)
								layoutContainerKeywords := container.NewVBox(layoutContainer)

								keywordsContent.Refresh()
								keywordsContent.Add(layoutContainerKeywords)
							}

							suggestKeywordsButton := widget.NewButtonWithIcon("Load Keywords", theme.DownloadIcon(), func() {
								keywordsContent.Objects = nil
								outputText.SetText("Processing keywords...")
								duration := 2 * time.Second
								time.Sleep(duration)
								outputText.SetText("")

								filePath := csvFilePathEntry.Text
								numChunks := 10 // 10 goroutines
								startRow, _ := strconv.Atoi(dataStartsAtRowEntry.Text)
								columnLetters := []string{categoryCellLetterEntry.Text, productCellLetterEntry.Text, sumOfSalesCellLetterEntry.Text} // Replace with corresponding column letters
								// Convert column letters to their corresponding indices (0-based).
								columnIndices := make([]int, len(columnLetters))
								for i, letter := range columnLetters {
									columnIndices[i] = core.ColLetterToIndex(letter)
								}
								keywordColumnLetter := categoryCellLetterEntry.Text
								var emptyKeywords []string
								var emptySkipKeywords []string

								dataChan, err := core.ParseCsv(filePath, numChunks, startRow, keywordColumnLetter, columnIndices, columnLetters, emptyKeywords, emptySkipKeywords)
								if err != nil {
									outputText.SetText("Error parsing CSV File. " + err.Error())
									duration := 2 * time.Second
									time.Sleep(duration)
									outputText.SetText("")
									return
								}

								var foundKeywords []string
								var keywordCounts map[string]int // Map to store keyword counts

								// Initialize keywordCounts map
								keywordCounts = make(map[string]int)

								if dataChan != nil {
									for chunk := range dataChan {
										for _, columnData := range chunk {
											foundKeywords = append(foundKeywords, columnData[categoryCellLetterEntry.Text])
											keyword := columnData[categoryCellLetterEntry.Text]
											keywordCounts[keyword]++ // Increment the count for the keyword
										}
									}
								}

								// Extract unique keywords from foundKeywords
								uniqueKeywords := core.UniqueStrings(foundKeywords)

								// Populate storeKeywordCache with unique keyword counts
								var storeKeywordCache []core.StoreKeywordCache
								for _, keyword := range uniqueKeywords {
									count := keywordCounts[keyword]
									storeKeywordCache = append(storeKeywordCache, core.StoreKeywordCache{
										Keyword: keyword,
										Count:   count,
									})
								}

								if len(uniqueKeywords) > 0 {

									keywordCacheFileName := "store_cache_keywords_" + filteredStores[index].StoreName + "_" + core.RemoveExtension(epocell_file_path)

									/**
									If cache doesnt exist, create it
									*/

									errKeywordsCache := core.UpdateStoreKeywordsCacheIfNecessary(keywordCacheFileName, storeKeywordCache)
									if errKeywordsCache != nil {
										// Handle the error here, if needed
										outputText.SetText("Unable to update/create store keywords cache: " + errKeywordsCache.Error())
										duration := 10 * time.Second
										time.Sleep(duration)
										outputText.SetText("")
									}

									outputText.SetText("Loading keywords...")

									duration := 2 * time.Second
									time.Sleep(duration)
									outputText.SetText("")

									searchEntryKeywords := widget.NewEntry()
									searchEntryKeywords.PlaceHolder = "Search keywords..."

									filteredKeywords := storeKeywordCache

									keywordList := widget.NewList(
										func() int {
											return len(filteredKeywords)
										},
										func() fyne.CanvasObject {
											buttonLabel := widget.NewLabel("")
											buttonCopy := widget.NewButtonWithIcon("Copy", theme.ContentCopyIcon(), func() {

											})
											contentContainer := container.New(layout.NewHBoxLayout(), buttonLabel, buttonCopy)
											return contentContainer
										},
										func(index int, item fyne.CanvasObject) {
											hBox := item.(*fyne.Container)
											keywordLabel := hBox.Objects[0].(*widget.Label)
											buttonCopy := hBox.Objects[1].(*widget.Button)

											keyword := filteredKeywords[index].Keyword
											productCount := filteredKeywords[index].Count

											keywordLabel.SetText(keyword + " : " + strconv.Itoa(productCount) + "")

											buttonCopy.OnTapped = func() {
												// Append keyword to input and separate with commas
												input := keywordsEntry // Replace with your input field reference
												currentText := keywordsEntry.Text
												if currentText != "" {
													currentText += ", "
												}
												currentText += keyword
												input.SetText(currentText)

												// Provide user feedback if needed
												buttonCopy.SetText("Copied")
												buttonCopy.Disable()
											}
										},
									)

									searchEntryKeywords.OnChanged = func(query string) {
										if query == "" {
											// If the query is empty, show all stores
											filteredKeywords = storeKeywordCache
										} else {
											// Otherwise, filter the stores based on the query
											filteredKeywords = core.FilterKeywords(query, storeKeywordCache)
										}

										// Define a function that returns the length of the filtered stores
										getFilteredKeywordsCount := func() int {
											return len(filteredKeywords)
										}

										// Assign the function to storeList.Length
										keywordList.Length = getFilteredKeywordsCount

										// Refresh the list
										keywordList.Refresh()
									}

									scrollContainerKeywords := container.NewVScroll(keywordList)
									scrollContainerKeywords.SetMinSize(fyne.NewSize(1000, 250)) // Adjust the height as needed

									searchKeywordsContainer := container.NewHBox(searchEntryKeywords)
									searchKeywordsContainer.Layout = layout.NewMaxLayout()

									layoutContainer := container.NewVBox(widget.NewLabelWithStyle("Suggested Keywords:", fyne.TextAlign(0), fyne.TextStyle{Bold: true}), searchKeywordsContainer, scrollContainerKeywords)

									//keywordScrollListContainer := container.NewVBox(searchKeywordsContainer, scrollContainer)
									layoutContainerKeywords := container.NewVBox(layoutContainer)

									keywordsContent.Refresh()
									keywordsContent.Add(layoutContainerKeywords)

								}
							})

							loadColumnsButton := widget.NewButtonWithIcon("Show columns", theme.InfoIcon(), func() {
								columns, err := core.MapColumnsWithLetters(csvFilePathEntry.Text)
								if err != nil {
									outputText.SetText("Unable to parse Csv File. " + err.Error())
									duration := 2 * time.Second
									time.Sleep(duration)
									outputText.SetText("")
									return
								}

								if len(columns) > 0 {

									columnsList := widget.NewList(
										func() int {
											return len(columns)
										},
										func() fyne.CanvasObject {
											buttonLabel := widget.NewLabel("")
											contentContainer := container.New(layout.NewHBoxLayout(), buttonLabel)
											return contentContainer
										},
										func(index int, item fyne.CanvasObject) {
											hBox := item.(*fyne.Container)
											keywordLabel := hBox.Objects[0].(*widget.Label)
											keywordLabel.SetText(columns[index].Name + " : " + columns[index].Letter)
										},
									)

									scrollContainercolumns := container.NewVScroll(columnsList)
									scrollContainercolumns.SetMinSize(fyne.NewSize(300, 300)) // Adjust the height as needed

									layoutContainerList := container.NewVBox(widget.NewLabelWithStyle("Found columns:", fyne.TextAlign(0), fyne.TextStyle{Bold: true}), scrollContainercolumns)

									//columnscrollListContainer := container.NewVBox(searchcolumnsContainer, scrollContainer)
									layoutContainerColumns := container.NewVBox(layoutContainerList)

									viewColumnsWindow := epocellApp.NewWindow("Column information for : " + csvFilePathEntry.Text)
									viewColumnsWindowIcon, _ := fyne.LoadResourceFromPath("icons/epocell-icon.png")
									viewColumnsWindow.SetIcon(viewColumnsWindowIcon)

									viewColumnsWindowContainer := container.NewVBox(
										layoutContainerColumns,
									)

									viewColumnsWindow.SetContent(viewColumnsWindowContainer)
									viewColumnsWindow.Resize(fyne.NewSize(300, 300))
									viewColumnsWindow.CenterOnScreen()
									viewColumnsWindow.Show()

								} else {
									outputText.SetText("Error parsing CSV File. " + err.Error())
									duration := 2 * time.Second
									time.Sleep(duration)
									outputText.SetText("")
									return
								}

							})

							suggestKeywordsButtonContainer := container.NewHBox(suggestKeywordsButton)
							suggestKeywordsButtonContainer.Layout = layout.NewMaxLayout()

							saveConfigButton := widget.NewButtonWithIcon("Save", theme.DocumentSaveIcon(), func() {
								csvFile := csvFilePathEntry.Text
								sheetId := sheetIdEntry.Text
								keywords := core.FilterNonEmptyKeywords(strings.Split(keywordsEntry.Text, ","))

								if len(keywords) > 0 {
									for i := range keywords {
										keywords[i] = strings.TrimSpace(keywords[i])
									}
								} else {
									// No keywords found, set the array as empty
									keywords = nil
								}

								skipKeywords := core.FilterNonEmptyKeywords(strings.Split(skipKeywordsEntry.Text, ","))

								if len(skipKeywords) > 0 {
									for i := range skipKeywords {
										skipKeywords[i] = strings.TrimSpace(skipKeywords[i])
									}
								} else {
									skipKeywords = nil
								}

								dataStartsAtRow := dataStartsAtRowEntry.Text
								categoryCellLetter := categoryCellLetterEntry.Text
								productCellLetter := productCellLetterEntry.Text
								sumOfSalesLetter := sumOfSalesCellLetterEntry.Text

								if csvFile != "" &&
									sheetId != "" &&
									dataStartsAtRow != "" &&
									categoryCellLetter != "" &&
									productCellLetter != "" &&
									sumOfSalesLetter != "" {

									// Read existing data from the JSON file
									existingData, err := core.ReadStoreCacheFromFile(cache_file_name)
									if err != nil {
										outputText.SetText("Error reading Store Cache : " + err.Error())
										duration := 1 * time.Second
										time.Sleep(duration)
										outputText.SetText("")
									}

									// Find the specific store data in the existing data

									selectedStoreIndex := originalIndex

									if len(existingData) < 1 {
										outputText.SetText("Error reading Store Cache : " + err.Error())
										duration := 1 * time.Second
										time.Sleep(duration)
										outputText.SetText("")
									}

									var selectedStore core.StoreCache

									for _, store := range existingData {
										if store.Index == selectedStoreIndex {
											selectedStore = store
											break // Exit the loop once the desired index is found
										}
									}

									for idx, store := range existingData {
										if store.Index == selectedStoreIndex {
											// Modify the selected store's data
											existingData[idx].CsvFile = csvFile
											existingData[idx].SheetId, _ = strconv.Atoi(sheetId)
											existingData[idx].Keywords = keywords
											existingData[idx].SkipKeywords = skipKeywords
											existingData[idx].DataStartsAtRow, _ = strconv.Atoi(dataStartsAtRow)
											existingData[idx].CategoryCellLetter = categoryCellLetter
											existingData[idx].ProductCellLetter = productCellLetter
											existingData[idx].SumOfSalesCellLetter = sumOfSalesLetter
											break
										}
									}

									// Update the selected store data
									selectedStore.CsvFile = csvFile
									selectedStore.SheetId, _ = strconv.Atoi(sheetId)
									selectedStore.Keywords = keywords
									selectedStore.SkipKeywords = skipKeywords
									selectedStore.DataStartsAtRow, _ = strconv.Atoi(dataStartsAtRow)
									selectedStore.CategoryCellLetter = categoryCellLetter
									selectedStore.ProductCellLetter = productCellLetter
									selectedStore.SumOfSalesCellLetter = sumOfSalesLetter

									// Write the updated data back to the JSON file
									err = core.UpdateStoreCacheToFile(cache_file_name, existingData)
									if err != nil {
										outputText.SetText("Error updating Store Cache : " + err.Error())
										duration := 1 * time.Second
										time.Sleep(duration)
										outputText.SetText("")
									} else {

										/**
										Refresh data in the form
										*/

										filteredStores[index].CsvFile = csvFile
										filteredStores[index].SheetId, _ = strconv.Atoi(sheetId)
										filteredStores[index].Keywords = keywords
										filteredStores[index].SkipKeywords = skipKeywords
										filteredStores[index].DataStartsAtRow, _ = strconv.Atoi(dataStartsAtRow)
										filteredStores[index].CategoryCellLetter = categoryCellLetter
										filteredStores[index].ProductCellLetter = productCellLetter
										filteredStores[index].SumOfSalesCellLetter = sumOfSalesLetter

										outputText.SetText("Configuration saved!")

										duration := 1 * time.Second
										time.Sleep(duration)
										outputText.SetText("")
									}

								} else {
									outputText.SetText("Error: All fields are required!")
									duration := 1 * time.Second
									time.Sleep(duration)
									outputText.SetText("")
								}

							})

							saveConfigButtonContainer := container.NewHBox(saveConfigButton)
							saveConfigButtonContainer.Layout = layout.NewMaxLayout()

							lastRowButtons := container.New(layout.NewGridLayoutWithColumns(3), suggestKeywordsButtonContainer, loadColumnsButton, saveConfigButtonContainer)

							configOneStoreWindowContainer := container.NewVBox(
								mainText,
								csvfilePickerContainer,
								sheetIdContainer,
								keywordsContent,
								keywordsContainer,
								skipKeywordsContainer,
								dataStartsAtRowContainer,
								categoryCellLetterContainer,
								productCellLetterContainer,
								sumOfSalesCellLetterContainer,
								lastRowButtons,
								outputText,
							)

							configOneStoreWindow.SetContent(configOneStoreWindowContainer)
							configOneStoreWindow.Resize(fyne.NewSize(1000, 250))
							configOneStoreWindow.CenterOnScreen()
							configOneStoreWindow.Show()
						}
					}(index) // Capture the index in the event handler closure

					buttonCheck.OnTapped = func(index int) func() {
						return func() {

							if isChecked {

								for idx, store := range existingData {
									if store.Index == selectedStoreIndex {
										// Modify the selected store's data
										existingData[idx].Status = false
										break
									}
								}

								buttonCheck.SetIcon(theme.CheckButtonCheckedIcon())
								buttonCheck.SetText("Checked")

							} else {

								for idx, store := range existingData {
									if store.Index == selectedStoreIndex {
										// Modify the selected store's data
										existingData[idx].Status = true
										break
									}
								}

								buttonCheck.SetIcon(theme.CheckButtonIcon())
								buttonCheck.SetText("Unchecked")
							}

							// Write the updated data back to the JSON file
							err = core.UpdateStoreCacheToFile(cache_file_name, existingData)
							if err != nil {
								fmt.Println("Error updating Store Cache : " + err.Error())
							}
						}
					}(index) // Capture the index in the event handler closure
				},
			)

			storeList.OnSelected = func(id int) {
				// Handler for item selection
				// Access the corresponding storeCache using the id
				// Show the store data using a dialog or another UI element
			}

			searchEntry.OnChanged = func(query string) {
				if query == "" {
					// If the query is empty, show all stores
					filteredStores = filteredStores
				} else {
					// Otherwise, filter the stores based on the query
					filteredStores = core.FilterStoreCache(query, storeCacheArray)
				}

				// Define a function that returns the length of the filtered stores
				getFilteredStoreCount := func() int {
					return len(filteredStores)
				}

				// Assign the function to storeList.Length
				storeList.Length = getFilteredStoreCount

				// Refresh the list
				storeList.Refresh()
			}

			scrollContainer := container.NewVScroll(storeList)
			scrollContainer.SetMinSize(fyne.NewSize(650, 600)) // Adjust the height as needed

			searchContainer := container.NewHBox(searchEntry)
			searchContainer.Layout = layout.NewMaxLayout()

			storeListLayoutContainer := container.NewVBox(searchContainer, scrollContainer)
			layoutContainer := container.NewVBox(widget.NewLabel("Store List"), storeListLayoutContainer)

			storeConfigWindow.SetContent(layoutContainer)
			storeConfigWindow.Resize(fyne.NewSize(800, 600))
			storeConfigWindow.CenterOnScreen()
			storeConfigWindow.Show()

		} else {
			statusLabel.SetText("No epocell filename provided.")
		}

	})

	storeConfigurationContainer := container.NewHBox(storeConfigurationButton)
	storeConfigurationContainer.Layout = layout.NewMaxLayout()

	computeButton := widget.NewButtonWithIcon("Compute", theme.ComputerIcon(), func() {
		if config.EpocellFile != "" {

			epocell_file_path := filepath.Base(config.EpocellFile)
			cache_file_name := "store_cache_" + core.RemoveExtension(epocell_file_path) + ".json"

			storeCacheArray, err := core.ReadStoreCacheFromFile(cache_file_name)
			if err != nil {
				checkFormStatus.SetText("Error reading store cache: " + err.Error())
				duration := 1 * time.Second
				time.Sleep(duration)
				checkFormStatus.SetText("")
				return
			}

			checkFormStatus.SetText("Computing, please wait...")

			toCompute := core.FilterCacheData(storeCacheArray)

			if toCompute != nil {

				currentTime := time.Now()
				timestamp := currentTime.Format("2006-01-02_15-04-05")
				columns := []string{"Store Name", "Sales Units", "Sales Units Deviation", "Sum of Sales", "Used Keywords", "Skipped Keywords"}

				fileName := "data_export_" + core.RemoveExtension(epocell_file_path) + "_" + timestamp + ".csv"

				var csvData [][]string

				for _, item := range toCompute {
					// Apply your filtering logic here

					filePath := item.CsvFile
					numChunks := 16 // 16 goroutines
					startRow := item.DataStartsAtRow
					columnLetters := []string{item.CategoryCellLetter, item.ProductCellLetter, item.SumOfSalesCellLetter} // Replace with corresponding column letters
					// Convert column letters to their corresponding indices (0-based).
					columnIndices := make([]int, len(columnLetters))
					for i, letter := range columnLetters {
						columnIndices[i] = core.ColLetterToIndex(letter)
					}
					keywordColumnLetter := item.CategoryCellLetter
					keywords := item.Keywords // Replace with your list of keywords
					skipKeywords := item.SkipKeywords

					dataChan, err := core.ParseCsv(filePath, numChunks, startRow, keywordColumnLetter, columnIndices, columnLetters, keywords, skipKeywords)
					if err != nil {
						checkFormStatus.SetText("Error parsing CSV File. " + err.Error())
						duration := 1 * time.Second
						time.Sleep(duration)
						checkFormStatus.SetText("")
						return
					}

					var productCsvData [][]string

					productExportColumns := []string{"Category Name", "Product Name", "Sum of Sales"}
					csvFilePath := filepath.Base(item.CsvFile)
					productExportFileName := "product_export_" + core.RemoveExtension(csvFilePath) + "_" + core.RemoveExtension(epocell_file_path) + "_" + timestamp + ".csv"

					if dataChan != nil {

						sum := 0
						for chunk := range dataChan {

							for _, columnData := range chunk {
								productExportRow := []string{columnData[item.CategoryCellLetter], columnData[item.ProductCellLetter], columnData[item.SumOfSalesCellLetter]}
								productCsvData = append(productCsvData, productExportRow)

								numberStr := columnData[item.SumOfSalesCellLetter]
								number, err := strconv.Atoi(numberStr)
								if err != nil {
									fmt.Printf("Unable to convert to int in order to calculate sum %s\n", err)
									continue
								}
								sum += number
							}

						}
						usedKeywords := strings.Join(item.Keywords, ",")
						skippedKeywords := strings.Join(item.SkipKeywords, ",")
						row := []string{item.StoreName, item.SalesUnits, item.SalesUnitsDeviation, strconv.Itoa(sum), usedKeywords, skippedKeywords}
						csvData = append(csvData, row)

						err = core.CreateAndSaveCSV(productExportColumns, "output/products/"+productExportFileName, productCsvData)

						if err != nil {
							checkFormStatus.SetText("Error saving CSV File. " + err.Error())
							duration := 4 * time.Second
							time.Sleep(duration)
							checkFormStatus.SetText("")
							return
						}
					}
				}

				err = core.CreateAndSaveCSV(columns, "output/result/"+fileName, csvData)
				if err != nil {
					checkFormStatus.SetText("Error saving CSV File. " + err.Error())
					duration := 2 * time.Second
					time.Sleep(duration)
					checkFormStatus.SetText("")
					return
				} else {
					checkFormStatus.SetText("CSV File Generated in output folder. ")
					duration := 3 * time.Second
					time.Sleep(duration)
					checkFormStatus.SetText("")

					if len(csvData) > 0 {

						//expectedRowLength := len(columns)

						// Iterate through csvData and ensure each row matches the expected length

						grid := container.NewGridWithColumns(len(columns))

						for _, colName := range columns {
							headerLabel := widget.NewLabel(colName)
							grid.Add(headerLabel)
						}

						for rowIndex, rowData := range csvData {
							for _, cellData := range rowData {
								input := widget.NewEntry()
								input.SetText(cellData)
								input.Resize(input.MinSize()) // Set the input size to its minimum size

								rowIndexCopy := rowIndex
								input.OnChanged = func(text string) {
									csvData[rowIndexCopy][len(columns)] = text
								}

								grid.Add(input)
							}
						}

						computeTable := container.NewVScroll(grid)
						computeTable.SetMinSize(fyne.NewSize(620, 300)) // Adjust the height as needed

						layoutContainerCompute := container.NewVBox(computeTable)

						computeResultsWindow := epocellApp.NewWindow("Compute Results")
						computeResultsWindowIcon, _ := fyne.LoadResourceFromPath("icons/epocell-icon.png")
						computeResultsWindow.SetIcon(computeResultsWindowIcon)

						computeResultsWindow.SetContent(container.NewVBox(
							layoutContainerCompute,
							widget.NewButtonWithIcon("Done", theme.ConfirmIcon(), func() {
								computeResultsWindow.Close()
							}),
						))
						computeResultsWindow.Resize(computeResultsWindow.Canvas().Size())
						computeResultsWindow.CenterOnScreen()

						computeResultsWindow.Show()
					}

				}

			} else {
				checkFormStatus.SetText("Error, no data found to compute. Make sure the CSV Files exist. ")
				duration := 2 * time.Second
				time.Sleep(duration)
				checkFormStatus.SetText("")
				return
			}

		}
	})

	computeButtonContainer := container.NewHBox(computeButton)
	computeButtonContainer.Layout = layout.NewMaxLayout()

	wipeDataButton := widget.NewButtonWithIcon("Wipe data", theme.DeleteIcon(), func() {
		wipeDataWindow := epocellApp.NewWindow("Wipe data")
		wipeDataWindowIcon, _ := fyne.LoadResourceFromPath("icons/epocell-icon.png")
		wipeDataWindow.SetIcon(wipeDataWindowIcon)

		mainText := widget.NewLabel("This action will allow you to wipe data generated by the program.\n" +
			"RESTART THE APP after this.")

		wipeElements := []string{"Cache", "Config", "Output"}

		checkItems := make([]*widget.Check, 0)
		checkValues := make([]binding.Bool, 0)

		for range wipeElements {
			check := widget.NewCheck("", nil)
			checkItems = append(checkItems, check)
			checkValues = append(checkValues, binding.NewBool())
		}

		listWipe := container.NewVBox()

		for i, element := range wipeElements {
			row := container.NewHBox(
				widget.NewLabel(element),
				checkItems[i],
			)
			listWipe.Add(row)

			// Capture checkbox state change
			index := i // To avoid closure issues in the callback
			checkItems[i].OnChanged = func(checked bool) {
				err := checkValues[index].Set(checked)
				if err != nil {
					return
				}
				//fmt.Printf("Checkbox #%d (%s) changed to %v\n", index+1, element, checked)
			}
		}

		wipeOutputData := widget.NewLabel("")

		proceedButton := widget.NewButtonWithIcon("Confirm selection", theme.ConfirmIcon(), func() {
			checkboxStatus := make(map[string]bool)
			for element, value := range checkValues {
				checkboxStatus[wipeElements[element]], _ = value.Get()
			}

			for item, value := range checkboxStatus {

				location := ""
				if item == "Cache" {
					location = "cache/"
				}

				if item == "Config" {
					location = "config/"
				}

				if item == "Output" {
					location = "output/"
				}

				if value {
					deletedFiles, err := core.DeleteFilesRecursively(location)

					if err != nil {
						wipeOutputData.SetText("Error: " + err.Error())
						duration := 1 * time.Second
						time.Sleep(duration)
					} else {

						for _, deletedFile := range deletedFiles {
							wipeOutputData.SetText("Deleted: " + deletedFile)
							duration := 300 * time.Millisecond
							time.Sleep(duration)
						}
					}

					wipeOutputData.SetText("Wiping " + item + " done!")

					duration := 1 * time.Second
					time.Sleep(duration)
					wipeOutputData.SetText("")
				}

			}
		})

		cancelButton := widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), func() {
			wipeDataWindow.Close()
		})

		restartButton := widget.NewButtonWithIcon("Restart app", theme.ViewRefreshIcon(), func() {
			epocellApp.Quit()
			core.RestartApp()
		})

		proceedButtonContainer := container.NewHBox(proceedButton)
		proceedButtonContainer.Layout = layout.NewMaxLayout()

		cancelButtonContainer := container.NewHBox(cancelButton)
		cancelButtonContainer.Layout = layout.NewMaxLayout()

		restartButtonContainer := container.NewHBox(restartButton)
		restartButtonContainer.Layout = layout.NewMaxLayout()

		buttonsRow := container.New(layout.NewGridLayoutWithColumns(3), proceedButtonContainer, restartButtonContainer, cancelButtonContainer)

		wipeDataWindow.SetContent(container.NewVBox(
			mainText,
			listWipe,
			buttonsRow,
			wipeOutputData,
		))
		wipeDataWindow.Resize(wipeDataWindow.Canvas().Size())
		wipeDataWindow.CenterOnScreen()

		wipeDataWindow.Show()
	})

	wipeDataButtonContainer := container.NewHBox(wipeDataButton)
	wipeDataButtonContainer.Layout = layout.NewMaxLayout()

	// Create the multi-select widget
	rowButtons := container.New(layout.NewGridLayoutWithColumns(3), saveConfigContainer, storeConfigurationContainer, wipeDataButtonContainer)
	lastRowButtons := container.New(layout.NewGridLayoutWithColumns(1), computeButtonContainer)

	copyrigthLabel := widget.NewLabel("Owner: alexanderdth")
	copyrigthLabelContainer := container.NewHBox(copyrigthLabel)
	copyrigthLabelContainer.Layout = layout.NewMaxLayout()

	programData := widget.NewLabelWithStyle("Epocell Go v1.4", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true})
	programDataContainer := container.NewHBox(programData)
	programDataContainer.Layout = layout.NewMaxLayout()

	footerContainer := container.New(layout.NewGridLayoutWithColumns(2), copyrigthLabelContainer, programDataContainer)

	windowContainer := container.NewVBox(
		label,
		filePickerContainer,
		epocellStartsAtRowContainer,
		storeNameContainer,
		storeSalesUnitsContainer,
		storeSalesUnitsDeviationContainer,
		rowButtons,
		lastRowButtons,
		checkFormStatus,
		footerContainer,
	)

	// container2 := container.NewVScroll(list)
	// container2.SetMinSize(fyne.NewSize(200, 200))

	layoutContainer := container.NewVBox(windowContainer)
	epocellWindow.SetContent(layoutContainer)
	epocellWindow.Resize(fyne.NewSize(550, 275))
	epocellWindow.CenterOnScreen()
	epocellWindow.ShowAndRun()
}
