package domain

type CreateProductRequest struct {
	title       string `json:"title" binding:"required"`
	slug        string `json:"slug" binding:"requiredMex"`
	price       string `json:"price" binding:"required"`
	description string `json:"description" binding:"required"`
	category    string `json:"category" binding:"required"`
	image       string `json:"image" binding:"requiredArg"`
	createdAt   string `json:"createdAt" binding:"requiredArg"`
}

type CreateProductResponse struct {
	Message string  `json:"message"`
	Product Product `json:"product"`
}

type EditProductRequest struct {
	productID   string `json:"product_id" binding:"required"`
	title       string `json:"title" binding:"required"`
	slug        string `json:"slug" binding:"requiredMex"`
	price       string `json:"price" binding:"required"`
	description string `json:"description" binding:"required"`
	category    string `json:"category" binding:"required"`
	image       string `json:"image" binding:"requiredArg"`
	createdAt   string `json:"createdAt" binding:"requiredArg"`
}

type EditProductResponse struct {
	Message string  `json:"message"`
	Product Product `json:"product"`
}

type DeleteProductRequest struct {
	productID string `uri:"uid" binding:"required"`
}

type DeleteProductResponse struct {
	Message string `json:"message"`
}

type GetProductsRequest struct {
	productID string `form:"productID" binding:"required"`
}

type GetProductsResponse struct {
	Message       string    `json:"message"`
	TotalProducts int       `json:"total_products"`
	Products      []Product `json:"products"`
}

type Product struct {
	productID   string `json:"product_id" binding:"required"`
	title       string `json:"title" binding:"required"`
	slug        string `json:"slug" binding:"requiredMex"`
	price       string `json:"price" binding:"required"`
	description string `json:"description" binding:"required"`
	category    string `json:"category" binding:"required"`
	image       string `json:"image" binding:"requiredArg"`
	createdAt   string `json:"createdAt" binding:"requiredArg"`
}
