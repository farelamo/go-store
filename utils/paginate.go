package utils

func Paginate(reqPage, reqPageSize int) (page, pageSize int) {
	page = reqPage
	if reqPage <= 0 {
		page = 1
	}
	
	pageSize = reqPageSize
	if reqPageSize <= 0 {
		pageSize = 10
	}

	return page, pageSize
}
