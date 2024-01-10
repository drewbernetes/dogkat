/*
Copyright 2024 EscherCloud.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package s3

import (
	"fmt"
	"github.com/rhnvrm/simples3"
	"io"
)

type S3 struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	s3Conn    *simples3.S3
}

// PutToS3 Pushes a file to an S3 bucket.
func (s *S3) PutToS3(contentType, key, fileName string, body io.ReadSeeker) error {
	s3Conn := simples3.New("us-east-1", s.AccessKey, s.SecretKey)
	s3Conn.SetEndpoint(s.Endpoint)

	// Put the file into S3.
	input := simples3.UploadInput{
		Bucket:      s.Bucket,
		ObjectKey:   key,
		ContentType: contentType,
		FileName:    fileName,
		Body:        body,
	}
	_, err := s.upload(input)
	if err != nil {
		return fmt.Errorf("failed to push file to S3: %v\n", err)
	}

	return nil
}

// PutToS3 Pushes a file to an S3 bucket.
func (s *S3) upload(input simples3.UploadInput) (simples3.PutResponse, error) {
	return s.s3Conn.FilePut(input)
}
