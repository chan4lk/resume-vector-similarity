package main

import (
	"context"
	"io"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/pgvector/pgvector-go"
)

type Job struct {
	Id          int64
	Title       string
	Company     string
	Description string
	Embeddings  pgvector.Vector
}

// add a new job
func addJob(title, company, filepath string) error {
	// open the job file and read the contents
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	job := Job{
		Title:       title,
		Company:     company,
		Description: string(data),
	}

	err = createJob(job)
	if err != nil {
		return err
	}
	return nil
}

func createJob(job Job) error {
	ctx := context.Background()
	embeddings, err := GetEmbeddings(job.Description)
	if err != nil {
		return err
	}

	conn, err := pgx.Connect(ctx, "postgres://localhost/jobs")
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	_, err = conn.Exec(ctx, "INSERT INTO jobs (title, company, description, "+
		"embeddings) VALUES ($1, $2, $3, $4)",
		job.Title, job.Company, job.Description,
		pgvector.NewVector(embeddings))
	if err != nil {
		return err
	}
	return nil
}

func main() {
	addJob("Software Developer", "MSI GLOBAL PRIVATE LIMITED", "./jobs/dev01.txt")
	addJob("Software Developer", "LMA RECRUITMENT SINGAPORE PTE. LTD.", "./jobs/dev02.txt")
	addJob("Software Developer", "GROCERY LOGISTICS OF SINGAPORE PTE LTD", "./jobs/dev03.txt")
	addJob("Nurse", "MSI GLOBAL PRIVATE LIMITED", "./jobs/nurse01.txt")
}
