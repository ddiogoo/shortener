package dto

type ShortenedRequest struct {
	Url string `json:"url"`
}

func NewShortenedRequest() ShortenedRequest {
	return ShortenedRequest{
		Url: "",
	}
}

type ShortenedResponse struct {
	Id           uint  `json:"id"`
	RowsAffected int64 `json:"rowsAffected"`
}

func NewShortenedResponse(id uint, rowsAffected int64) ShortenedResponse {
	return ShortenedResponse{
		Id:           id,
		RowsAffected: rowsAffected,
	}
}
