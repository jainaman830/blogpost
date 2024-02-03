// client.go
package main

import (
	"context"
	"fmt"
	"log"

	blog "project/blogpost/blog/blogs"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	client := blog.NewBlogPostClient(conn)
	post := &blog.Post{
		Title:   "Post title",
		Content: "Some Dummy Content.",
		Author:  "Aman Jain",
		Tags:    []string{"tag 1", "tag2"},
	}

	createdPost, err := client.CreatePost(context.Background(), post)
	if err != nil {
		log.Fatalf("Failed to create post: %v", err)
	}
	fmt.Println("Post created: ", createdPost)

	readPost := &blog.ReadPostRequest{
		PostID: createdPost.PostID,
	}

	outputPost, err := client.ReadPost(context.Background(), readPost)
	if err != nil {
		log.Fatalf("Failed to read post: %v", err)
	}

	fmt.Println("Post read output: ", outputPost)
	updatedPost := &blog.UpdatePostRequest{
		PostID:  outputPost.PostID,
		Title:   "New title",
		Content: "Updated content",
		Author:  outputPost.Author,
		Tags:    []string{"tag 1", "tag 2", "tag 3"},
	}

	postAfterUpdate, err := client.UpdatePost(context.Background(), updatedPost)
	if err != nil {
		log.Fatalf("Failed to update post: %v", err)
	}
	fmt.Println("Post after update: ", postAfterUpdate)

	deletePost := &blog.DeletePostRequest{
		PostID: postAfterUpdate.PostID,
	}
	deletePostResponse, err := client.DeletePost(context.Background(), deletePost)
	if err != nil {
		log.Fatalf("Failed to delete post: %v", err)
	}
	if deletePostResponse.Success {
		fmt.Println(deletePostResponse.Message)
	} else {
		fmt.Println("Error in delete: ", deletePostResponse.Message)
	}

	outputPost, err = client.ReadPost(context.Background(), readPost)
	if err != nil {
		log.Fatalf("Failed to read post: %v", err)
	}

	fmt.Println("Post read output: ", outputPost)
}
