package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// AdminStatsResponse represents the dashboard statistics
type AdminStatsResponse struct {
	GalleryWorksCount    int64                  `json:"gallery_works_count"`
	ShopItemsCount       int64                  `json:"shop_items_count"`
	NewsCount            int64                  `json:"news_count"`
	UsersCount           int64                  `json:"users_count"`
	CoursesCount         int64                  `json:"courses_count"`
	MasterClassesCount   int64                  `json:"master_classes_count"`
	ReviewsTotal         int64                  `json:"reviews_total"`
	ReviewsPending       int64                  `json:"reviews_pending"`
	PurchasesTotal       int64                  `json:"purchases_total"`
	RevenueTotal         int64                  `json:"revenue_total"`
	RevenueMonth         int64                  `json:"revenue_month"`
	ActiveChats          int64                  `json:"active_chats"`
	UsersRegisteredMonth int64                  `json:"users_registered_month"`
	SalesByMonth         []SalesByMonthEntry    `json:"sales_by_month"`
	PopularCategories    []PopularCategoryEntry `json:"popular_categories"`
}

type SalesByMonthEntry struct {
	Month   string `json:"month"`
	Count   int64  `json:"count"`
	Revenue int64  `json:"revenue"`
}

type PopularCategoryEntry struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

// RecentActivityItem represents a single activity entry
type RecentActivityItem struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Text      string `json:"text"`
	Time      string `json:"time"`
	CreatedAt string `json:"created_at"`
}

// AdminGetStats returns dashboard statistics
func AdminGetStats(c *gin.Context) {
	stats := AdminStatsResponse{}

	// Gallery works count
	err := db.QueryRow("SELECT COUNT(*) FROM works").Scan(&stats.GalleryWorksCount)
	if err != nil {
		log.Printf("Error counting works: %v", err)
	}

	// Shop items count (sales)
	err = db.QueryRow("SELECT COUNT(*) FROM sales").Scan(&stats.ShopItemsCount)
	if err != nil {
		log.Printf("Error counting sales: %v", err)
	}

	// News count
	err = db.QueryRow("SELECT COUNT(*) FROM news").Scan(&stats.NewsCount)
	if err != nil {
		log.Printf("Error counting news: %v", err)
	}

	// Users count
	err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&stats.UsersCount)
	if err != nil {
		log.Printf("Error counting users: %v", err)
	}

	// Courses count
	err = db.QueryRow("SELECT COUNT(*) FROM products WHERE type = 'course'").Scan(&stats.CoursesCount)
	if err != nil {
		log.Printf("Error counting courses: %v", err)
	}

	// Master classes count
	err = db.QueryRow("SELECT COUNT(*) FROM products WHERE type = 'masterclass'").Scan(&stats.MasterClassesCount)
	if err != nil {
		log.Printf("Error counting master classes: %v", err)
	}

	// Reviews total
	err = db.QueryRow("SELECT COUNT(*) FROM reviews").Scan(&stats.ReviewsTotal)
	if err != nil {
		log.Printf("Error counting reviews: %v", err)
	}

	// Reviews pending (not approved)
	err = db.QueryRow("SELECT COUNT(*) FROM reviews WHERE is_approved = FALSE").Scan(&stats.ReviewsPending)
	if err != nil {
		log.Printf("Error counting pending reviews: %v", err)
	}

	// Purchases total
	err = db.QueryRow("SELECT COUNT(*) FROM purchases").Scan(&stats.PurchasesTotal)
	if err != nil {
		log.Printf("Error counting purchases: %v", err)
	}

	// Revenue total (sum of price_paid in kopecks, convert to rubles)
	var revenueTotalKopecks int64
	err = db.QueryRow("SELECT COALESCE(SUM(price_paid), 0) FROM purchases WHERE status = 'active'").Scan(&revenueTotalKopecks)
	if err != nil {
		log.Printf("Error calculating total revenue: %v", err)
	}
	stats.RevenueTotal = revenueTotalKopecks / 100

	// Revenue this month (in rubles)
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	var revenueMonthKopecks int64
	err = db.QueryRow(
		"SELECT COALESCE(SUM(price_paid), 0) FROM purchases WHERE status = 'active' AND created_at >= ?",
		monthStart.Format("2006-01-02 15:04:05"),
	).Scan(&revenueMonthKopecks)
	if err != nil {
		log.Printf("Error calculating monthly revenue: %v", err)
	}
	stats.RevenueMonth = revenueMonthKopecks / 100

	// Active chats (unresolved threads)
	err = db.QueryRow("SELECT COUNT(*) FROM chat_threads WHERE is_resolved = FALSE").Scan(&stats.ActiveChats)
	if err != nil {
		log.Printf("Error counting active chats: %v", err)
	}

	// Users registered this month
	err = db.QueryRow(
		"SELECT COUNT(*) FROM users WHERE created_at >= ?",
		monthStart.Format("2006-01-02 15:04:05"),
	).Scan(&stats.UsersRegisteredMonth)
	if err != nil {
		log.Printf("Error counting new users: %v", err)
	}

	// Sales by month (last 12 months)
	salesByMonthRows, err := db.Query(`
		SELECT
			strftime('%Y-%m', created_at) as month,
			COUNT(*) as count,
			COALESCE(SUM(price_paid), 0) as revenue
		FROM purchases
		WHERE status = 'active'
			AND created_at >= date('now', '-12 months')
		GROUP BY strftime('%Y-%m', created_at)
		ORDER BY month ASC
	`)
	if err == nil {
		defer salesByMonthRows.Close()
		for salesByMonthRows.Next() {
			var entry SalesByMonthEntry
			var revenueKopecks int64
			err := salesByMonthRows.Scan(&entry.Month, &entry.Count, &revenueKopecks)
			if err != nil {
				log.Printf("Error scanning sales by month: %v", err)
				continue
			}
			entry.Revenue = revenueKopecks / 100
			stats.SalesByMonth = append(stats.SalesByMonth, entry)
		}
	} else {
		log.Printf("Error querying sales by month: %v", err)
	}

	// Popular categories (by product count)
	catRows, err := db.Query(`
		SELECT pc.name_ru, COUNT(p.id) as cnt
		FROM product_categories pc
		LEFT JOIN products p ON p.category_id = pc.id
		GROUP BY pc.id, pc.name_ru
		ORDER BY cnt DESC
		LIMIT 5
	`)
	if err == nil {
		defer catRows.Close()
		for catRows.Next() {
			var entry PopularCategoryEntry
			err := catRows.Scan(&entry.Name, &entry.Count)
			if err != nil {
				log.Printf("Error scanning popular categories: %v", err)
				continue
			}
			stats.PopularCategories = append(stats.PopularCategories, entry)
		}
	} else {
		log.Printf("Error querying popular categories: %v", err)
	}

	c.JSON(http.StatusOK, stats)
}

// AdminGetRecentActivity returns recent activity across all entities
func AdminGetRecentActivity(c *gin.Context) {
	type activityRow struct {
		id        int64
		itemType  string
		text      string
		createdAt string
	}

	var activities []RecentActivityItem

	// Recent works (gallery)
	workRows, err := db.Query(`
		SELECT id, 'gallery' as type,
			CASE WHEN title_ru != '' THEN title_ru ELSE title_en END as text,
			created_at
		FROM works
		ORDER BY created_at DESC LIMIT 5
	`)
	if err == nil {
		defer workRows.Close()
		for workRows.Next() {
			var row activityRow
			err := workRows.Scan(&row.id, &row.itemType, &row.text, &row.createdAt)
			if err != nil {
				log.Printf("Error scanning work activity: %v", err)
				continue
			}
			activities = append(activities, RecentActivityItem{
				ID:        formatActivityID("gallery", row.id),
				Type:      "gallery",
				Text:      "Добавлена работа: " + row.text,
				Time:      formatTimeAgo(row.createdAt),
				CreatedAt: row.createdAt,
			})
		}
	} else {
		log.Printf("Error querying recent works: %v", err)
	}

	// Recent shop items
	saleRows, err := db.Query(`
		SELECT id, 'shop' as type,
			CASE WHEN title_ru != '' THEN title_ru ELSE title_en END as text,
			created_at
		FROM sales
		ORDER BY created_at DESC LIMIT 5
	`)
	if err == nil {
		defer saleRows.Close()
		for saleRows.Next() {
			var row activityRow
			err := saleRows.Scan(&row.id, &row.itemType, &row.text, &row.createdAt)
			if err != nil {
				log.Printf("Error scanning sale activity: %v", err)
				continue
			}
			activities = append(activities, RecentActivityItem{
				ID:        formatActivityID("shop", row.id),
				Type:      "shop",
				Text:      "Добавлен товар: " + row.text,
				Time:      formatTimeAgo(row.createdAt),
				CreatedAt: row.createdAt,
			})
		}
	} else {
		log.Printf("Error querying recent sales: %v", err)
	}

	// Recent news
	newsRows, err := db.Query(`
		SELECT id, 'news' as type,
			CASE WHEN title_ru != '' THEN title_ru ELSE title_en END as text,
			created_at
		FROM news
		ORDER BY created_at DESC LIMIT 5
	`)
	if err == nil {
		defer newsRows.Close()
		for newsRows.Next() {
			var row activityRow
			err := newsRows.Scan(&row.id, &row.itemType, &row.text, &row.createdAt)
			if err != nil {
				log.Printf("Error scanning news activity: %v", err)
				continue
			}
			activities = append(activities, RecentActivityItem{
				ID:        formatActivityID("news", row.id),
				Type:      "news",
				Text:      "Опубликована новость: " + row.text,
				Time:      formatTimeAgo(row.createdAt),
				CreatedAt: row.createdAt,
			})
		}
	} else {
		log.Printf("Error querying recent news: %v", err)
	}

	// Recent reviews
	reviewRows, err := db.Query(`
		SELECT r.id, 'review' as type,
			COALESCE(u.username, 'пользователь') as username,
			SUBSTR(COALESCE(r.comment_ru, r.comment_en, ''), 1, 80) as comment,
			r.created_at
		FROM reviews r
		JOIN users u ON r.user_id = u.id
		ORDER BY r.created_at DESC LIMIT 5
	`)
	if err == nil {
		defer reviewRows.Close()
		for reviewRows.Next() {
			var id int64
			var itemType, username, comment, createdAt string
			err := reviewRows.Scan(&id, &itemType, &username, &comment, &createdAt)
			if err != nil {
				log.Printf("Error scanning review activity: %v", err)
				continue
			}
			text := "Новый отзыв от " + username
			if comment != "" {
				text += ": \"" + comment + "\""
			}
			activities = append(activities, RecentActivityItem{
				ID:        formatActivityID("review", id),
				Type:      "review",
				Text:      text,
				Time:      formatTimeAgo(createdAt),
				CreatedAt: createdAt,
			})
		}
	} else {
		log.Printf("Error querying recent reviews: %v", err)
	}

	// Recent purchases
	purchaseRows, err := db.Query(`
		SELECT p.id, 'purchase' as type,
			COALESCE(u.username, 'пользователь') as username,
			COALESCE(pr.title_ru, pr.title_en, 'продукт') as product_name,
			p.created_at
		FROM purchases p
		JOIN users u ON p.user_id = u.id
		JOIN products pr ON p.product_id = pr.id
		ORDER BY p.created_at DESC LIMIT 5
	`)
	if err == nil {
		defer purchaseRows.Close()
		for purchaseRows.Next() {
			var id int64
			var itemType, username, productName, createdAt string
			err := purchaseRows.Scan(&id, &itemType, &username, &productName, &createdAt)
			if err != nil {
				log.Printf("Error scanning purchase activity: %v", err)
				continue
			}
			activities = append(activities, RecentActivityItem{
				ID:        formatActivityID("purchase", id),
				Type:      "purchase",
				Text:      "Покупка: " + username + " приобрёл(а) \"" + productName + "\"",
				Time:      formatTimeAgo(createdAt),
				CreatedAt: createdAt,
			})
		}
	} else {
		log.Printf("Error querying recent purchases: %v", err)
	}

	// Recent user registrations
	userRows, err := db.Query(`
		SELECT id, 'user' as type, username, created_at
		FROM users
		ORDER BY created_at DESC LIMIT 5
	`)
	if err == nil {
		defer userRows.Close()
		for userRows.Next() {
			var id int64
			var itemType, username, createdAt string
			err := userRows.Scan(&id, &itemType, &username, &createdAt)
			if err != nil {
				log.Printf("Error scanning user activity: %v", err)
				continue
			}
			activities = append(activities, RecentActivityItem{
				ID:        formatActivityID("user", id),
				Type:      "user",
				Text:      "Зарегистрирован новый пользователь: " + username,
				Time:      formatTimeAgo(createdAt),
				CreatedAt: createdAt,
			})
		}
	} else {
		log.Printf("Error querying recent users: %v", err)
	}

	// Sort all activities by created_at descending and take top 15
	sortActivitiesByTime(activities)
	if len(activities) > 15 {
		activities = activities[:15]
	}

	c.JSON(http.StatusOK, activities)
}

func formatActivityID(prefix string, id int64) string {
	return prefix + "-" + fmt.Sprintf("%d", id)
}

func formatTimeAgo(dateStr string) string {
	date, err := time.Parse("2006-01-02 15:04:05", dateStr)
	if err != nil {
		// Try other formats
		date, err = time.Parse(time.RFC3339, dateStr)
		if err != nil {
			return dateStr
		}
	}
	now := time.Now()
	diff := now.Sub(date)

	if diff < time.Minute {
		return "только что"
	}
	if diff < time.Hour {
		mins := int(diff.Minutes())
		return fmt.Sprintf("%d мин. назад", mins)
	}
	if diff < 24*time.Hour {
		hours := int(diff.Hours())
		return fmt.Sprintf("%d ч. назад", hours)
	}
	if diff < 7*24*time.Hour {
		days := int(diff.Hours() / 24)
		return fmt.Sprintf("%d дн. назад", days)
	}
	return date.Format("02.01.2006")
}

func sortActivitiesByTime(activities []RecentActivityItem) {
	for i := 0; i < len(activities); i++ {
		for j := i + 1; j < len(activities); j++ {
			if activities[i].CreatedAt < activities[j].CreatedAt {
				activities[i], activities[j] = activities[j], activities[i]
			}
		}
	}
}
