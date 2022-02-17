package message

import (
	"context"
	"otus_dialog_go/internal/otusdb"
)

type MessagesListResponse struct {
	Items []*Message
}

func (r *MessagesListResponse) ToResponse() map[string]interface{} {
	var items interface{}
	if len(r.Items) > 0 {
		items = MessagesToResponse(r.Items)
	} else {
		items = make([]string, 0)
	}
	return map[string]interface{}{
		"items": items,
	}
}

func GetMessage(ctx context.Context, id string, shardFactor string, message *Message) error {
	err := otusdb.Db.QueryRowContext(
		ctx,
		"SELECT BIN_TO_UUID(id), message, author_id, chat_id, shard_factor, created_at FROM messages WHERE id = UUID_TO_BIN(?) and shard_factor = '"+shardFactor+"'",
		id,
	).Scan(&message.Id, &message.Message, &message.AuthorId, &message.ChatId, &message.ShardFactor, &message.CreatedAt)
	return err
}

func GetMessagesByChatId(ctx context.Context, chatId string, shardFactor string, result *MessagesListResponse) error {
	rows, err := otusdb.Db.QueryContext(
		ctx,
		"SELECT BIN_TO_UUID(id), message, author_id, chat_id, shard_factor, created_at FROM messages WHERE shard_factor = '"+shardFactor+"' AND chat_id=?"+
			" ORDER BY created_at DESC",
		chatId,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	count := 0
	for rows.Next() {
		count++
		message := &Message{}
		err := rows.Scan(&message.Id, &message.Message, &message.AuthorId, &message.ChatId, &message.ShardFactor, &message.CreatedAt)
		if err != nil {
			return err
		}
		result.Items = append(result.Items, message)
	}
	return nil
}
