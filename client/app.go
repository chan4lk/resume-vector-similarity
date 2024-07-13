package main

import (
	"context"
	"fmt"
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
	embeddings, err := getEmbeddings(job.Description)
	if err != nil {
		return err
	}

	conn, err := pgx.Connect(ctx, "postgres://localhost/jobsdb")
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

func getJobs(cv string) ([]Job, []float64, error) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, "postgres://localhost/jobsdb")
	if err != nil {
		return []Job{}, []float64{}, err
	}
	defer conn.Close(ctx)

	embeddings, err := getEmbeddings(cv)
	if err != nil {
		return []Job{}, []float64{}, err
	}
	rows, err := conn.Query(ctx,
		"SELECT id, title, company, description, "+
			"(1 - (embeddings <=> $1)) as cosine_distance "+
			"FROM jobs ORDER BY cosine_distance DESC LIMIT 5",
		pgvector.NewVector(embeddings))
	if err != nil {
		return []Job{}, []float64{}, err
	}
	defer rows.Close()
	var jobs []Job
	var distances []float64
	for rows.Next() {
		var job Job
		var distance float64
		err = rows.Scan(&job.Id, &job.Title, &job.Company, &job.Description,
			&distance)
		if err != nil {
			return []Job{}, []float64{}, err
		}
		jobs = append(jobs, job)
		distances = append(distances, distance)
	}
	if rows.Err() != nil {
		return []Job{}, []float64{}, err
	}
	return jobs, distances, nil
}

func addJobs() {
	addJob("Software Developer", "MSI GLOBAL PRIVATE LIMITED", "./jobs/dev01.txt")
	addJob("Software Developer", "LMA RECRUITMENT SINGAPORE PTE. LTD.", "./jobs/dev02.txt")
	addJob("Software Developer", "GROCERY LOGISTICS OF SINGAPORE PTE LTD", "./jobs/dev03.txt")
	addJob("Nurse", "MSI GLOBAL PRIVATE LIMITED", "./jobs/nurse01.txt")
}

func main() {
	file, _ := os.Open("./cvs/dev01.txt")
	cv, _ := io.ReadAll(file)
	jobs, distances, _ := getJobs(string(cv))

	for i, job := range jobs {
		fmt.Printf("%s (%s), %.3f\n", job.Title, job.Company, distances[i])
	}
}
