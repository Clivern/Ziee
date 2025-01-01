// Copyright 2026 Actx0. All rights reserved.
// License can be found in the LICENSE file.

package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	smithy "github.com/aws/smithy-go"
)

type S3Store struct {
	api    *awss3.Client
	config Config
}

// NewS3Store creates an S3-backed document store.
func NewS3Store(config Config) (*S3Store, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(
		context.Background(),
		awsconfig.WithRegion(config.S3.Region),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				config.S3.AccessKeyId,
				config.S3.SecretAccessKey,
				"",
			),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("s3 storage load config: %w", err)
	}

	client := awss3.NewFromConfig(awsCfg, func(options *awss3.Options) {
		options.BaseEndpoint = aws.String(config.S3.Endpoint)
		options.UsePathStyle = config.S3.UsePathStyle
	})

	return &S3Store{
		api:    client,
		config: config,
	}, nil
}

// Config returns the store configuration.
func (s *S3Store) Config() Config {
	return s.config
}

// PutDocument writes document bytes to storage.
func (s *S3Store) PutDocument(ctx context.Context, id string, text string) (string, error) {
	key := s.objectKey(id)

	_, err := s.api.PutObject(ctx, &awss3.PutObjectInput{
		Bucket:      aws.String(s.config.S3.Bucket),
		Key:         aws.String(key),
		Body:        strings.NewReader(text),
		ContentType: aws.String("text/plain; charset=utf-8"),
	})
	if err != nil {
		return "", fmt.Errorf("s3 storage put document: %w", err)
	}

	return s.documentURI(key)
}

// GetDocument reads document bytes from storage.
func (s *S3Store) GetDocument(ctx context.Context, id string) (string, error) {
	key := s.objectKey(id)

	output, err := s.api.GetObject(ctx, &awss3.GetObjectInput{
		Bucket: aws.String(s.config.S3.Bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		if isS3NotFound(err) {
			return "", ErrDocumentNotFound
		}
		return "", fmt.Errorf("s3 storage get document: %w", err)
	}
	defer output.Body.Close()

	data, err := io.ReadAll(output.Body)
	if err != nil {
		return "", fmt.Errorf("s3 storage read document: %w", err)
	}

	return string(data), nil
}

// DeleteDocument removes a document from storage.
func (s *S3Store) DeleteDocument(ctx context.Context, id string) error {
	key := s.objectKey(id)

	_, err := s.api.DeleteObject(ctx, &awss3.DeleteObjectInput{
		Bucket: aws.String(s.config.S3.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if isS3NotFound(err) {
			return ErrDocumentNotFound
		}
		return fmt.Errorf("s3 storage delete document: %w", err)
	}

	return nil
}

// DocumentURI returns the URI for a stored document.
func (s *S3Store) DocumentURI(id string) (string, error) {
	return s.documentURI(s.objectKey(id))
}

// objectKey returns the S3 object key for a document.
func (s *S3Store) objectKey(id string) string {
	return s.config.S3.Prefix + "/" + DocumentKey(id)
}

// documentURI builds the URI for a stored document.
func (s *S3Store) documentURI(key string) (string, error) {
	return fmt.Sprintf("s3://%s/%s", s.config.S3.Bucket, key), nil
}

// isS3NotFound reports whether an error is an S3 not-found response.
func isS3NotFound(err error) bool {
	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.ErrorCode() {
		case "NoSuchKey", "NotFound", "404":
			return true
		}
	}
	return false
}
