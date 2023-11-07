package core

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/tealeg/xlsx"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

func readChunk(filePath string, sheetIndex int, startRow, endRow int, columnIndices []int, columnLetters []string, wg *sync.WaitGroup, dataChan chan<- []map[string]string) {
	defer wg.Done()

	xlFile, err := xlsx.OpenFile(filePath)
	if err != nil {
		fmt.Printf("Error opening the file: %s\n", err)
		return
	}

	sheet := xlFile.Sheets[sheetIndex]
	chunkData := make([]map[string]string, 0)

	// Create a map to track processed rows
	processedRows := make(map[string]bool)

	for rowIdx := startRow; rowIdx < endRow; rowIdx++ {
		row := sheet.Rows[rowIdx]
		rowData := make(map[string]string)

		// Collect data for each column
		uniqueID := "" // Initialize an empty string for the unique ID

		for i, colIdx := range columnIndices {
			cell := row.Cells[colIdx]
			columnLetter := columnLetters[i]
			rowData[columnLetter] = cell.String()

			// Build the unique identifier using column values
			uniqueID += cell.String() + "_"
		}

		// Check if the row has already been processed
		if !processedRows[uniqueID] {
			// If not, add it to the chunkData and mark it as processed
			chunkData = append(chunkData, rowData)
			processedRows[uniqueID] = true
		}
	}

	dataChan <- chunkData
}

func ParseEpocell(filePath string, sheetIndex int, numChunks, startRow int, columnIndices []int, columnLetters []string) (chan []map[string]string, error) {
	xlFile, err := xlsx.OpenFile(filePath)
	if err != nil {
		fmt.Printf("Error opening the file: %s\n", err)
		return nil, nil
	}

	sheet := xlFile.Sheets[sheetIndex]
	if len(sheet.Rows) == 0 {
		fmt.Println("The sheet is empty.")
		return nil, nil
	}

	totalRows := len(sheet.Rows)

	// Ensure that startRow is within the valid range of rows.
	if startRow >= totalRows {
		fmt.Println("Start row exceeds the total number of rows.")
		return nil, nil
	}

	// Calculate the number of rows to read in each chunk.
	rowsPerChunk := totalRows / numChunks
	extraRows := totalRows % numChunks

	var wg sync.WaitGroup
	dataChan := make(chan []map[string]string, numChunks)

	for i := 0; i < numChunks; i++ {
		chunkStartRow := startRow + i*rowsPerChunk
		chunkEndRow := chunkStartRow + rowsPerChunk

		// Add extra rows to the last chunk.
		if i == numChunks-1 {
			chunkEndRow += extraRows
		}

		// Ensure that chunkEndRow does not exceed the totalRows.
		if chunkEndRow > totalRows {
			chunkEndRow = totalRows
		}

		wg.Add(1)
		go readChunk(filePath, sheetIndex, chunkStartRow, chunkEndRow, columnIndices, columnLetters, &wg, dataChan)
	}

	go func() {
		wg.Wait()
		close(dataChan)
	}()

	return dataChan, nil
}
func fileExists(filePath string) bool {
	if _, err := os.Stat(filePath); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		return false
	}
}

func FilterCacheData(data []StoreCache) []StoreCache {
	remaining := []StoreCache{}

	for _, item := range data {
		// Apply your filtering logic here
		if fileExists(item.CsvFile) {
			remaining = append(remaining, item)
		}
	}
	return remaining
}

func readCsvChunk(filePath string, startRow, endRow int, keywordColumnLetter string, columnIndices []int, columnLetters []string, keywords []string, skipKeywords []string, dataChan chan<- []map[string]string, wg *sync.WaitGroup) {
	defer wg.Done()

	csvFile, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening the file: %s\n", err)
		return
	}
	defer func(csvFile *os.File) {
		err := csvFile.Close()
		if err != nil {
			fmt.Printf("Error closing the file: %s\n", err)
		}
	}(csvFile)

	reader := csv.NewReader(csvFile)
	rows, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("Error reading CSV: %s\n", err)
		return
	}

	keywordColumnIdx := ColLetterToIndex(keywordColumnLetter)
	if keywordColumnIdx == -1 {
		fmt.Printf("Keyword column letter not found\n")
		return
	}

	chunkData := make([]map[string]string, 0)

	for rowIdx := startRow; rowIdx < endRow && rowIdx < len(rows); rowIdx++ {
		row := rows[rowIdx]

		keywordCellContent := strings.TrimSpace(row[keywordColumnIdx])

		if len(keywords) > 0 {
			if len(skipKeywords) > 0 {
				if containsKeyword(keywordCellContent, keywords) && doesntContainKeywords(keywordCellContent, skipKeywords) {
					rowData := make(map[string]string)

					for i, colIdx := range columnIndices {
						if colIdx >= 0 && colIdx < len(row) {
							cellContent := strings.TrimSpace(row[colIdx])
							rowData[columnLetters[i]] = cellContent
						}
					}

					chunkData = append(chunkData, rowData)
				}
			} else {
				if containsKeyword(keywordCellContent, keywords) {
					rowData := make(map[string]string)

					for i, colIdx := range columnIndices {
						if colIdx >= 0 && colIdx < len(row) {
							cellContent := strings.TrimSpace(row[colIdx])
							rowData[columnLetters[i]] = cellContent
						}
					}

					chunkData = append(chunkData, rowData)
				}
			}
		} else {
			rowData := make(map[string]string)

			for i, colIdx := range columnIndices {
				if colIdx >= 0 && colIdx < len(row) {
					cellContent := strings.TrimSpace(row[colIdx])
					rowData[columnLetters[i]] = cellContent
				}
			}
			chunkData = append(chunkData, rowData)
		}

	}

	dataChan <- chunkData
}

func UniqueStrings(slice []string) []string {
	encountered := map[string]bool{}
	result := []string{}

	for _, str := range slice {
		if !encountered[str] {
			encountered[str] = true
			result = append(result, str)
		}
	}

	return result
}

func containsKeyword(s string, keywords []string) bool {
	s = strings.ToLower(s) // Convert cell content to lowercase for case-insensitive matching

	for _, kw := range keywords {
		if strings.Contains(s, strings.ToLower(kw)) {
			return true
		}
	}
	return false
}

func doesntContainKeywords(s string, keywords []string) bool {
	lowerS := strings.ToLower(s) // Convert cell content to lowercase for case-insensitive matching

	for _, kw := range keywords {
		if strings.Contains(lowerS, strings.ToLower(kw)) {
			return false
		}
	}
	return true
}

func CreateAndSaveCSV(columnNames []string, filename string, data [][]string) error {
	// Create a new CSV file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write column names to the CSV file
	if err := writer.Write(columnNames); err != nil {
		return err
	}

	// Write data to the CSV file
	for _, row := range data {
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}
func ParseCsv(filePath string, numChunks, startRow int, keywordColumnLetter string, columnIndices []int, columnLetters []string, keywords []string, skipKeywords []string) (chan []map[string]string, error) {

	if len(skipKeywords) > 0 {
		fmt.Println("Found skipped keywords")
	} else {
		fmt.Println("Didnt find skipped keywords")
	}

	csvFile, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening the file: %s\n", err)
		return nil, err
	}
	defer func(csvFile *os.File) {
		err := csvFile.Close()
		if err != nil {

		}
	}(csvFile)

	reader := csv.NewReader(csvFile)
	rows, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("Error reading CSV: %s\n", err)
		return nil, err
	}

	totalRows := len(rows)

	// Ensure that startRow is within the valid range of rows.
	if startRow >= totalRows {
		fmt.Println("Start row exceeds the total number of rows.")
		return nil, nil
	}

	// Calculate the number of rows to read in each chunk.
	rowsPerChunk := totalRows / numChunks
	extraRows := totalRows % numChunks

	var wg sync.WaitGroup
	dataChan := make(chan []map[string]string, numChunks)

	for i := 0; i < numChunks; i++ {
		chunkStartRow := startRow + i*rowsPerChunk
		chunkEndRow := chunkStartRow + rowsPerChunk

		// Add extra rows to the last chunk.
		if i == numChunks-1 {
			chunkEndRow += extraRows
		}

		// Ensure that chunkEndRow does not exceed the totalRows.
		if chunkEndRow > totalRows {
			chunkEndRow = totalRows
		}

		wg.Add(1)
		go readCsvChunk(filePath, chunkStartRow, chunkEndRow, keywordColumnLetter, columnIndices, columnLetters, keywords, skipKeywords, dataChan, &wg)
	}

	go func() {
		wg.Wait()
		close(dataChan)
	}()

	return dataChan, nil
}

func ColLetterToIndex(letter string) int {
	letter = strings.ToUpper(letter)
	base := int('A') - 1 // 'A' should be mapped to index 0.
	colIdx := 0

	for i := 0; i < len(letter); i++ {
		colIdx = colIdx*26 + int(letter[i]) - base
	}

	return colIdx - 1 // Adjust for 0-based index.
}

func getColumnLetter(columnIndex int) string {
	// Convert the column index to Excel-like column letter
	columnLetter := ""
	for columnIndex >= 0 {
		mod := columnIndex % 26
		columnLetter = string(rune('A'+mod)) + columnLetter
		columnIndex = (columnIndex-mod)/26 - 1
	}
	return columnLetter
}

func MapColumnsWithLetters(csvFilePath string) ([]ColumnInfo, error) {
	var columns []ColumnInfo

	// Open the CSV file
	file, err := os.Open(csvFilePath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Unable to close file")
		}
	}(file)

	// Create a CSV reader
	reader := csv.NewReader(file)

	// Read the first row
	firstRow, err := reader.Read()
	if err != nil {
		return nil, err
	}

	// Populate the columns slice
	for i, columnName := range firstRow {
		columnLetter := getColumnLetter(i)
		column := ColumnInfo{
			Letter: columnLetter,
			Name:   columnName,
		}
		columns = append(columns, column)
	}

	return columns, nil
}

func DeleteFilesRecursively(dirPath string) ([]string, error) {
	deletedFiles := []string{}
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return deletedFiles, err
	}

	for _, file := range files {
		filePath := filepath.Join(dirPath, file.Name())
		if file.IsDir() {
			subDeletedFiles, err := DeleteFilesRecursively(filePath) // Recurse into subdirectories
			if err != nil {
				return deletedFiles, err
			}
			deletedFiles = append(deletedFiles, subDeletedFiles...)
		} else {
			err := os.Remove(filePath)
			if err != nil {
				return deletedFiles, err
			}
			deletedFiles = append(deletedFiles, filePath)
		}
	}
	return deletedFiles, nil
}
func FilterNonEmptyKeywords(keywords []string) []string {
	validKeywords := make([]string, 0)

	for _, kw := range keywords {
		trimmedKeyword := strings.TrimSpace(kw)
		if trimmedKeyword != "" {
			validKeywords = append(validKeywords, trimmedKeyword)
		}
	}

	return validKeywords
}
func ReadStoreCacheFromFile(fileName string) ([]StoreCache, error) {
	filePath := "cache/" + fileName
	// Read the file contents
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data into a slice of StoreCache structs
	var storeCache []StoreCache
	err = json.Unmarshal(data, &storeCache)
	if err != nil {
		return nil, err
	}

	return storeCache, nil
}

func UpdateStoreCacheToFile(filename string, data []StoreCache) error {
	updatedData, err := json.MarshalIndent(
		data,
		"",
		"    ",
	)
	if err != nil {
		return err
	}

	err = os.WriteFile("cache/"+filename, updatedData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func ReadStoreKeywordsCache(filename string) ([]StoreKeywordCache, error) {
	filePath := "cache/" + filename + ".json"

	// Read the file contents
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data into a slice of StoreKeywordCache
	var keywordCache []StoreKeywordCache
	err = json.Unmarshal(data, &keywordCache)
	if err != nil {
		return nil, err
	}

	return keywordCache, nil
}

func updateStoreKeywordCache(filename string, data []StoreKeywordCache) error {
	filePath := "cache/" + filename + ".json"

	updatedData, err := json.MarshalIndent(
		data,
		"",
		"    ",
	)
	if err != nil {
		return err
	}

	// Remove the existing file if it exists
	if _, err := os.Stat(filePath); err == nil {
		if err := os.Remove(filePath); err != nil {
			return err
		}
	}

	// Create and write to the new file
	err = os.WriteFile(filePath, updatedData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func UpdateStoreKeywordsCacheIfNecessary(cacheFileName string, storeKeywordCache []StoreKeywordCache) error {
	filePath := "cache/" + cacheFileName + ".json"

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		err := saveStoreKeywordCache(storeKeywordCache, cacheFileName)
		if err != nil {
			return err
		}
	} else if err == nil {
		err := updateStoreKeywordCache(cacheFileName, storeKeywordCache)
		if err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

func saveStoreKeywordCache(data []StoreKeywordCache, filename string) error {
	dataJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	cacheFolder := "cache"
	if _, err := os.Stat(cacheFolder); os.IsNotExist(err) {
		err := os.Mkdir(cacheFolder, os.ModePerm)
		if err != nil {
			return err
		}
	}

	filePath := fmt.Sprintf("%s/%s.json", cacheFolder, filename)

	// Check if the file already exists
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("a JSON file with the same name already exists")
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Unable to close file.")
		}
	}(file)

	_, err = file.Write(dataJSON)
	return err
}

func LoadJsonConfig() (Config, error) {
	configFile, err := os.Open("config/config.json")
	if err != nil {
		return Config{}, fmt.Errorf("error opening epocell configuration file: %w", err)
	}
	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {
			fmt.Println("Unable to close file.")
		}
	}(configFile)

	var configuration Config

	decoder := json.NewDecoder(configFile)
	if err := decoder.Decode(&configuration); err != nil {
		return Config{}, fmt.Errorf("error decoding JSON: %w", err)
	}

	return configuration, nil
}

func formatStoreName(storeName string) string {
	formatted := strings.ReplaceAll(storeName, " ", "_")
	reg := regexp.MustCompile("[^a-zA-Z0-9_]")
	formatted = reg.ReplaceAllString(formatted, "")
	return formatted
}

func SaveStoreCache(storeCaches []StoreCache, filename string) error {
	data, err := json.MarshalIndent(storeCaches, "", "  ")
	if err != nil {
		return err
	}

	cacheFolder := "cache"
	if _, err := os.Stat(cacheFolder); os.IsNotExist(err) {
		err := os.Mkdir(cacheFolder, os.ModePerm)
		if err != nil {
			return err
		}
	}

	filePath := fmt.Sprintf("%s/%s.json", cacheFolder, filename)

	// Check if the file already exists
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		return fmt.Errorf("a JSON file with the same name already exists")
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Unable to close file.")
		}
	}(file)

	_, err = file.Write(data)
	return err
}

func RemoveExtension(filename string) string {
	return filename[:len(filename)-len(filepath.Ext(filename))]
}

func FilterStoreCache(query string, storeCacheArray []StoreCache) []StoreCache {
	filteredStores := []StoreCache{}
	for _, store := range storeCacheArray {
		if strings.Contains(strings.ToLower(store.StoreName), strings.ToLower(query)) {
			filteredStores = append(filteredStores, store)
		}
	}
	return filteredStores
}

func FilterKeywords(query string, keywordCache []StoreKeywordCache) []StoreKeywordCache {
	var filteredKeywordCache []StoreKeywordCache
	lowerQuery := strings.ToLower(query)

	for _, keyword := range keywordCache {
		if strings.Contains(strings.ToLower(keyword.Keyword), lowerQuery) {
			filteredKeywordCache = append(filteredKeywordCache, keyword)
		}
	}

	return filteredKeywordCache
}

func RestartApp() {
	// Get the current executable path
	executablePath, _ := os.Executable()

	// Restart the application using exec.Command
	cmd := exec.Command(executablePath)
	err := cmd.Start()
	if err != nil {
		fmt.Println("Error restarting app:", err)
	}
}
