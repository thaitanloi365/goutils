package gormutils

import (
	"math"

	"github.com/jinzhu/gorm"
)

// CountResult result
type CountResult struct {
	Count int `json:"count"`
}

// Pagination ...
type Pagination struct {
	HasNext     bool        `json:"has_next"`
	HasPrev     bool        `json:"has_prev"`
	PerPage     int         `json:"per_page"`
	NextPage    int         `json:"next_page"`
	Page        int         `json:"current_page"`
	PrevPage    int         `json:"prev_page"`
	Offset      int         `json:"offset"`
	Records     interface{} `json:"records"`
	TotalRecord int         `json:"total_record"`
	TotalPage   int         `json:"total_page"`
}

// PaginationParam ...
type PaginationParam struct {
	DB          *gorm.DB
	RawSQL      *gorm.DB
	CountRawSQL *gorm.DB
	Page        int
	Limit       int
	OrderBy     []string
}

// Paging paging list
func Paging(p *PaginationParam, result interface{}) *Pagination {
	db := p.DB

	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit == 0 {
		p.Limit = 10
	}
	if len(p.OrderBy) > 0 {
		for _, o := range p.OrderBy {
			db = db.Order(o)
			if p.RawSQL != nil {
				p.RawSQL = p.RawSQL.Order(o)
			}
		}

	}

	done := make(chan bool, 1)
	var pagination Pagination
	var count int
	var offset int

	if p.CountRawSQL != nil {
		go func() {
			var result CountResult
			p.CountRawSQL.Scan(&result)
			count = result.Count
			done <- true
		}()
	} else {
		go countRecords(db, result, done, &count)
	}

	if p.Page == 1 {
		offset = 0
	} else {
		offset = (p.Page - 1) * p.Limit
	}

	if p.RawSQL != nil {
		p.RawSQL.Limit(p.Limit).Offset(offset).Scan(result)
	} else {
		db.Limit(p.Limit).Offset(offset).Find(result)
	}

	<-done

	pagination.TotalRecord = count
	pagination.Records = result
	pagination.Page = p.Page

	pagination.Offset = offset
	pagination.PerPage = p.Limit
	pagination.TotalPage = int(math.Ceil(float64(count) / float64(p.Limit)))

	if p.Page > 1 {
		pagination.PrevPage = p.Page - 1
	} else {
		pagination.PrevPage = p.Page
	}

	if p.Page == pagination.TotalPage {
		pagination.NextPage = p.Page
	} else {
		pagination.NextPage = p.Page + 1
	}

	pagination.HasNext = pagination.TotalPage > pagination.Page
	pagination.HasPrev = pagination.Page > 1

	return &pagination
}

func countRecords(db *gorm.DB, anyType interface{}, done chan bool, count *int) {
	db.Model(anyType).Count(count)
	done <- true
}
