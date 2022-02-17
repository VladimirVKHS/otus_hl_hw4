package create_message_handler

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"io/ioutil"
	"net/http"
	httpHelper "otus_dialog_go/internal/helpers/http"
	"otus_dialog_go/internal/logger"
	"otus_dialog_go/internal/models/message"
)

func CreateMessageHandler(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "BadRequest", http.StatusBadRequest)
		return
	}
	r.Body.Close()

	var request message.MessageCreateRequest
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

	createdMessage, err := request.CreateMessage(r.Context())

	if err != nil {
		logger.Error(err.Error())
		httpHelper.InternalServerErrorResponse(w)
		return
	}

	httpHelper.JsonResponse(w, createdMessage.ToResponse())
}
