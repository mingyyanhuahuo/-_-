package service

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"sync"
	"time"

	"lostfound/dao"
	"lostfound/dto"
	"lostfound/model"
	"lostfound/pkg/errcode"
	"lostfound/pkg/logger"

	"go.uber.org/zap"
)

const (
	matchCacheTTL    = 10 * time.Minute
	matchDefaultSize = 10
	matchMaxSize     = 20
	matchMinScore    = 0.3
)

type matchCacheEntry struct {
	items     []dto.MatchItem
	expiredAt time.Time
}

var matchCache sync.Map

func GetItemMatches(itemID uint, req *dto.MatchRequest) ([]dto.MatchItem, error) {
	item, err := getItemOrFail(itemID)
	if err != nil {
		return nil, err
	}
	limit := req.Limit
	if limit <= 0 {
		limit = matchDefaultSize
	}
	if limit > matchMaxSize {
		limit = matchMaxSize
	}
	cacheKey := fmt.Sprintf("%d:%d", itemID, limit)
	if v, ok := matchCache.Load(cacheKey); ok {
		if entry, ok := v.(matchCacheEntry); ok && time.Now().Before(entry.expiredAt) {
			return entry.items, nil
		}
	}
	items, err := computeMatches(item, limit)
	if err != nil {
		return nil, err
	}
	matchCache.Store(cacheKey, matchCacheEntry{items: items, expiredAt: time.Now().Add(matchCacheTTL)})
	return items, nil
}

func computeMatches(item *model.Item, limit int) ([]dto.MatchItem, error) {
	oppositeType := model.ItemTypeFound
	if item.Type == model.ItemTypeFound {
		oppositeType = model.ItemTypeLost
	}
	candidates, err := dao.ListItemsForMatch(oppositeType, item.ID)
	if err != nil {
		logger.Logger.Error("查询匹配候选失败", zap.Uint("itemId", item.ID), zap.Error(err))
		return nil, errcode.ErrInternalServer
	}
	results := make([]dto.MatchItem, 0, len(candidates))
	for i := range candidates {
		c := &candidates[i]
		titleSim := textSimilarity(item.Title, c.Title)
		categoryScore := 0.0
		if item.CategoryID == c.CategoryID {
			categoryScore = 1.0
		}
		locSim := textSimilarity(item.Location, c.Location)
		timeSim := timeProximity(item.LostTime, c.LostTime)
		score := titleSim*0.4 + categoryScore*0.25 + locSim*0.20 + timeSim*0.15
		if score < matchMinScore {
			continue
		}
		reasons := make([]string, 0, 4)
		if titleSim >= 0.5 {
			reasons = append(reasons, "标题高度相似")
		}
		if categoryScore > 0 {
			reasons = append(reasons, "同一分类")
		}
		if locSim >= 0.5 {
			reasons = append(reasons, "地点接近")
		}
		if timeSim >= 0.5 {
			reasons = append(reasons, "时间接近")
		}
		results = append(results, dto.MatchItem{
			ItemBrief:   toItemBrief(c),
			MatchScore:  math.Round(score*100) / 100,
			MatchReason: reasons,
		})
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].MatchScore > results[j].MatchScore
	})
	if len(results) > limit {
		results = results[:limit]
	}
	if results == nil {
		results = []dto.MatchItem{}
	}
	return results, nil
}

// textSimilarity 基于字符二元组 Jaccard 计算文本相似度(0-1)。
func textSimilarity(a, b string) float64 {
	ba := bigrams(a)
	bb := bigrams(b)
	if len(ba) == 0 && len(bb) == 0 {
		return 0
	}
	inter := 0
	for k := range ba {
		if _, ok := bb[k]; ok {
			inter++
		}
	}
	union := len(ba) + len(bb) - inter
	if union == 0 {
		return 0
	}
	return float64(inter) / float64(union)
}

func bigrams(s string) map[string]struct{} {
	r := []rune(strings.ToLower(s))
	m := make(map[string]struct{}, len(r))
	for i := 0; i+1 < len(r); i++ {
		m[string(r[i:i+2])] = struct{}{}
	}
	return m
}

// timeProximity 时间接近度，30 天内线性衰减，越接近分值越高。
func timeProximity(a, b time.Time) float64 {
	days := math.Abs(a.Sub(b).Hours()) / 24
	if days > 30 {
		return 0
	}
	return 1 - days/30
}