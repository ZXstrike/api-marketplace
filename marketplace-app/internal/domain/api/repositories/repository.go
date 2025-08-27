package repositories

import (
	"strings"
	"time"

	"github.com/ZXstrike/shared/pkg/models"
	"gorm.io/gorm"
)

type Repository interface {
	GetUserByID(id string) (*models.User, error)
	GetApiTotalNumbers() (int64, error)
	GetUserTotalNumbers() (int64, error)
	GetTransaction24HNumbers() (int64, error)
	GetRecentAllTopUps() ([]models.PaymentTransaction, error)
	CreateAPI(api models.API, pricePercall float64) (string, error)
	GetAPIByID(id string) (*models.API, error)
	GetAllAPI(page int, length int, query string) ([]models.API, error)
	GetAllAPIByUserID(userID string) ([]models.API, error)
	UpdateAPI(api models.API) error
	DeleteAPI(id string) error
	GetCategoryBySlug(slug string) (*models.Category, error)
	GetAllCategories() ([]models.Category, error)
	GetAPIVersionByID(id string) (*models.APIVersion, error)
	CreateAPIEndpoint(apiEndpoint models.Endpoint) error
	GetAPIEndpointByID(id string) (*models.Endpoint, error)
	UpdateAPIEndpoint(apiEndpoint models.Endpoint) error
	DeleteAPIEndpoint(id string) error
	GetAllEndpointsByAPIVersionID(apiVersionID string) ([]models.Endpoint, error)
	GetConsumerOverview(userId string) (map[string]interface{}, error)
	GetProviderOverview(userId string) (map[string]interface{}, error)
}

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetUserByID(id string) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&user).Association("Roles").Find(&user.Roles); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repository) GetApiTotalNumbers() (int64, error) {
	var count int64
	if err := r.db.Model(&models.API{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *repository) GetUserTotalNumbers() (int64, error) {
	var count int64
	if err := r.db.Model(&models.User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *repository) GetTransaction24HNumbers() (int64, error) {
	var count int64
	cutoff := time.Now().Add(-24 * time.Hour)

	tx := r.db.Model(&models.PaymentTransaction{}).
		Where("created_at >= ?", cutoff).
		Count(&count)

	if tx.Error != nil {
		// Handle missing column gracefully
		if strings.Contains(tx.Error.Error(), `column "created_at" does not exist`) {
			// Attempt to add the column (requires model definition to have CreatedAt)
			if err := r.db.Migrator().AddColumn(&models.PaymentTransaction{}, "CreatedAt"); err == nil {
				// Retry with the time filter after adding column
				if retryErr := r.db.Model(&models.PaymentTransaction{}).
					Where("created_at >= ?", cutoff).
					Count(&count).Error; retryErr == nil {
					return count, nil
				}
			}
		}
		return 0, tx.Error
	}

	return count, nil
}

func (r *repository) GetRecentAllTopUps() ([]models.PaymentTransaction, error) {
	var transactions []models.PaymentTransaction

	tx := r.db.Preload("User").Preload("Subscription").Model(&models.PaymentTransaction{}).Limit(100).Order("created_at desc").Find(&transactions)

	if tx.Error != nil {
		// Handle missing column gracefully
		if strings.Contains(tx.Error.Error(), `column "created_at" does not exist`) {
			// Attempt to add the column (requires model definition to have CreatedAt)
			if err := r.db.Migrator().AddColumn(&models.PaymentTransaction{}, "CreatedAt"); err == nil {
				// Retry with the time filter after adding column
				if retryErr := r.db.Preload("User").Preload("Subscription").Model(&models.PaymentTransaction{}).
					Limit(50).Error; retryErr == nil {

					return transactions, nil
				}
			}

		}
		return nil, tx.Error
	}

	return transactions, nil
}

func (r *repository) CreateAPI(api models.API, pricePercall float64) (string, error) {

	if err := r.db.Create(&api).Error; err != nil {
		return "", err
	}

	apiVersion := models.APIVersion{
		APIID:         api.ID,
		API:           api,
		VersionString: "v1.0.0",
		PricePerCall:  pricePercall,
	}

	if err := r.db.Create(&apiVersion).Error; err != nil {
		return "", err
	}

	if err := r.db.Model(&api).Association("Versions").Append(&apiVersion); err != nil {
		return "", err
	}

	return apiVersion.ID, nil
}

func (r *repository) GetAPIByID(id string) (*models.API, error) {
	var api models.API
	if err := r.db.First(&api, "id = ?", id).Error; err != nil {
		return nil, err
	}
	// Preload the Provider and Categories associations.
	if err := r.db.Model(&api).Preload("Provider").Preload("Categories").First(&api).Error; err != nil {
		return nil, err
	}
	// Preload the Versions association.
	if err := r.db.Model(&api).Preload("Versions").First(&api).Error; err != nil {
		return nil, err
	}

	return &api, nil
}

func (r *repository) GetAllAPI(page int, length int, query string) ([]models.API, error) {
	if page < 1 {
		page = 1
	}
	if length <= 0 {
		length = 10
	}

	var apis []models.API
	db := r.db.
		Preload("Provider").
		Preload("Categories").
		Preload("Versions.Endpoints")

	if query = strings.TrimSpace(query); query != "" {
		// Postgres: ILIKE. If not Postgres, switch to LOWER(name) LIKE LOWER(?)
		db = db.Where("name ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%")
	}

	err := db.
		Offset((page - 1) * length).
		Limit(length).
		Find(&apis).Error
	if err != nil {
		return nil, err
	}

	// Return empty slice (usual pattern).
	return apis, nil
}

func (r *repository) GetAllAPIByUserID(userID string) ([]models.API, error) {
	var apis []models.API
	if err := r.db.Where("provider_id = ?", userID).Find(&apis).Error; err != nil {
		return nil, err
	}
	return apis, nil
}

func (r *repository) UpdateAPI(api models.API) error {
	return r.db.Save(&api).Error
}

func (r *repository) DeleteAPI(id string) error {
	return r.db.Delete(&models.API{}, "id = ?", id).Error
}

func (r *repository) GetCategoryBySlug(slug string) (*models.Category, error) {
	var category models.Category
	if err := r.db.First(&category, "slug = ?", slug).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *repository) GetAllCategories() ([]models.Category, error) {
	var categories []models.Category
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *repository) GetAPIVersionByID(id string) (*models.APIVersion, error) {
	var apiVersion models.APIVersion
	if err := r.db.First(&apiVersion, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &apiVersion, nil
}

func (r *repository) CreateAPIEndpoint(apiEndpoint models.Endpoint) error {
	return r.db.Create(&apiEndpoint).Error
}

func (r *repository) GetAPIEndpointByID(id string) (*models.Endpoint, error) {
	var apiEndpoint models.Endpoint
	if err := r.db.First(&apiEndpoint, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &apiEndpoint, nil
}

func (r *repository) UpdateAPIEndpoint(apiEndpoint models.Endpoint) error {
	return r.db.Save(&apiEndpoint).Error
}

func (r *repository) DeleteAPIEndpoint(id string) error {
	return r.db.Delete(&models.Endpoint{}, "id = ?", id).Error
}

func (r *repository) GetAllEndpointsByAPIVersionID(apiVersionID string) ([]models.Endpoint, error) {
	var endpoints []models.Endpoint
	if err := r.db.Where("api_version_id = ?", apiVersionID).Find(&endpoints).Error; err != nil {
		return nil, err
	}
	return endpoints, nil
}

func (r *repository) GetConsumerOverview(userId string) (map[string]interface{}, error) {
	overview := make(map[string]interface{})

	// 1. Active Subscriptions Count
	var activeSubscriptionsCount int64
	if err := r.db.Model(&models.Subscription{}).
		Where("consumer_user_id = ?", userId).
		Count(&activeSubscriptionsCount).Error; err != nil {
		return nil, err
	}
	overview["active_subscriptions_count"] = activeSubscriptionsCount

	// 2. Total Monthly Cost from all active subscriptions
	var totalMonthlyCost float64
	if err := r.db.Model(&models.UsageLog{}).
		Joins("JOIN subscriptions ON usage_logs.subscription_id = subscriptions.id").
		Joins("JOIN api_versions ON subscriptions.api_version_id = api_versions.id").
		Where("subscriptions.consumer_user_id = ?", userId).
		Select("COALESCE(SUM(api_versions.price_per_call), 0)").
		Scan(&totalMonthlyCost).Error; err != nil {
		return nil, err
	}
	overview["total_monthly_cost"] = totalMonthlyCost

	// 3. Total requests in the last 30 days
	var requestsLast30Days int64
	if err := r.db.Model(&models.UsageLog{}).
		Joins("JOIN subscriptions ON usage_logs.subscription_id = subscriptions.id").
		Where("subscriptions.consumer_user_id = ? AND usage_logs.request_timestamp >= ?", userId, time.Now().AddDate(0, -1, 0)).
		Count(&requestsLast30Days).Error; err != nil {
		return nil, err
	}
	overview["requests_last_30_days"] = requestsLast30Days

	// 4. Total requests in the last 7 days for all subscribed APIs
	type dailyUsage struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}
	var requestsLast7Days []dailyUsage
	now := time.Now()

	for i := 6; i >= 0; i-- {
		var count int64
		day := now.AddDate(0, 0, -i)
		dayStart := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, now.Location())
		dayEnd := dayStart.AddDate(0, 0, 1)

		if err := r.db.Model(&models.UsageLog{}).
			Joins("JOIN subscriptions ON usage_logs.subscription_id = subscriptions.id").
			Where("subscriptions.consumer_user_id = ? AND usage_logs.request_timestamp >= ? AND usage_logs.request_timestamp < ?", userId, dayStart, dayEnd).
			Count(&count).Error; err != nil {
			return nil, err
		}

		dateLabel := day.Format("January 2")
		if i == 0 {
			dateLabel = "Today"
		}

		requestsLast7Days = append(requestsLast7Days, dailyUsage{
			Date:  dateLabel,
			Count: count,
		})
	}
	overview["requests_last_7_days"] = requestsLast7Days

	return overview, nil
}

func (r *repository) GetProviderOverview(userId string) (map[string]interface{}, error) {
	overview := make(map[string]interface{})

	// 1. Active Subscriber Count
	var activeSubscriberCount int64
	if err := r.db.Model(&models.Subscription{}).
		Joins("JOIN api_versions ON subscriptions.api_version_id = api_versions.id").
		Joins("JOIN apis ON api_versions.api_id = apis.id").
		Where("apis.provider_id = ?", userId).
		Count(&activeSubscriberCount).Error; err != nil {
		return nil, err
	}
	overview["active_subscriber_count"] = activeSubscriberCount

	// 2. Total Monthly Revenue from all active subscriptions
	var totalMonthlyRevenue float64
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	if err := r.db.Model(&models.PaymentTransaction{}).
		Joins("JOIN wallets ON payment_transactions.user_id = wallets.user_id").
		Where("payment_transactions.user_id = ? AND wallets.wallet_type = ? AND payment_transactions.amount > 0", userId, models.WalletTypeProvider).
		Select("COALESCE(SUM(payment_transactions.amount), 0)").
		Scan(&totalMonthlyRevenue).Error; err != nil {
		return nil, err
	}
	overview["total_revenue"] = totalMonthlyRevenue

	// 3. Total requests in the last 30 days
	var requestsLast30Days int64
	if err := r.db.Model(&models.UsageLog{}).
		Joins("JOIN subscriptions ON usage_logs.subscription_id = subscriptions.id").
		Joins("JOIN api_versions ON subscriptions.api_version_id = api_versions.id").
		Joins("JOIN apis ON api_versions.api_id = apis.id").
		Where("apis.provider_id = ? AND usage_logs.request_timestamp >= ?", userId, thirtyDaysAgo).
		Count(&requestsLast30Days).Error; err != nil {
		return nil, err
	}
	overview["requests_last_30_days"] = requestsLast30Days

	// 4. Total requests in the last 4 weeks for all subscribed APIs
	var requestsLast4Weeks []map[string]interface{}
	now := time.Now()
	for i := 3; i >= 0; i-- {
		var count int64
		weekStart := now.AddDate(0, 0, -7*i)
		weekStart = time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, now.Location())
		weekEnd := weekStart.AddDate(0, 0, 7)
		if err := r.db.Model(&models.UsageLog{}).
			Joins("JOIN subscriptions ON usage_logs.subscription_id = subscriptions.id").
			Joins("JOIN api_versions ON subscriptions.api_version_id = api_versions.id").
			Joins("JOIN apis ON api_versions.api_id = apis.id").
			Where("apis.provider_id = ? AND usage_logs.request_timestamp >= ? AND usage_logs.request_timestamp < ?", userId, weekStart, weekEnd).
			Count(&count).Error; err != nil {
			return nil, err
		}
		weekLabel := weekStart.Format("Jan 2")
		if i == 0 {
			weekLabel = "This Week"
		}
		requestsLast4Weeks = append(requestsLast4Weeks, map[string]interface{}{
			"week":  weekLabel,
			"count": count,
		})
	}
	overview["requests_last_4_weeks"] = requestsLast4Weeks

	return overview, nil
}
