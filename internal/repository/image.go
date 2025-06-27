package repository

import (
	"context"

	"imageBot/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type imageRepository struct {
	db *pgxpool.Pool
}

func NewImageRepository(db *pgxpool.Pool) *imageRepository {
	return &imageRepository{
		db: db,
	}
}

func (r *imageRepository) SaveImage(image *model.Image) error {
	query := `
			INSERT INTO 
			    images (id,text, content, prompt) 
			VALUES ($1,$2,$3,$4)
			`
	_, err := r.db.Exec(context.Background(), query, image.ID, image.Text, image.Content, image.Prompt)
	if err != nil {
		return err

	}
	return nil
}
func (r *imageRepository) SaveImageMessage(message_id int, image_id string) error {
	query := `
			INSERT INTO 
				image_messages(message_id, image_id)
			VALUES ($1,$2)
			`
	_, err := r.db.Exec(context.Background(), query, message_id, image_id)
	if err != nil {
		return err
	}
	return nil

}

func (r *imageRepository) GetImage(delta int) ([]model.Image, error) {
	query := `
			SELECT id, message_id FROM images
			          JOIN public.image_messages im on images.id = im.image_id
			WHERE created_at > NOW()- ($1 *24* interval '1  hours') 
				`
	if delta > 365 {
		query += " AND status=='month_top'"
	} else if delta > 30 {
		query += " AND status=='week_top'"
	}
	rows, err := r.db.Query(context.Background(), query, delta)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var images []model.Image
	for rows.Next() {
		var image model.Image
		if err := rows.Scan(&image.ID, &image.MessageId); err != nil {
			return nil, err
		}
		images = append(images, image)
	}
	return images, nil
}
