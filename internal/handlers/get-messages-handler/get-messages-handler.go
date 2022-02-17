package get_messages_handler

import (
	"net/http"
	"otus_dialog_go/internal/helpers/chat"
	httpHelper "otus_dialog_go/internal/helpers/http"
	"otus_dialog_go/internal/logger"
	"otus_dialog_go/internal/models/message"
	"strconv"
)

func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	user1Id, _ := strconv.Atoi(r.URL.Query().Get("user1_id"))
	if user1Id == 0 {
		httpHelper.ValidationErrorResponse(w, "user1_id not provided")
		return
	}
	if user1Id < 0 {
		httpHelper.ValidationErrorResponse(w, "user1_id is invalid")
		return
	}
	user2Id, _ := strconv.Atoi(r.URL.Query().Get("user2_id"))
	if user2Id == 0 {
		httpHelper.ValidationErrorResponse(w, "user2_id not provided")
		return
	}
	if user2Id < 0 {
		httpHelper.ValidationErrorResponse(w, "user2_id is invalid")
		return
	}

	chatId := chat.GenerateChatId(user1Id, user2Id)
	shardFactor := chat.GenerateShardFactor(user1Id, user2Id)

	result := &message.MessagesListResponse{}
	if err := message.GetMessagesByChatId(r.Context(), chatId, shardFactor, result); err != nil {
		logger.Error(err.Error())
		httpHelper.InternalServerErrorResponse(w)
		return
	}

	httpHelper.JsonResponse(w, result.ToResponse())
}
