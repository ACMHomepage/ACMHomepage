package crawler

import (
	"context"
	"net/url"
	"strconv"

	"github.com/skogkatt-org/ACMHomepage/backend/storage"
)

func UpdateOjUser(db storage.DB, ctx context.Context, ojUser *storage.OjUser) {
	oj := new(storage.Oj)
	db.GetOjByID(ctx, oj)
	if oj.OjName == "codeforces" {
		// 获取rating
		res := apiPost("https://codeforces.com/api/user.info", url.Values{"handles": {ojUser.Handle}})
		if res["status"] == "OK" {
			info := res["result"].([]interface{})[0].(map[string]interface{})
			ojUser.CurrentRating = int64(info["rating"].(float64))
			ojUser.MaxRating = int64(info["maxRating"].(float64))
		}

		// 更新submission
		UpdateCfSubmission(db, ctx, ojUser.ID, ojUser.Handle)

		// 获取AcceptCount
		db.GetAcceptCount(ctx, ojUser)

		ojUser.Link = "https://codeforces.com/profile/" + ojUser.Handle

		db.SetOjUser(ctx, ojUser)
	}
}

func UpdateCfSubmission(db storage.DB, ctx context.Context, ojUserID int64, handle string) {
	res := apiPost("https://codeforces.com/api/user.status", url.Values{"handle": {handle}})
	if res["status"] == "OK" {
		result := res["result"].([]map[string]interface{})
		for _, info := range result {
			submission := new(storage.Submission)
			submission.OjUserID = ojUserID

			problemInfo := info["problem"].(map[string]interface{})
			contestId := problemInfo["contestId"].(int64)
			contestIndex := problemInfo["index"].(string)

			// 通过link查询ProblemID，没有则加入数据库
			problem := new(storage.Problem)
			problem.Link = "https://codeforces.com/contest/" + strconv.FormatInt(contestId, 10) + "/problem/" + contestIndex
			exists, _ := db.GetProblemByLink(ctx, problem)
			if !exists {
				problem.Rating = int64(problemInfo["points"].(float64))
				problem.ProblemName = problemInfo["name"].(string)
				// TODO
				// 处理tag problem_tag
			}
			submission.ProblemID = problem.ID

			verdict := new(storage.Verdict)
			verdict.VerdictName = info["verdict"].(string)
			db.GetVerdictID(ctx, verdict)
			submission.VerdictID = verdict.ID

			submission.SubmitTime = info["creationTimeSeconds"].(int64) // unix-format time

			submissionId := info["id"].(int64)
			submission.Link = "https://codeforces.com/contest/" + strconv.FormatInt(contestId, 10) + "/submission/" + strconv.FormatInt(submissionId, 10)
		}
	}
}
