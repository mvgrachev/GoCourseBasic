package server

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"GoCourseBasic/homework-4/models"

	"github.com/go-chi/chi"
)

// getAllPosts - возвращает все посты
func (serv *Server) getAllPosts(w http.ResponseWriter, r *http.Request) {
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
		serv.SendInternalErr(w, err)
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	templ, err := template.New("Page").Parse(string(data))
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	posts, err := models.GetAllPostItems(serv.db)
	if err != nil {
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

	post, err := models.GetPost(serv.db, postId)
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
	data, _ := ioutil.ReadAll(r.Body)

	post := models.PostItem{}
	_ = json.Unmarshal(data, &post)

	if err := post.Insert(serv.db); err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	data, _ = json.Marshal(post)
	w.Write(data)
}

// deletePostHandler - удаляет пост
func (serv *Server) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	postId := chi.URLParam(r, "id")

	post := models.PostItem{Id: postId}
	if err := post.Delete(serv.db); err != nil {
		serv.SendInternalErr(w, err)
		return
	}
}

// putPostHandler - обновляет задачу
func (serv *Server) putPostHandler(w http.ResponseWriter, r *http.Request) {
	postId := chi.URLParam(r, "id")

	data, _ := ioutil.ReadAll(r.Body)

	post := models.PostItem{}
	_ = json.Unmarshal(data, &post)
	post.Id = postId

	if err := post.Update(serv.db); err != nil {
		serv.SendInternalErr(w, err)
		return
	}
}
