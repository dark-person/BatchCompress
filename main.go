package main

import (
	"archive/zip"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
)

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func getInputFolder() (folder string, filenames []string) {
	var result_file_list []string

	// Get input folder from user
	var input_folder string
	fmt.Println("Enter the folder to compress..")
	fmt.Scanln(&input_folder)

	// Verify user input by open folder
	_, err_open := os.Open(input_folder)
	if err_open != nil {
		fmt.Println("Error opening input folder, please enter valid folder.")
		return getInputFolder()
	}

	filelist, err1 := os.ReadDir(input_folder)
	if err1 != nil {
		fmt.Println("Error opening input folder, please enter valid folder.")
		return getInputFolder()
	}

	// Put file into string list
	for _, file := range filelist {
		result_file_list = append(result_file_list, filepath.Join(input_folder, file.Name()))
	}

	return input_folder, result_file_list
}

func getFileNumber() int {
	var number int
	fmt.Println("Enter the number of file to contain in one zip:")
	fmt.Scanln(&number)

	return number
}

func getStringByInput(message string) string {
	var temp string
	fmt.Println("Enter the folder to compress..")
	fmt.Scanln(&temp)
	return temp
}

func ConfirmPrompt(message string) bool {
	var decision string
	fmt.Print(message, "(Y/n) ")
	fmt.Scanln(&decision)

	if decision == "Y" || decision == "y" {
		return true
	} else if decision == "N" || decision == "n" {
		return false
	} else {
		return ConfirmPrompt("Invalid Input. Please Enter Again.\n\n" + message)
	}
}

// func appendFiles(filename string, zipw *zip.Writer) error {
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		return fmt.Errorf("Failed to open %s: %s", filename, err)
// 	}
// 	defer file.Close()

// 	wr, err := zipw.Create(filename)
// 	if err != nil {
// 		msg := "Failed to create entry for %s in zip file: %s"
// 		return fmt.Errorf(msg, filename, err)
// 	}

// 	if _, err := io.Copy(wr, file); err != nil {
// 		return fmt.Errorf("Failed to write %s to zip: %s", filename, err)
// 	}

// 	return nil
// }

func createFlatZip(w io.Writer, files ...string) error {
	z := zip.NewWriter(w)
	for _, file := range files {
		src, err := os.Open(file)
		if err != nil {
			return err
		}
		info, err := src.Stat()
		if err != nil {
			return err
		}
		hdr, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		hdr.Name = filepath.Base(file) // Write only the base name in the header
		dst, err := z.CreateHeader(hdr)
		if err != nil {
			return err
		}
		_, err = io.Copy(dst, src)
		if err != nil {
			return err
		}
		src.Close()
	}
	return z.Close()
}

func main() {
	// fmt.Println("Hello, world!")

	// ========================== User Input ===========================

	// Get the Input folder by user
	_, filelist := getInputFolder()
	// fmt.Println("Test Filelist : ", filelist)

	// Get number of file in one zip file from user
	number := getFileNumber()

	// Get Archive Name
	archive_name := getStringByInput("Input the zip filename (e.g. \"test\" will output \"test_1.zip\")")

	// ========================== User Input ===========================

	// ========================== User Confirm ========================

	// Check how many compress file are needed
	compress_needed := int(math.Ceil(float64(len(filelist)) / float64(number)))
	// fmt.Println("Compress File will be Generated: ", compress_needed)
	temp := ConfirmPrompt(fmt.Sprintf("There will be %d compress files generated, will be named as format of %s_1.zip", compress_needed, archive_name))
	if !temp {
		os.Exit(99)
	}

	// ========================== User Confirm ========================

	// Looping Start
	previous := 0
	for i := 0; i < compress_needed; i++ {
		archive_name := "test_" + strconv.Itoa(i) + ".zip"
		archive, err := os.Create(archive_name)
		if err != nil {
			fmt.Printf("Failed to generate archive %s: %v\n", archive_name, err)
		}
		defer archive.Close()

		high_bound := Min((i+1)*number, len(filelist))
		fmt.Println("Test", i, ": ", filelist[previous:high_bound])

		err2 := createFlatZip(archive, filelist[previous:high_bound]...)
		if err2 != nil {
			fmt.Printf("Failed to put file into archive %s: %v\n", archive_name, err)
		}

		previous = (i + 1) * number
	}

	// Report user by terminal message
	fmt.Println("All file has archived successfully.")

	// todo Delete the file that compressed

}
