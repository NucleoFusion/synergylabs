package main

import (
	"fmt"
	"net/http"

	"assn.com/db"
	"assn.com/routes"
)

func main() {
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

	http.ListenAndServe(":4000", nil)
}
