package httpsvc

import (
	"net/http"

	"github.com/mashingan/smapping"
	"github.com/miun173/autograd/model"

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

	assignment := &model.Assignment{}
	err = smapping.FillStruct(assignment, smapping.MapFields(assignmentReq))
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	err = s.assignmentUsecase.Create(c.Request().Context(), assignment)
	if err != nil {
		logrus.Error(err)
		return responseError(c, err)
	}

	return c.JSON(http.StatusOK, assignmentModelToResponse(assignment))
}
