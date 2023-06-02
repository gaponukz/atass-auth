package controller

import (
	"auth/src/entities"
	"auth/src/storage"
	"net/http"
)

type Controller struct {
	Storage storage.IStorage[entities.User]
}

func (contr *Controller) Singin(responseWriter http.ResponseWriter, request *http.Request) {}

func (contr *Controller) Refresh(responseWriter http.ResponseWriter, request *http.Request) {}

func (contr *Controller) Logout(responseWriter http.ResponseWriter, request *http.Request) {}

func (contr *Controller) Welcome(responseWriter http.ResponseWriter, request *http.Request) {}
