package main

import (
	"net/http"

	"cloud.google.com/go/datastore"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

type Post struct {
	Body string
}

func createPost(body string) error {
	ctx := appengine.NewContext()

	post := Post{
		Body: body,
	}

	key := datastore.NewIncompleteKey(ctx, "Post", nil)
	_, err := datastore.Put(ctx, key, &post)
	if err != nil {
		return err
	}

	return nil
}

func listPosts() ([]Post, error) {
	ctx := appengine.NewContext()

	var posts []Post
	query := datastore.NewQuery("Post")
	_, err := query.GetAll(ctx, &posts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return posts, nil
	}

	return posts, nil
}
