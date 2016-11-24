package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/context"

	"github.com/speedland/go/iterator/slice"
	vr "github.com/speedland/go/services/watson/visualrecognition"
	"github.com/speedland/go/x/xarchive/xzip"
	"github.com/speedland/go/x/ximage"
)

// collectSources collects source files for classifier inputs.
// the directory structure should be:
//
//     {rootDir}/
//         {class_name}/
//              {img}.png
//              {img}.jpeg
//         {class_name}/
//              {img}.png
//              {img}.jpeg
//              something.txt // <- ignored
//         something.txt      // <- ignored
//
// and the returned map will be a map from {class_name} to [{img}, {img}, ...]
func collectSources(rootDir string) (map[string][]string, error) {
	sources := make(map[string][]string)
	log.Printf("Collecting source files from %s...\n", rootDir)
	classDirs, err := ioutil.ReadDir(rootDir)
	if err != nil {
		return nil, err
	}
	for _, classDirInfo := range classDirs {
		classDirPath := filepath.Join(rootDir, classDirInfo.Name())
		if !classDirInfo.IsDir() {
			log.Printf("S: %s\n", classDirPath)
			continue
		} else {
			var filePathList []string
			files, err := ioutil.ReadDir(classDirPath)
			if err != nil {
				log.Printf("E: %s - %s\n", classDirPath, err)
				continue
			}
			for _, fileInfo := range files {
				fpath := filepath.Join(classDirPath, fileInfo.Name())
				if strings.HasPrefix(mime.TypeByExtension(filepath.Ext(fpath)), "image/") {
					filePathList = append(filePathList, fpath)
				} else {
					log.Printf("S: %s\n", fpath)
				}
			}
			if len(filePathList) > 0 {
				sources[classDirInfo.Name()] = filePathList
			}
		}
	}
	return sources, nil
}

func extractFacesFromSources(sources map[string][]string, tempdir string, client *vr.Client, size int) (map[string][]string, error) {
	newSources := make(map[string][]string)
	totalFiles := 0
	for _, files := range sources {
		totalFiles += len(files)
	}
	processedFiles := 0
	for class, files := range sources {
		var faceFiles []string
		dstDir := filepath.Join(tempdir, class)
		if err := os.MkdirAll(dstDir, os.FileMode(0755)); err != nil {
			return nil, err
		}
		_list := slice.SplitByLength(files, size)
		for _, list := range _list.([][]string) {
			resp, err := func() (*vr.FaceDetectResponse, error) {
				var rsList []*xzip.RawSource
				for _, filepath := range list {
					rs, err := xzip.NewRawSourceFromFile(filepath)
					if err != nil {
						log.Printf("E: %s - %v", filepath, err)
						continue
					}
					defer rs.Close()
					log.Printf("A: %s", filepath)
					rsList = append(rsList, rs)
				}
				// call the API
				log.Printf("[%d/%d] Requesting files to detect faces ...", processedFiles, totalFiles)
				processedFiles += len(rsList)
				return client.DetectFacesOnImages(context.Background(), xzip.NewArchiver(rsList...))
			}()
			if err != nil { // unrecoverable API error so return.
				return nil, err
			}
			for i, image := range resp.Images {
				srcPath := list[i]
				if len(image.Faces) == 0 {
					log.Printf("S: %s - no faces detected", srcPath)
					continue
				}
				if len(image.Faces) > 1 {
					log.Printf("S: %s - multiple faces detected", srcPath)
					continue
				}
				dstPath := filepath.Join(dstDir, filepath.Base(srcPath))
				err := func() error {
					loc := image.Faces[0].FaceLocation
					src, err := os.Open(srcPath)
					if err != nil {
						return err
					}
					dst, err := os.Create(dstPath)
					if err != nil {
						return err
					}
					return ximage.Crop(
						src, dst, ximage.TypeByExtension(filepath.Ext(srcPath)),
						int(loc.Left), int(loc.Top), int(loc.Left+loc.Width), int(loc.Top+loc.Height),
					)
				}()
				if err != nil {
					log.Printf("E: %s - error creating a face file: %v", srcPath, err)
					continue
				}
				log.Printf("C: %s", dstPath)
				faceFiles = append(faceFiles, dstPath)
			}
		}
		newSources[class] = faceFiles
	}
	return newSources, nil
}

// buildExamples builds the example files from sources
func buildExamples(sources map[string][]string) (map[string]*xzip.Archiver, *xzip.Archiver) {
	positives := make(map[string]*xzip.Archiver)
	var negative *xzip.Archiver
	for class, files := range sources {
		var rawSources []*xzip.RawSource
		for _, file := range files {
			if rs, err := xzip.NewRawSourceFromFile(file); err != nil {
				log.Printf("E: %s - %v", file, err)
			} else {
				log.Printf("A: %s", file)
				rawSources = append(rawSources, rs)
			}
		}
		if class == "negative" {
			negative = xzip.NewArchiver(rawSources...)
		} else {
			positives[class] = xzip.NewArchiver(rawSources...)
		}
	}
	return positives, negative
}

// prepareSources to prepare the raw sources from file path or url.
// For urls, this downloads them into tempdir and handle as local file on it.
func prepareSources(tempdir string, args ...string) ([]*xzip.RawSource, []*sourceInfo) {
	var rawSources []*xzip.RawSource
	var sources []*sourceInfo
	for _, arg := range args {
		var err error
		var rs *xzip.RawSource
		var info *sourceInfo
		if strings.HasPrefix(arg, "http://") || strings.HasPrefix(arg, "https://") {
			var file string
			fmt.Printf("Downloading %s....\n", arg)
			file, err = download(arg, tempdir)
			if err == nil {
				rs, err = xzip.NewRawSourceFromFile(file)
				info = &sourceInfo{
					orig: arg,
					temp: file,
				}
			}
		} else {
			rs, err = xzip.NewRawSourceFromFile(arg)
			info = &sourceInfo{
				orig: arg,
			}
		}
		if err != nil {
			log.Printf("E: %s - %s\n", arg, err)
			continue
		}
		log.Printf("A: %s", arg)
		rawSources = append(rawSources, rs)
		sources = append(sources, info)
	}
	return rawSources, sources
}
