package media

import (
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"weigo/api/database"
	"weigo/api/helper"

	"github.com/labstack/echo"
	"github.com/nfnt/resize"
	uuid "github.com/satori/go.uuid"
)

const (
	maxUploadSize = int64(500 * 1024 * 1024) // 5MB

)

var (
	fileExtRegex = regexp.MustCompile(`(png|jpe?g|gif)$`)
)

func storeJPEG(nameUUID string, resizedImage image.Image) (string, error) {
	outFile, err := os.Create(filepath.Join("../ui/uploaded/media/" + nameUUID + ".jpg"))
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer outFile.Close()
	opt := &jpeg.Options{
		Quality: 70,
	}
	err = jpeg.Encode(outFile, resizedImage, opt)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return nameUUID + ".jpg", nil
}

func storeGIF(nameUUID string, resizedImage image.Image) (string, error) {
	outFile, err := os.Create(filepath.Join("../ui/uploaded/media/" + nameUUID + ".gif"))
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer outFile.Close()
	err = gif.Encode(outFile, resizedImage, nil)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return nameUUID + ".gif", nil
}

func postMedia(c echo.Context) error {
	if c.Request().ContentLength > maxUploadSize {
		return c.JSON(http.StatusRequestEntityTooLarge, helper.BuildErrorResponse("the media is too large"))
	}
	defer c.Request().Body.Close()
	if err := c.Request().ParseMultipartForm(1024); err != nil {
		return err
	}
	form := c.Request().MultipartForm
	defer form.RemoveAll()
	fileH, _, err := c.Request().FormFile("media")
	if err != nil {
		return err
	}
	defer fileH.Close()
	fImage := fileH.(*os.File)
	img, ext, err := image.Decode(fImage)
	if err != nil {
		return c.JSON(http.StatusNotAcceptable, helper.BuildErrorResponse("invalid image"))
	}
	if !fileExtRegex.MatchString(ext) {
		return c.JSON(http.StatusNotAcceptable, helper.BuildErrorResponse("invalid file type"))
	}

	var outName = ""
	var errProcessImage error
	switch strings.ToLower(ext) {
	case "jpeg", "png":
		resizedImage := resize.Thumbnail(1000, 1000, img, resize.Lanczos3)
		outName, errProcessImage = storeJPEG(uuid.NewV4().String(), resizedImage)
	case "gif":
		outName, errProcessImage = storeGIF(uuid.NewV4().String(), img)
	default:
		return errors.New("unknow file")
	}
	if errProcessImage != nil {
		return errProcessImage
	}
	p := c.Get("user").(*database.Profile)
	if _, err := database.MediaCreated(p, outName); err != nil {
		log.Println(err)
		return err
	}
	return c.JSON(http.StatusOK, map[string]string{
		"id": outName,
	})
}

// Init .
func Init(group *echo.Group) {
	group.Use(helper.AuthMiddleware(), helper.GetUserMiddleware)
	group.POST("", postMedia)
}
