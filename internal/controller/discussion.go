package controller

import (
    "time"
    "strconv"
    "net/http"
    "github.com/labstack/echo/v4"

    "user-management-system/internal/model"
    "user-management-system/internal/utils"
    "user-management-system/internal/controller/param"
    "user-management-system/internal/controller/context"
)

func PostDiscussion(c echo.Context) error {
    claims := c.Get("claims").(*utils.Claims)
    req := new(param.PostDiscussionRequest)
    if err := context.BindAndVali(c, req); err != nil {
        return echo.ErrBadRequest
    }
    discussion := model.Discussion{}
    discussion.Title = req.Title
    discussion.Content = req.Content
    discussion.UserId = claims.UserId
    discussion.Time = time.Now().Unix()
    if err := model.CreateDiscussion(&discussion); err != nil {
        return echo.ErrInternalServerError
    }
    var resp struct {
        Message string  `json:"message"`
    }
    resp.Message = "Success!"
    return c.JSON(http.StatusOK, &resp)
}

func DiscussionInfo(c echo.Context) error {
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        return echo.ErrNotFound
    }
    discussion, err := model.FindDiscussionById(id)
    if err == model.ErrDiscussionNotFound {
        return echo.ErrNotFound
    } else if err != nil {
        return echo.ErrInternalServerError
    }
    comments, err := model.FindCommentsByDiscId(id)
    if err != nil {
        return echo.ErrInternalServerError
    }
    type Comment struct {
        UserId      uint64  `json:"userid"`
        Content     string  `json:"content"`
        PostTime    int64   `json:"posttime`
    }
    var resp struct {
        Title       string      `json:"title"`
        Content     string      `json:"content"`
        PostTime    int64       `json:"posttime"`
        Comments    []Comment   `json:"comments"`
    }
    resp.Title = discussion.Title
    resp.Content = discussion.Content
    resp.PostTime = discussion.Time
    for _, comment := range comments {
        resp.Comments = append(resp.Comments, Comment{
            UserId: comment.UserId,
            Content: comment.Content,
            PostTime: comment.Time,
        })
    }
    return c.JSON(http.StatusOK, &resp)
}

func PostComment(c echo.Context) error {
    claims := c.Get("claims").(*utils.Claims)
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        return echo.ErrNotFound
    }
    req := new(param.PostCommentRequest)
    if err := context.BindAndVali(c, req); err != nil {
        return echo.ErrBadRequest
    }
    comment := model.Comment{}
    comment.DiscId = id
    comment.UserId = claims.UserId
    comment.Content = req.Content
    comment.Time = time.Now().Unix()
    if err := model.CreateComment(&comment); err != nil {
        return echo.ErrInternalServerError
    }
    var resp struct {
        Message string  `json:"message"`
    }
    resp.Message = "Success!"
    return c.JSON(http.StatusOK, &resp)
}