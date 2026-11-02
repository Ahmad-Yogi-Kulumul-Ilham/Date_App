package requests

import (
	"github.com/goravel/framework/contracts/http"

	"leochat/app/models"
)

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

const maxReportDetailsLength = 500

type ReportRequest struct {
	UserID  int    `form:"user_id" json:"user_id"`
	Reason  string `form:"reason" json:"reason"`
	Details string `form:"details" json:"details"`
}

func (r *ReportRequest) Authorize(http.Context) error { return nil }

func (r *ReportRequest) Rules(http.Context) map[string]any {
	return map[string]any{
		"user_id": "required|integer|min:1",
		"reason":  "required|in:fake_profile,inappropriate,harassment,spam,underage,other",
		"details": "sometimes|string|max:" + itoa(maxReportDetailsLength),
	}
}

func (r *ReportRequest) Target() uint { return uint(r.UserID) }

func (r *ReportRequest) ReportReason() models.ReportReason { return models.ReportReason(r.Reason) }
