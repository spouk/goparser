package library

import (
	. "github.com/kkdai/youtube"
	"strings"
	"fmt"
	"os"
	"io"
	"net/http"
	"errors"
	"log"
	"path/filepath"
	"path"
)

func (d *Parser) DownloaderVideo(nameGoroutine, pathSaveDir, linkYoutube string, pool bool) (error) {
	if pool {
		defer func() {
			d.Done()
		}()
	}
	log.Printf("[DownloadVideo] starting `%s`\n", linkYoutube)

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

func (d *Parser) DownloadImage(pathToSaveWithFilename, link string, pool bool) (error) {
	if pool {
		defer func() {
			d.Done()
		}()
	}
	log.Printf("[DownloadImage] starting `%s`\n", link)
	resp, err := http.Get(link)
	if err != nil {
		log.Printf("[DownloadImage] [%s] error: %v\n", link, err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("http error: %v\n", http.StatusForbidden))
	}
	var fname = path.Base(link)
	file, err := os.Create(strings.Join([]string{pathToSaveWithFilename,fname}, "/"))
	if err != nil {
		log.Printf("[DownloadImage] [%s] error: %v\n", link, err)
		return err
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Printf("[DownloadImage] [%s] error: %v\n", link, err)
		return err
	}
	file.Close()
	log.Printf("%s success download\n", pathToSaveWithFilename)
	return nil
}

