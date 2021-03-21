package vk

// WallGetResponse - wall.get response
type WallGetResponse struct {
	Response struct {
		Count int `json:"count"`
		Items []struct {
			ID          int    `json:"id"`
			FromID      int    `json:"from_id"`
			OwnerID     int    `json:"owner_id"`
			Date        int    `json:"date"`
			Type        string `json:"type"`
			MarkedAsAds int    `json:"marked_as_ads"`
			PostType    string `json:"post_type"`
			Text        string `json:"text"`
			SignerID    int    `json:"signer_id"`
			IsPinned    int    `json:"is_pinned,omitempty"`
			Attachments []struct {
				Type  string `json:"type"`
				Photo struct {
					ID      int `json:"id"`
					AlbumID int `json:"album_id"`
					OwnerID int `json:"owner_id"`
					UserID  int `json:"user_id"`
					Sizes   []struct {
						Type   string `json:"type"`
						URL    string `json:"url"`
						Width  int    `json:"width"`
						Height int    `json:"height"`
					} `json:"sizes"`
					Text      string `json:"text"`
					Date      int    `json:"date"`
					PostID    int    `json:"post_id"`
					AccessKey string `json:"access_key"`
				} `json:"photo"`
			} `json:"attachments"`
			PostSource struct {
				Type string `json:"type"`
			} `json:"post_source"`
			Comments struct {
				Count         int  `json:"count"`
				CanPost       int  `json:"can_post"`
				GroupsCanPost bool `json:"groups_can_post"`
			} `json:"comments"`
			Likes struct {
				Count      int `json:"count"`
				UserLikes  int `json:"user_likes"`
				CanLike    int `json:"can_like"`
				CanPublish int `json:"can_publish"`
			} `json:"likes"`
			Reposts struct {
				Count        int `json:"count"`
				UserReposted int `json:"user_reposted"`
			} `json:"reposts"`
			Views struct {
				Count int `json:"count"`
			} `json:"views"`
			IsFavorite     bool `json:"is_favorite"`
			AdsEasyPromote struct {
				Type int `json:"type"`
			} `json:"ads_easy_promote"`
		} `json:"items"`
		NextFrom string `json:"next_from"`
	} `json:"response"`
}
