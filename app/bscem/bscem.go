package main

import (
	"encoding/csv"
	"fmt"
	"github.com/docopt/docopt-go"
	"github.com/jesand/stats/dist"
	"github.com/jesand/stats/model"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	USAGE = `bscem - Run Expectation Maximization to infer parameters and truth
        for a binary symmetric channel

Usage:
  bscem <prefs> <qrel> <research_task> <topic> [--pair]

Options:
  <prefs>          A CSV file containing channel output
  <qrel>           The QREL containing gold standard assessments
  <research_task>  The research task to assess
  <topic>          The topic to assess
  --pair           Use the BSCPair model
`

	FldTaskStr    = "research_task"
	FldTopicStr   = "topic"
	FldTopicIdStr = "topic_id"
	FldStatusStr  = "assn_status"
	FldWorkerStr  = "worker_id"
	FldDoc1Str    = "left_doc"
	FldDoc2Str    = "right_doc"
	FldVoteStr    = "result"

	ValApproved = "Approved"
	ValLt       = "win"
)

func main() {

	// Parse the command line
	args, _ := docopt.Parse(USAGE, nil, true, "1.0", false)
	var (
		file        = args["<prefs>"].(string)
		qrelPath    = args["<qrel>"].(string)
		taskFilter  = args["<research_task>"].(string)
		topicFilter = args["<topic>"].(string)
		pairModel   = args["--pair"].(bool)
		fields      []string
		rows        [][]string
	)

	// Initialize
	if f, err := os.Open(file); err != nil {
		fmt.Println("Could not open ", file, ":", err)
		return
	} else {
		r := csv.NewReader(f)

		fmt.Println("Loading", file)
		if fields, err = r.Read(); err != nil {
			fmt.Println("Could not read CSV field names - ", err)
			return
		} else if rows, err = r.ReadAll(); err != nil {
			fmt.Println("Could not read CSV contents - ", err)
			return
		}
		f.Close()
	}

	// Find our fields
	var (
		fldTask    int
		fldTopic   int
		fldTopicId int
		fldStatus  int
		fldWorker  int
		fldDoc1    int
		fldDoc2    int
		fldVote    int

		topicIdFilter string
	)
	for i, field := range fields {
		if field == FldTaskStr {
			fldTask = i
		} else if field == FldTopicStr {
			fldTopic = i
		} else if field == FldTopicIdStr {
			fldTopicId = i
		} else if field == FldStatusStr {
			fldStatus = i
		} else if field == FldWorkerStr {
			fldWorker = i
		} else if field == FldDoc1Str {
			fldDoc1 = i
		} else if field == FldDoc2Str {
			fldDoc2 = i
		} else if field == FldVoteStr {
			fldVote = i
		}
	}
	for _, row := range rows {
		if row[fldTopic] == topicFilter {
			topicIdFilter = row[fldTopicId]
			break
		}
	}

	// Load the QREL
	var qrel = make(map[string]int)
	if f, err := os.Open(qrelPath); err != nil {
		fmt.Println("Could not open ", qrelPath, ":", err)
		return
	} else {
		fmt.Println("Loading", qrelPath)
		r := csv.NewReader(f)
		r.Comma = ' '
		r.FieldsPerRecord = 4
		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println("Could not parse", qrelPath, ":", err)
				return
			}
			if record[0] == topicIdFilter {
				qrel[record[2]], _ = strconv.Atoi(record[3])
			}
		}
		f.Close()
		fmt.Println("Found assessments for", len(qrel), "documents")
	}

	// Prepare a model for the data
	fmt.Println("Building model...")
	const initNoise = 1e-3
	var (
		mod1           = model.NewMultipleBSCModel(1, 1)
		mod2           = model.NewMultipleBSCPairModel(1, 1, 1, 1)
		numPos, numNeg int
	)
	for _, row := range rows {
		if row[fldTask] != taskFilter || row[fldTopic] != topicFilter ||
			row[fldStatus] != ValApproved || row[fldVote] != ValLt {
			continue
		}
		var (
			winner, loser = row[fldDoc1], row[fldDoc2]
			small, big    = winner, loser
			isLt          bool
		)
		if big < small {
			small, big = big, small
		}
		if small == winner {
			isLt = true
			numPos++
		} else {
			isLt = false
			numNeg++
		}
		var (
			question = fmt.Sprintf("%s %s", small, big)
			worker   = row[fldWorker]
		)
		if pairModel {
			if !mod2.HasChannel(question, worker) {
				mod2.AddChannel(question, initNoise, worker, initNoise)
			}
			mod2.AddObservation(question, question, worker, isLt)
		} else {
			if !mod1.HasChannel(worker) {
				mod1.AddChannel(worker, initNoise)
			}
			mod1.AddObservation(question, worker, isLt)
		}
	}
	fmt.Println(numPos, "Pos", numNeg, "Neg")

	if pairModel {

		// Define an evaluation method
		eval := func(mod *model.MultipleBSCPairModel, round int, stage string) {
			var correct, total float64
			for doc, value := range mod.Inputs {
				var (
					docs = strings.Split(doc, " ")
					rel0 = qrel[docs[0]]
					rel1 = qrel[docs[1]]
					lt   = dist.BooleanSpace.BoolValue(value.Outcome())
				)
				if rel0 != rel1 {
					total++
					if lt && rel0 == 1 {
						correct++
					}
				}
			}

			if round == 0 {
				fmt.Printf("%s score: %f accuracy: %d/%d = %f\n", stage,
					mod.Score(), int(correct), int(total), correct/total)
			} else {
				fmt.Printf("Round %d %s score: %f accuracy: %d/%d = %f\n", round, stage,
					mod.Score(), int(correct), int(total), correct/total)
			}

			if stage == "beta1" {
				fmt.Printf("Beta1 parameters: alpha=%f beta=%f mean=%f variance=%f\n",
					mod.Noise1Dist.Alpha, mod.Noise1Dist.Beta, mod.Noise1Dist.Mean(), mod.Noise1Dist.Variance())
			} else if stage == "beta2" {
				fmt.Printf("Beta2 parameters: alpha=%f beta=%f mean=%f variance=%f\n",
					mod.Noise2Dist.Alpha, mod.Noise2Dist.Beta, mod.Noise2Dist.Mean(), mod.Noise2Dist.Variance())
			}
		}

		// Run expectation maximization on the model
		// mod2.UpdateBeta1 = false
		// mod2.UpdateBeta2 = false
		mod2.EM(0, 1e-3, eval)
	} else {

		// Define an evaluation method
		eval := func(mod *model.MultipleBSCModel, round int, stage string) {
			var correct, total float64
			for doc, value := range mod.Inputs {
				var (
					docs = strings.Split(doc, " ")
					rel0 = qrel[docs[0]]
					rel1 = qrel[docs[1]]
					lt   = dist.BooleanSpace.BoolValue(value.Outcome())
				)
				if rel0 != rel1 {
					total++
					if lt && rel0 == 1 {
						correct++
					}
				}
			}

			if round == 0 {
				fmt.Printf("%s score: %f accuracy: %d/%d = %f\n", stage,
					mod.Score(), int(correct), int(total), correct/total)
			} else {
				fmt.Printf("Round %d %s score: %f accuracy: %d/%d = %f\n", round, stage,
					mod.Score(), int(correct), int(total), correct/total)
			}
			if stage == "beta" {
				fmt.Printf("Beta parameters: alpha=%f beta=%f mean=%f variance=%f\n",
					mod.NoiseDist.Alpha, mod.NoiseDist.Beta, mod.NoiseDist.Mean(), mod.NoiseDist.Variance())
			}
		}

		// Run expectation maximization on the model
		mod1.EM(0, 1e-3, eval)
	}
}
