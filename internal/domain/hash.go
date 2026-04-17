package domain

type HashDto struct {
	Hash string `uri:"hash" binding:"required"`
}
