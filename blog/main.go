package main

import (
	"context"
	"log"
	"net"
	"project/blogpost/blog/blogs"

	"google.golang.org/grpc"
)

type blogdetail struct {
	blogs.UnimplementedBlogPostServer
}

func CreatePost(ctx context.Context, post *blogs.Post) (*blogs.Post, error) {
	return &blogs.Post{
		PostID:           "Post1",
		Title:            "Post 1",
		Content:          "Some content to test",
		Author:           "Aman",
		PublicationDate:  "03-02-2024",
		ModificationDate: "03-02-2024",
		Tags:             []string{"tag1", "tag2"},
	}, nil
}
func main() {
	listner, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Error in listner: %s", err)
	}
	registrar := grpc.NewServer()
	service := &blogdetail{}
	blogs.RegisterBlogPostServer(registrar, service)
	err = registrar.Serve(listner)
	if err != nil {
		log.Fatalf("Error in serve: %f", err)
	}
}
