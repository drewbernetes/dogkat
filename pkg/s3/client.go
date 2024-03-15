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
	s3 "github.com/drewbernetes/simple-s3"
	"os"
)

type S3 struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
}

// PutToS3 Pushes a file to an S3 bucket.
func (s *S3) PutToS3(key string, body *os.File) error {
	s3Conn, err := s3.New(s.Endpoint, s.AccessKey, s.SecretKey, s.Bucket, "us-east-1")
	if err != nil {
		return err
	}

	// Put the file into S3.
	err = s3Conn.Put(key, body)
	if err != nil {
		return fmt.Errorf("failed to push file to S3: %v\n", err)
	}

	return nil
}
