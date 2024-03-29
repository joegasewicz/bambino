package main

import (
	"encoding/json"
	"errors"
	"fmt"
	entityfileuploader "github.com/joegasewicz/entity-file-uploader"
	"github.com/joegasewicz/gomek"
	"log"
	"net/http"
)

type FileView struct{}
type UserView struct{}
type HealthView struct{}

var agnosticUploader *entityfileuploader.FileManager

// file
func (f *FileView) Get(w http.ResponseWriter, r *http.Request, d *gomek.Data) {
	id, err := gomek.GetParams(r, "id")
	if err != nil {
		return
		gomek.JSON(w, nil, http.StatusBadRequest)
	}
	var fileModel FileModel
	DB.First(&fileModel, "id = ?", id[0])
	agnosticUploader = NewFilesManager(fileModel.EntityName, "")
	fileURL := agnosticUploader.Get(fileModel.FileName, fileModel.ID)
	data := struct{ URL string }{URL: fileURL}
	gomek.JSON(w, data, http.StatusOK)
}

func (f *FileView) Post(w http.ResponseWriter, r *http.Request, d *gomek.Data) {
	var options OptionSchema
	var fileResp FileRespSchema

	optionsStr, err := gomek.GetParams(r, "options")
	fileRespSlices := make([]FileRespSchema, 0)
	if err != nil {
		log.Println(err.Error())
		gomek.JSON(w, nil, http.StatusBadRequest)
		return
	}
	if optionsStr == nil {
		err := errors.New("no options send with request")
		log.Println(err.Error())
		gomek.JSON(w, nil, http.StatusBadRequest)
		return
	}
	err = json.Unmarshal([]byte(optionsStr[0]), &options)
	if err != nil {
		log.Println(err.Error())
		gomek.JSON(w, nil, http.StatusBadRequest)
		return
	}
	fileManager := NewFilesManager(options.EntityName, "")

	// TODO https://github.com/joegasewicz/bambino/issues/8
	for i, optionsFileName := range options.Files {
		fileName := optionsFileName
		if err != nil {
			log.Println(err.Error())
			gomek.JSON(w, nil, http.StatusBadRequest)
			return
		}
		data, err := json.Marshal(options.Data)
		if err != nil {
			log.Println(err.Error())
			gomek.JSON(w, nil, http.StatusBadRequest)
			return
		}
		fileModel := FileModel{
			Name:       optionsFileName,
			FileName:   fileName,
			Data:       string(data),
			EntityName: options.EntityName,
		}
		result := DB.Create(&fileModel)
		if result.Error != nil {
			log.Printf("error saving file %s\n", result.Error.Error())
			gomek.JSON(w, nil, http.StatusInternalServerError)
			return
		}
		if result.RowsAffected == 0 {
			log.Printf("unable to save file with name: %s", fileName)
		}
		_, err = fileManager.Upload(w, r, fileModel.ID, optionsFileName)
		if err != nil {
			// Handle file uploads over http

			err = fileManager.ReceiveMultiPartFormDataAndSaveToDir(r, options.Fields[i], fileModel.ID)
			if err != nil {
				log.Println(err)
				log.Printf("unable to store file on server with name: %s", fileName)
				gomek.JSON(w, nil, http.StatusInternalServerError)
				return
			}
		}

		path := fmt.Sprintf("%d/%s", fileModel.ID, optionsFileName)
		url := fmt.Sprintf("%s/%s/%s/%s", AppConfig.GetUrl(), fileManager.UploadDir, options.EntityName, path)
		fileResp = FileRespSchema{
			ID:         fileModel.ID,
			FileName:   fileName,
			Name:       optionsFileName,
			Data:       options.Data,
			EntityName: fileModel.EntityName,
			Url:        url,
			Path:       fmt.Sprintf("%s/%s/%s", fileManager.UploadDir, options.EntityName, path),
			CreatedOn:  fileModel.CreatedAt.String(),
		}
		fileRespSlices = append(fileRespSlices, fileResp)
	}
	gomek.JSON(w, fileRespSlices, http.StatusOK)
}

func (f *FileView) Put(w http.ResponseWriter, r *http.Request, d *gomek.Data) {
	panic("implement me")
}
func (f *FileView) Delete(w http.ResponseWriter, r *http.Request, d *gomek.Data) {
	panic("implement me")
}

// health
func (h *HealthView) Get(w http.ResponseWriter, r *http.Request, d *gomek.Data) {
	templateData := make(gomek.Data)
	templateData["health"] = map[string]string{
		"status": "OK",
	}
	gomek.JSON(w, templateData, http.StatusOK)
}
func (h *HealthView) Post(w http.ResponseWriter, r *http.Request, d *gomek.Data) {
	panic("implement me")
}
func (h *HealthView) Put(w http.ResponseWriter, r *http.Request, d *gomek.Data) {
	panic("implement me")
}
func (h *HealthView) Delete(w http.ResponseWriter, r *http.Request, d *gomek.Data) {
	panic("implement me")
}

// user
func (u *UserView) Get(w http.ResponseWriter, r *http.Request, d *gomek.Data) {
	panic("implement me")
}
func (u *UserView) Post(w http.ResponseWriter, r *http.Request, d *gomek.Data) {
	panic("implement me")
}
func (u *UserView) Put(w http.ResponseWriter, r *http.Request, d *gomek.Data) {
	panic("implement me")
}
func (u *UserView) Delete(w http.ResponseWriter, r *http.Request, d *gomek.Data) {
	panic("implement me")
}
