package controller

import (
    "time"
    "strconv"
    "github.com/labstack/echo/v4"
    "github.com/golang-jwt/jwt/v4"

    "user-management-system/internal/config"
    "user-management-system/internal/model"
    "user-management-system/internal/utils"
)

func PostDiscussion(c echo.Context) error {
    claims := c.Get("claims").(*utils.Claims)
    var req struct {
        Title   string    `json:"title"`
        Content string    `json:"content"`
    }
    if err := c.Bind(&req); err != nil {
        return echo.ErrBadRequest
    }
    discussion := model.Discussion{}
    discussion.Title = req.Title
    discussion.Content = req.Content
    discussion.UserId = claims.User.ID
    discussion.Time = time.Now().Unix()
    if err := model.CreateDiscussion(discussion); err != nil {
        return echo.ErrInternalServerError
    }
}

func DiscussionInfo(c echo.Context) error {
    claims := c.Get("claims").(*utils.Claims)
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
        UserId      string  `json:"userid"`
        Content     string  `json:"content"`
        PostTime    string  `json:"posttime`
    }
    var resp struct {
        Title       string      `json:"title"`
        Content     string      `json:"content"`
        PostTime    time.Time   `json:"posttime"`
        Comments    []Comment   `json:"comments"`
    }
    resp.Title = discussion.Title
    resp.Content = discussion.Content
    resp.PostTime = discussion.Time
    for _, comment := range comments {
        resp.Comments = append(resp.Comments, Comment{
            UserId: comment.UserId,
            Content: comment.Content,
            PostTime: comment.PostTime,
        })
    }
    return c.JSON(200, &resp)
}

func PostComment(c echo.Context) error {
    claims := c.Get("claims").(*utils.Claims)
    id, strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        return echo.ErrNotFound
    }
    var req struct {
        Content    string    `json:"content"`
    }
    if err := c.Bind(&req); err != nil {
        return echo.ErrBadRequest
    }
    comment := model.Comment{}
    comment.DiscId = id
    comment.UserId = claims.User.ID
    comment.Content = req.Content
    comment.Time = time.Now().Unix()
    if err := model.CreateComment(comment); err != nil {
        return echo.ErrInternalServerError
    }
    var resp struct {
        Message string  `json:"message"`
    }
    resp.Message = "Success!"
    return c.JSON(200, &resp)
}