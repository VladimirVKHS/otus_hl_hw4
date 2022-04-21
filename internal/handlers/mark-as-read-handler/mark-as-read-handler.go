package mark_as_read_handler

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"io/ioutil"
	"net/http"
	"otus_dialog_go/internal/constants"
	httpHelper "otus_dialog_go/internal/helpers/http"
	"otus_dialog_go/internal/logger"
	"otus_dialog_go/internal/models/message"
)

func MarkAsReadHandler(w http.ResponseWriter, r *http.Request) {

	requestId := r.Context().Value(constants.RequestIDKey).(string)
	logger.Info(fmt.Sprintf("Mark as read. Request ID = %s", requestId))

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "BadRequest", http.StatusBadRequest)
		return
	}
	r.Body.Close()

	var request message.MarkAsReadRequest
	if err := json.Unmarshal(data, &request); err != nil {
		http.Error(w, "BadRequest", http.StatusBadRequest)
		return
	}

	validate := validator.New()
	validationErrs := validate.Struct(request)
	if validationErrs != nil {
		httpHelper.ValidationErrorResponse(w, validationErrs.Error())
		return
	}

	affected, err2 := request.MarkAsRead(r.Context())
	if err2 != nil {
		logger.Error(err.Error())
		httpHelper.InternalServerErrorResponse(w)
		return
	}

	httpHelper.JsonResponse(w, map[string]interface{}{
		"AffectedMessages": affected,
	})
}
