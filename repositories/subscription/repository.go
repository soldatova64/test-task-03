package subscription

import (
	sql "github.com/jmoiron/sqlx"
	"log"
	"test-task-03/entity"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetList() ([]entity.Subscription, error) {
	rows, err := r.db.Query(`
		SELECT id, service_name, price, user_id, start_date, end_date, is_deleted 
		FROM subscription 
		WHERE is_deleted = FALSE`,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []entity.Subscription

	for rows.Next() {
		var s entity.Subscription
		err = rows.Scan(&s.ID, &s.ServiceName, &s.Price, &s.UserID, &s.StartDate, &s.EndDate, &s.IsDeleted)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, s)
	}

	return subscriptions, nil
}

func (r *Repository) Get(id int) (entity.Subscription, error) {
	var s entity.Subscription

	query := `
		SELECT id, service_name, price, user_id, start_date, end_date, is_deleted 
		FROM subscription 
		WHERE id = :id 
		  AND is_deleted = FALSE`
	query, params, err := sql.Named(query, map[string]interface{}{
		"id": id,
	})

	query = r.db.Rebind(query)

	err = r.db.Get(&s, query, params...)

	if err != nil {
		log.Println(err)
		return s, err
	}

	return s, nil
}

func (r *Repository) Post(subscription entity.Subscription) (entity.Subscription, error) {
	query := `
		INSERT INTO subscription (service_name, price, user_id, start_date, end_date) 
		VALUES (:service_name, :price, :user_id, :start_date, :end_date) RETURNING id`

	query, params, err := sql.Named(query, map[string]interface{}{
		"service_name": subscription.ServiceName,
		"price":        subscription.Price,
		"user_id":      subscription.UserID,
		"start_date":   subscription.StartDate,
		"end_date":     subscription.EndDate,
	})

	query = r.db.Rebind(query)

	err = r.db.QueryRowx(query, params...).Scan(&subscription.ID)

	if err != nil {
		log.Println(err)
		return subscription, err
	}

	return subscription, nil
}

func (r *Repository) Put(subscription entity.Subscription) (entity.Subscription, error) {
	query := `
		UPDATE subscription 
		SET service_name = :service_name,
			price = :price,
			user_id = :user_id,
			start_date = :start_date,
			end_date = :end_date
    	WHERE id = :id 
      	  AND is_deleted = FALSE`

	query, params, err := sql.Named(query, map[string]interface{}{
		"id":           subscription.ID,
		"service_name": subscription.ServiceName,
		"price":        subscription.Price,
		"user_id":      subscription.UserID,
		"start_date":   subscription.StartDate,
		"end_date":     subscription.EndDate,
	})

	query = r.db.Rebind(query)

	_, err = r.db.Queryx(query, params...)

	if err != nil {
		log.Println(err)
		return subscription, err
	}

	return subscription, nil
}

func (r *Repository) Delete(id int) error {
	query := `UPDATE subscription SET is_deleted = TRUE WHERE id = :id`

	query, params, err := sql.Named(query, map[string]interface{}{
		"id": id,
	})

	query = r.db.Rebind(query)

	_, err = r.db.Queryx(query, params...)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (r *Repository) GetSubscriptionSum(userID, serviceName, dateFrom, dateTo string) (int, error) {
	var sum int

	query := `SELECT COALESCE(SUM(price), 0) FROM subscription WHERE is_deleted = FALSE`
	args := map[string]interface{}{}

	if dateFrom != "" {
		query += ` AND end_date >= :date_from`
		args["date_from"] = dateFrom
	}

	if dateTo != "" {
		query += ` AND start_date <= :date_to`
		args["date_to"] = dateTo
	}

	if userID != "" {
		query += ` AND user_id = :user_id`
		args["user_id"] = userID
	}

	if serviceName != "" {
		query += ` AND service_name = :service_name`
		args["service_name"] = serviceName
	}

	query, params, err := sql.Named(query, args)

	query = r.db.Rebind(query)

	err = r.db.Get(&sum, query, params...)

	if err != nil {
		log.Println(err)
		return 0, err
	}

	return sum, nil
}
