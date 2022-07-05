package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/HackerDom/letoctf-taskbot-2022-tasks/filestore/pkg/storage"
)

const maxFileSize = 100

func (s *server) uploadFile(ctx context.Context, _ http.ResponseWriter, r *http.Request) (interface{}, *errApi) {
	if r.Method != http.MethodPut {
		return nil, newApiErr(http.StatusMethodNotAllowed, "method not allowed")
	}

	userId, err := getUserIdFrom(ctx)
	if err != nil {
		return nil, internalErr("get user id failed", err)
	}
	if userId == uuid.Nil {
		return nil, newApiErr(http.StatusUnauthorized, "unauthorized")
	}

	file, fileHandler, err := r.FormFile("file")
	if err != nil {
		return nil, newApiErr(http.StatusBadRequest, err.Error())
	}

	if fileHandler.Size > maxFileSize {
		return nil, newApiErr(http.StatusBadRequest, fmt.Sprintf("file must be less than %v bytes", maxFileSize))
	}

	// для того чтобы нельзя было перезаписать файл с флагом
	_, err = s.filestore.Stat(fileHandler.Filename)
	switch {
	case errors.Is(err, storage.ErrNotExist):
	case err == nil:
		return nil, newApiErr(http.StatusBadRequest, "file already exist")
	default:
		return nil, internalErr("get file stat failed", err)
	}

	var encrypted bool
	encryptedValue := r.Form.Get("encrypted")
	switch encryptedValue {
	case "true":
		encrypted = true
	case "false":
		encrypted = false
	default:
		return nil, newApiErr(http.StatusBadRequest, "form field \"encrypted\" must be true or false")
	}

	opts := &storage.FileOptions{
		Name:      fileHandler.Filename,
		Size:      fileHandler.Size,
		OwnerId:   userId,
		Encrypted: encrypted,
	}
	if err := s.filestore.Update(ctx, file, opts); err != nil {
		return nil, internalErr("update file failed", err)
	}

	return nil, nil
}

func (s *server) listFiles(_ context.Context, _ http.ResponseWriter, r *http.Request) (interface{}, *errApi) {
	if r.Method != http.MethodGet {
		return nil, newApiErr(http.StatusMethodNotAllowed, "method not allowed")
	}

	filesInfo := s.filestore.List()
	return filesInfo, nil
}

func (s *server) getFileOwner(_ context.Context, _ http.ResponseWriter, r *http.Request) (interface{}, *errApi) {
	if r.Method != http.MethodGet {
		return nil, newApiErr(http.StatusMethodNotAllowed, "method not allowed")
	}

	filename := r.URL.Query().Get("filename")
	if filename == "" {
		return nil, newApiErr(http.StatusBadRequest, "filename query param required")
	}

	stat, err := s.filestore.Stat(filename)
	if err != nil {
		return nil, internalErr("get owner failed", err)
	}

	return stat.OwnerId, nil
}

func (s *server) getFile(ctx context.Context, _ http.ResponseWriter, r *http.Request) (interface{}, *errApi) {
	if r.Method != http.MethodGet {
		return nil, newApiErr(http.StatusMethodNotAllowed, "method not allowed")
	}

	filename := r.URL.Query().Get("filename")
	if filename == "" {
		return nil, newApiErr(http.StatusBadRequest, "filename query param required")
	}

	stat, err := s.filestore.Stat(filename)
	if err != nil {
		return nil, internalErr("get owner id failed", err)
	}

	file, err := s.filestore.Get(ctx, filename)
	if err != nil {
		return nil, internalErr("get file failed", err)
	}

	userId, err := getUserIdFrom(ctx)
	if err != nil {
		return nil, internalErr("get user id failed", err)
	}

	if userId != uuid.Nil && stat.Encrypted && stat.OwnerId == userId {
		key, err := stat.OwnerId.MarshalBinary()
		if err != nil {
			return nil, internalErr("internal error", err)
		}
		file, err = decrypt(key, file)
		if err != nil {
			return nil, internalErr("decryption failed", err)
		}
	}

	return file, nil
}
