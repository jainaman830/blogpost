package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	blog "project/blogpost/blog/blogs"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type blogServer struct {
	blog.UnimplementedBlogPostServer
	posts map[string]*blog.Post
}

func (s *blogServer) CreatePost(ctx context.Context, req *blog.Post) (*blog.Post, error) {
	post := &blog.Post{
		PostID:           uuid.New().String(),
		Title:            req.Title,
		Content:          req.Content,
		Author:           req.Author,
		PublicationDate:  time.Now().String(),
		ModificationDate: time.Now().String(),
		Tags:             req.Tags,
	}

	s.posts[post.PostID] = post
	return post, nil
}

func (s *blogServer) ReadPost(ctx context.Context, req *blog.ReadPostRequest) (*blog.Post, error) {
	post, ok := s.posts[req.PostID]
	if !ok {
		return nil, fmt.Errorf("Post not found")
	}
	return post, nil
}

func (s *blogServer) UpdatePost(ctx context.Context, request *blog.UpdatePostRequest) (*blog.Post, error) {
	post, ok := s.posts[request.PostID]
	if !ok {
		return nil, fmt.Errorf("Post not found")
	}

	post.Title = request.Title
	post.Content = request.Content
	post.Author = request.Author
	post.Tags = request.Tags
	post.ModificationDate = time.Now().String()
	s.posts[request.PostID] = post
	return post, nil
}

func (s *blogServer) DeletePost(ctx context.Context, request *blog.DeletePostRequest) (*blog.DeletePostResponse, error) {
	if _, ok := s.posts[request.PostID]; !ok {
		return &blog.DeletePostResponse{Success: false, Message: "Post not found"}, nil
	}
	delete(s.posts, request.PostID)
	return &blog.DeletePostResponse{Success: true, Message: "Post deleted successfully"}, nil

}
func (s *blogServer) ReadAllPosts(ctx context.Context, req *blog.ReadAllPost) (*blog.AllPosts, error) {
	var allPosts []*blog.Post
	fmt.Println("s.posts", s.posts)
	for _, post := range s.posts {
		allPosts = append(allPosts, post)
	}
	output := &blog.AllPosts{
		Allpost: allPosts,
	}
	return output, nil
}
func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	blog.RegisterBlogPostServer(server, &blogServer{posts: make(map[string]*blog.Post)})
	fmt.Println("server is running on localhost:8080")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
