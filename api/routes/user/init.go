package user

import (
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
	"weigo/api/database"
	"weigo/api/helper"

	// Add images .
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/labstack/echo"
	"github.com/nfnt/resize"
	"github.com/satori/go.uuid"
)

const (
	// MaxAvatarSize .
	MaxAvatarSize = int64(2 * 1024 * 1024)
)

func generateToken(p *database.Profile) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = uint64(p.ID)
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	return token.SignedString([]byte(helper.Secret))
}

func login(c echo.Context) error {
	var data struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.Bind(&data); err != nil {
		log.Println(err)
		return err
	}
	p, err := database.ProfileGetByUsername(data.Username)
	if err == database.ErrNotFound {
		return c.JSON(http.StatusNotFound, helper.BuildErrorResponse("User not found"))
	} else if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Error"))
	}
	if p.Authenticate(data.Password) {
		t, err := generateToken(p)
		if err != nil {
			log.Println(err)
			return err
		}
		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	}
	return echo.ErrUnauthorized
}

func loginByKey(c echo.Context) error {
	key := c.Param("key")
	op, err := database.OauthProfile{}.GetByKey(key)
	if err != nil {
		log.Println(err)
		return err
	}
	if op == nil {
		return c.JSON(http.StatusNotFound, helper.BuildErrorResponse("Token not found"))
	}
	var profile *database.Profile
	var getProfileErr error
	if op.Type == "github" {
		profile, getProfileErr = database.ProfileGetByGithubID(op.AccountID)
	} else if op.Type == "google" {
		profile, getProfileErr = database.ProfileGetByGoogleID(op.AccountID)
	}
	if getProfileErr != nil {
		log.Println(getProfileErr)
		return c.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Error"))
	} else if profile == nil {
		return c.JSON(http.StatusNotFound, helper.BuildErrorResponse("User not found"))
	}
	t, err := generateToken(profile)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Error"))
	}
	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

func loginSetupGet(c echo.Context) error {
	key := c.Param("key")
	op, err := database.OauthProfile{}.GetByKey(key)
	if err != nil {
		log.Println(err)
		return err
	}
	if op == nil {
		return c.JSON(http.StatusNotFound, helper.BuildErrorResponse("Token not found"))
	}
	return c.JSON(http.StatusOK, *op)
}

func loginSetupPost(c echo.Context) error {
	key := c.Param("key")
	op, err := database.OauthProfile{}.GetByKey(key)
	if err != nil {
		log.Println(err)
		return err
	}
	if op == nil {
		return c.JSON(http.StatusNotFound, helper.BuildErrorResponse("Token not found"))
	}
	var data struct {
		Username string `json:"username"`
	}
	if err := c.Bind(&data); err != nil {
		return err
	}
	err = validation.ValidateStruct(
		&data,
		validation.Field(&data.Username, validation.Required),
	)
	if err != nil {
		errs := err.(validation.Errors)
		outErrors := []string{}
		for f, e := range errs {
			outErrors = append(outErrors, fmt.Sprintf("%s %s", f, e))
		}
		return c.JSON(http.StatusBadRequest, helper.BuildErrorsResponse(outErrors))
	}
	var createProfileError error
	if op.Type == "github" {
		_, createProfileError = database.ProfileCreateFromGitbub(data.Username, op.AccountID)
	} else if op.Type == "google" {
		_, createProfileError = database.ProfileCreateFromGoogle(data.Username, op.AccountID)
	}
	if createProfileError == database.ErrDuplicated {
		return c.JSON(http.StatusBadRequest, helper.BuildErrorResponse("This username is already taken"))
	} else if createProfileError != nil {
		return c.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("error"))
	}
	return c.String(http.StatusNoContent, "")
}

func register(c echo.Context) error {
	var data struct {
		Email           string `json:"email" valid:"email~Invalid email,required"`
		Username        string `json:"username" valid:"required"`
		Password        string `json:"password" valid:"required,minLength(6)~password too small (min 6 characters)"`
		PasswordConfirm string `json:"passwordConfirm" valid:"required,minLength(10)~passwordConfirm too small (min 6 characters)"`
	}
	if err := c.Bind(&data); err != nil {
		return err
	}
	err := validation.ValidateStruct(
		&data,
		validation.Field(&data.Email, validation.Required, is.Email),
		validation.Field(&data.Username, validation.Required),
		validation.Field(&data.Password, validation.Required, validation.Match(regexp.MustCompile(`^.{6,}$`)).Error("too small")),
		validation.Field(&data.PasswordConfirm, validation.Required, validation.Match(regexp.MustCompile(`^.{6,}$`)).Error("too small")),
	)
	if err != nil {
		errs := err.(validation.Errors)
		outErrors := []string{}
		for f, e := range errs {
			outErrors = append(outErrors, fmt.Sprintf("%s %s", f, e))
		}
		return c.JSON(http.StatusBadRequest, helper.BuildErrorsResponse(outErrors))
	}

	_, err = database.ProfileCreate(data.Email, data.Username, data.Password)
	if err == database.ErrDuplicated {
		return c.JSON(http.StatusBadRequest, helper.BuildErrorResponse("User or email already exists"))
	} else if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"ok": "ok",
	})
}

func user(c echo.Context) error {
	return c.JSON(http.StatusOK, c.Get("user"))
}

func userByID(c echo.Context) error {
	profileID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}
	p, err := database.ProfileGetByID(profileID)
	if err != nil {
		return nil
	}
	return c.JSON(http.StatusOK, p)
}

func searchUser(c echo.Context) error {
	p := c.Get("user").(*database.Profile)
	var data struct {
		Content string `json:"contains"`
	}
	if err := c.Bind(&data); err != nil {
		return err
	}
	profiles, err := database.ProfileSearchByName(p, data.Content)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, profiles)
}

func addFriend(c echo.Context) error {
	p := c.Get("user").(*database.Profile)
	profileID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}
	err = database.ProfileAddFriendByID(p, profileID)
	if err != nil {
		return err
	}
	return c.String(http.StatusNoContent, "")
}

func removeFriend(c echo.Context) error {
	p := c.Get("user").(*database.Profile)
	profileID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}
	err = database.ProfileRemoveFriendByID(p, profileID)
	if err != nil {
		return err
	}
	return c.String(http.StatusNoContent, "")
}

func listFriends(c echo.Context) error {
	p := c.Get("user").(*database.Profile)
	profiles, err := database.ProfileListFriends(p)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, profiles)
}

func updatePassword(c echo.Context) error {
	p := c.Get("user").(*database.Profile)
	var data struct {
		OldPassword        string `json:"oldPassword"`
		NewPassword        string `json:"newPassword"`
		NewPasswordConfirm string `json:"newPasswordConfirm"`
	}
	if err := c.Bind(&data); err != nil {
		return nil
	}
	err := validation.ValidateStruct(
		&data,
		validation.Field(
			&data.NewPassword,
			validation.Required,
			validation.Match(regexp.MustCompile(`^.{6,}$`)).Error("too small (6 char min.)"),
		),
		validation.Field(
			&data.NewPasswordConfirm,
			validation.Required,
			validation.Match(regexp.MustCompile(`^.{6,}$`)).Error("too small  (6 char min.)"),
		),
	)
	if err != nil {
		if err != nil {
			errs := err.(validation.Errors)
			outErrors := []string{}
			for f, e := range errs {
				outErrors = append(outErrors, fmt.Sprintf("%s %s", f, e))
			}
			return c.JSON(http.StatusBadRequest, helper.BuildErrorsResponse(outErrors))
		}
	}
	if data.NewPassword != data.NewPasswordConfirm {
		return c.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Password confirmation mismatch"))
	}
	if !p.Authenticate(data.OldPassword) {
		return c.JSON(http.StatusForbidden, helper.BuildErrorResponse("Your old password doesn't match"))
	}
	return database.ProfileUpdatePassword(p, data.NewPassword)
}

func readPart(part *multipart.Part, p *database.Profile) error {
	defer part.Close()
	if part.FormName() != "avatar" {
		return nil
	}

	os.Remove("../ui/uploaded/avatar/" + strconv.FormatUint(p.ID, 10) + ".jpeg")
	picture, _ := os.Create("../ui/uploaded/avatar/" + strconv.FormatUint(p.ID, 10) + ".jpeg")
	defer picture.Close()
	io.Copy(picture, part)
	return nil
}

func userAvatar(c echo.Context) error {
	p := c.Get("user").(*database.Profile)
	if c.Request().ContentLength > MaxAvatarSize {
		c.Request().Body.Close()
		return c.String(http.StatusRequestEntityTooLarge, "")
	}
	c.Request().Body = http.MaxBytesReader(c.Response(), c.Request().Body, MaxAvatarSize)
	defer c.Request().Body.Close()
	err := c.Request().ParseMultipartForm(1024)
	if err != nil {
		log.Println(err)
		return err
	}
	defer c.Request().MultipartForm.RemoveAll()
	fi, info, err := c.Request().FormFile("avatar")
	if err != nil {
		fmt.Println(err)
		return err
	}
	fi.Close()
	fimage, err := info.Open()
	if err != nil {
		log.Println(err)
		return err
	}
	defer fimage.Close()
	img, _, err := image.Decode(fimage.(*os.File))
	if err != nil {
		log.Println(err)
		return err
	}
	m := resize.Thumbnail(500, 500, img, resize.Lanczos3)
	imageID := uuid.NewV4().String()
	out, err := os.Create(filepath.Join("..", "ui", "uploaded", "avatar", fmt.Sprintf("%d-%s.jpg", p.ID, imageID)))
	if err != nil {
		log.Println(err)
	}
	defer out.Close()
	jpeg.Encode(out, m, nil)
	return p.ChangeAvatar(imageID)
}

func updateInfo(c echo.Context) error {
	p := c.Get("user").(*database.Profile)
	var data struct {
		Username string `json:"username"`
	}
	if err := c.Bind(&data); err != nil {
		return nil
	}
	err := validation.ValidateStruct(
		&data,
		validation.Field(
			&data.Username,
			validation.Required,
			validation.Match(regexp.MustCompile(`^[a-zA-Z0-9_]+$`)).Error("invalid characters"),
		),
	)
	if err != nil {
		if err != nil {
			errs := err.(validation.Errors)
			outErrors := []string{}
			for f, e := range errs {
				outErrors = append(outErrors, fmt.Sprintf("%s %s", f, e))
			}
			return c.JSON(http.StatusBadRequest, helper.BuildErrorsResponse(outErrors))
		}
	}
	err = database.ProfileUpdateUsername(p, data.Username)
	if err == database.ErrDuplicated {
		return c.JSON(http.StatusBadRequest, helper.BuildErrorResponse("this username already exists"))
	}
	return nil
}

// Init .
func Init(group *echo.Group) {
	group.GET("/login/setup/:key", loginSetupGet)
	group.POST("/login/setup/:key", loginSetupPost)
	group.GET("/login/:key", loginByKey)
	group.POST("/login", login)
	group.POST("/register", register)
	group.GET("", user, helper.AuthMiddleware(), helper.GetUserMiddleware)
	group.PUT("/password", updatePassword, helper.AuthMiddleware(), helper.GetUserMiddleware)
	group.PUT("/info", updateInfo, helper.AuthMiddleware(), helper.GetUserMiddleware)
	group.POST("/avatar", userAvatar, helper.AuthMiddleware(), helper.GetUserMiddleware)
	group.POST("/search", searchUser, helper.AuthMiddleware(), helper.GetUserMiddleware)
	group.POST("/friendship/:id", addFriend, helper.AuthMiddleware(), helper.GetUserMiddleware)
	group.DELETE("/friendship/:id", removeFriend, helper.AuthMiddleware(), helper.GetUserMiddleware)
	group.GET("/friendship", listFriends, helper.AuthMiddleware(), helper.GetUserMiddleware)
	group.GET("/:id", userByID, helper.AuthMiddleware(), helper.GetUserMiddleware)
	group.GET("/oauth/google/callback", googleOauth)
	group.GET("/oauth/github/callback", githubOauth)
}
