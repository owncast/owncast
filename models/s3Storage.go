package models

type S3 struct {
	Enabled         bool   `json:"enabled"`
	Endpoint        string `json:"endpoint,omitempty"`
	ServingEndpoint string `json:"servingEndpoint,omitempty"`
	AccessKey       string `json:"accessKey,omitempty"`
	Secret          string `json:"secret,omitempty"`
	Bucket          string `json:"bucket,omitempty"`
	Region          string `json:"region,omitempty"`
	ACL             string `json:"acl,omitempty"`
}
