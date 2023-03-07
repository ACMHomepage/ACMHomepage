package crawler

import (
	"context"
	"net/url"
	"strconv"

	"github.com/skogkatt-org/ACMHomepage/backend/storage"
)

func UpdateOjUser(db storage.DB, ctx context.Context, ojUser *storage.OjUser) {
	if ojUser.OjName == "codeforces" {
		res := apiPost("https://codeforces.com/api/user.info", url.Values{"handles": {ojUser.Handle}})
		if res["status"] == "OK" {
			info := res["result"].([]interface{})[0].(map[string]interface{})
			ojUser.CurrentRating = int64(info["rating"].(float64))
			ojUser.MaxRating = int64(info["maxRating"].(float64))
		}

		// Update user's codeforces submissions.
		UpdateCfSubmission(db, ctx, ojUser)

		ojUser.AcceptCount, _ = db.GetAcceptCount(ctx, ojUser.ID)
		ojUser.Link = "https://codeforces.com/profile/" + ojUser.Handle
		db.UpdateOjUser(ctx, ojUser)
	}
}

func UpdateCfSubmission(db storage.DB, ctx context.Context, ojUser *storage.OjUser) {
	res := apiPost("https://codeforces.com/api/user.status", url.Values{"handle": {ojUser.Handle}})
	if res["status"] == "OK" {
		result := res["result"].([]interface{})
		for _, _info := range result {
			info := _info.(map[string]interface{})
			problemInfo := info["problem"].(map[string]interface{})
			if problemInfo["contestId"] == nil {
				// NOTE: Currently only problems in official problemset are processed.
				//       Gym and SGU problemset have not been stored in our db.
				continue
			}
			contestId := problemInfo["contestId"].(float64)
			submissionId := info["id"].(float64)
			submissionLink := "https://codeforces.com/contest/" + strconv.FormatFloat(contestId, 'f', 0, 64) + "/submission/" + strconv.FormatFloat(submissionId, 'f', 0, 64)
			var submission storage.Submission
			if err := db.GetSubmissionByLink(ctx, &submission, submissionLink); err != nil {
				submission.OjUserID = ojUser.ID
				submission.Link = submissionLink
				contestIndex := problemInfo["index"].(string)
				problemLink := "https://codeforces.com/contest/" + strconv.FormatFloat(contestId, 'f', 0, 64) + "/problem/" + contestIndex
				var problem storage.Problem
				if err := db.GetProblemByLink(ctx, &problem, problemLink); err != nil {
					problem.OjName = ojUser.OjName
					problem.Name = problemInfo["name"].(string)
					problem.Link = problemLink
					if problemInfo["rating"] != nil {
						problem.Rating = int64(problemInfo["rating"].(float64))
					}
					if err := db.CreateProblem(ctx, &problem); err != nil {
						panic(err)
					}
					tags := problemInfo["tags"].([]interface{})
					for _, _tagName := range tags {
						tagName := _tagName.(string)
						var tag storage.Tag
						if err := db.GetTag(ctx, &tag, tagName); err != nil {
							tag.Name = tagName
							if err := db.CreateTag(ctx, &tag); err != nil {
								panic(err)
							}
						}
						var problemToTag storage.ProblemToTag
						if err := db.GetProblemToTag(ctx, &problemToTag, problem.ID, tag.Name); err != nil {
							problemToTag.ProblemID = problem.ID
							problemToTag.TagName = tag.Name
							if err := db.CreateProblemToTag(ctx, &problemToTag); err != nil {
								panic(err)
							}
						}
					}
				}
				submission.ProblemID = problem.ID
				verdictName := info["verdict"].(string)
				var verdict storage.Verdict
				if err := db.GetVerdict(ctx, &verdict, verdictName); err != nil {
					verdict.Name = verdictName
					if err := db.CreateVerdict(ctx, &verdict); err != nil {
						panic(err)
					}
				}
				submission.VerdictName = verdict.Name
				submission.SubmitTime = info["creationTimeSeconds"].(float64)

				if err := db.CreateSubmission(ctx, &submission); err != nil {
					panic(err)
				}
			} else {
				break
			}
		}
	}
}
