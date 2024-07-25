package handlers

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BlogHandler struct {
	DB *mongo.Database
}

type Blog struct {
	ID      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title   string             `json:"title" bson:"title"`
	Body    string             `json:"body" bson:"body"`
	Author  string             `json:"author" bson:"author"`
	ImageID string             `json:"image_id" bson:"image_id"`
}

func (h *BlogHandler) CreateBlog(c *gin.Context) {
	var blog Blog
	if err := c.ShouldBindJSON(&blog); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	_, err := h.DB.Collection("blogs").InsertOne(c.Request.Context(), blog)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Blog Created Succesfully"})
}

func (h *BlogHandler) UpdateBlog(c *gin.Context) {
	var blog Blog
	if err := c.ShouldBindJSON(&blog); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	id := c.Param("id")
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	_, err = h.DB.Collection("blogs").UpdateOne(c.Request.Context(), bson.M{"_id": objId}, bson.M{"$set": blog})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Blog Updated Successfully"})
}

func (h *BlogHandler) GetBlogs(c *gin.Context) {
	var blogs []Blog
	cursor, err := h.DB.Collection("blogs").Find(c.Request.Context(), bson.M{})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var b Blog
		cursor.Decode(&b)
		blogs = append(blogs, b)
	}
	if err := cursor.Err(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, blogs)
}

func (h *BlogHandler) GetBlog(c *gin.Context) {
	id := c.Param("id")
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}
	var blog Blog
	err = h.DB.Collection("blogs").FindOne(c.Request.Context(), bson.M{"_id": objId}).Decode(&blog)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, blog)
}

func (h *BlogHandler) GetUserBlog(c *gin.Context) {
	var request struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var blogs []Blog
	cursor, err := h.DB.Collection("blogs").Find(c.Request.Context(), bson.M{"author": request.Email})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	for cursor.Next(context.Background()) {
		var b Blog
		cursor.Decode(&b)
		blogs = append(blogs, b)
	}
	if err := cursor.Err(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, blogs)
}

func (h *BlogHandler) DeleteBlog(c *gin.Context) {
	id := c.Param("id")
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	_, err = h.DB.Collection("blogs").DeleteOne(c.Request.Context(), bson.M{"_id": objId})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Blog Deleted Successfully"})
}
