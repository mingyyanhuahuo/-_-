package model

type OverviewResponse struct {
	TotalClaims   int64 `json:"totalClaims"`
	LostCount     int64 `json:"lostCount"`
	FoundCount    int64 `json:"foundCount"`
	PendingAudit  int64 `json:"pendingAudit"`
	ClaimedCount  int64 `json:"claimedCount"`
	ClosedCount   int64 `json:"closedCount"`
	UserCount     int64 `json:"userCount"`
	TodayNewItems int64 `json:"todayNewItems"`
}

type CategoryResponse struct {
	CategoryId   uint   `json:"categoryId"`
	CategoryName string `json:"categoryName"`
	Count        int64  `json:"count"`
}

type TrendResponse struct {
	Date       string `json:"date"`
	LostCount  int64  `json:"lostCount"`
	FoundCount int64  `json:"foundCount"`
}

type ClaimRateResponse struct {
	TotalClaims     int64   `json:"totalClaims"`
	ApprovedClaims  int64   `json:"approvedClaims"`
	RejectedClaims  int64   `json:"rejectedClaims"`
	CancelledClaims int64   `json:"cancelledClaims"`
	ClaimRate       float64 `json:"claimRate"`
}
