package requests

import "github.com/goravel/framework/contracts/http"

type BlockRequest struct {
	UserID int `form:"user_id" json:"user_id"`
}

func (r *BlockRequest) Authorize(http.Context) error { return nil }

func (r *BlockRequest) Rules(http.Context) map[string]any {
	return map[string]any{
		"user_id": "required|integer|min:1",
	}
}

func (r *BlockRequest) Target() uint { return uint(r.UserID) }
