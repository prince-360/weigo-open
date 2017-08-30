package stream

import (
	"log"
	"net/http"
	"strconv"
	"weigo/api/database"
	"weigo/api/helper"

	"github.com/labstack/echo"
)

func newPost(c echo.Context) error {
	p := c.Get("user").(*database.Profile)
	var data struct {
		Content string   `json:"content"`
		Medias  []string `json:"medias"`
	}
	if err := c.Bind(&data); err != nil {
		return err
	}
	if len(data.Medias) > 4 {
		data.Medias = data.Medias[0:4]
	}
	post, err := database.PostCreate(p, data.Content, data.Medias)
	if err != nil {
		log.Println(err)
		return nil
	}
	return c.JSON(http.StatusOK, post)
}

func listPosts(c echo.Context) error {
	p := c.Get("user").(*database.Profile)
	var fromString = c.QueryParam("from")
	var from *uint64
	if f, err := strconv.ParseUint(fromString, 10, 64); err == nil {
		from = &f
	}
	posts, err := database.PostListIDs(p, from)
	if err != nil {
		return nil
	}
	return c.JSON(http.StatusOK, posts)
}

func getPost(c echo.Context) error {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}
	post, err := database.PostGetByID(postID)
	if err != nil {
		return nil
	}
	if post == nil {
		return c.String(http.StatusNotFound, "")
	}
	return c.JSON(http.StatusOK, post)
}

func likePost(c echo.Context) error {
	p := c.Get("user").(*database.Profile)
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}
	err = database.PostLikeByID(postID, p.ID)
	if err != nil {
		return nil
	}
	return c.String(http.StatusNoContent, "")
}

func unlikePost(c echo.Context) error {
	p := c.Get("user").(*database.Profile)
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}
	err = database.PostUnlikeByID(postID, p.ID)
	if err != nil {
		return nil
	}
	return c.String(http.StatusNoContent, "")
}

func newComment(c echo.Context) error {
	p := c.Get("user").(*database.Profile)
	var data struct {
		Content string `json:"content"`
	}
	if err := c.Bind(&data); err != nil {
		return err
	}
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}
	commentID, err := database.CommentCreate(postID, p.ID, data.Content)
	if err != nil {
		return nil
	}
	return c.JSON(http.StatusOK, map[string]uint64{
		"comment_id": commentID,
	})
}

func listComment(c echo.Context) error {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}
	items, err := database.CommentListByPostID(postID)
	if err != nil {
		return nil
	}
	return c.JSON(http.StatusOK, items)
}

// Init .
func Init(group *echo.Group) {
	group.Use(helper.AuthMiddleware(), helper.GetUserMiddleware)
	group.POST("/post", newPost)
	group.GET("/post", listPosts)
	group.GET("/post/:id", getPost)
	group.GET("/post/:id/like", likePost)
	group.GET("/post/:id/unlike", unlikePost)
	group.POST("/post/:id/comment", newComment)
	group.GET("/post/:id/comment", listComment)

}
