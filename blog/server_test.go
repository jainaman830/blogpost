package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"testing"
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
func setupTestServer() *blogServer {
	return &blogServer{
		posts: make(map[string]*blog.Post),
	}
}

func TestBlogService(t *testing.T) {
	t.Run("CreatePost", testCreatePost)
	t.Run("ReadPost", testReadPost)
	t.Run("UpdatePost", testUpdatePost)
	t.Run("DeletePost", testDeletePost)
}

func testCreatePost(t *testing.T) {
	s := setupTestServer()

	req := &blog.Post{
		Title:   "Testing Post",
		Content: "Some Testing content.",
		Author:  "Test Author",
		Tags:    []string{"tag 1", "tag 2"},
	}

	createdPost, err := s.CreatePost(context.Background(), req)
	if err != nil {
		t.Fatalf("CreatePost failed: %v", err)
	}

	if createdPost.PostID == "" {
		t.Errorf("Expected PostID to be set, got empty string")
	}
}

func testReadPost(t *testing.T) {
	s := setupTestServer()
	req := &blog.Post{
		Title:   "Testing Post 2",
		Content: "Some Testing content 2.",
		Author:  "Test Author 2",
		Tags:    []string{"tag 1", "tag 2"},
	}

	createdPost, err := s.CreatePost(context.Background(), req)
	if err != nil {
		t.Fatalf("CreatePost failed: %v", err)
	}
	readReq := &blog.ReadPostRequest{PostID: createdPost.PostID}
	readPost, err := s.ReadPost(context.Background(), readReq)
	if err != nil {
		t.Fatalf("ReadPost failed: %v", err)
	}

	if readPost == nil || readPost.PostID != createdPost.PostID {
		t.Errorf("Expected to read the created post, got different or nil post")
	}
}

func testUpdatePost(t *testing.T) {
	s := setupTestServer()

	// Create a dummy post for testing
	req := &blog.Post{
		Title:   "Testing Post 3",
		Content: "Some Testing content 3.",
		Author:  "Test Author 3",
		Tags:    []string{"tag 1", "tag 2"},
	}

	createdPost, err := s.CreatePost(context.Background(), req)
	if err != nil {
		t.Fatalf("CreatePost failed: %v", err)
	}

	// Update the created post
	updateReq := &blog.UpdatePostRequest{
		PostID: createdPost.PostID,
		Title:  "Updated Post Title",
		Tags:   []string{"tag3"},
	}

	updatedPost, err := s.UpdatePost(context.Background(), updateReq)
	if err != nil {
		t.Fatalf("UpdatePost failed: %v", err)
	}

	if updatedPost.Title != updateReq.Title || updatedPost.Tags[0] != updateReq.Tags[0] {
		t.Errorf("Expected post to be updated, got different post")
	}
}

func testDeletePost(t *testing.T) {
	s := setupTestServer()
	req := &blog.Post{
		Title:   "Testing Post 4",
		Content: "Some Testing content 4.",
		Author:  "Test Author 4",
		Tags:    []string{"tag 1", "tag 2"},
	}

	createdPost, err := s.CreatePost(context.Background(), req)
	if err != nil {
		t.Fatalf("CreatePost failed: %v", err)
	}

	// Delete the created post
	deleteReq := &blog.DeletePostRequest{PostID: createdPost.PostID}
	deleteResponse, err := s.DeletePost(context.Background(), deleteReq)
	if err != nil {
		t.Fatalf("DeletePost failed: %v", err)
	}

	if !deleteResponse.Success {
		t.Errorf("Expected post to be deleted successfully, got failure response")
	}

	// Ensure the post is not present after deletion
	readReq := &blog.ReadPostRequest{PostID: createdPost.PostID}
	readPost, err := s.ReadPost(context.Background(), readReq)
	if readPost != nil {
		t.Errorf("Expected post to be deleted, but it is still present")
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	blogServer := setupTestServer()
	blog.RegisterBlogPostServer(server, blogServer)
	fmt.Println("server is running on localhost:8080")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
