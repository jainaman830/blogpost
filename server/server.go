package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"time"

	blog "project/blogpost/blog/blogs"

	"google.golang.org/grpc"
)

type blogServer struct {
	blog.UnimplementedBlogPostServer
	posts map[string]*blog.Post
}

func (s *blogServer) CreatePost(ctx context.Context, req *blog.Post) (*blog.Post, error) {
	postID, err := GenerateRandomID(16)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	post := &blog.Post{
		PostID:           postID,
		Title:            req.Title,
		Content:          req.Content,
		Author:           req.Author,
		PublicationDate:  time.Now().String(),
		ModificationDate: time.Now().String(),
		Tags:             req.Tags,
	}

	s.posts[postID] = post
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
func GenerateRandomID(length int) (string, error) {
	if length%2 != 0 {
		return "", fmt.Errorf("length must be even")
	}

	randomBytes := make([]byte, length/2)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(randomBytes), nil
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
