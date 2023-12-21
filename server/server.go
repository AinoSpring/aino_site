package server

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"aino-spring.com/aino_site/config"
	"aino-spring.com/aino_site/database"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Server struct {
  Router *gin.Engine
  Config *config.Config
  Database *database.Connection
}

func NewServer(db *database.Connection, conf *config.Config) *Server {
  server := new(Server)
  server.Database = db
  server.Config = conf

  server.Router = gin.Default()

  store := cookie.NewStore([]byte(conf.SessionSecret))
  server.Router.Use(sessions.Sessions("ainosession", store))

  server.Router.SetFuncMap(template.FuncMap{
    "dict": func(values ...interface{}) (map[string]interface{}, error) {
      if len(values) % 2 != 0 {
        return nil, errors.New("Invalid dict call")
      }
      dict := make(map[string]interface{}, len(values) / 2)
      for i := 0; i < len(values); i += 2 {
        key, ok := values[i].(string)
        if !ok {
          return nil, errors.New("Dict keys must be strings")
        }
        dict[key] = values[i + 1]
      }
      return dict, nil
    },
  })

  server.Router.LoadHTMLGlob("templates/**/*")

  pages, err := db.FetchPages()
  if err != nil {
    log.Panic(err)
  }
  pager := NewPagerFromDBPages(pages)
  server.LoadPager(pager)

  server.Router.Static("/static", "static")
  server.Router.StaticFile("/favicon.ico", "favicon.ico")
  server.Router.NoRoute(server.GetHandler(http.StatusNotFound, "not-found", gin.H{}))

  return server
}

func (server *Server) SetupManualPages() {
  server.Router.GET("/posts", func (c *gin.Context) {
    posts, err := server.Database.FetchPosts()
    if err != nil {
      log.Panic(err)
    }
    sort.Slice(posts, func (i, j int) bool {
      return posts[i].Date.Unix() > posts[j].Date.Unix()
    })
    c.HTML(http.StatusOK, "posts", server.GetValues("posts", c, gin.H{"posts": posts}))
  })

  server.Router.GET("/posts/:id", func (c *gin.Context) {
    id := c.Param("id")
    post, err := server.Database.FetchPost(id)
    if err != nil {
      c.Redirect(http.StatusTemporaryRedirect, "/posts")
      return
    }
    postMap := gin.H{"id": id, "title": post.Title, "date": post.Date.Format("January 02, 2006"), "abstract": post.Abstract, "contents": template.HTML(post.Contents), "public": post.Public}
    c.HTML(http.StatusOK, "post", server.GetValues("post", c, gin.H{"post": postMap}))
  })

  server.Router.GET("/posts/:id/edit", func (c *gin.Context) {
    id := c.Param("id")
    if !server.IsAuthed(c) {
      c.Redirect(http.StatusTemporaryRedirect, "/posts/" + id)
      return
    }
    post, err := server.Database.FetchPost(id)
    if err != nil {
      c.Redirect(http.StatusTemporaryRedirect, "/posts")
      return
    }
    postMap := gin.H{"id": id, "title": post.Title, "date": post.Date.Format("January 02, 2006"), "abstract": post.Abstract, "contents": template.HTML(post.Contents), "public": post.Public}
    c.HTML(http.StatusOK, "edit-post", server.GetValues("edit-post", c, gin.H{"post": postMap}))
  })

  server.Router.GET("/new-post", func (c *gin.Context) {
    if !server.IsAuthed(c) {
      c.Redirect(http.StatusTemporaryRedirect, "/posts")
      return
    }
    c.HTML(http.StatusOK, "new-post", server.GetValues("new-post", c, gin.H{"date": time.Now().Format("January 02, 2006")}))
  })

  server.Router.GET("/logout", func (c *gin.Context) {
    redirect := c.DefaultQuery("redirect", "/")
    server.Logout(c)
    c.Redirect(http.StatusTemporaryRedirect, redirect)
  })
}

func (server *Server) IsAuthed(c *gin.Context) bool {
  queryPassword := c.Query("password")
  session := sessions.Default(c)
  sessionPassword := session.Get("password")
  if queryPassword == "" && sessionPassword == nil {
    return false
  }
  if queryPassword != "" {
    isAuthed := queryPassword == server.Config.AdminPassword
    session.Set("password", queryPassword)
    session.Save()
    return isAuthed
  }
  return sessionPassword == server.Config.AdminPassword
}

func (server *Server) Logout(c *gin.Context) {
  session := sessions.Default(c)
  session.Clear()

  session.Options(sessions.Options{MaxAge: -1})
  session.Save()
}

func (server *Server) GetValues(template string, c *gin.Context, values gin.H) gin.H {
  caser := cases.Title(language.English)
  title := strings.ReplaceAll(template, "-", " ")
  title = strings.ReplaceAll(title, "_", " ")
  title = caser.String(title)
  values["title"] = title
  values["template"] = template
  values["authed"] = server.IsAuthed(c)
  return values
}

func (server *Server) GetHandler(status int, template string, values gin.H) func (*gin.Context) {
  return func (c *gin.Context) {
    c.HTML(http.StatusOK, template, server.GetValues(template, c, values))
  }
}

func (server *Server) GetAdminHandler(status int, template string, values gin.H) func (*gin.Context) {
  return func (c *gin.Context) {
    if server.IsAuthed(c) {
      c.HTML(http.StatusOK, template, server.GetValues(template, c, values))
      return
    }
    c.Redirect(http.StatusTemporaryRedirect, "/login")
  }
}

func (server *Server) LoadPager(pager *Pager) {
  for _, path := range pager.GetPaths() {
    var handler func (*gin.Context)
    if pager.IsAdmin(path) {
      handler = server.GetAdminHandler(http.StatusOK, pager.GetTemplate(path), gin.H{})
    } else {
      handler = server.GetHandler(http.StatusOK, pager.GetTemplate(path), gin.H{})
    }
    server.Router.GET(path, handler)
  }
}

func (server *Server) Run(address string) {
  server.Router.Run(address)
}

