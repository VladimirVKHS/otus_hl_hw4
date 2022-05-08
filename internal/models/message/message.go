package message

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"otus_dialog_go/internal/helpers/chat"
	"otus_dialog_go/internal/otusdb"
	"strings"
)

type Message struct {
	Id          string
	AuthorId    int    `validate:"required"`
	ChatId      string `validate:"lte=255,required"`
	Message     string `validate:"lte=4096,required"`
	ShardFactor string
	CreatedAt   string
	IsRead      bool
}

type MessageCreateRequest struct {
	AuthorId    int    `validate:"required"`
	AddresseeId int    `validate:"required"`
	Message     string `validate:"lte=4096,required"`
}

type MarkAsReadRequest struct {
	AuthorId    int      `validate:"required"`
	AddresseeId int      `validate:"required"`
	MessageIds  []string `validate:"required,gt=0,lte=1000,dive,uuid"`
}

func (m *MessageCreateRequest) GenerateChatId() string {
	return chat.GenerateChatId(m.AuthorId, m.AddresseeId)
}

func (m *MessageCreateRequest) GenerateShardFactor() string {
	return chat.GenerateShardFactor(m.AuthorId, m.AddresseeId)
}

func (m *MessageCreateRequest) CreateMessage(ctx context.Context) (*Message, error) {
	id := uuid.New().String()
	shardFactor := m.GenerateShardFactor()
	_, err := otusdb.Db.ExecContext(
		ctx,
		"INSERT INTO messages SET id = UUID_TO_BIN(?), message = ?, chat_id = ?, author_id = ?, shard_factor = '"+shardFactor+"'",
		id, m.Message, m.GenerateChatId(), m.AuthorId,
	)
	if err != nil {
		return nil, err
	}
	message := &Message{}
	if err := GetMessage(ctx, id, shardFactor, message); err != nil {
		return nil, err
	}

	return message, nil
}

func (m *MarkAsReadRequest) GenerateShardFactor() string {
	return chat.GenerateShardFactor(m.AuthorId, m.AddresseeId)
}

func (m *MarkAsReadRequest) MarkAsRead(ctx context.Context) (int, error) {
	var idsArray []string
	for _, id := range m.MessageIds {
		idsArray = append(idsArray, fmt.Sprintf("UUID_TO_BIN('%s')", id))
	}
	shardFactor := m.GenerateShardFactor()
	res, err := otusdb.Db.ExecContext(
		ctx,
		"UPDATE messages SET is_read = true WHERE shard_factor = '"+shardFactor+"' AND id IN ("+strings.Join(idsArray, ", ")+") AND author_id <> ?",
		m.AddresseeId,
	)
	affected, _ := res.RowsAffected()
	return int(affected), err
}

func (m *Message) Refresh(ctx context.Context) error {
	return GetMessage(ctx, m.Id, m.ShardFactor, m)
}

func (m *Message) ToResponse() map[string]interface{} {
	return map[string]interface{}{
		"Id":        m.Id,
		"Message":   m.Message,
		"AuthorId":  m.AuthorId,
		"CreatedAt": m.CreatedAt,
		"IsRead":    m.IsRead,
	}
}

func MessagesToResponse(messages []*Message) []map[string]interface{} {
	var result []map[string]interface{}
	for _, m := range messages {
		result = append(result, m.ToResponse())
	}
	return result
}
