package repositories

import (
	"senkou-catalyst-be/app/dtos"
	"senkou-catalyst-be/app/models"
	"senkou-catalyst-be/utils/converter"
	"senkou-catalyst-be/utils/query"
	"time"

	"gorm.io/gorm"
)

type ProductInteractionRepository interface {
	StoreProductInteractionLog(interaction *models.ProductMetric) error
	GetMerchantProductsMetric(merchantID string, params *query.QueryParams) ([]dtos.ProductReport, error)
	GetMerchantProductsMetricStats(merchantID string, params *query.QueryParams) (*dtos.ProductMetricStats, error)
}

type ProductInteractionRepositoryInstance struct {
	db *gorm.DB
}

func NewProductInteractionRepository(db *gorm.DB) ProductInteractionRepository {
	return &ProductInteractionRepositoryInstance{
		db: db,
	}
}

// StoreProductInteractionLog stores a new product interaction log in the database
func (r *ProductInteractionRepositoryInstance) StoreProductInteractionLog(interaction *models.ProductMetric) error {
	return r.db.Create(interaction).Error
}

func (r *ProductInteractionRepositoryInstance) GetMerchantProductsMetric(merchantID string, params *query.QueryParams) ([]dtos.ProductReport, error) {
	var productReports []dtos.ProductReport

	currentDateFrom := converter.ParseDate(params.DateFrom, time.Now().AddDate(0, 0, -30)).Format("2006-01-02")
	currentDateTo := converter.ParseDate(params.DateTo, time.Now()).Format("2006-01-02")

	query := r.db.Model(&models.Product{}).
		Select(`
            products.title as name,
            COUNT(CASE WHEN pm.interaction = 'view' AND DATE(pm.created_at) >= ? AND DATE(pm.created_at) <= ? THEN 1 END) as total_views,
            COUNT(CASE WHEN pm.interaction = 'click' AND DATE(pm.created_at) >= ? AND DATE(pm.created_at) <= ? THEN 1 END) as total_clicks
        `,
			currentDateFrom,
			currentDateTo,
			currentDateFrom,
			currentDateTo,
		).
		Joins("LEFT JOIN product_metrics pm ON pm.product_id = products.id").
		Where("products.merchant_id = ? AND DATE(pm.created_at) BETWEEN ? AND ?", merchantID, currentDateFrom, currentDateTo).
		Group("products.id, products.title").
		Having("COUNT(pm.id) > 0")

	rows, err := query.Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var report dtos.ProductReport
		err := rows.Scan(&report.Name, &report.TotalViews, &report.TotalClicks)
		if err != nil {
			return nil, err
		}
		productReports = append(productReports, report)
	}

	return productReports, nil
}

func (r *ProductInteractionRepositoryInstance) GetMerchantProductsMetricStats(merchantID string, params *query.QueryParams) (*dtos.ProductMetricStats, error) {
	productMetricStats := new(dtos.ProductMetricStats)

	currentDateFrom := converter.ParseDate(params.DateFrom, time.Now().AddDate(0, 0, -30))
	currentDateTo := converter.ParseDate(params.DateTo, time.Now())

	query := `
		SELECT
			COUNT(CASE WHEN pm.interaction = 'view' THEN 1 END) AS total_views,
			COUNT(CASE WHEN pm.interaction = 'click' THEN 1 END) AS total_clicks
		FROM product_metrics pm
			LEFT OUTER JOIN products p ON p.id = pm.product_id
			LEFT JOIN merchants m ON m.id = p.merchant_id
			WHERE m.id = ? AND DATE(pm.created_at) BETWEEN ? AND ?
	`
	if err := r.db.Raw(query, merchantID, currentDateFrom.Format("2006-01-01"), currentDateTo.Format("2006-01-02")).Scan(productMetricStats).Error; err != nil {
		return nil, err
	}

	return productMetricStats, nil
}
