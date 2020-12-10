package httpsvc

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (s *Server) handleCreateAssignment(c echo.Context) error {
	assignmentReq := &assignmentRequest{}
	err := c.Bind(assignmentReq)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	assignment := assignmentRequestToModel(assignmentReq)
	err = s.assignmentUsecase.Create(c.Request().Context(), assignment)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	return c.JSON(http.StatusOK, assignmentModelToResponse(assignment))
}

func (s *Server) handleUpdateAssignment(c echo.Context) error {
	assignmentReq := &assignmentUpdateRequest{}
	err := c.Bind(assignmentReq)
	if err != nil {

		return responseError(c, err)
	}

	assignment := assigmentUpdateReqToModel(assignmentReq)
	err = s.assignmentUsecase.Update(c.Request().Context(), assignment)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	return c.JSON(http.StatusOK, assignmentModelToResponse(assignment))

}
