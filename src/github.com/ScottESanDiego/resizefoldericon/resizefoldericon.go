// Resize JPG files to 500x500
// Use to make thumbnails of album art
//
// Takes arguments to specify directory to use as base (e.g., /data/media/mp3)
//
// 2017-12-18	ScottE	Initial version
// 2017-12-20	ScottE	Lots of refactoring, added command line arguments

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"image"
	"image/jpeg"
	"github.com/nfnt/resize"
	"io/ioutil"
	"flag"
)

// Setup our constants here
const extension string = ".jpg" 		// Only work on files with this extension
const targetsize int = 500 				// This is used as both width and height
const aspecttolerance float32 = 0.98 	// This is how out-of-square we can be, and still get resized to be square

// Default to actually replacing files
var dryrun bool = false

// Compare passed in file's image dimensions with desired dimensions, and return true/false if this image needs resizing
func checkCandidateSize(CandidateFilePath string) (bool, int, int, error) {
	file, err := os.Open(CandidateFilePath)
	if err != nil {
		return false,0,0,err
	}
	defer file.Close()
	
	image, err := jpeg.DecodeConfig(file)
	if err != nil {
		return false,0,0,err
	}
	
	if (image.Width != targetsize) || (image.Height != targetsize) {
		aspect := float32(image.Width)/float32(image.Height)
		if (aspect < aspecttolerance) || (aspect > 1+(1-aspecttolerance)) {
			fmt.Println("Skipping",CandidateFilePath,": check Aspect Ratio")
			return false,0,0,nil
		}
		return true,image.Width,image.Height,nil
	}
	
	return false,0,0,nil
}

// Given a file passed in, resize it to the desired dimensions
func resizeJpeg(JpegPath string, OldWidth int, OldHeight int) (error) {
	fmt.Println("Resizing",JpegPath,"from",OldWidth,"x",OldHeight)

	file, err := os.Open(JpegPath)
	if err != nil {
		return err
	}
	defer file.Close()

	decodedimage, err := jpeg.Decode(file)
	if err != nil {
		return err
	}
	imgResized := resize.Resize(uint(targetsize),uint(targetsize),decodedimage,resize.Bicubic)
	
	err = saveResizedJpeg(JpegPath, imgResized)
	if err != nil {
		return err
	}
	
	return nil
}

// Given an Image and a path, save the Image to the path
func saveResizedJpeg(JpegPath string, ResizedImage image.Image) (error) {	
	tmpfile, err := ioutil.TempFile(filepath.Dir(JpegPath),"resize")
	if err != nil {
		return err
	}

	jpeg.Encode(tmpfile, ResizedImage, nil)
		
	if dryrun {
		defer os.Remove(tmpfile.Name())
	} else {

		tmpfile.Close()
		
		err := os.Rename(tmpfile.Name(), JpegPath)
		if err != nil {
			return err
		}
		
		// Do this just to clean up temporarily
		defer os.Remove(tmpfile.Name())
	}
		
	return nil
}

// Check if the file matches our extension, and then if it has a size other than what we want
func findCandidates(path string, info os.FileInfo, err error) (error) {
	// the filepath.Walk() call could have passed in an error - if so, just abort
	if err != nil {
		return err
	}
	
	r, err := regexp.MatchString(extension, info.Name())
	if err != nil {
		return err
	}

	if r {
		needResize, width, height, err := checkCandidateSize(path)
		if err != nil {
			return err
		}
		
		if needResize {
			err := resizeJpeg(path,width,height)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func main() {
	// Don't print date and time in log message
	log.SetFlags(log.Lshortfile)
	
	dir := flag.String("d", "", "Directory to recurse through")
	flag.BoolVar(&dryrun,"n", false, "Do a dry run and don't actually replace files")
	flag.Parse()
	
	if *dir == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	
	if dryrun {
		fmt.Println("***Dry run - not actually replacing files***")
	}

	err := filepath.Walk(*dir, findCandidates)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

