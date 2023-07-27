package models

import (
	"bytes"
	"io"
	"plex_monitor/internal/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// createBucket creates a new GridFS bucket for storing raw request wires.
func createBucket(bucketName string) (*gridfs.Bucket, error) {
	opts := options.GridFSBucket().SetName(bucketName)
	bucket, err := gridfs.NewBucket(database.DB, opts)
	if err != nil {
		return nil, err
	}

	return bucket, nil
}

// AddFileToBucket adds a file to the GridFS bucket.
func AddFileToBucket(bucketName string, filename string, file []byte, metadata bson.M) (*primitive.ObjectID, error) {
	// Create a new bucket or get the existing one
	bucket, err := createBucket(bucketName)
	if err != nil {
		return nil, err
	}

	// Add metadata to the file
	uploadOpts := options.GridFSUpload().SetMetadata(metadata)

	// Upload the file to the bucket
	objectID, err := bucket.UploadFromStream(filename, io.NopCloser(bytes.NewReader(file)), uploadOpts)
	if err != nil {
		return nil, err
	}

	return &objectID, nil
}

// CountFilesInBucket counts the number of files in the GridFS bucket.
func CountFilesInBucket(bucketName string, query bson.M) (int64, error) {
	// Create a new bucket or get the existing one
	bucket, err := createBucket(bucketName)
	if err != nil {
		return 0, err
	}

	// Count the number of files in the bucket
	count, err := bucket.GetFilesCollection().CountDocuments(database.Ctx, query)
	if err != nil {
		return 0, err
	}

	return count, nil
}
