package main

import (
	"fmt"
	"net/http"
	"os"

	"assn.com/db"
	"assn.com/routes"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	client, err := db.ConnectToDB()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	userColl := db.GetCollection("users", client)
	jobsColl := db.GetCollection("jobs", client)

	signup := routes.SignupHandle{Collection: userColl}
	http.Handle("/signup", &signup)

	login := routes.LoginHandle{Collection: userColl}
	http.Handle("/login", &login)

	upload := routes.UploadResume{Collection: userColl}
	http.Handle("/uploadResume", &upload)

	createJob := routes.JobCreateHandle{Collection: jobsColl}
	http.Handle("/admin/job", &createJob)

	getJobById := routes.GetJobByIdHandle{Collection: jobsColl}
	http.Handle("/admin/job/{id}", &getJobById)

	getJobs := routes.GetJobsHandle{Collection: jobsColl}
	http.Handle("/jobs", &getJobs)

	getApplicants := routes.GetApplicants{Collection: userColl}
	http.Handle("/admin/applicants", &getApplicants)

	getApplicantById := routes.GetApplicantById{Collection: userColl}
	http.Handle("/admin/applicant/{id}", &getApplicantById)

	apply := routes.ApplyToJob{Collection: jobsColl}
	http.Handle("/jobs/apply", &apply)

	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
