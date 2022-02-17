package chat

import "strconv"

func GenerateChatId(user1Id int, user2Id int) string {
	if user1Id < user2Id {
		return strconv.Itoa(user1Id) + "_" + strconv.Itoa(user2Id)
	}
	return strconv.Itoa(user2Id) + "_" + strconv.Itoa(user1Id)
}

func GenerateShardFactor(user1Id int, user2Id int) string {
	return strconv.Itoa(user1Id%10 + user2Id%10)
}
