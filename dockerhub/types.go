package dockerhub

import (
	"github.com/docker/distribution/manifest/schema1"
)

type (
	Error struct {
		Code    ErrorCode   `json:"code"`
		Message string      `json:"message"`
		Detail  interface{} `json:"detail"`
	}

	ErrorCode string

	// Response 通用的响应报文
	Response struct {
		Name                string `json:"name,omitempty"`
		*ResponseRepository `json:",inline"`
		*ResponseTag        `json:",inline"`
		*ResponseManifest   `json:",inline"`
		Errors              []Error `json:"errors,omitempty"`
	}

	// ResponseTag 查询tag接口响应报文
	ResponseTag struct {
		Tags []string `json:"tags,omitempty"`
	}

	ResponseManifest struct {
		Tag          string            `json:"tag,omitempty"`
		Architecture string            `json:"architecture,omitempty"`
		FSLayers     []schema1.FSLayer `json:"fsLayers,omitempty"`
		History      []schema1.History `json:"history,omitempty"`
	}

	// ResponseRepository 查询存储库响应报文
	ResponseRepository struct {
		Repository []string `json:"repositories,omitempty"`
	}

	ListTagsOption struct {
		Number int
	}
)

const (
	// https://docs.docker.com/registry/spec/api/#errors-2
	ErrBlobUnknown         ErrorCode = "BLOB_UNKNOWN"
	ErrBlobUploadInvalid   ErrorCode = "BLOB_UPLOAD_INVALID"
	ErrBlobUploadUnknown   ErrorCode = "BLOB_UPLOAD_UNKNOWN"
	ErrDigestInvalid       ErrorCode = "DIGEST_INVALID"
	ErrManifestBlobUnknown ErrorCode = "MANIFEST_BLOB_UNKNOWN"
	ErrManifestInvalid     ErrorCode = "MANIFEST_INVALID"
	ErrManifestUnknown     ErrorCode = "MANIFEST_UNKNOWN"
	ErrManifestUnverified  ErrorCode = "MANIFEST_UNVERIFIED"
	ErrNameInvalid         ErrorCode = "NAME_INVALID"
	ErrNameUnknown         ErrorCode = "NAME_UNKNOWN"
	ErrSizeInvalid         ErrorCode = "SIZE_INVALID"
	ErrTagInvalid          ErrorCode = "TAG_INVALID"
	ErrUnauthorized        ErrorCode = "UNAUTHORIZED"
	ErrDenied              ErrorCode = "UNAUTHORIZED"
	ErrSupported           ErrorCode = "UNSUPPORTED"
)
