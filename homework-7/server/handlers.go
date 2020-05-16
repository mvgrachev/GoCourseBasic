package server

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"GoCourseBasic/homework-7/models"
	"github.com/go-chi/chi"
	"strings"
)

// getAllPosts - возвращает все посты
func (serv *Server) getAllPosts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	templateName := chi.URLParam(r, "template")

	if templateName == "" {
		templateName = serv.indexTemplate
	}

	file, err := os.Open(path.Join(serv.templatesDir, templateName))
	if err != nil {
		if err == os.ErrNotExist {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		serv.lg.Debug("00000000000000")
		serv.SendInternalErr(w, err)
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		serv.lg.Debug("222222222")
		serv.SendInternalErr(w, err)
		return
	}

	templ, err := template.New("Page").Parse(string(data))
	if err != nil {
		serv.lg.Debug("33333333")
		serv.SendInternalErr(w, err)
		return
	}

	posts, err := models.GetAllPostItems(ctx, serv.db)
	if err != nil {
		serv.lg.Debug("4444444444")
		serv.lg.Debug(serv.db)
		serv.SendInternalErr(w, err)
		return
	}

	serv.Page.Posts = posts

	if err := templ.Execute(w, serv.Page.Posts); err != nil {
		serv.SendInternalErr(w, err)
		return
	}
}


// getPost - возвращает шаблон
func (serv *Server) getPost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	postId := chi.URLParam(r, "Id")

	if postId == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	file, err := os.Open(path.Join(serv.templatesDir, "post.html"))
	if err != nil {
		if err == os.ErrNotExist {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		serv.SendInternalErr(w, err)
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	templ, err := template.New("Post").Parse(string(data))
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	post, err := models.GetPost(ctx, serv.db, postId)
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	if err := templ.Execute(w, post); err != nil {
		serv.SendInternalErr(w, err)
		return
	}
}

// postPostHandler - добавляет новый post
func (serv *Server) postPostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data, _ := ioutil.ReadAll(r.Body)

	post := &models.Post{}
	err := json.Unmarshal(data, &post)
	var body []string
	for _, value := range post.Body.([]interface{}) {
		body = append(body, value.(string))
	}
	post.Body = strings.Join(body, "\n")

	if err = post.Create(ctx, serv.db); err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	data, _ = json.Marshal(post)
	w.Write(data)
}

// deletePostHandler - удаляет пост
func (serv *Server) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	postId := chi.URLParam(r, "id")

	post := models.Post{}
	post.Id = postId
	if err := post.Delete(ctx, serv.db); err != nil {
		serv.SendInternalErr(w, err)
		return
	}
}

// putPostHandler - обновляет пост
func (serv *Server) putPostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	postId := chi.URLParam(r, "id")
	//postId = postId.(string)
	data, _ := ioutil.ReadAll(r.Body)

	//post, _ := models.GetPost(ctx, serv.db, postId)
	post := models.Post{}
	post.Id = postId
	err := json.Unmarshal(data, &post)
	if err != nil {
		panic(err)
	}
	var body []string
	for _, value := range post.Body.([]interface{}) {
		body = append(body, value.(string))
	}
	post.Body = strings.Join(body, "\n")

	if err := post.Update(ctx, serv.db); err != nil {
		serv.SendInternalErr(w, err)
		return
	}
}
