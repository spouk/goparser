package library

import (
	. "github.com/kkdai/youtube"
	"strings"
	"fmt"
	"os"
	"io"
	"sync"
	"net/http"
	"errors"
	"log"
)

type Downloader struct {
	sync.WaitGroup
}
func (d *Downloader) DownloaderVideo(nameGoroutine, pathSaveDir, linkYoutube string, pool bool) (error) {
	if pool {
		defer func() {
			d.Done()
		}()
	}

	//check exists save path
	_, err := os.Stat(pathSaveDir)
	if err != nil {
		//create new dir
		err := os.Mkdir(pathSaveDir, os.ModeDir|os.ModePerm)
		if err != nil {
			return err
		}
	}
	//download instance
	y := NewYoutube(true)
	err = y.DecodeURL(linkYoutube)
	if err != nil {
		return err
	}
	//extract filename + ext + make outputfilename
	var tt = strings.Split(y.StreamList[0]["type"], ";")[0]
	var typeVideo = strings.Split(tt, "/")[1]                                    //ex: mp4
	var name = strings.Join([]string{y.StreamList[0]["author"], typeVideo}, ".") //fname
	//var fname = strings.Join([]string{name, typeVideo}, ".")

	//download video
	ps := strings.Join([]string{pathSaveDir, name}, "/")
	fmt.Printf("PATH SAVE: %v\n", ps)
	//go showPercent(nameGoroutine, ps, y.DownloadPercent)
	//y.StartDownload("/tmp/tester.mp4")
	y.StartDownload(ps)

	return nil
}

func (d *Downloader) DownloadImage(pathToSaveWithFilename, link string, pool bool) (error) {
	if pool {
		defer func() {
			d.Done()
		}()
	}
	resp, err := http.Get(link)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("http error: %v\n", http.StatusForbidden))
	}
	file, err := os.Create(pathToSaveWithFilename)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}
	file.Close()
	log.Printf("%s success download\n", pathToSaveWithFilename)
	return nil
}

