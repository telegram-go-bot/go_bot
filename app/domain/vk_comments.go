package domain

import (
	"encoding/json"
	"log"
	"strconv"
)

// CommentItem : Internal representation of comment item from json
type CommentItem struct {
	text          string
	likes         int
	attachmentURL []struct { // might be nil
		URL string
	}
}

//GetComments : returns only comments with Likes > @minLikes
// and items with attachments are ignored
// Return: bool - true means @count of comments were processed, and we have some more.
//				 You are free to call me again with offset += count
func GetComments(messageID int, ownerID int, count int, offset int, minLikes int) (bool, []CommentItem) {
	parameters := make(map[string]string)
	parameters["count"] = strconv.Itoa(count)
	parameters["offset"] = strconv.Itoa(offset)
	parameters["post_id"] = strconv.Itoa(messageID)
	parameters["owner_id"] = strconv.Itoa(ownerID)
	parameters["need_likes"] = "1"

	resp, err := VkRequest("wall.getComments", parameters, authResponse)
	if err != nil {
		log.Printf("wall.get request failed: %s\n", err)
		return false, nil
	}

	var commentsResponse GetCommentsResponse
	err = json.Unmarshal(resp, &commentsResponse)
	if err != nil {
		log.Printf(string(resp))
		log.Printf("Error unmarshalling json response: %s\n", err)
		return false, nil
	}

	items := commentsResponse.Response.Items

	res := make([]CommentItem, 0, len(items))
	for _, item := range items {
		if item.Likes.Count <= minLikes || len(item.Text) == 0 {
			continue
		}

		//if item.Attachments != nil {
		// TODO: fill attachment item with most largest photo URL.
		// skip other types instewad of "photo"
		//}

		res = append(res, CommentItem{item.Text, item.Likes.Count, nil})
	}

	if len(res) == 0 {
		return false, nil
	}

	return len(items) == count, res
}

// GetTopComment : returns true on success. If false - string is undefined
func GetTopComment(messageID int, ownerID int, commentsCount int, minLikes int) string {
	offset := 0
	const blockSize = 100

	if commentsCount == 0 {
		return ""
	}

	var res CommentItem
	for {
		hasMoreComments, commentsPart := GetComments(messageID, ownerID, blockSize, offset, minLikes)
		for _, comment := range commentsPart {
			if comment.likes > res.likes {
				res = comment
			}
		}

		if hasMoreComments == false || commentsCount <= 0 {
			break
		}

		//TODO: prefer offset or commentsCount ...
		offset += blockSize
		commentsCount -= blockSize
	}
	return res.text
}
