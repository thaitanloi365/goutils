package gormpaging

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

// PaginationCountOnlyParam ...
type PaginationCountOnlyParam struct {
	CountSQL    *gorm.DB
	CountRawSQL *gorm.DB
	Page        int
	Limit       int
}

func pagingCountOnly(p *PaginationCountOnlyParam, result interface{}) *Pagination {
	db := p.CountSQL

	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit == 0 {
		p.Limit = 10
	}

	done := make(chan bool, 1)
	var paginator Pagination
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

	<-done

	paginator.TotalRecord = count
	paginator.Records = result
	paginator.Page = p.Page

	paginator.Offset = offset
	paginator.PerPage = p.Limit
	paginator.TotalPage = int(math.Ceil(float64(count) / float64(p.Limit)))

	if p.Page > 1 {
		paginator.PrevPage = p.Page - 1
	} else {
		paginator.PrevPage = p.Page
	}

	if p.Page == paginator.TotalPage {
		paginator.NextPage = p.Page
	} else {
		paginator.NextPage = p.Page + 1
	}

	paginator.HasNext = paginator.TotalPage > paginator.Page
	paginator.HasPrev = paginator.Page > 1

	return &paginator
}

// Paging paging list
func paging(p *PaginationParam, result interface{}) *Pagination {
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
	var paginator Pagination
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

	paginator.TotalRecord = count
	paginator.Records = result
	paginator.Page = p.Page

	paginator.Offset = offset
	paginator.PerPage = p.Limit
	paginator.TotalPage = int(math.Ceil(float64(count) / float64(p.Limit)))

	if p.Page > 1 {
		paginator.PrevPage = p.Page - 1
	} else {
		paginator.PrevPage = p.Page
	}

	if p.Page == paginator.TotalPage {
		paginator.NextPage = p.Page
	} else {
		paginator.NextPage = p.Page + 1
	}

	paginator.HasNext = paginator.TotalPage > paginator.Page
	paginator.HasPrev = paginator.Page > 1

	return &paginator
}

func countRecords(db *gorm.DB, anyType interface{}, done chan bool, count *int) {
	db.Model(anyType).Count(count)
	done <- true
}
